package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetInternalInitiateFleetSetupResponse wraps the response from a fleet.internal.initiate_fleet_setup call
type FleetInternalInitiateFleetSetupResponse struct {
	StatusCode int
	Body       *FleetInternalInitiateFleetSetupResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetInternalInitiateFleetSetupResponseBody struct {
	IsInitialized  bool `json:"isInitialized"`
	NonFatalErrors []struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"nonFatalErrors"`
}

// newFleetInternalInitiateFleetSetup returns a function that performs POST /api/fleet/setup API requests
func (api *API) newFleetInternalInitiateFleetSetup() func(context.Context, ...RequestOption) (*FleetInternalInitiateFleetSetupResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*FleetInternalInitiateFleetSetupResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.internal.initiate_fleet_setup")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/setup"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "fleet.internal.initiate_fleet_setup")
			if reader := instrument.RecordRequestBody(ctx, "fleet.internal.initiate_fleet_setup", httpReq.Body); reader != nil {
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
		resp := &FleetInternalInitiateFleetSetupResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetInternalInitiateFleetSetupResponseBody

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
