package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetListAgentPolicies wraps the response from a FleetListAgentPolicies call
type FleetListAgentPoliciesResponse struct {
	StatusCode int
	Body       *GetFleetAgentPoliciesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type GetFleetAgentPoliciesResponseBody struct {
	Items   []AgentPolicy `json:"items"`
	Page    float32       `json:"page"`
	PerPage float32       `json:"perPage"`
	Total   float32       `json:"total"`
}

type FleetAgentPoliciesRequest struct {
	Params FleetAgentPoliciesRequestParams
}

// GetFleetAgentPoliciesRequest defines parameters for GetFleetAgentPolicies.
type FleetAgentPoliciesRequestParams struct {
	Page         *float32 `form:"page,omitempty" json:"page,omitempty"`
	PerPage      *float32 `form:"perPage,omitempty" json:"perPage,omitempty"`
	Kuery        *string  `form:"kuery,omitempty" json:"kuery,omitempty"`
	NoAgentCount *bool    `form:"noAgentCount,omitempty" json:"noAgentCount,omitempty"`
	// Full get full policies with package policies populated
	Full   *bool   `form:"full,omitempty" json:"full,omitempty"`
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

// newFleetAgentListPoliciesFunc returns a function that performs GET /api/fleet/agent_policies API requests
func (api *API) newFleetAgentListPoliciesFunc() func(context.Context, *FleetAgentPoliciesRequest, ...RequestOption) (*FleetListAgentPoliciesResponse, error) {
	return func(ctx context.Context, req *FleetAgentPoliciesRequest, opts ...RequestOption) (*FleetListAgentPoliciesResponse, error) {
		if req == nil {
			req = &FleetAgentPoliciesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_policies"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.FormatFloat(float64(*req.Params.Page), 'f', -1, 32)
		}
		if req.Params.PerPage != nil {
			params["perPage"] = strconv.FormatFloat(float64(*req.Params.PerPage), 'f', -1, 32)
		}
		if req.Params.Kuery != nil {
			params["kuery"] = *req.Params.Kuery
		}
		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
		}
		if req.Params.Full != nil {
			params["full"] = strconv.FormatBool(*req.Params.Full)
		}
		if req.Params.NoAgentCount != nil {
			params["noAgentCount"] = strconv.FormatBool(*req.Params.NoAgentCount)
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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.list", httpReq.Body); reader != nil {
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
		resp := &FleetListAgentPoliciesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result GetFleetAgentPoliciesResponseBody

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

				errorMessage := string(bodyBytes)
				if jsonMap, ok := errorObj.(map[string]interface{}); ok {
					if msg, exists := jsonMap["message"]; exists {
						if msgStr, ok := msg.(string); ok {
							errorMessage = msgStr
						}
					}
				}

				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
