package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// FleetPackagePoliciesDeleteResponse wraps the response from a <todo> call
type FleetPackagePoliciesDeleteResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesDeleteResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesDeleteResponseBody struct {
	ID string `json:"id"`
}

type FleetPackagePoliciesDeleteRequest struct {
	PackagePolicyId string
	Params          FleetPackagePoliciesDeleteRequestParams
}

type FleetPackagePoliciesDeleteRequestParams struct {
	Force *bool
}

// newFleetPackagePoliciesDelete returns a function that performs DELETE /api/fleet/package_policies/{packagePolicyId} API requests
func (api *API) newFleetPackagePoliciesDelete() func(context.Context, *FleetPackagePoliciesDeleteRequest, ...RequestOption) (*FleetPackagePoliciesDeleteResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesDeleteRequest, opts ...RequestOption) (*FleetPackagePoliciesDeleteResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.delete")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/package_policies/%s", req.PackagePolicyId)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Force != nil {
			params["force"] = strconv.FormatBool(*req.Params.Force)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.delete")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.delete", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesDeleteResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesDeleteResponseBody

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
