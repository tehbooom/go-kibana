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

// FleetEPMBulkInstallPackagesResponse wraps the response from a FleetEPMBulkInstallPackages  call
type FleetEPMBulkInstallPackagesResponse struct {
	StatusCode int
	Body       *FleetEPMBulkInstallPackagesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMBulkInstallPackagesResponseBody struct {
}

type PostFleetEpmPackagesBulkResponse struct {
	Items []FleetEPMBulkInstallPackagesResponseItems `json:"items"`
}

type FleetEPMBulkInstallPackagesResponseItems struct {
	Name   string `json:"name"`
	Result struct {
		Assets        *[]Package  `json:"assets,omitempty"`
		Error         interface{} `json:"error"`
		InstallSource *string     `json:"installSource,omitempty"`
		InstallType   string      `json:"installType"`
		Status        *string     `json:"status,omitempty"`
	} `json:"result"`
	Version string `json:"version"`
}

// FleetEPMBulkInstallPackagesRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMBulkInstallPackagesRequest struct {
	Body   FleetEPMBulkInstallPackagesRequestBody
	Params FleetEPMBulkInstallPackagesRequestParams
}

type FleetEPMBulkInstallPackagesRequestParams struct {
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

type FleetEPMBulkInstallPackagesRequestBody struct {
	Force    *bool             `json:"force,omitempty"`
	Packages []json.RawMessage `json:"packages"`
}

// PostFleetEpmPackagesBulkJSONBodyPackages1 defines packages for FleetEPMBulkInstallPackagesRequestBody .
type PostFleetEpmPackagesBulkJSONBodyPackages struct {
	Name       string `json:"name"`
	Prerelease *bool  `json:"prerelease,omitempty"`
	Version    string `json:"version"`
}

// AddStringPackage adds a string package to the packages slice
func (j *FleetEPMBulkInstallPackagesRequestBody) AddStringPackage(val string) {
	data, _ := json.Marshal(val)
	j.Packages = append(j.Packages, data)
}

// AddObjectPackage adds an object package to the packages slice
func (j *FleetEPMBulkInstallPackagesRequestBody) AddObjectPackage(val PostFleetEpmPackagesBulkJSONBodyPackages) {
	data, _ := json.Marshal(val)
	j.Packages = append(j.Packages, data)
}

// newFleetEPMBulkInstallPackages returns a function that performs POST /api/fleet/epm/packages/_bulk API requests
func (api *API) newFleetEPMBulkInstallPackages() func(context.Context, *FleetEPMBulkInstallPackagesRequest, ...RequestOption) (*FleetEPMBulkInstallPackagesResponse, error) {
	return func(ctx context.Context, req *FleetEPMBulkInstallPackagesRequest, opts ...RequestOption) (*FleetEPMBulkInstallPackagesResponse, error) {
		if req == nil {
			req = &FleetEPMBulkInstallPackagesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.bulk.install_packages")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}
		// Build query parameters
		params := make(map[string]string)

		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}

		path := "/api/fleet/epm/packages/_bulk"

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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.bulk.install_packages")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.bulk.install_packages", httpReq.Body); reader != nil {
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
		resp := &FleetEPMBulkInstallPackagesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMBulkInstallPackagesResponseBody

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
