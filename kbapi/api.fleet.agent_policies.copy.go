package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetPostAgentPolicyResponse wraps the response from a FleetCopyAgentPolicy call
type FleetCopyAgentPolicyResponse struct {
	StatusCode int
	Body       *FleetCopyAgentPolicyResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetCopyAgentPolicyResponseBody struct {
	Item AgentPolicy `json:"item"`
}

// FleetPostAgentPolicyRequest is the request for newFleetCopyAgentPolicy
type FleetCopyAgentPolicyRequest struct {
	// ID of agent policy
	ID     string
	Params FleetAgentPolicyCopyRequestParams
	Body   FleetAgentPolicyCopyRequestBody
}

type FleetAgentPolicyCopyRequestParams struct {
	// Values are simplified or legacy
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

type FleetAgentPolicyCopyRequestBody struct {
	Description *string `json:"description,omitempty"`
	Name        string  `json:"name"`
}

// newFleetCopyAgentPolicy returns a function that performs POST /api/fleet/agent_policies/{agentPolicyId}/copy API requests
func (api *API) newFleetCopyAgentPolicy() func(context.Context, *FleetCopyAgentPolicyRequest, ...RequestOption) (*FleetCopyAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetCopyAgentPolicyRequest, opts ...RequestOption) (*FleetCopyAgentPolicyResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Required Agent Policy ID and Body is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.copy")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agent_policies/%s/copy", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Format != nil {
			params["format"] = *StrPtr(*req.Params.Format)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.copy")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.copy", httpReq.Body); reader != nil {
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
		resp := &FleetCopyAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetCopyAgentPolicyResponseBody

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
