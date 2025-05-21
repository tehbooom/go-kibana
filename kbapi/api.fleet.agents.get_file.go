package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkGetAgentPolicies wraps the response from a FleetBulkGetAgentPolicies call
type FleetGetAgentFileResponse struct {
	StatusCode int
	Body       *FleetGetAgentFileResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

// TODO: Test this API to ensure it returns items
type FleetGetAgentFileResponseBody struct {
	Items interface{} `json:"items"`
}

type FleetGetAgentFileRequest struct {
	FileID   string
	FileName string
}

// newFleetGetAgentPolicy returns a function that performs get /api/fleet/agent/{agentID} API requests
func (api *API) newFleetGetAgentFile() func(context.Context, *FleetGetAgentFileRequest, ...RequestOption) (*FleetGetAgentFileResponse, error) {
	return func(ctx context.Context, req *FleetGetAgentFileRequest, opts ...RequestOption) (*FleetGetAgentFileResponse, error) {
		if req.FileID == "" || req.FileName == "" {
			return nil, fmt.Errorf("File name or ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.get_file")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/files/%s/%s", req.FileID, req.FileName)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agents.get_file")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.get_file", httpReq.Body); reader != nil {
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
		resp := &FleetGetAgentFileResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetGetAgentFileResponseBody

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
