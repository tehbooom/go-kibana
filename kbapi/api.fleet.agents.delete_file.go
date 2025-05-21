package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetUpdateAgentResponse wraps the response from a FleetUpdateAgent call
type FleetDeleteFileResponse struct {
	StatusCode int
	Body       *FleetDeleteFileResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetDeleteFileResponseBody struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

type FleetDeleteFileRequest struct {
	FileID string
}

// newFleetUpdateAgent returns a function that performs DELETE /api/fleet/agents/files/{fileId} API requests
func (api *API) newFleetDeleteFile() func(context.Context, *FleetDeleteFileRequest, ...RequestOption) (*FleetDeleteFileResponse, error) {
	return func(ctx context.Context, req *FleetDeleteFileRequest, opts ...RequestOption) (*FleetDeleteFileResponse, error) {
		if req.FileID == "" {
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
			newCtx = instrument.Start(ctx, "fleet.agents.delete_file")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/files/%s", req.FileID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.agents.delete_file")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.delete_file", httpReq.Body); reader != nil {
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
		resp := &FleetDeleteFileResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetDeleteFileResponseBody

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
