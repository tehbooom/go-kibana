package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityAIAssistantCreateConversationResponse wraps the response from a <todo> call
type SecurityAIAssistantCreateConversationResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantConversationResponse
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantCreateConversationRequest struct {
	Body SecurityAIAssistantCreateConversationRequestBody
}

type SecurityAIAssistantCreateConversationRequestBody struct {
	ApiConfig *SecurityAIAssistantAPIConfig `json:"apiConfig,omitempty"`
	// Category The conversation category.
	// Values are assistant or insights.
	Category                           *string                       `json:"category,omitempty"`
	ExcludeFromLastConversationStorage *bool                         `json:"excludeFromLastConversationStorage,omitempty"`
	ID                                 *string                       `json:"id,omitempty"`
	Messages                           *[]SecurityAIAssistantMessage `json:"messages,omitempty"`
	Replacements                       *map[string]string            `json:"replacements,omitempty"`
	Title                              string                        `json:"title"`
}

// newSecurityAIAssistantCreateConversation returns a function that performs POST /api/security_ai_assistant/current_user/conversations API requests
func (api *API) newSecurityAIAssistantCreateConversation() func(context.Context, *SecurityAIAssistantCreateConversationRequest, ...RequestOption) (*SecurityAIAssistantCreateConversationResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantCreateConversationRequest, opts ...RequestOption) (*SecurityAIAssistantCreateConversationResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.create_conversation")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/current_user/conversations"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.create_conversation")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.create_conversation", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantCreateConversationResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantConversationResponse

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
