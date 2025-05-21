package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// AlertingHealthResponse wraps the response from a <todo> call
type AlertingHealthResponse struct {
	StatusCode int
	Body       *AlertingHealthResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type AlertingHealthResponseBody struct {
	// AlertingFrameworkHealth Three substates identify the health of the alerting framework: `decryption_health`, `execution_health`, and `read_health`.
	AlertingFrameworkHealth *struct {
		// DecryptionHealth The timestamp and status of the rule decryption.
		DecryptionHealth *struct {
			Status    *string `json:"status,omitempty"`
			Timestamp *string `json:"timestamp,omitempty"`
		} `json:"decryption_health,omitempty"`
		// ExecutionHealth The timestamp and status of the rule run.
		ExecutionHealth *struct {
			Status    *string `json:"status,omitempty"`
			Timestamp *string `json:"timestamp,omitempty"`
		} `json:"execution_health,omitempty"`
		// ReadHealth The timestamp and status of the rule reading events.
		ReadHealth *struct {
			Status    *string `json:"status,omitempty"`
			Timestamp *string `json:"timestamp,omitempty"`
		} `json:"read_health,omitempty"`
	} `json:"alerting_framework_health,omitempty"`
	// HasPermanentEncryptionKey If `false`, the encrypted saved object plugin does not have a permanent encryption key.
	HasPermanentEncryptionKey *bool `json:"has_permanent_encryption_key,omitempty"`
	// IsSufficientlySecure If `false`, security is enabled but TLS is not.
	IsSufficientlySecure *bool `json:"is_sufficiently_secure,omitempty"`
}

// newAlertingHealth returns a function that performs GET /api/alerting/_health API requests
func (api *API) newAlertingHealth() func(context.Context, ...RequestOption) (*AlertingHealthResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*AlertingHealthResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "alerting.health")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/alerting/_health"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "alerting.health")
			if reader := instrument.RecordRequestBody(ctx, "alerting.health", httpReq.Body); reader != nil {
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
		resp := &AlertingHealthResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result AlertingHealthResponseBody

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
