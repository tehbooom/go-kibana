package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMGetPackageFileResponse wraps the response from a FleetEPMGetPackageFile call
type FleetEPMGetPackageFileResponse struct {
	StatusCode int
	Body       []byte
	Error      interface{}
	RawBody    io.ReadCloser
}

// FleetEPMGetPackageFileRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMGetPackageFileRequest struct {
	PackageName    string
	PackageVersion string
	FilePath       string
}

// newFleetEPMGetPackageFile returns a function that performs GET /api/fleet/epm/packages/{pkgName}/{pkgVersion}/{filePath} API requests
func (api *API) newFleetEPMGetPackageFile() func(context.Context, *FleetEPMGetPackageFileRequest, ...RequestOption) (*FleetEPMGetPackageFileResponse, error) {
	return func(ctx context.Context, req *FleetEPMGetPackageFileRequest, opts ...RequestOption) (*FleetEPMGetPackageFileResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.epm.get_package_file")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		var path string

		path = fmt.Sprintf("/api/fleet/epm/packages/%s/%s/%s", req.PackageName, req.PackageVersion, req.FilePath)

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

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.get_package_file")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_package_file", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetPackageFileResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		if httpResp.StatusCode == 200 {
			resp.Body = bodyBytes
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
