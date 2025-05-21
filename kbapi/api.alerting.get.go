package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// AlertingGetResponse wraps the response from a <todo> call
type AlertingGetResponse struct {
	StatusCode int
	Body       *AlertingGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type AlertingGetResponseBody struct {
	AlertingResponseBase
	ActiveSnoozes  []string         `json:"active_snoozes,omitempty"`
	AlertDelay     *AlertDelay      `json:"alert_delay,omitempty"`
	Flapping       *Flapping        `json:"flapping,omitempty"`
	IsSnoozedUntil *string          `json:"is_snoozed_until,omitempty"`
	MappedParams   map[string]any   `json:"mapped_params,omitempty"`
	Monitoring     *Monitoring      `json:"monitoring,omitempty"`
	NotifyWhen     *string          `json:"notify_when,omitempty"`
	SnoozeSchedule []SnoozeSchedule `json:"snooze_schedule,omitempty"`
	ViewInAppURL   *string          `json:"view_in_app_relative_url,omitempty"`
}

type AlertingGetRequest struct {
	ID string
}

// newAlertingGet returns a function that performs GET /api/alerting/rule/{id} API requests
func (api *API) newAlertingGet() func(context.Context, *AlertingGetRequest, ...RequestOption) (*AlertingGetResponse, error) {
	return func(ctx context.Context, req *AlertingGetRequest, opts ...RequestOption) (*AlertingGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "alerting.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/alerting/rule/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

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
			instrument.BeforeRequest(httpReq, "alerting.get")
			if reader := instrument.RecordRequestBody(ctx, "alerting.get", httpReq.Body); reader != nil {
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
		resp := &AlertingGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result AlertingGetResponseBody

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
