package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMCreateCustomIntegrationResponse wraps the response from a FleetEPMCreateCustomIntegration  call
type FleetEPMCreateCustomIntegrationResponse struct {
	StatusCode int
	Body       *FleetEPMCreateCustomIntegrationResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMCreateCustomIntegrationResponseBody struct {
	Meta struct {
		InstallSource string `json:"install_source"`
	} `json:"_meta"`
	Items []FleetEPMCreateCustomIntegrationResponseItems `json:"items"`
}

type FleetEPMCreateCustomIntegrationResponseItems struct {
	ID       string `json:"id"`
	OriginID string `json:"originId"`
	Type     string `json:"type"`
}

// FleetEPMCreateCustomIntegrationRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMCreateCustomIntegrationRequest struct {
	Body FleetEPMCreateCustomIntegrationRequestBody
}

type FleetEPMCreateCustomIntegrationRequestBody struct {
	Datasets []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"datasets"`
	Force           *bool  `json:"force,omitempty"`
	IntegrationName string `json:"integrationName"`
}

// newFleetEPMCreateCustomIntegration returns a function that performs GET /api/fleet/epm/categories API requests
func (api *API) newFleetEPMCreateCustomIntegration() func(context.Context, *FleetEPMCreateCustomIntegrationRequest, ...RequestOption) (*FleetEPMCreateCustomIntegrationResponse, error) {
	return func(ctx context.Context, req *FleetEPMCreateCustomIntegrationRequest, opts ...RequestOption) (*FleetEPMCreateCustomIntegrationResponse, error) {
		if req == nil {
			req = &FleetEPMCreateCustomIntegrationRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.create_custom_integration")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/epm/categories"

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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.create_custom_integration")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.create_custom_integration", httpReq.Body); reader != nil {
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
		resp := &FleetEPMCreateCustomIntegrationResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMCreateCustomIntegrationResponseBody

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
