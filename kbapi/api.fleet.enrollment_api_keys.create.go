package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEnrollmentAPIKeysCreateResponse  wraps the response from a FleetBulkGetAgentPolicies call
type FleetEnrollmentAPIKeysCreateResponse struct {
	StatusCode int
	Body       *FleetEnrollmentAPIKeysCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEnrollmentAPIKeysCreateResponseBody struct {
	Action string           `json:"action"`
	Item   EnrollmentApiKey `json:"item"`
}

// FleetEnrollmentAPIKeysCreateRequest  defines parameters for GetFleetAgents.
type FleetEnrollmentAPIKeysCreateRequest struct {
	Body FleetEnrollmentAPIKeysCreateRequestBody
}

type FleetEnrollmentAPIKeysCreateRequestBody struct {
	Expiration *string `json:"expiration"`
	Name       *string `json:"name"`
	PolicyID   string  `json:"policy_id"`
}

// newFleetEnrollmentAPIKeysCreate returns a function that performs POST /api/fleet/enrollment_api_keys API requests
func (api *API) newFleetEnrollmentAPIKeysCreate() func(context.Context, *FleetEnrollmentAPIKeysCreateRequest, ...RequestOption) (*FleetEnrollmentAPIKeysCreateResponse, error) {
	return func(ctx context.Context, req *FleetEnrollmentAPIKeysCreateRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysCreateResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("PolicyID is not defined in request")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.enrollment_api_keys.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/enrollment_api_keys"

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

		jsonBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.enrollment_api_keys.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.enrollment_api_keys.create", httpReq.Body); reader != nil {
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
		resp := &FleetEnrollmentAPIKeysCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEnrollmentAPIKeysCreateResponseBody

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
