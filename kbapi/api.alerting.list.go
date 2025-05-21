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

// TODO: Update the call
// AlertingListResponse wraps the response from a <todo> call
type AlertingListResponse struct {
	StatusCode int
	Body       *AlertingListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type AlertingListResponseBody struct {
	Data    []AlertingResponseBase `json:"data"`
	Page    int                    `json:"page"`
	PerPage int                    `json:"per_page"`
	Total   int                    `json:"total"`
}

type AlertingListRequest struct {
	Params AlertingListRequestParams
}

type AlertingListRequestParams struct {
	// PerPage The number of rules to return per page.
	PerPage *int `form:"per_page,omitempty" json:"per_page,omitempty"`
	// Page The page number to return.
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// Search An Elasticsearch simple_query_string query that filters the objects in the response.
	Search *string `form:"search,omitempty" json:"search,omitempty"`
	// DefaultSearchOperator The default operator to use for the simple_query_string.
	DefaultSearchOperator *string `form:"default_search_operator,omitempty" json:"default_search_operator,omitempty"`
	// SearchFields The fields to perform the simple_query_string parsed query against.
	SearchFields *[]string `form:"search_fields,omitempty" json:"search_fields,omitempty"`
	// SortField Determines which field is used to sort the results. The field must exist in the `attributes` key of the response.
	SortField *string `form:"sort_field,omitempty" json:"sort_field,omitempty"`
	// SortOrder Determines the sort order.
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
	// HasReference Filters the rules that have a relation with the reference objects with a specific type and identifier.
	HasReference *struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `form:"has_reference,omitempty" json:"has_reference,omitempty"`
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
	// Filter A KQL string that you filter with an attribute from your saved object. It should look like `savedObjectType.attributes.title: "myTitle"`. However, if you used a direct attribute of a saved object, such as `updatedAt`, you must define your filter, for example, `savedObjectType.updatedAt > 2018-12-22`.
	Filter          *string   `form:"filter,omitempty" json:"filter,omitempty"`
	FilterConsumers *[]string `form:"filter_consumers,omitempty" json:"filter_consumers,omitempty"`
}

// newAlertingList returns a function that performs GET /api/alerting/rules/_find API requests
func (api *API) newAlertingList() func(context.Context, *AlertingListRequest, ...RequestOption) (*AlertingListResponse, error) {
	return func(ctx context.Context, req *AlertingListRequest, opts ...RequestOption) (*AlertingListResponse, error) {
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
			newCtx = instrument.Start(ctx, "alerting.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/alerting/rules/_find"

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
		if req.Params.Search != nil {
			params["search"] = *req.Params.Search
		}
		if req.Params.DefaultSearchOperator != nil {
			params["default_search_operator"] = *req.Params.DefaultSearchOperator
		}
		if req.Params.SearchFields != nil && len(*req.Params.SearchFields) > 0 {
			searchFields := strings.Join(*req.Params.SearchFields, ",")
			params["search_fields"] = searchFields
		}
		if req.Params.FilterConsumers != nil && len(*req.Params.FilterConsumers) > 0 {
			filterConsumers := strings.Join(*req.Params.FilterConsumers, ",")
			params["filter_consumers"] = filterConsumers
		}
		if req.Params.HasReference != nil {
			hasRefJSON, err := json.Marshal(req.Params.HasReference)
			if err != nil {
				return nil, err
			}
			params["has_reference"] = string(hasRefJSON)
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
			instrument.BeforeRequest(httpReq, "alerting.list")
			if reader := instrument.RecordRequestBody(ctx, "alerting.list", httpReq.Body); reader != nil {
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
		resp := &AlertingListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result AlertingListResponseBody

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
