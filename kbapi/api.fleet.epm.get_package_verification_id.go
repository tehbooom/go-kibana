package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMGetPackageSignatureVerificationIDResponse wraps the response from a FleetEPMGetPackageSignatureVerificationID  call
type FleetEPMGetPackageSignatureVerificationIDResponse struct {
	StatusCode int
	Body       *FleetEPMGetPackageSignatureVerificationIDResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMGetPackageSignatureVerificationIDResponseBody struct {
	ID string `json:"id"`
}

// newFleetEPMGetPackageSignatureVerificationID returns a function that performs GET /api/fleet/epm/verification_key_id API requests
func (api *API) newFleetEPMGetPackageSignatureVerificationID() func(context.Context, ...RequestOption) (*FleetEPMGetPackageSignatureVerificationIDResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*FleetEPMGetPackageSignatureVerificationIDResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.get_package_verification_id")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/epm/verification_key_id"

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
			instrument.BeforeRequest(httpReq, "fleet.epm.get_package_verification_id")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_package_verification_id", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetPackageSignatureVerificationIDResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMGetPackageSignatureVerificationIDResponseBody

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
