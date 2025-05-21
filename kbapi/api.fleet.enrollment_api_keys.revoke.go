package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEnrollmentAPIKeysRevokeResponse  wraps the response from a FleetEnrollmentAPIKeysRevoke call
type FleetEnrollmentAPIKeysRevokeResponse struct {
	StatusCode int
	Body       *FleetEnrollmentAPIKeysRevokeResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEnrollmentAPIKeysRevokeResponseBody struct {
	Action string `json:"action"`
}

// FleetEnrollmentAPIKeysRevokeRequest  defines parameters for GetFleetAgents.
type FleetEnrollmentAPIKeysRevokeRequest struct {
	KeyID string
}

// newFleetEnrollmentAPIKeysRevoke returns a function that performs DELETE /api/fleet/enrollment_api_keys/{keyId} API requests
func (api *API) newFleetEnrollmentAPIKeysRevoke() func(context.Context, *FleetEnrollmentAPIKeysRevokeRequest, ...RequestOption) (*FleetEnrollmentAPIKeysRevokeResponse, error) {
	return func(ctx context.Context, req *FleetEnrollmentAPIKeysRevokeRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysRevokeResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Key ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.enrollment_api_keys.revoke")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/enrollment_api_keys/%s", req.KeyID)

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

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.enrollment_api_keys.revoke")
			if reader := instrument.RecordRequestBody(ctx, "fleet.enrollment_api_keys.revoke", httpReq.Body); reader != nil {
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
		resp := &FleetEnrollmentAPIKeysRevokeResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEnrollmentAPIKeysRevokeResponseBody

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
