package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityAIAssistantGetKnowledgeBaseResponse wraps the response from a <todo> call
type SecurityAIAssistantGetKnowledgeBaseResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantGetKnowledgeBaseResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantGetKnowledgeBaseResponseBody struct {
	ELSERExists                *bool   `json:"elser_exists,omitempty"`
	IsSetupAvailable           *bool   `json:"is_setup_available,omitempty"`
	IsSetupInProgress          *bool   `json:"is_setup_in_progress,omitempty"`
	ProductDocumentationStatus *string `json:"product_documentation_status,omitempty"`
	SecurityLabsExists         *bool   `json:"security_labs_exists,omitempty"`
	UserDataExists             *bool   `json:"user_data_exists,omitempty"`
}

type SecurityAIAssistantGetKnowledgeBaseRequest struct {
	// Resource The KnowledgeBase resource value.
	Resource string
}

// newSecurityAIAssistantGetKnowledgeBase returns a function that performs GET /api/security_ai_assistant/knowledge_base/{resource} API requests
func (api *API) newSecurityAIAssistantGetKnowledgeBase() func(context.Context, *SecurityAIAssistantGetKnowledgeBaseRequest, ...RequestOption) (*SecurityAIAssistantGetKnowledgeBaseResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantGetKnowledgeBaseRequest, opts ...RequestOption) (*SecurityAIAssistantGetKnowledgeBaseResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.get_knowledge_base")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/security_ai_assistant/knowledge_base/%s", req.Resource)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.get_knowledge_base")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.get_knowledge_base", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantGetKnowledgeBaseResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantGetKnowledgeBaseResponseBody

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
