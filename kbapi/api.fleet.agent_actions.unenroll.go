package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetUnenrollAgentResponse wraps the response from a FleetUnenrollAgent call
type FleetUnenrollAgentResponse struct {
	StatusCode int
	Body       *FleetUnenrollAgentResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetUnenrollAgentResponseBody struct {
}

type FleetUnenrollAgentRequest struct {
	AgentID string
	Body    FleetUnenrollAgentRequestBody
}

type FleetUnenrollAgentRequestBody struct {
	Force  *bool `json:"force,omitempty"`
	Revoke *bool `json:"revoke,omitempty"`
}

// newFleetUnenrollAgent returns a function that performs POST /api/fleet/agent/{agentID}/unenroll API requests
func (api *API) newFleetUnenrollAgent() func(context.Context, *FleetUnenrollAgentRequest, ...RequestOption) (*FleetUnenrollAgentResponse, error) {
	return func(ctx context.Context, req *FleetUnenrollAgentRequest, opts ...RequestOption) (*FleetUnenrollAgentResponse, error) {
		if req.AgentID == "" {
			return nil, fmt.Errorf("Agent ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.unenroll")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/%s/unenroll", req.AgentID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agents.unenroll")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.unenroll", httpReq.Body); reader != nil {
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
		resp := &FleetUnenrollAgentResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetUnenrollAgentResponseBody

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

				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
