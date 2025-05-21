package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// ConnectorsGetResponse wraps the response from a <todo> call
type ConnectorsGetResponse struct {
	StatusCode int
	Body       *ConnectorsGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type ConnectorsGetResponseBody struct {
	// Id The identifier for the connector.
	ID     string                 `json:"id"`
	Config map[string]interface{} `json:"config"`
	// Name  The name of the rule.
	Name string `json:"name"`
	// ConnectorTypeId The connector type identifier.
	ConnectorTypeID string `json:"connector_type_id"`
	// IsDeprecated Indicates whether the connector is deprecated.
	IsDeprecated bool `json:"is_deprecated"`
	// IsMissingSecrets Indicates whether the connector is missing secrets.
	IsMissingSecrets bool `json:"is_missing_secrets"`
	// IsPreconfigured Indicates whether the connector is preconfigured. If true, the `config` and `is_missing_secrets` properties are omitted from the response.
	IsPreconfigured bool `json:"is_preconfigured"`
	// IsSystemActionType Indicates whether the connector is used for system actions.
	IsSystemActionType bool `json:"is_system_action_type"`
}

type ConnectorsGetRequest struct {
	ID string
}

// newConnectorsGet returns a function that performs GET /api/actions/connector/{id} API requests
func (api *API) newConnectorsGet() func(context.Context, *ConnectorsGetRequest, ...RequestOption) (*ConnectorsGetResponse, error) {
	return func(ctx context.Context, req *ConnectorsGetRequest, opts ...RequestOption) (*ConnectorsGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "connectors.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/actions/connector/%s", req.ID)

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
			instrument.BeforeRequest(httpReq, "connectors.get")
			if reader := instrument.RecordRequestBody(ctx, "connectors.get", httpReq.Body); reader != nil {
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
		resp := &ConnectorsGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result ConnectorsGetResponseBody

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
