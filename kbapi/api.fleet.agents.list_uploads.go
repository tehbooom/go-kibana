package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetListAgentUploadsResponse wraps the response from a FleetListAgentUploads call
type FleetListAgentUploadsResponse struct {
	StatusCode int
	Body       *FleetListAgentUploadsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetListAgentUploadsResponseBody struct {
	Items []struct {
		ActionId   string  `json:"actionId"`
		CreateTime string  `json:"createTime"`
		Error      *string `json:"error,omitempty"`
		FilePath   string  `json:"filePath"`
		Id         string  `json:"id"`
		Name       string  `json:"name"`
		Status     string  `json:"status"`
	} `json:"item"`
}

type FleetListAgentUploadsRequest struct {
	AgentID string
}

// newFleetListAgentUploads returns a function that performs GET /api/fleet/agent/{agentID}/uploads API requests
func (api *API) newFleetListAgentUploads() func(context.Context, *FleetListAgentUploadsRequest, ...RequestOption) (*FleetListAgentUploadsResponse, error) {
	return func(ctx context.Context, req *FleetListAgentUploadsRequest, opts ...RequestOption) (*FleetListAgentUploadsResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.agents.list_uploads")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/%s/uploads", req.AgentID)

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
			instrument.BeforeRequest(httpReq, "fleet.agents.list_uploads")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.list_uploads", httpReq.Body); reader != nil {
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
		resp := &FleetListAgentUploadsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetListAgentUploadsResponseBody

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
