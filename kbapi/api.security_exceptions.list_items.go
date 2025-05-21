package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// SecurityExceptionsListItemsResponse wraps the response from a ListItems call
type SecurityExceptionsListItemsResponse struct {
	StatusCode int
	Body       *SecurityExceptionsListItemsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityExceptionsListItemsResponseBody struct {
	Data    []SecurityExceptionsItem `json:"data"`
	Page    int                      `json:"page"`
	PerPage int                      `json:"per_page"`
	Total   int                      `json:"total"`
}

type SecurityExceptionsListItemsRequest struct {
	Params SecurityExceptionsListItemsRequestParams
}

type SecurityExceptionsListItemsRequestParams struct {
	// Filter Filters the returned results according to the value of the specified field,
	// using the <field name>:<field value> syntax.
	// Minimum length of each is 1. Default value is [] (empty).
	Filter *string
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string
	// SortField Determines which field is used to sort the results.
	SortField *string
	// SortOrder Values are asc or desc.
	SortOrder *string
	// Page Page number
	// Minimum value is 1. Default value is 1.
	Page *int
	// PerPage AnonymizationFields per page
	// Minimum value is 0. Default value is 20.
	PerPage *int
	Search  *string
	// ListID The list_ids of the items to fetch.
	ListID *[]string
}

// newSecurityExceptionsListItems returns a function that performs GET /api/exception_lists/items/_find API requests
func (api *API) newSecurityExceptionsListItems() func(context.Context, *SecurityExceptionsListItemsRequest, ...RequestOption) (*SecurityExceptionsListItemsResponse, error) {
	return func(ctx context.Context, req *SecurityExceptionsListItemsRequest, opts ...RequestOption) (*SecurityExceptionsListItemsResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_exceptions.list_items")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/exception_lists/items/_find"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Filter != nil {
			params["filter"] = *req.Params.Filter
		}
		if req.Params.Search != nil {
			params["search"] = *req.Params.Search
		}
		if req.Params.ListID != nil {
			params["list_id"] = strings.Join(*req.Params.ListID, ",")
		}
		if req.Params.NamespaceType != nil {
			params["namespace_type"] = *req.Params.NamespaceType
		}
		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PerPage != nil {
			params["per_page"] = strconv.Itoa(*req.Params.PerPage)
		}
		if req.Params.SortField != nil {
			params["sort_field"] = *req.Params.SortField
		}
		if req.Params.SortOrder != nil {
			params["sort_order"] = *req.Params.SortOrder
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
			instrument.BeforeRequest(httpReq, "security_exceptions.list_items")
			if reader := instrument.RecordRequestBody(ctx, "security_exceptions.list_items", httpReq.Body); reader != nil {
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
		resp := &SecurityExceptionsListItemsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityExceptionsListItemsResponseBody

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
