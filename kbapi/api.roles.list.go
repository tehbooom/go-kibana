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
// RolesListResponse wraps the response from a <todo> call
type RolesListResponse struct {
	StatusCode int
	Body       *RolesListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type RolesListResponseBody struct {
	Roles []Role
}

type RolesListRequest struct {
	Params RolesListRequestParams
}

type RolesListRequestParams struct {
	// ReplaceDeprecatedPrivileges If `true` and the response contains any privileges that are associated with deprecated features, they are omitted in favor of details about the appropriate replacement feature privileges.
	ReplaceDeprecatedPrivileges *bool `form:"replaceDeprecatedPrivileges,omitempty" json:"replaceDeprecatedPrivileges,omitempty"`
}

// newRolesList returns a function that performs GET /api/security/role API requests
func (api *API) newRolesList() func(context.Context, *RolesListRequest, ...RequestOption) (*RolesListResponse, error) {
	return func(ctx context.Context, req *RolesListRequest, opts ...RequestOption) (*RolesListResponse, error) {
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
			newCtx = instrument.Start(ctx, "roles.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security/role"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ReplaceDeprecatedPrivileges != nil {
			params["replaceDeprecatedPrivileges"] = strconv.FormatBool(*req.Params.ReplaceDeprecatedPrivileges)
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
			instrument.BeforeRequest(httpReq, "roles.list")
			if reader := instrument.RecordRequestBody(ctx, "roles.list", httpReq.Body); reader != nil {
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
		resp := &RolesListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result RolesListResponseBody

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
