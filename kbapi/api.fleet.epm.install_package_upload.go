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

// FleetEPMInstallPackageUploadResponse wraps the response from a FleetEPMInstallPackageUpload  call
type FleetEPMInstallPackageUploadResponse struct {
	StatusCode int
	Body       *FleetEPMInstallPackageUploadResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMInstallPackageUploadResponseBody struct {
	Meta struct {
		InstallSource string `json:"install_source"`
	} `json:"_meta"`
	Items []Package `json:"items"`
}

// FleetEPMInstallPackageUploadRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMInstallPackageUploadRequest struct {
	Package []byte
	Params  FleetEPMInstallPackageUploadRequestParams
}

type FleetEPMInstallPackageUploadRequestParams struct {
	IgnoreMappingUpdateErrors *bool `form:"ignoreMappingUpdateErrors,omitempty" json:"ignoreMappingUpdateErrors,omitempty"`
	SkipDataStreamRollover    *bool `form:"skipDataStreamRollover,omitempty" json:"skipDataStreamRollover,omitempty"`
}

// newFleetEPMInstallPackageUpload returns a function that performs POST /api/fleet/epm/packages API requests
func (api *API) newFleetEPMInstallPackageUpload() func(context.Context, *FleetEPMInstallPackageUploadRequest, ...RequestOption) (*FleetEPMInstallPackageUploadResponse, error) {
	return func(ctx context.Context, req *FleetEPMInstallPackageUploadRequest, opts ...RequestOption) (*FleetEPMInstallPackageUploadResponse, error) {
		if req == nil {
			req = &FleetEPMInstallPackageUploadRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.install_package_upload")
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

		path := "/api/fleet/epm/packages"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(req.Package))
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

		httpReq.Header.Set("Content-Type", "application/gzip; application/zip")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.install_package_upload")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.install_package_upload", httpReq.Body); reader != nil {
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
		resp := &FleetEPMInstallPackageUploadResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMInstallPackageUploadResponseBody

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
