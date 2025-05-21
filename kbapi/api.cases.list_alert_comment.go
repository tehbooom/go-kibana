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
// CasesListCommentsAlertsResponse wraps the response from a <todo> call
type CasesListCommentsAlertsResponse struct {
	StatusCode int
	Body       *CasesListCommentsAlertsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesListCommentsAlertsResponseBody struct {
	Comments []json.RawMessage `json:"comments"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
	Total    int               `json:"total"`
}

// Method to get all comments with their appropriate types
func (resp *CasesListCommentsAlertsResponseBody) GetAllTypedComments() ([]interface{}, error) {
	if len(resp.Comments) == 0 {
		return []interface{}{}, nil
	}

	comments := make([]interface{}, 0, len(resp.Comments))
	for _, rawComment := range resp.Comments {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(rawComment, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "user":
			var comment UserCommentResponse
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		case "alert":
			var comment AlertCommentResponse
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		default:
			var comment BaseComment
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

type CasesListCommentsAlertsRequest struct {
	ID     string
	Params CasesListCommentsAlertsRequestParams
}

type CasesListCommentsAlertsRequestParams struct {
	// PerPage The number of rules to return per page.
	// Maximum value is 100. Default value is 20.
	PerPage *int `form:"per_page,omitempty" json:"per_page,omitempty"`
	// Page The page number to return.
	// Default value is 1
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// SortOrder Determines the sort order.
	// Values are asc or desc. Default value is desc.
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
}

// newCasesListCommentsAlerts returns a function that performs GET /api/cases/{caseId}/comments/_find API requests
func (api *API) newCasesListCommentsAlerts() func(context.Context, *CasesListCommentsAlertsRequest, ...RequestOption) (*CasesListCommentsAlertsResponse, error) {
	return func(ctx context.Context, req *CasesListCommentsAlertsRequest, opts ...RequestOption) (*CasesListCommentsAlertsResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.list_alert_comment")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/cases/%s/comments/_find", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PerPage != nil {
			params["per_page"] = strconv.Itoa(*req.Params.PerPage)
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
			instrument.BeforeRequest(httpReq, "cases.list_alert_comment")
			if reader := instrument.RecordRequestBody(ctx, "cases.list_alert_comment", httpReq.Body); reader != nil {
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
		resp := &CasesListCommentsAlertsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesListCommentsAlertsResponseBody

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
