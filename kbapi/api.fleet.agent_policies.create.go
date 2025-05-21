package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetCreateAgentPolicy wraps the response from a FleetCreateAgentPolicy call
type FleetCreateAgentPolicyResponse struct {
	StatusCode int
	Header     http.Header
	Body       *FleetCreateAgentPolicyResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetCreateAgentPolicyResponseBody struct {
	Item AgentPolicy `json:"item"`
}

// PostFleetAgentPoliciesJSONBody defines parameters for PostFleetAgentPolicies.
type FleetCreateAgentPolicyRequest struct {
	Body FleetCreateAgentPolicyRequestBody
}

// newFleetCreateAgentPolicyFunc returns a function that performs POST /api/fleet/agent_policies API requests
func (api *API) newFleetCreateAgentPolicyFunc() func(context.Context, *FleetCreateAgentPolicyRequest, ...RequestOption) (*FleetCreateAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetCreateAgentPolicyRequest, opts ...RequestOption) (*FleetCreateAgentPolicyResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Agent Policy not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_policies"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.create", httpReq.Body); reader != nil {
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
		resp := &FleetCreateAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			Header:     httpResp.Header,
			RawBody:    httpResp.Body,
		}

		var result FleetCreateAgentPolicyResponseBody

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
