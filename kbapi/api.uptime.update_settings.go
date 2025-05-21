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
// UptimeUpdateSettingsResponse wraps the response from a <todo> call
type UptimeUpdateSettingsResponse struct {
	StatusCode int
	Body       *UptimeUpdateSettingsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type UptimeUpdateSettingsResponseBody struct {
}

type UptimeUpdateSettingsRequest struct {
	Body UptimeUpdateSettingsRequestBody
}

type UptimeUpdateSettingsRequestBody struct {
	// CertAgeThreshold The number of days after a certificate is created to trigger an alert.
	CertAgeThreshold *float32 `json:"certAgeThreshold,omitempty"`

	// CertExpirationThreshold The number of days before a certificate expires to trigger an alert.
	CertExpirationThreshold *float32 `json:"certExpirationThreshold,omitempty"`

	// DefaultConnectors A list of connector IDs to be used as default connectors for new alerts.
	DefaultConnectors *[]string `json:"defaultConnectors,omitempty"`

	// DefaultEmail The default email configuration for new alerts.
	DefaultEmail *UptimeDefaultEmail `json:"defaultEmail,omitempty"`

	// HeartbeatIndices An index pattern string to be used within the Uptime app and alerts to query Heartbeat data.
	HeartbeatIndices *string `json:"heartbeatIndices,omitempty"`
}

type UptimeDefaultEmail struct {
	Bcc *[]string `json:"bcc,omitempty"`
	Cc  *[]string `json:"cc,omitempty"`
	To  *[]string `json:"to,omitempty"`
}

// newUptimeUpdateSettings returns a function that performs PUT /api/uptime/settings API requests
func (api *API) newUptimeUpdateSettings() func(context.Context, *UptimeUpdateSettingsRequest, ...RequestOption) (*UptimeUpdateSettingsResponse, error) {
	return func(ctx context.Context, req *UptimeUpdateSettingsRequest, opts ...RequestOption) (*UptimeUpdateSettingsResponse, error) {
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
			newCtx = instrument.Start(ctx, "uptime.update_settings")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/uptime/settings"

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
			instrument.BeforeRequest(httpReq, "uptime.update_settings")
			if reader := instrument.RecordRequestBody(ctx, "uptime.update_settings", httpReq.Body); reader != nil {
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
		resp := &UptimeUpdateSettingsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result UptimeUpdateSettingsResponseBody

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
