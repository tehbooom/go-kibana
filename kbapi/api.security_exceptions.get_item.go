package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityExceptionsGetItemResponse wraps the response from a <todo> call
type SecurityExceptionsGetItemResponse struct {
	StatusCode int
	Body       *SecurityExceptionsItem
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityExceptionsGetItemRequest struct {
	Params SecurityExceptionsGetItemRequestParams
}

type SecurityExceptionsGetItemRequestParams struct {
	// ID Exception list's identifier.
	// Either id or list_id must be specified.
	ID *string
	// ItemID Human readable exception list string identifier, e.g. trusted-linux-processes.
	// Either id or list_id must be specified.
	ItemID *string
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string
}

// newSecurityExceptionsGetItem returns a function that performs GET /api/exception_lists/items API requests
func (api *API) newSecurityExceptionsGetItem() func(context.Context, *SecurityExceptionsGetItemRequest, ...RequestOption) (*SecurityExceptionsGetItemResponse, error) {
	return func(ctx context.Context, req *SecurityExceptionsGetItemRequest, opts ...RequestOption) (*SecurityExceptionsGetItemResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_exceptions.get_item")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/exception_lists/items"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ID != nil {
			params["id"] = *req.Params.ID
		}
		if req.Params.ItemID != nil {
			params["item_id"] = *req.Params.ItemID
		}
		if req.Params.NamespaceType != nil {
			params["namespace_type"] = *req.Params.NamespaceType
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
			instrument.BeforeRequest(httpReq, "security_exceptions.get_item")
			if reader := instrument.RecordRequestBody(ctx, "security_exceptions.get_item", httpReq.Body); reader != nil {
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
		resp := &SecurityExceptionsGetItemResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityExceptionsItem

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
