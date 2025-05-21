package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMGetPackagesLimitedResponse wraps the response from a FleetEPMGetPackagesLimited  call
type FleetEPMGetPackagesLimitedResponse struct {
	StatusCode int
	Body       *FleetEPMGetPackagesLimitedResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMGetPackagesLimitedResponseBody struct {
	Items []string `json:"items"`
}

// newFleetEPMGetPackagesLimited returns a function that performs GET /api/fleet/epm/packages/limited API requests
func (api *API) newFleetEPMGetPackagesLimited() func(context.Context, ...RequestOption) (*FleetEPMGetPackagesLimitedResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*FleetEPMGetPackagesLimitedResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.get_packages_limited")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/epm/packages/limited"

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
			instrument.BeforeRequest(httpReq, "fleet.epm.get_packages_limited")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_packages_limited", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetPackagesLimitedResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMGetPackagesLimitedResponseBody

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
