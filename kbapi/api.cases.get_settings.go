package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: Update the call
// CasesGetSettingsResponse wraps the response from a <todo> call
type CasesGetSettingsResponse struct {
	StatusCode int
	Body       *[]CasesSettingsResponse
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesGetSettingsRequest struct {
	Params CasesGetSettingsRequestParams
}

type CasesGetSettingsRequestParams struct {
	// Owner a filter to limit the response to a specific set of applications.
	// If this parameter is omitted, the response contains information about all the cases that the user has access to read.
	Owner *[]string `form:"owner,omitempty" json:"owner,omitempty"`
}

// newCasesGetSettings returns a function that performs GET /api/cases/configure API requests
func (api *API) newCasesGetSettings() func(context.Context, *CasesGetSettingsRequest, ...RequestOption) (*CasesGetSettingsResponse, error) {
	return func(ctx context.Context, req *CasesGetSettingsRequest, opts ...RequestOption) (*CasesGetSettingsResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.get_settings")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/cases/configure"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Owner != nil {
			params["owner"] = strings.Join(*req.Params.Owner, ",")
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "cases.get_settings")
			if reader := instrument.RecordRequestBody(ctx, "cases.get_settings", httpReq.Body); reader != nil {
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
		resp := &CasesGetSettingsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result []CasesSettingsResponse

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
