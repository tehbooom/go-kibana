package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityAIAssistantDeleteKnowledgeBaseEntryResponse wraps the response from a <todo> call
type SecurityAIAssistantDeleteKnowledgeBaseEntryResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantDeleteKnowledgeBaseEntryResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantDeleteKnowledgeBaseEntryResponseBody struct {
	ID string `json:"id"`
}

type SecurityAIAssistantDeleteKnowledgeBaseEntryRequest struct {
	// ID The Knowledge Base Entry's id value.
	ID string
}

// newSecurityAIAssistantDeleteKnowledgeBaseEntry returns a function that performs DELETE /api/security_ai_assistant/knowledge_base/entries/{id} API requests
func (api *API) newSecurityAIAssistantDeleteKnowledgeBaseEntry() func(context.Context, *SecurityAIAssistantDeleteKnowledgeBaseEntryRequest, ...RequestOption) (*SecurityAIAssistantDeleteKnowledgeBaseEntryResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantDeleteKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantDeleteKnowledgeBaseEntryResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.delete_knowledge_base_entry")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/security_ai_assistant/knowledge_base/entries/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.delete_knowledge_base_entry")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.delete_knowledge_base_entry", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantDeleteKnowledgeBaseEntryResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantDeleteKnowledgeBaseEntryResponseBody

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
