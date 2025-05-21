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
// SecurityDetectionsListRulesResponse wraps the response from a <todo> call
type SecurityDetectionsListRulesResponse struct {
	StatusCode int
	Body       *SecurityDetectionsListRulesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsListRulesResponseBody struct {
	Data    []json.RawMessage `json:"data"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
	Total   int               `json:"total"`
}

type SecurityDetectionsListRulesRequest struct {
	Params SecurityDetectionsListRulesRequestParams
}

type SecurityDetectionsListRulesRequestParams struct {
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
	// Filter Search query
	Filter *string `form:"filter,omitempty" json:"filter,omitempty"`
	// SortField Values are created_at, createdAt, enabled, execution_summary.last_execution.date,
	// execution_summary.last_execution.metrics.execution_gap_duration_s, execution_summary.last_execution.metrics.total_indexing_duration_ms,
	// execution_summary.last_execution.metrics.total_search_duration_ms, execution_summary.last_execution.status, name,
	// risk_score, riskScore, severity, updated_at, or updatedAt.
	SortField *string `form:"sort_field,omitempty" json:"sort_field,omitempty"`
	// SortOrder Values are asc or desc.
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
	// Page Page number
	// Minimum value is 1. Default value is 1.
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// PerPage AnonymizationFields per page
	// Minimum value is 0. Default value is 20.
	PerPage        *int    `form:"per_page,omitempty" json:"per_page,omitempty"`
	GapsRangeStart *string `form:"gaps_range_start,omitempty" json:"gaps_range_start,omitempty"`
	GapsRangeEnd   *string `form:"gaps_range_end,omitempty" json:"gaps_range_end,omitempty"`
}

// newSecurityDetectionsListRules returns a function that performs GET /api/detection_engine/rules/_find API requests
func (api *API) newSecurityDetectionsListRules() func(context.Context, *SecurityDetectionsListRulesRequest, ...RequestOption) (*SecurityDetectionsListRulesResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsListRulesRequest, opts ...RequestOption) (*SecurityDetectionsListRulesResponse, error) {
		if req == nil {
			req = &SecurityDetectionsListRulesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_detections.list_rules")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/_find"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Filter != nil {
			params["filter"] = *req.Params.Filter
		}
		if req.Params.Fields != nil {
			params["fields"] = strings.Join(*req.Params.Fields, ",")
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
		if req.Params.GapsRangeStart != nil {
			params["gaps_range_start"] = *req.Params.GapsRangeStart
		}
		if req.Params.GapsRangeEnd != nil {
			params["gaps_range_end"] = *req.Params.GapsRangeEnd
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
			instrument.BeforeRequest(httpReq, "security_detections.list_rules")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.list_rules", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsListRulesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsListRulesResponseBody

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
