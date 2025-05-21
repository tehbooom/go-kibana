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
// EndpointExceptionsListItemsResponse wraps the response from a <todo> call
type EndpointExceptionsListItemsResponse struct {
	StatusCode int
	Body       *EndpointExceptionsListItemsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type EndpointExceptionsListItemsResponseBody struct {
	Data    []EndpointExceptionsListItem `json:"data"`
	Page    int                          `json:"page"`
	PerPage int                          `json:"per_page"`
	Pit     string                       `json:"pit,omitempty"`
	Total   int                          `json:"total"`
}

type EndpointExceptionsListItemsRequest struct {
	Params EndpointExceptionsListItemsRequestParams
}

type EndpointExceptionsListItemsRequestParams struct {
	// Filter Filters the returned results according to the value of the specified field,
	// using the `<field name>:<field value>` syntax.
	Filter *string `form:"filter,omitempty" json:"filter,omitempty"`
	// Page The page number to return
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// PerPage The number of exception list items to return per page
	PerPage *int `form:"per_page,omitempty" json:"per_page,omitempty"`
	// SortField Determines which field is used to sort the results
	SortField *string `form:"sort_field,omitempty" json:"sort_field,omitempty"`
	// SortOrder Determines the sort order, which can be `desc` or `asc`
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
}

// newEndpointExceptionsListItems returns a function that performs GET /api/endpoint_list/items/_find API requests
func (api *API) newEndpointExceptionsListItems() func(context.Context, *EndpointExceptionsListItemsRequest, ...RequestOption) (*EndpointExceptionsListItemsResponse, error) {
	return func(ctx context.Context, req *EndpointExceptionsListItemsRequest, opts ...RequestOption) (*EndpointExceptionsListItemsResponse, error) {
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
			newCtx = instrument.Start(ctx, "endpoint.exceptions.list_items")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/endpoint_list/items/_find"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Filter != nil {
			params["filter"] = *req.Params.Filter
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
			instrument.BeforeRequest(httpReq, "endpoint.exceptions.list_items")
			if reader := instrument.RecordRequestBody(ctx, "endpoint.exceptions.list_items", httpReq.Body); reader != nil {
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
		resp := &EndpointExceptionsListItemsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result EndpointExceptionsListItemsResponseBody

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
