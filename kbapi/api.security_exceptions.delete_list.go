package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityExceptionsDeleteListResponse wraps the response from a <todo> call
type SecurityExceptionsDeleteListResponse struct {
	StatusCode int
	Body       *SecurityExceptionsList
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityExceptionsDeleteListRequest struct {
	Params SecurityExceptionsDeleteListRequestParams
}

type SecurityExceptionsDeleteListRequestParams struct {
	// ID Exception list's identifier.
	// Either id or list_id must be specified.
	ID *string
	// ListID Human readable exception list string identifier, e.g. trusted-linux-processes.
	// Either id or list_id must be specified.
	ListID *string
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string
}

// newSecurityExceptionsDeleteList returns a function that performs DELETE /api/exception_lists API requests
func (api *API) newSecurityExceptionsDeleteList() func(context.Context, *SecurityExceptionsDeleteListRequest, ...RequestOption) (*SecurityExceptionsDeleteListResponse, error) {
	return func(ctx context.Context, req *SecurityExceptionsDeleteListRequest, opts ...RequestOption) (*SecurityExceptionsDeleteListResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_exceptions.delete_list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/exception_lists"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ID != nil {
			params["id"] = *req.Params.ID
		}
		if req.Params.ListID != nil {
			params["list_id"] = *req.Params.ListID
		}
		if req.Params.NamespaceType != nil {
			params["namespace_type"] = *req.Params.NamespaceType
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_exceptions.delete_list")
			if reader := instrument.RecordRequestBody(ctx, "security_exceptions.delete_list", httpReq.Body); reader != nil {
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
		resp := &SecurityExceptionsDeleteListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityExceptionsList

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
