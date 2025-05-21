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
// SecurityAIAssistantCreateKnowledgeBaseResponse wraps the response from a <todo> call
type SecurityAIAssistantCreateKnowledgeBaseResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantCreateKnowledgeBaseResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantCreateKnowledgeBaseResponseBody struct {
	Success bool `json:"success"`
}

type SecurityAIAssistantCreateKnowledgeBaseRequest struct {
	// Resource The KnowledgeBase resource value.
	Resource string
	Params   SecurityAIAssistantCreateKnowledgeBaseRequestParams
	Body     SecurityAIAssistantCreateKnowledgeBaseRequestBody
}

type SecurityAIAssistantCreateKnowledgeBaseRequestParams struct {
	// ModelID Optional ELSER modelId to use when setting up the Knowledge Base
	ModelID *string
	// IgnoreSecurityLabs Indicates whether we should or should not install Security Labs docs when setting up the Knowledge Base
	// Default value is false.
	IgnoreSecurityLabs *bool
}

type SecurityAIAssistantCreateKnowledgeBaseRequestBody struct {
}

// newSecurityAIAssistantCreateKnowledgeBase returns a function that performs POST /api/security_ai_assistant/knowledge_base/{resource} API requests
func (api *API) newSecurityAIAssistantCreateKnowledgeBase() func(context.Context, *SecurityAIAssistantCreateKnowledgeBaseRequest, ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantCreateKnowledgeBaseRequest, opts ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.create_knowledge_base")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/security_ai_assistant/knowledge_base/%s", req.Resource)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ModelID != nil {
			params["modelId"] = *req.Params.ModelID
		}
		if req.Params.IgnoreSecurityLabs != nil {
			params["ignoreSecurityLabs"] = strconv.FormatBool(*req.Params.IgnoreSecurityLabs)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.create_knowledge_base")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.create_knowledge_base", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantCreateKnowledgeBaseResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantCreateKnowledgeBaseResponseBody

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
