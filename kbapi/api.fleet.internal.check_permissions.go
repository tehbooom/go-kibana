package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// FleetInternalCheckPermissionsResponse wraps the response from a  call
type FleetInternalCheckPermissionsResponse struct {
	StatusCode int
	Body       *FleetInternalCheckPermissionsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetInternalCheckPermissionsResponseBody struct {
	// Values are MISSING_SECURITY, MISSING_PRIVILEGES, or MISSING_FLEET_SERVER_SETUP_PRIVILEGES.
	Error   *string `json:"error,omitempty"`
	Success bool    `json:"success"`
}

type FleetInternalCheckPermissionsRequest struct {
	Params FleetInternalCheckPermissionsRequestParams
}

type FleetInternalCheckPermissionsRequestParams struct {
	FleetServerSetup *bool `form:"fleetServerSetup,omitempty" json:"fleetServerSetup,omitempty"`
}

// newFleetInternalCheckPermissions returns a function that performs GET /api/fleet/check-permissions API requests
func (api *API) newFleetInternalCheckPermissions() func(context.Context, *FleetInternalCheckPermissionsRequest, ...RequestOption) (*FleetInternalCheckPermissionsResponse, error) {
	return func(ctx context.Context, req *FleetInternalCheckPermissionsRequest, opts ...RequestOption) (*FleetInternalCheckPermissionsResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.internal.check_permissions")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/check-permissions "

		// Build query parameters
		params := make(map[string]string)

		if req.Params.FleetServerSetup != nil {
			params["fleetServerSetup"] = strconv.FormatBool(*req.Params.FleetServerSetup)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.internal.check_permissions")
			if reader := instrument.RecordRequestBody(ctx, "fleet.internal.check_permissions", httpReq.Body); reader != nil {
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
		resp := &FleetInternalCheckPermissionsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetInternalCheckPermissionsResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
