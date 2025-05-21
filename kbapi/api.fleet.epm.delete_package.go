package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMDeletePackageResponse wraps the response from a FleetEPMDeletePackage  call
type FleetEPMDeletePackageResponse struct {
	StatusCode int
	Body       *FleetEPMDeletePackageResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMDeletePackageResponseBody struct {
	Items []Package `json:"items"`
}

// FleetEPMDeletePackageRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMDeletePackageRequest struct {
	PackageName    string
	PackageVersion *string
	Params         FleetEPMDeletePackageRequestParams
}

type FleetEPMDeletePackageRequestParams struct {
	Force *bool
}

// newFleetEPMDeletePackage returns a function that performs DELETE /api/fleet/epm/packages/{pkgName}/{pkgVersion} API requests
func (api *API) newFleetEPMDeletePackage() func(context.Context, *FleetEPMDeletePackageRequest, ...RequestOption) (*FleetEPMDeletePackageResponse, error) {
	return func(ctx context.Context, req *FleetEPMDeletePackageRequest, opts ...RequestOption) (*FleetEPMDeletePackageResponse, error) {
		if req.PackageName == "" {
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
			newCtx = instrument.Start(ctx, "fleet.epm.delete_package")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		var path string

		if req.PackageVersion != nil {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s/%s", req.PackageName, *req.PackageVersion)
		} else {
			path = fmt.Sprintf("/api/fleet/epm/packages/%s", req.PackageName)
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Force != nil {
			params["force"] = strconv.FormatBool(*req.Params.Force)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.epm.delete_package")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.delete_package", httpReq.Body); reader != nil {
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
		resp := &FleetEPMDeletePackageResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMDeletePackageResponseBody

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
