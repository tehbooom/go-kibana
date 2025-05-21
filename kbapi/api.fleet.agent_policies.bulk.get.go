package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkGetAgentPolicies wraps the response from a FleetBulkGetAgentPolicies call
type FleetBulkGetAgentPoliciesResponse struct {
	StatusCode int
	Body       *FleetBulkGetAgentPoliciesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

// FleetBulkGetAgentPoliciesRequest is the request for newFleetBulkGetAgentPolicies
type FleetBulkGetAgentPoliciesRequest struct {
	Params FleetBulkGetAgentPoliciesRequestParams
	Body   FleetBulkGetAgentPoliciesRequestBody
}

type FleetBulkGetAgentPoliciesRequestParams struct {
	// Values are simplified or legacy
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

type FleetBulkGetAgentPoliciesRequestBody struct {
	// Full get full policies with package policies populated
	Full *bool `json:"full,omitempty"`
	// IDs list of package policy IDs
	IDs           []string `json:"ids"`
	IgnoreMissing *bool    `json:"ignoreMissing,omitempty"`
}

// newFleetBulkGetAgentPolicies returns a function that performs POST /api/fleet/agent_policies/_bulk_get API requests
func (api *API) newFleetBulkGetAgentPolicies() func(context.Context, *FleetBulkGetAgentPoliciesRequest, ...RequestOption) (*FleetBulkGetAgentPoliciesResponse, error) {
	return func(ctx context.Context, req *FleetBulkGetAgentPoliciesRequest, opts ...RequestOption) (*FleetBulkGetAgentPoliciesResponse, error) {
		if req.Body.IDs == nil {
			return nil, fmt.Errorf("Required Agent Policy IDs is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.bulk.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_policies/_bulk_get"

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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.bulk.get")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.bulk.get", httpReq.Body); reader != nil {
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
		resp := &FleetBulkGetAgentPoliciesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBulkGetAgentPoliciesResponseBody

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
