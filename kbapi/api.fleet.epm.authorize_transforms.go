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

// FleetEPMAuthorizeTransformsResponse wraps the response from a FleetEPMAuthorizeTransforms  call
type FleetEPMAuthorizeTransformsResponse struct {
	StatusCode int
	Body       *FleetEPMAuthorizeTransformsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMAuthorizeTransformsResponseBody []TransformResult

type TransformResult struct {
	Success     bool   `json:"success"`
	TransformID string `json:"transformId"`
}

// FleetEPMAuthorizeTransformsRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMAuthorizeTransformsRequest struct {
	PackageName    string
	PackageVersion *string
	Body           FleetEPMAuthorizeTransformsRequestBody
	Params         FleetEPMAuthorizeTransformsRequestParams
}

type FleetEPMAuthorizeTransformsRequestParams struct {
	Prerelease *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
}

type FleetEPMAuthorizeTransformsRequestBody struct {
	Transforms []struct {
		TransformId string `json:"transformId"`
	} `json:"transforms"`
}

// newFleetEPMAuthorizeTransforms returns a function that performs POST /api/fleet/epm/packages/{pkgName}/{pkgVersion}/transforms/authorize API requests
func (api *API) newFleetEPMAuthorizeTransforms() func(context.Context, *FleetEPMAuthorizeTransformsRequest, ...RequestOption) (*FleetEPMAuthorizeTransformsResponse, error) {
	return func(ctx context.Context, req *FleetEPMAuthorizeTransformsRequest, opts ...RequestOption) (*FleetEPMAuthorizeTransformsResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.epm.authorize_transforms")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/epm/packages/%s/%s", req.PackageName, *req.PackageVersion)

		params := make(map[string]string)

		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.authorize_transforms")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.authorize_transforms", httpReq.Body); reader != nil {
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
		resp := &FleetEPMAuthorizeTransformsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMAuthorizeTransformsResponseBody

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
