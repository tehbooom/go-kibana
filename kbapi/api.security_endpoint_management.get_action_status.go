package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityEndpointManagementGetActionStatusResponse wraps the response from a <todo> call
type SecurityEndpointManagementGetActionStatusResponse struct {
	StatusCode int
	Body       *SecurityEndpointManagementGetActionStatusResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityEndpointManagementGetActionStatusResponseBody struct {
	Body struct {
		Data struct {
			AgentID        string `json:"agent_id"`
			PendingActions struct {
				Execute          *int `json:"execute,omitempty"`
				GetFile          *int `json:"get-file,omitempty"`
				Isolate          *int `json:"isolate,omitempty"`
				KillProcess      *int `json:"kill-process,omitempty"`
				RunningProcesses *int `json:"running-processes,omitempty"`
				Scan             *int `json:"scan,omitempty"`
				SuspendProcess   *int `json:"suspend-process,omitempty"`
				Unisolate        *int `json:"unisolate,omitempty"`
				Upload           *int `json:"upload,omitempty"`
			} `json:"pending_actions"`
		} `json:"data"`
	} `json:"body"`
}

type SecurityEndpointManagementGetActionStatusRequest struct {
	Params SecurityEndpointManagementGetActionStatusRequestParams
}

type SecurityEndpointManagementGetActionStatusRequestParams struct {
	Query *SecurityEndpointManagementGetActionStatusRequestQueryParam
}

// newSecurityEndpointManagementGetActionStatus returns a function that performs GET /api/endpoint/action_status API requests
func (api *API) newSecurityEndpointManagementGetActionStatus() func(context.Context, *SecurityEndpointManagementGetActionStatusRequest, ...RequestOption) (*SecurityEndpointManagementGetActionStatusResponse, error) {
	return func(ctx context.Context, req *SecurityEndpointManagementGetActionStatusRequest, opts ...RequestOption) (*SecurityEndpointManagementGetActionStatusResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_endpoint_management.get_action_status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/endpoint/action_status"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Query != nil {
			queryJSON, err := json.Marshal(req.Params.Query)
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			params["query"] = string(queryJSON)
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
			instrument.BeforeRequest(httpReq, "security_endpoint_management.get_action_status")
			if reader := instrument.RecordRequestBody(ctx, "security_endpoint_management.get_action_status", httpReq.Body); reader != nil {
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
		resp := &SecurityEndpointManagementGetActionStatusResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityEndpointManagementGetActionStatusResponseBody

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
