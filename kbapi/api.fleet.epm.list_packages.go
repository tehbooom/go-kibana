package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMListPackagesResponse wraps the response from a FleetEPMListPackages  call
type FleetEPMListPackagesResponse struct {
	StatusCode int
	Body       *FleetEPMListPackagesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMListPackagesResponseBody struct {
	Items []PackageInfoList `json:"items"`
}

// FleetEPMListPackagesRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMListPackagesRequest struct {
	Params FleetEPMListPackagesRequestParams
}

type FleetEPMListPackagesRequestParams struct {
	Category                 *string `form:"category,omitempty" json:"category,omitempty"`
	Prerelease               *bool   `form:"prerelease,omitempty" json:"prerelease,omitempty"`
	ExcludeInstallStatus     *bool   `form:"excludeInstallStatus,omitempty" json:"excludeInstallStatus,omitempty"`
	WithPackagePoliciesCount *bool   `form:"withPackagePoliciesCount,omitempty" json:"withPackagePoliciesCount,omitempty"`
}

// newFleetEPMListPackages returns a function that performs GET /api/fleet/epm/packages API requests
func (api *API) newFleetEPMListPackages() func(context.Context, *FleetEPMListPackagesRequest, ...RequestOption) (*FleetEPMListPackagesResponse, error) {
	return func(ctx context.Context, req *FleetEPMListPackagesRequest, opts ...RequestOption) (*FleetEPMListPackagesResponse, error) {
		if req == nil {
			req = &FleetEPMListPackagesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.list_packages")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}
		if req.Params.WithPackagePoliciesCount != nil {
			params["withPackagePoliciesCount"] = strconv.FormatBool(*req.Params.WithPackagePoliciesCount)
		}
		if req.Params.ExcludeInstallStatus != nil {
			params["excludeInstallStatus"] = strconv.FormatBool(*req.Params.ExcludeInstallStatus)
		}
		if req.Params.Category != nil {
			params["category"] = *req.Params.Category
		}

		path := "/api/fleet/epm/packages"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.list_packages")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.list_packages", httpReq.Body); reader != nil {
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
		resp := &FleetEPMListPackagesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMListPackagesResponseBody

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
