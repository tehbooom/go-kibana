package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMGetInputsTemplateResponse wraps the response from a FleetEPMGetInputsTemplate  call
type FleetEPMGetInputsTemplateResponse struct {
	StatusCode int
	Body       *FleetEPMGetInputsTemplateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMGetInputsTemplateResponseBody struct {
	Inputs []FleetEPMGetInputsTemplateResponseBodyInput `json:"inputs"`
}

type FleetEPMGetInputsTemplateResponseBodyInput struct {
	ID      string                                        `json:"id"`
	Type    string                                        `json:"type"`
	Streams []FleetEPMGetInputsTemplateResponseBodyStream `json:"streams"`
}

type FleetEPMGetInputsTemplateResponseBodyDataStream struct {
	Dataset string `json:"dataset"`
	Type    string `json:"type"`
}

type FleetEPMGetInputsTemplateResponseBodyStream struct {
	ID         string                                          `json:"id"`
	DataStream FleetEPMGetInputsTemplateResponseBodyDataStream `json:"data_stream"`
}

// FleetEPMGetInputsTemplateRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMGetInputsTemplateRequest struct {
	PackageName    string
	PackageVersion string
	Params         FleetEPMGetInputsTemplateRequestParams
}

type FleetEPMGetInputsTemplateRequestParams struct {
	IgnoreUnverified *bool `form:"ignoreUnverified,omitempty" json:"ignoreUnverified,omitempty"`
	Prerelease       *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
	// Values are json, yml, or yaml. Default value is json.
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

// newFleetEPMGetInputsTemplate returns a function that performs GET /api/fleet/epm/templates/{pkgName}/{pkgVersion}/inputs API requests
func (api *API) newFleetEPMGetInputsTemplate() func(context.Context, *FleetEPMGetInputsTemplateRequest, ...RequestOption) (*FleetEPMGetInputsTemplateResponse, error) {
	return func(ctx context.Context, req *FleetEPMGetInputsTemplateRequest, opts ...RequestOption) (*FleetEPMGetInputsTemplateResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Required package name or version is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.get_inputs_template")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		var path string

		path = fmt.Sprintf("/api/fleet/epm/templates/%s/%s/inputs", req.PackageName, req.PackageVersion)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.IgnoreUnverified != nil {
			params["ignoreUnverified"] = strconv.FormatBool(*req.Params.IgnoreUnverified)
		}
		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}
		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
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
			instrument.BeforeRequest(httpReq, "fleet.epm.get_inputs_template")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_inputs_template", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetInputsTemplateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMGetInputsTemplateResponseBody

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
