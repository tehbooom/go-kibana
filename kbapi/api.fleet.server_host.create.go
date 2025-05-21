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
// FleetServerHostCreateResponse wraps the response from a <todo> call
type FleetServerHostCreateResponse struct {
	StatusCode int
	Body       *FleetServerHostCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetServerHostCreateResponseBody struct {
	Item FleetServerHostItem `json:"item"`
}

type FleetServerHostCreateRequest struct {
	Body FleetServerHostCreateRequestBody
}

type FleetServerHostCreateRequestBody struct {
	HostURLs        []string `json:"host_urls"`
	ID              *string  `json:"id"`
	IsDefault       *bool    `json:"is_default,omitempty"`
	IsInternal      *bool    `json:"is_internal,omitempty"`
	IsPreconfigured *bool    `json:"is_preconfigured,omitempty"`
	Name            string   `json:"name"`
	ProxyID         *string  `json:"proxy_id"`
}

// newFleetServerHostCreate returns a function that performs POST /api/fleet/fleet_server_hosts API requests
func (api *API) newFleetServerHostCreate() func(context.Context, *FleetServerHostCreateRequest, ...RequestOption) (*FleetServerHostCreateResponse, error) {
	return func(ctx context.Context, req *FleetServerHostCreateRequest, opts ...RequestOption) (*FleetServerHostCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.server_host.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/fleet_server_hosts"

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
			instrument.BeforeRequest(httpReq, "fleet.server_host.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.server_host.create", httpReq.Body); reader != nil {
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
		resp := &FleetServerHostCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetServerHostCreateResponseBody

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
