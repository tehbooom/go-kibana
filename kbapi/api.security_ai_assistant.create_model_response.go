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
// SecurityAIAssistantCreateModelResponseResponse wraps the response from a <todo> call
type SecurityAIAssistantCreateModelResponseResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantCreateModelResponseResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantCreateModelResponseResponseBody struct{}

type SecurityAIAssistantCreateModelResponseRequest struct {
	Body SecurityAIAssistantCreateModelResponseRequestBody
}

type SecurityAIAssistantCreateModelResponseRequestBody struct {
	ConnectorID      string                           `json:"connectorId"`
	ConversationID   *string                          `json:"conversationId,omitempty"`
	IsStream         *bool                            `json:"isStream,omitempty"`
	LangSmithAPIKey  *string                          `json:"langSmithApiKey,omitempty"`
	LangSmithProject *string                          `json:"langSmithProject,omitempty"`
	Messages         []SecurityAIAssistantChatMessage `json:"messages"`
	Model            *string                          `json:"model,omitempty"`
	Persist          bool                             `json:"persist"`
	PromptID         *string                          `json:"promptId,omitempty"`
	ResponseLanguage *string                          `json:"responseLanguage,omitempty"`
}

// newSecurityAIAssistantCreateModelResponse returns a function that performs POST /api/security_ai_assistant/chat/complete API requests
func (api *API) newSecurityAIAssistantCreateModelResponse() func(context.Context, *SecurityAIAssistantCreateModelResponseRequest, ...RequestOption) (*SecurityAIAssistantCreateModelResponseResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantCreateModelResponseRequest, opts ...RequestOption) (*SecurityAIAssistantCreateModelResponseResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.create_model_response")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/chat/complete"

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.create_model_response")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.create_model_response", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantCreateModelResponseResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantCreateModelResponseResponseBody

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
