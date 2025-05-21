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
// FleetServerHostUpdateResponse wraps the response from a <todo> call
type FleetServerHostUpdateResponse struct {
	StatusCode int
	Body       *FleetServerHostUpdateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetServerHostUpdateResponseBody struct {
	Item FleetServerHostItem `json:"item"`
}

type FleetServerHostUpdateRequest struct {
	ID   string
	Body FleetServerHostUpdateRequestBody
}

type FleetServerHostUpdateRequestBody struct {
	HostURLs        []string `json:"host_urls"`
	ID              *string  `json:"id"`
	IsDefault       *bool    `json:"is_default,omitempty"`
	IsInternal      *bool    `json:"is_internal,omitempty"`
	IsPreconfigured *bool    `json:"is_preconfigured,omitempty"`
	Name            string   `json:"name"`
	ProxyID         *string  `json:"proxy_id"`
}

// newFleetServerHostUpdate returns a function that performs PUT /api/fleet/fleet_server_hosts/{itemId} API requests
func (api *API) newFleetServerHostUpdate() func(context.Context, *FleetServerHostUpdateRequest, ...RequestOption) (*FleetServerHostUpdateResponse, error) {
	return func(ctx context.Context, req *FleetServerHostUpdateRequest, opts ...RequestOption) (*FleetServerHostUpdateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.server_host.update")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/fleet_server_hosts/%s", req.ID)

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
			instrument.BeforeRequest(httpReq, "fleet.server_host.update")
			if reader := instrument.RecordRequestBody(ctx, "fleet.server_host.update", httpReq.Body); reader != nil {
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
		resp := &FleetServerHostUpdateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetServerHostUpdateResponseBody

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
