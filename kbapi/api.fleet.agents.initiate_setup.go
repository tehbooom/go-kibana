package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetInitiateSetupResponse  wraps the response from a FleetInitiateSetup call
type FleetInitiateSetupResponse struct {
	StatusCode int
	Body       *FleetInitiateSetupResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetInitiateSetupResponseBody struct {
	IsInitialized  bool `json:"isInitialized"`
	NonFatalErrors []struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"nonFatalErrors"`
}

// FleetInitiateSetupRequest  defines parameters for FleetInitiateSetup.
type FleetInitiateSetupRequest struct {
	Body FleetInitiateSetupRequestBody
}

type FleetInitiateSetupRequestBody struct {
	AdminPassword string `json:"admin_password"`
	AdminUsername string `json:"admin_username"`
}

// newFleetInitiateSetup returns a function that performs POST /api/fleet/agents/setup API requests
func (api *API) newFleetInitiateSetup() func(context.Context, *FleetInitiateSetupRequest, ...RequestOption) (*FleetInitiateSetupResponse, error) {
	return func(ctx context.Context, req *FleetInitiateSetupRequest, opts ...RequestOption) (*FleetInitiateSetupResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Username or password is not set")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.initiate_setup")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/setup"

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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agents.initiate_setup")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.initiate_setup", httpReq.Body); reader != nil {
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
		resp := &FleetInitiateSetupResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetInitiateSetupResponseBody

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
