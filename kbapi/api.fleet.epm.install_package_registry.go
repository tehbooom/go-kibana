package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMInstallPackageRegistryResponse wraps the response from a FleetEPMInstallPackageRegistry  call
type FleetEPMInstallPackageRegistryResponse struct {
	StatusCode int
	Body       *FleetEPMInstallPackageRegistryResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMInstallPackageRegistryResponseBody struct {
	Meta struct {
		InstallSource string `json:"install_source"`
	} `json:"_meta"`
	Items []Package `json:"items"`
}

// FleetEPMInstallPackageRegistryRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMInstallPackageRegistryRequest struct {
	PackageName    string
	PackageVersion *string
	Params         FleetEPMInstallPackageRegistryRequestParams
	Body           FleetEPMInstallPackageRegistryRequestBody
}

type FleetEPMInstallPackageRegistryRequestParams struct {
	Prerelease                *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
	IgnoreMappingUpdateErrors *bool `form:"ignoreMappingUpdateErrors,omitempty" json:"ignoreMappingUpdateErrors,omitempty"`
	SkipDataStreamRollover    *bool `form:"skipDataStreamRollover,omitempty" json:"skipDataStreamRollover,omitempty"`
}

type FleetEPMInstallPackageRegistryRequestBody struct {
	Force             *bool `json:"force,omitempty"`
	IgnoreConstraints *bool `json:"ignore_constraints,omitempty"`
}

// newFleetEPMInstallPackageRegistry returns a function that performs POST /api/fleet/epm/packages/{pkgName}/{pkgVersion} API requests
func (api *API) newFleetEPMInstallPackageRegistry() func(context.Context, *FleetEPMInstallPackageRegistryRequest, ...RequestOption) (*FleetEPMInstallPackageRegistryResponse, error) {
	return func(ctx context.Context, req *FleetEPMInstallPackageRegistryRequest, opts ...RequestOption) (*FleetEPMInstallPackageRegistryResponse, error) {
		if req == nil {
			req = &FleetEPMInstallPackageRegistryRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.install_package_registry")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}
		// Build query parameters
		params := make(map[string]string)

		if req.Params.IgnoreMappingUpdateErrors != nil {
			params["ignoreMappingUpdateErrors"] = strconv.FormatBool(*req.Params.IgnoreMappingUpdateErrors)
		}
		if req.Params.SkipDataStreamRollover != nil {
			params["skipDataStreamRollover"] = strconv.FormatBool(*req.Params.SkipDataStreamRollover)
		}
		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}

		var path string
		if req.PackageVersion != nil {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s/%s", req.PackageName, *req.PackageVersion)
		} else {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s", req.PackageName)
		}

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

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		jsonBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.install_package_registry")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.install_package_registry", httpReq.Body); reader != nil {
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
		resp := &FleetEPMInstallPackageRegistryResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMInstallPackageRegistryResponseBody

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
