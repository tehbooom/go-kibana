package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// CasesGetConnectorsResponse wraps the response from a <todo> call
type CasesGetConnectorsResponse struct {
	StatusCode int
	Body       *CasesGetConnectorsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesGetConnectorsResponseBody []struct {
	ID     string `json:"id"`
	Config struct {
		APIURL     string `json:"apiUrl"`
		ProjectKey string `json:"projectKey"`
	} `json:"config"`
	Name string `json:"name"`
	// The type of connector.
	// Values are .cases-webhook, .jira, .none, .resilient, .servicenow, .servicenow-sir, or .swimlane.
	ActionTypeID      string `json:"actionTypeId"`
	IsDeprecated      bool   `json:"isDeprecated"`
	IsMissingSecrets  bool   `json:"isMissingSecrets"`
	IsPreconfigured   bool   `json:"isPreconfigured"`
	ReferencedByCount int    `json:"referencedByCount"`
}

// newCasesGetConnectors returns a function that performs GET /api/cases/configure/connectors/_find API requests
func (api *API) newCasesGetConnectors() func(context.Context, ...RequestOption) (*CasesGetConnectorsResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*CasesGetConnectorsResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "cases.get_connectors")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/cases/configure/connectors/_find"

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
			instrument.BeforeRequest(httpReq, "cases.get_connectors")
			if reader := instrument.RecordRequestBody(ctx, "cases.get_connectors", httpReq.Body); reader != nil {
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
		resp := &CasesGetConnectorsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesGetConnectorsResponseBody

		if httpResp.StatusCode < 299 {
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
			// For all non-success responses
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
