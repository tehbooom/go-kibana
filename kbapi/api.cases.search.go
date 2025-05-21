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
// CasesSearchResponse wraps the response from a <todo> call
type CasesSearchResponse struct {
	StatusCode int
	Body       *CasesSearchResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesSearchResponseBody struct {
	Cases                *[]CasesObjectResponse `json:"cases,omitempty"`
	CountClosedCases     *int                   `json:"count_closed_cases,omitempty"`
	CountInProgressCases *int                   `json:"count_in_progress_cases,omitempty"`
	CountOpenCases       *int                   `json:"count_open_cases,omitempty"`
	Page                 *int                   `json:"page,omitempty"`
	PerPage              *int                   `json:"per_page,omitempty"`
	Total                *int                   `json:"total,omitempty"`
}

type CasesSearchRequest struct {
	Params CasesSearchRequestParams
}

type CasesSearchRequestParams struct {
	// Assignees Filters the returned cases by assignees. Valid values are `none` or unique identifiers for the user profiles.
	// These identifiers can be found by using the suggest user profile API.
	Assignees *[]string `form:"assignees,omitempty" json:"assignees,omitempty"`

	// Category Filters the returned cases by category.
	Category *[]string `form:"category,omitempty" json:"category,omitempty"`

	// DefaultSearchOperator he default operator to use for the simple_query_string.
	DefaultSearchOperator *string `form:"defaultSearchOperator,omitempty" json:"defaultSearchOperator,omitempty"`

	// From Returns only cases that were created after a specific date. The date must be specified as a KQL data range or date match expression.
	From *string `form:"from,omitempty" json:"from,omitempty"`

	// Owner A filter to limit the response to a specific set of applications.
	// If this parameter is omitted, the response contains information about all the cases that the user has access to read.
	Owner *[]string `form:"owner,omitempty" json:"owner,omitempty"`

	// Page The page number to return.
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// PerPage The number of items to return. Limited to 100 items.
	PerPage *int `form:"perPage,omitempty" json:"perPage,omitempty"`

	// Reporters Filters the returned cases by the user name of the reporter.
	Reporters *[]string `form:"reporters,omitempty" json:"reporters,omitempty"`

	// Search An Elasticsearch simple_query_string query that filters the objects in the response.
	Search *string `form:"search,omitempty" json:"search,omitempty"`

	// SearchFields The fields to perform the simple_query_string parsed query against.
	SearchFields *[]string `form:"searchFields,omitempty" json:"searchFields,omitempty"`

	// Severity The severity of the case.
	Severity *string `form:"severity,omitempty" json:"severity,omitempty"`

	// SortField Determines which field is used to sort the results.
	SortField *string `form:"sortField,omitempty" json:"sortField,omitempty"`

	// SortOrder Determines the sort order.
	SortOrder *string `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`

	// Status Filters the returned cases by state.
	Status *string `form:"status,omitempty" json:"status,omitempty"`

	// Tags Filters the returned cases by tags.
	Tags *[]string `form:"tags,omitempty" json:"tags,omitempty"`

	// To Returns only cases that were created before a specific date. The date must be specified as a KQL data range or date match expression.
	To *string `form:"to,omitempty" json:"to,omitempty"`
}

// newCasesSearch returns a function that performs GET /api/cases/_find API requests
func (api *API) newCasesSearch() func(context.Context, *CasesSearchRequest, ...RequestOption) (*CasesSearchResponse, error) {
	return func(ctx context.Context, req *CasesSearchRequest, opts ...RequestOption) (*CasesSearchResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.search")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/cases/_find"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Assignees != nil {
			params["assignees"] = strings.Join(*req.Params.Assignees, ",")
		}
		if req.Params.Category != nil {
			params["category"] = strings.Join(*req.Params.Category, ",")
		}
		if req.Params.DefaultSearchOperator != nil {
			params["defaultSearchOperator"] = *req.Params.DefaultSearchOperator
		}
		if req.Params.From != nil {
			params["from"] = *req.Params.From
		}
		if req.Params.Owner != nil {
			params["owner"] = strings.Join(*req.Params.Owner, ",")
		}
		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PerPage != nil {
			params["perPage"] = strconv.Itoa(*req.Params.PerPage)
		}
		if req.Params.Reporters != nil {
			params["reporters"] = strings.Join(*req.Params.Reporters, ",")
		}
		if req.Params.Search != nil {
			params["search"] = *req.Params.Search
		}
		if req.Params.SearchFields != nil {
			params["searchFields"] = strings.Join(*req.Params.SearchFields, ",")
		}
		if req.Params.Severity != nil {
			params["severity"] = *req.Params.Severity
		}
		if req.Params.SortField != nil {
			params["sortField"] = *req.Params.SortField
		}
		if req.Params.SortOrder != nil {
			params["sortOrder"] = *req.Params.SortOrder
		}
		if req.Params.Status != nil {
			params["status"] = *req.Params.Status
		}
		if req.Params.Tags != nil {
			params["tags"] = strings.Join(*req.Params.Tags, ",")
		}
		if req.Params.To != nil {
			params["to"] = *req.Params.To
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
			instrument.BeforeRequest(httpReq, "cases.search")
			if reader := instrument.RecordRequestBody(ctx, "cases.search", httpReq.Body); reader != nil {
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
		resp := &CasesSearchResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesSearchResponseBody

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
