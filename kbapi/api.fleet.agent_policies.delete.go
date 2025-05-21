package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetDeleteAgentPolicyResponse wraps the response from a FleetDeleteAgentPolicy call
type FleetDeleteAgentPolicyResponse struct {
	StatusCode int
	Body       *FleetDeleteAgentPolicyResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetDeleteAgentPolicyResponseBody struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// FleetDeleteAgentPolicyRequest is the request for newFleetDeleteAgentPolicy
type FleetDeleteAgentPolicyRequest struct {
	Body FleetDeleteAgentPolicyRequestBody
}

type FleetDeleteAgentPolicyRequestBody struct {
	AgentPolicyId string `json:"agentPolicyId"`

	// Force bypass validation checks that can prevent agent policy deletion
	Force *bool `json:"force,omitempty"`
}

// newFleetDeleteAgentPolicy returns a function that performs POST /api/fleet/agent_policies/delete API requests
func (api *API) newFleetDeleteAgentPolicy() func(context.Context, *FleetDeleteAgentPolicyRequest, ...RequestOption) (*FleetDeleteAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetDeleteAgentPolicyRequest, opts ...RequestOption) (*FleetDeleteAgentPolicyResponse, error) {
		if req.Body.AgentPolicyId == "" {
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
			newCtx = instrument.Start(ctx, "fleet.agent_policies.delete")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_policies/delete"

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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.delete")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.delete", httpReq.Body); reader != nil {
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
		resp := &FleetDeleteAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetDeleteAgentPolicyResponseBody

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
