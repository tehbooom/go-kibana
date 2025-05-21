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
// RolesCreateOrUpdateMultiResponse wraps the response from a <todo> call
type RolesCreateOrUpdateMultiResponse struct {
	StatusCode int
	Body       *RolesCreateOrUpdateMultiResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type RolesCreateOrUpdateMultiResponseBody struct{}

type RolesCreateOrUpdateMultiRequest struct {
	Body RolesCreateOrUpdateMultiRequestBody
}

type RolesCreateOrUpdateMultiRequestBody struct {
	Roles map[string]Role `json:"roles"`
}

// newRolesCreateOrUpdateMulti returns a function that performs PUT /api/security/roles API requests
func (api *API) newRolesCreateOrUpdateMulti() func(context.Context, *RolesCreateOrUpdateMultiRequest, ...RequestOption) (*RolesCreateOrUpdateMultiResponse, error) {
	return func(ctx context.Context, req *RolesCreateOrUpdateMultiRequest, opts ...RequestOption) (*RolesCreateOrUpdateMultiResponse, error) {
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
			newCtx = instrument.Start(ctx, "roles.create_update_multi")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security/roles"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "roles.create_update_multi")
			if reader := instrument.RecordRequestBody(ctx, "roles.create_update_multi", httpReq.Body); reader != nil {
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
		resp := &RolesCreateOrUpdateMultiResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result RolesCreateOrUpdateMultiResponseBody

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
