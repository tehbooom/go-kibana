package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetInternalCheckFleetServerHealthResponse wraps the response from a fleet.internal.check_fleet_server_health call
type FleetInternalCheckFleetServerHealthResponse struct {
	StatusCode int
	Body       *FleetInternalCheckFleetServerHealthResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetInternalCheckFleetServerHealthResponseBody struct {
	HostId *string `json:"host_id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Status string  `json:"status"`
}

type FleetInternalCheckFleetServerHealthRequest struct {
	Body FleetInternalCheckFleetServerHealthRequestBody
}

type FleetInternalCheckFleetServerHealthRequestBody struct {
	ID string `json:"id"`
}

// newFleetInternalCheckFleetServerHealth returns a function that performs POST /api/fleet/health_check API requests
func (api *API) newFleetInternalCheckFleetServerHealth() func(context.Context, *FleetInternalCheckFleetServerHealthRequest, ...RequestOption) (*FleetInternalCheckFleetServerHealthResponse, error) {
	return func(ctx context.Context, req *FleetInternalCheckFleetServerHealthRequest, opts ...RequestOption) (*FleetInternalCheckFleetServerHealthResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.internal.check_fleet_server_health")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/health_check "

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

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
			instrument.BeforeRequest(httpReq, "fleet.internal.check_fleet_server_health")
			if reader := instrument.RecordRequestBody(ctx, "fleet.internal.check_fleet_server_health", httpReq.Body); reader != nil {
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
		resp := &FleetInternalCheckFleetServerHealthResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetInternalCheckFleetServerHealthResponseBody

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
