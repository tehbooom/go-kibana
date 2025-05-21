package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMUpdatePackageSettingsResponse wraps the response from a FleetEPMUpdatePackageSettings  call
type FleetEPMUpdatePackageSettingsResponse struct {
	StatusCode int
	Body       *FleetEPMUpdatePackageSettingsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMUpdatePackageSettingsResponseBody struct {
	Item PackageInfo `json:"item"`
}

// FleetEPMUpdatePackageSettingsRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMUpdatePackageSettingsRequest struct {
	PackageName    string
	PackageVersion *string
	Body           FleetEPMUpdatePackageSettingsRequestBody
}

type FleetEPMUpdatePackageSettingsRequestBody struct {
	KeepPoliciesUpToDate bool `json:"keepPoliciesUpToDate"`
}

// newFleetEPMUpdatePackageSettings returns a function that performs PUT /api/fleet/epm/packages/{pkgName}/{pkgVersion} API requests
func (api *API) newFleetEPMUpdatePackageSettings() func(context.Context, *FleetEPMUpdatePackageSettingsRequest, ...RequestOption) (*FleetEPMUpdatePackageSettingsResponse, error) {
	return func(ctx context.Context, req *FleetEPMUpdatePackageSettingsRequest, opts ...RequestOption) (*FleetEPMUpdatePackageSettingsResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.epm.update_package_settings")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		var path string

		if req.PackageVersion != nil {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s/%s", req.PackageName, *req.PackageVersion)
		} else {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s", req.PackageName)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.epm.update_package_settings")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.update_package_settings", httpReq.Body); reader != nil {
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
		resp := &FleetEPMUpdatePackageSettingsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMUpdatePackageSettingsResponseBody

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
