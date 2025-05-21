package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMGetPackageResponse wraps the response from a FleetEPMGetPackage  call
type FleetEPMGetPackageResponse struct {
	StatusCode int
	Body       *FleetEPMGetPackageResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMGetPackageResponseBody struct {
	Item PackageInfo `json:"item"`
}

// FleetEPMGetPackageRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMGetPackageRequest struct {
	PackageName    string
	PackageVersion *string
	Params         FleetEPMGetPackageRequestParams
}

type FleetEPMGetPackageRequestParams struct {
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`
	Prerelease       *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
	Full             *bool `form:"full,omitempty" json:"full,omitempty"`
	WithMetadata     *bool `form:"withMetadata,omitempty" json:"withMetadata,omitempty"`
}

// newFleetEPMGetPackage returns a function that performs GET /api/fleet/epm/packages/{pkgName}/{pkgVersion} API requests
func (api *API) newFleetEPMGetPackage() func(context.Context, *FleetEPMGetPackageRequest, ...RequestOption) (*FleetEPMGetPackageResponse, error) {
	return func(ctx context.Context, req *FleetEPMGetPackageRequest, opts ...RequestOption) (*FleetEPMGetPackageResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Required package name is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.get_package")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		var path string

		if req.PackageVersion != nil {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s/%s", req.PackageName, *req.PackageVersion)
		} else {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s", req.PackageName)
		}
		//
		// Build query parameters
		params := make(map[string]string)

		if req.Params.IgnoreUnverified != nil {
			params["ignoreUnverified"] = strconv.FormatBool(*req.Params.IgnoreUnverified)
		}
		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}
		if req.Params.Full != nil {
			params["full"] = strconv.FormatBool(*req.Params.Full)
		}
		if req.Params.WithMetadata != nil {
			params["withMetadata"] = strconv.FormatBool(*req.Params.WithMetadata)
		}

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
			instrument.BeforeRequest(httpReq, "fleet.epm.get_package")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_package", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetPackageResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMGetPackageResponseBody

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
