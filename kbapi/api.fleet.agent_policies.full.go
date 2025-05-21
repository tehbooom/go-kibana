package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetGetFullAgentPolicyResponse wraps the response from a FleetGetFullAgentPolicy call
type FleetGetFullAgentPolicyResponse struct {
	StatusCode int
	Error      interface{}
	RawBody    io.ReadCloser
	rawJSON    []byte
}

// FleetDownloadAgentPolicyRequest is the request for newFleetDownloadAgentPolicy
type FleetGetFullAgentPolicyRequest struct {
	// ID of agent policy
	ID     string
	Params FleetGetFullAgentPolicyRequestParams
}

type FleetGetFullAgentPolicyRequestParams struct {
	Download   *bool `form:"download,omitempty" json:"download,omitempty"`
	Standalone *bool `form:"standalone,omitempty" json:"standalone,omitempty"`
	Kubernetes *bool `form:"kubernetes,omitempty" json:"kubernetes,omitempty"`
}

// newFleetGetFullAgentPolicy returns a function that performs GET /api/fleet/agent_policies/{agentPolicyId}/full API requests
func (api *API) newFleetGetFullAgentPolicy() func(context.Context, *FleetGetFullAgentPolicyRequest, ...RequestOption) (*FleetGetFullAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetGetFullAgentPolicyRequest, opts ...RequestOption) (*FleetGetFullAgentPolicyResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.agent_policies.full")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agent_policies/%s/full", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Download != nil {
			params["download"] = strconv.FormatBool(*req.Params.Download)
		}
		if req.Params.Standalone != nil {
			params["standalone"] = strconv.FormatBool(*req.Params.Standalone)
		}
		if req.Params.Kubernetes != nil {
			params["kubernetes"] = strconv.FormatBool(*req.Params.Kubernetes)
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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.full")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.full", httpReq.Body); reader != nil {
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
		resp := &FleetGetFullAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		httpResp.Body.Close()

		if httpResp.StatusCode == 200 {
			resp.rawJSON = bodyBytes
			return resp, nil
		} else {
			// For all non-200 responses

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

// AsKubnernetes tries to parse the response as kubernetes yaml
func (f *FleetGetFullAgentPolicyResponse) AsKubnernetes() (string, error) {
	var wrapper struct {
		Item json.RawMessage `json:"item"`
	}

	if err := json.Unmarshal(f.rawJSON, &wrapper); err != nil {
		return "", fmt.Errorf("failed to unmarshal wrapper: %w", err)
	}

	// Try to unmarshal as a string first
	var yamlString string
	err := json.Unmarshal(wrapper.Item, &yamlString)
	if err == nil {
		return yamlString, nil
	}

	// Not a string, so it's not YAML
	return "", fmt.Errorf("response is not in YAML format")
}

// AsJSON tries to parse the response as the JSON structure
// func (f *FleetGetFullAgentPolicyResponse) AsJSON() (*GetFleetAgentPoliciesAgentpolicyidFull200Item1, error) {
// 	var wrapper struct {
// 		Item json.RawMessage `json:"item"`
// 	}
//
// 	if err := json.Unmarshal(f.rawJSON, &wrapper); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal wrapper: %w", err)
// 	}
//
// 	// Try to unmarshal as the JSON structure
// 	var jsonStruct GetFleetAgentPoliciesAgentpolicyidFull200Item1
// 	err := json.Unmarshal(wrapper.Item, &jsonStruct)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal as JSON structure: %w", err)
// 	}
//
// 	return &jsonStruct, nil
// }
