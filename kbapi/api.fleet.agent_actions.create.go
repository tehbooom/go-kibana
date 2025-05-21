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
// FleetAgentActionsCreateResponse wraps the response from a <todo> call
type FleetAgentActionsCreateResponse struct {
	StatusCode int
	Body       *FleetAgentActionsCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetAgentActionsCreateResponseBody struct {
	Item struct {
		Agents    *[]string `json:"agents,omitempty"`
		CreatedAt string    `json:"created_at"`
		Data      *struct {
			LogLevel string `json:"log_level,omitempty"`
		} `json:"data,omitempty"`
		Expiration               *string   `json:"expiration,omitempty"`
		ID                       string    `json:"id"`
		MinimumExecutionDuration *float32  `json:"minimum_execution_duration,omitempty"`
		Namespaces               *[]string `json:"namespaces,omitempty"`
		RolloutDurationSeconds   *float32  `json:"rollout_duration_seconds,omitempty"`
		SentAt                   *string   `json:"sent_at,omitempty"`
		SourceURI                *string   `json:"source_uri,omitempty"`
		StartTime                *string   `json:"start_time,omitempty"`
		Total                    *float32  `json:"total,omitempty"`
		Type                     string    `json:"type"`
	} `json:"item"`
}

type FleetAgentActionsCreateRequest struct {
	ID   string
	Body json.RawMessage
}

type FleetAgentActionsCreateRequestStandardBody struct {
	Action FleetAgentActionsCreateRequestStandardBodyAction `json:"action"`
}

type FleetAgentActionsCreateRequestStandardBodyAction struct {
	// Values are UNENROLL, UPGRADE, or POLICY_REASSIGN.
	Type string `json:"type"`
}

type FleetAgentActionsCreateRequestBodySettings struct {
	Action FleetAgentActionsCreateBodySettingsAction `json:"action"`
}

type FleetAgentActionsCreateBodySettingsAction struct {
	Data FleetAgentActionsCreateBodyData `json:"data"`
	// Value is SETTINGS.
	Type string `json:"type"`
}

type FleetAgentActionsCreateBodyData struct {
	// Values are debug, info, warning, or error.
	LogLevel string `json:"log_level"`
}

// SettingsRequest sets the Body field using FleetAgentActionsCreateRequestBodySettings
func (req *FleetAgentActionsCreateRequest) SettingsRequest(jsonBody *FleetAgentActionsCreateRequestBodySettings) error {
	data, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}
	req.Body = data
	return nil
}

// StandardRequest sets the Body field using FleetAgentActionsCreateRequestStandardBody
func (req *FleetAgentActionsCreateRequest) StandardRequest(jsonBody *FleetAgentActionsCreateRequestStandardBody) error {
	data, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}
	req.Body = data
	return nil
}

// newFleetAgentActionsCreate returns a function that performs POST /api/fleet/agents/{agentId}/actions API requests
func (api *API) newFleetAgentActionsCreate() func(context.Context, *FleetAgentActionsCreateRequest, ...RequestOption) (*FleetAgentActionsCreateResponse, error) {
	return func(ctx context.Context, req *FleetAgentActionsCreateRequest, opts ...RequestOption) (*FleetAgentActionsCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.agent_actions.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/%s/actions", req.ID)

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
			instrument.BeforeRequest(httpReq, "fleet.agent_actions.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_actions.create", httpReq.Body); reader != nil {
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
		resp := &FleetAgentActionsCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetAgentActionsCreateResponseBody

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
