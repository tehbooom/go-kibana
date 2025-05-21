package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetUninstallTokensGetDecryptedResponse wraps the response from a <todo> call
type FleetUninstallTokensGetDecryptedResponse struct {
	StatusCode int
	Body       *FleetUninstallTokensGetDecryptedResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetUninstallTokensGetDecryptedResponseBody struct {
	Item struct {
		CreatedAt  string    `json:"created_at"`
		ID         string    `json:"id"`
		Namespaces *[]string `json:"namespaces,omitempty"`
		PolicyID   string    `json:"policy_id"`
		PolicyName *string   `json:"policy_name"`
		Token      string    `json:"token"`
	} `json:"item"`
}

type FleetUninstallTokensGetDecryptedRequest struct {
	ID string
}

// newFleetUninstallTokensGetDecrypted returns a function that performs GET /api/fleet/uninstall_tokens/{uninstallTokenId} API requests
func (api *API) newFleetUninstallTokensGetDecrypted() func(context.Context, *FleetUninstallTokensGetDecryptedRequest, ...RequestOption) (*FleetUninstallTokensGetDecryptedResponse, error) {
	return func(ctx context.Context, req *FleetUninstallTokensGetDecryptedRequest, opts ...RequestOption) (*FleetUninstallTokensGetDecryptedResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.uninstall_tokens.get_decrypted")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/uninstall_tokens/%s", req.ID)

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
			instrument.BeforeRequest(httpReq, "fleet.uninstall_tokens.get_decrypted")
			if reader := instrument.RecordRequestBody(ctx, "fleet.uninstall_tokens.get_decrypted", httpReq.Body); reader != nil {
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
		resp := &FleetUninstallTokensGetDecryptedResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetUninstallTokensGetDecryptedResponseBody

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
