package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetGetAgentPolicyResponse wraps the response from a FleetGetAgentPolicy call
type FleetGetAgentPolicyResponse struct {
	StatusCode int
	Body       *FleetAgentPoliciesGetAgentPolicyResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetAgentPoliciesGetAgentPolicyResponseBody struct {
	Item AgentPolicy `json:"item"`
}

// FleetGetAgentPolicyRequest is the request for newFleetGetAgentPolicy
type FleetGetAgentPolicyRequest struct {
	// ID of Agent Policy
	ID     string
	Params FleetAgentPolicyRequestParams
}

type FleetAgentPolicyRequestParams struct {
	// Values are simplified or legacy.
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

// newFleetGetAgentPolicy returns a function that performs GET /api/fleet/agent_policies/{agentPolicyId} API requests
func (api *API) newFleetGetAgentPolicy() func(context.Context, *FleetGetAgentPolicyRequest, ...RequestOption) (*FleetGetAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetGetAgentPolicyRequest, opts ...RequestOption) (*FleetGetAgentPolicyResponse, error) {
		if req.ID == "" {
			return nil, fmt.Errorf("Required Agent Policy ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agent_policies/%s", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Format != nil {
			params["format"] = *StrPtr(*req.Params.Format)
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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.get")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.get", httpReq.Body); reader != nil {
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
		resp := &FleetGetAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetAgentPoliciesGetAgentPolicyResponseBody

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
