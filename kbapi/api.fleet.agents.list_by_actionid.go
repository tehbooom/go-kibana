package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkGetAgentPolicies wraps the response from a FleetBulkGetAgentPolicies call
type FleetListAgentsByActionIDResponse struct {
	StatusCode int
	Body       *FleetListAgentsByActionIDResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetListAgentsByActionIDResponseBody struct {
	Items []string `json:"items"`
}

// PostFleetAgentsJSONBody defines parameters for PostFleetAgents.
type FleetListAgentsByActionIDRequest struct {
	ActionIds []string `json:"actionIds"`
}

// newFleetListAgentsByActionID returns a function that performs POST /api/fleet/agents API requests
func (api *API) newFleetListAgentsByActionID() func(context.Context, *FleetListAgentsByActionIDRequest, ...RequestOption) (*FleetListAgentsByActionIDResponse, error) {
	return func(ctx context.Context, req *FleetListAgentsByActionIDRequest, opts ...RequestOption) (*FleetListAgentsByActionIDResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Action IDs not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.list_by_actionid")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents"

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

		jsonBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agents.list_by_actionid")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.list_by_actionid", httpReq.Body); reader != nil {
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
		resp := &FleetListAgentsByActionIDResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetListAgentsByActionIDResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
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
