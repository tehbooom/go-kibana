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
// SecurityAIAssistantListConversationsResponse wraps the response from a <todo> call
type SecurityAIAssistantListConversationsResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantListConversationsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantListConversationsResponseBody struct {
	Data    []SecurityAIAssistantConversationResponse `json:"data"`
	Page    int                                       `json:"page"`
	PerPage int                                       `json:"per_page"`
	Total   int                                       `json:"total"`
}

type SecurityAIAssistantListConversationsRequest struct {
	Params SecurityAIAssistantListConversationsRequestParams
}

type SecurityAIAssistantListConversationsRequestParams struct {
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
	// Filter Search query
	Filter *string `form:"filter,omitempty" json:"filter,omitempty"`
	// SortField Values are created_at, is_default, title, or updated_at.
	SortField *string `form:"sort_field,omitempty" json:"sort_field,omitempty"`
	// SortOrder Values are asc or desc.
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
	// Page Page number
	// Minimum value is 1. Default value is 1.
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// PerPage AnonymizationFields per page
	// Minimum value is 0. Default value is 20.
	PerPage *int `form:"per_page,omitempty" json:"per_page,omitempty"`
}

// newSecurityAIAssistantListConversations returns a function that performs GET /api/security_ai_assistant/current_user/conversations/_find API requests
func (api *API) newSecurityAIAssistantListConversations() func(context.Context, *SecurityAIAssistantListConversationsRequest, ...RequestOption) (*SecurityAIAssistantListConversationsResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantListConversationsRequest, opts ...RequestOption) (*SecurityAIAssistantListConversationsResponse, error) {
		if req == nil {
			req = &SecurityAIAssistantListConversationsRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_ai_assistant.list_conversations")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/current_user/conversations/_find"

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.list_conversations")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.list_conversations", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantListConversationsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantListConversationsResponseBody

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
