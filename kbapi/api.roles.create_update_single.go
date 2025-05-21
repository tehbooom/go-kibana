package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// RolesCreateUpdateSingleRoleResponse wraps the response from a <todo> call
type RolesCreateUpdateSingleRoleResponse struct {
	StatusCode int
	Body       *RolesCreateUpdateSingleRoleResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type RolesCreateUpdateSingleRoleResponseBody struct{}

type RolesCreateUpdateSingleRoleRequest struct {
	Name   string
	Params RolesCreateUpdateSingleRoleRequestParams
	Body   RolesCreateUpdateSingleRoleRequestBody
}

type RolesCreateUpdateSingleRoleRequestParams struct {
	CreateOnly *bool
}

type RolesCreateUpdateSingleRoleRequestBody struct {
	Kibana        []KibanaPermission     `json:"kibana"`
	Metadata      Metadata               `json:"metadata"`
	Description   string                 `json:"description"`
	Elasticsearch ElasticsearchPrivilege `json:"elasticsearch"`
}

// newRolesCreateUpdateSingleRole returns a function that performs PUT /api/security/role/{name} API requests
func (api *API) newRolesCreateUpdateSingleRole() func(context.Context, *RolesCreateUpdateSingleRoleRequest, ...RequestOption) (*RolesCreateUpdateSingleRoleResponse, error) {
	return func(ctx context.Context, req *RolesCreateUpdateSingleRoleRequest, opts ...RequestOption) (*RolesCreateUpdateSingleRoleResponse, error) {
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
			newCtx = instrument.Start(ctx, "roles.create_update_single")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/security/role/%s", req.Name)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.CreateOnly != nil {
			params["createOnly"] = strconv.FormatBool(*req.Params.CreateOnly)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "roles.create_update_single")
			if reader := instrument.RecordRequestBody(ctx, "roles.create_update_single", httpReq.Body); reader != nil {
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
		resp := &RolesCreateUpdateSingleRoleResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result RolesCreateUpdateSingleRoleResponseBody

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
