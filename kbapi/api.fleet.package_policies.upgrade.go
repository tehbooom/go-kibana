package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetPackagePoliciesUpgradeResponse wraps the response from a <todo> call
type FleetPackagePoliciesUpgradeResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesUpgradeResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesUpgradeResponseBody []struct {
	Body struct {
		Message string `json:"message"`
	} `json:"body"`
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	StatusCode float64 `json:"statusCode"`
	Success    bool    `json:"success"`
}

type FleetPackagePoliciesUpgradeRequest struct {
	Body FleetPackagePoliciesUpgradeRequestBody
}

type FleetPackagePoliciesUpgradeRequestBody struct {
	PackagePolicyIDs []string ` json:"packagePolicyIds"`
}

// newFleetPackagePoliciesUpgrade returns a function that performs POST /api/fleet/package_policies/upgrade API requests
func (api *API) newFleetPackagePoliciesUpgrade() func(context.Context, *FleetPackagePoliciesUpgradeRequest, ...RequestOption) (*FleetPackagePoliciesUpgradeResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesUpgradeRequest, opts ...RequestOption) (*FleetPackagePoliciesUpgradeResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.upgrade")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/package_policies/upgrade"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.upgrade")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.upgrade", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesUpgradeResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesUpgradeResponseBody

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
