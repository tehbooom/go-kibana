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
// SecurityAIAssistantBulkActionPromptsResponse wraps the response from a <todo> call
type SecurityAIAssistantBulkActionPromptsResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantBulkActionPromptsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantBulkActionPromptsResponseBody struct {
	PromptsCount *int `json:"prompts_count,omitempty"`
	Attributes   struct {
		Errors  *[]SecurityAIAssistantBulkActionPromptsErrors `json:"errors,omitempty"`
		Results SecurityAIAssistantBulkActionKBResults        `json:"results"`
		Summary SecurityAIAssistantBulkActionSummary          `json:"summary"`
	} `json:"attributes"`
	Message    *string `json:"message,omitempty"`
	StatusCode *int    `json:"status_code,omitempty"`
	Success    *bool   `json:"success,omitempty"`
}

type SecurityAIAssistantBulkActionPromptsRequest struct {
	Body SecurityAIAssistantBulkActionPromptsRequestBody
}

type SecurityAIAssistantBulkActionPromptsRequestBody struct {
	Create []SecurityAIAssistantPromptRequest `json:"create,omitempty"`
	Delete SecurityAIAssistantDeleteObject    `json:"delete,omitempty"`
	Update []SecurityAIAssistantPromptRequest `json:"update,omitempty"`
}

// newSecurityAIAssistantBulkActionPrompts returns a function that performs POST /api/security_ai_assistant/prompts/_bulk_action API requests
func (api *API) newSecurityAIAssistantBulkActionPrompts() func(context.Context, *SecurityAIAssistantBulkActionPromptsRequest, ...RequestOption) (*SecurityAIAssistantBulkActionPromptsResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantBulkActionPromptsRequest, opts ...RequestOption) (*SecurityAIAssistantBulkActionPromptsResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.bulk_action_prompts")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/prompts/_bulk_action"

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.bulk_action_prompts")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.bulk_action_prompts", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantBulkActionPromptsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantBulkActionPromptsResponseBody

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
