package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetBulkGetAgentPolicies wraps the response from a FleetBulkGetAgentPolicies call
type FleetGetAgentResponse struct {
	StatusCode int
	Body       *FleetGetAgentResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetGetAgentRequest struct {
	AgentID string
	Params  FleetGetAgentRequestParams
}

type FleetGetAgentRequestParams struct {
	WithMetrics *bool `form:"withMetrics,omitempty" json:"withMetrics,omitempty"`
}

// newFleetGetAgent returns a function that performs GET /api/fleet/agent/{agentID} API requests
func (api *API) newFleetGetAgent() func(context.Context, *FleetGetAgentRequest, ...RequestOption) (*FleetGetAgentResponse, error) {
	return func(ctx context.Context, req *FleetGetAgentRequest, opts ...RequestOption) (*FleetGetAgentResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.agents.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agents/%s", req.AgentID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.WithMetrics != nil {
			params["withMetrics"] = strconv.FormatBool(*req.Params.WithMetrics)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agents.get")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.get", httpReq.Body); reader != nil {
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
		resp := &FleetGetAgentResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetGetAgentResponseBody

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
