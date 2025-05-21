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
// FleetMessageSigningServiceRotateResponse wraps the response from a <todo> call
type FleetMessageSigningServiceRotateResponse struct {
	StatusCode int
	Body       *FleetMessageSigningServiceRotateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetMessageSigningServiceRotateResponseBody struct {
	Message string `json:"message"`
}

type FleetMessageSigningServiceRotateRequest struct {
	Params FleetMessageSigningServiceRotateRequestParams
}

type FleetMessageSigningServiceRotateRequestParams struct {
	Acknowledge *bool
}

// newFleetMessageSigningServiceRotate returns a function that performs POST /api/fleet/message_signing_service/rotate_key_pair API requests
func (api *API) newFleetMessageSigningServiceRotate() func(context.Context, *FleetMessageSigningServiceRotateRequest, ...RequestOption) (*FleetMessageSigningServiceRotateResponse, error) {
	return func(ctx context.Context, req *FleetMessageSigningServiceRotateRequest, opts ...RequestOption) (*FleetMessageSigningServiceRotateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.message_sisgning_service.rotate")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/message_signing_service/rotate_key_pair "

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Acknowledge != nil {
			params["acknowledge"] = strconv.FormatBool(*req.Params.Acknowledge)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.message_sisgning_service.rotate")
			if reader := instrument.RecordRequestBody(ctx, "fleet.message_sisgning_service.rotate", httpReq.Body); reader != nil {
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
		resp := &FleetMessageSigningServiceRotateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetMessageSigningServiceRotateResponseBody

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
