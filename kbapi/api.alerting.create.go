package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// AlertingCreateResponse wraps the response from a <todo> call
type AlertingCreateResponse struct {
	StatusCode int
	Body       *AlertingCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type AlertingCreateResponseBody struct {
	ID                  string           `json:"id"`
	Name                string           `json:"name"`
	Tags                []string         `json:"tags"`
	Params              map[string]any   `json:"params"`
	Actions             []ActionResponse `json:"actions"`
	Enabled             bool             `json:"enabled"`
	Running             bool             `json:"running"`
	Consumer            string           `json:"consumer"`
	LastRun             *LastRun         `json:"last_run,omitempty"`
	MuteAll             bool             `json:"mute_all"`
	NextRun             *string          `json:"next_run,omitempty"`
	Revision            int              `json:"revision"`
	Schedule            Schedule         `json:"schedule"`
	Throttle            *string          `json:"throttle"`
	CreatedAt           string           `json:"created_at"`
	CreatedBy           string           `json:"created_by"`
	UpdatedAt           string           `json:"updated_at"`
	UpdatedBy           string           `json:"updated_by"`
	AlertDelay          *AlertDelay      `json:"alert_delay,omitempty"`
	NotifyWhen          *string          `json:"notify_when"`
	RuleTypeID          string           `json:"rule_type_id"`
	APIKeyOwner         string           `json:"api_key_owner"`
	MutedAlertIDs       []string         `json:"muted_alert_ids"`
	ExecutionStatus     ExecutionStatus  `json:"execution_status"`
	ScheduledTaskID     string           `json:"scheduled_task_id"`
	APIKeyCreatedByUser bool             `json:"api_key_created_by_user"`
}

type AlertingCreateRequest struct {
	ID   string
	Body AlertingCreateRequestBody
}

type AlertingCreateRequestBody struct {
	Name       string         `json:"name"`
	Tags       []string       `json:"tags,omitempty"`
	Params     map[string]any `json:"params"`
	Actions    []ActionCreate `json:"actions,omitempty"`
	Consumer   string         `json:"consumer"`
	Schedule   Schedule       `json:"schedule"`
	AlertDelay *AlertDelay    `json:"alert_delay,omitempty"`
	RuleTypeID string         `json:"rule_type_id"`
}

// newAlertingCreate returns a function that performs POST /api/alerting/rule/{id} API requests
func (api *API) newAlertingCreate() func(context.Context, *AlertingCreateRequest, ...RequestOption) (*AlertingCreateResponse, error) {
	return func(ctx context.Context, req *AlertingCreateRequest, opts ...RequestOption) (*AlertingCreateResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "alerting.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/alerting/rule/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "alerting.create")
			if reader := instrument.RecordRequestBody(ctx, "alerting.create", httpReq.Body); reader != nil {
				httpReq.Body = reader
			}
		}

		// Execute request
		httpResp, err := api.transport.Perform(httpReq)

		if instrument != nil {
			instrument.AfterRequest(httpReq, "kibana", path)
		}

		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		// Prepare response
		resp := &AlertingCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result AlertingCreateResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
