package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetPackagePoliciesGetResponse wraps the response from a <todo> call
type FleetPackagePoliciesGetResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesGetResponseBody struct {
	Item PackagePolicy `json:"item"`
}

type FleetPackagePoliciesGetRequest struct {
	PackagePolicyId string
	Params          FleetPackagePoliciesGetRequestParams
}

type FleetPackagePoliciesGetRequestParams struct {
	// Values are simplified or legacy
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

// newFleetPackagePoliciesGet returns a function that performs GET /api/fleet/package_policies/{packagePolicyId} API requests
func (api *API) newFleetPackagePoliciesGet() func(context.Context, *FleetPackagePoliciesGetRequest, ...RequestOption) (*FleetPackagePoliciesGetResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesGetRequest, opts ...RequestOption) (*FleetPackagePoliciesGetResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.package_policies.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/package_policies/%s", req.PackagePolicyId)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
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
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.package_policies.get")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.get", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesGetResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
