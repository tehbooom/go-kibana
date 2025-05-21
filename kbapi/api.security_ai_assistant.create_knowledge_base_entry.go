package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// To get the body use one of the following methods for SecurityAIAssistantCreateKnowledgeBaseEntryResponse
// GetDocumentEntry or GetIndexEntry
// TODO: Update the call
// SecurityAIAssistantCreateKnowledgeBaseEntryResponse wraps the response from a <todo> call
type SecurityAIAssistantCreateKnowledgeBaseEntryResponse struct {
	StatusCode int
	Body       json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

func (resp *SecurityAIAssistantCreateKnowledgeBaseEntryResponse) GetDocumentEntry() (*SecurityAIAssistantDocumentEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var docResponse SecurityAIAssistantDocumentEntryResponse
	err := json.Unmarshal(resp.Body, &docResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document entry response: %v", err)
	}

	return &docResponse, nil
}

func (resp *SecurityAIAssistantCreateKnowledgeBaseEntryResponse) GetIndexEntry() (*SecurityAIAssistantIndexEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var docResponse SecurityAIAssistantIndexEntryResponse
	err := json.Unmarshal(resp.Body, &docResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document entry response: %v", err)
	}

	return &docResponse, nil
}

// To set the body use one of the following methods for SecurityAIAssistantCreateKnowledgeBaseEntryRequest
// SetDocumentEntry or SetIndexEntry
type SecurityAIAssistantCreateKnowledgeBaseEntryRequest struct {
	Body json.RawMessage
}

func (body *SecurityAIAssistantCreateKnowledgeBaseEntryRequest) SetDocumentEntry(entry SecurityAIAssistantDocumentEntryRequest) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	body.Body = data
	return nil
}

func (body *SecurityAIAssistantCreateKnowledgeBaseEntryRequest) SetIndexEntry(entry SecurityAIAssistantIndexEntryResponse) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	body.Body = data
	return nil
}

// newSecurityAIAssistantCreateKnowledgeBaseEntry returns a function that performs POST /api/security_ai_assistant/knowledge_base/entries API requests
func (api *API) newSecurityAIAssistantCreateKnowledgeBaseEntry() func(context.Context, *SecurityAIAssistantCreateKnowledgeBaseEntryRequest, ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseEntryResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantCreateKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseEntryResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.create_knowledge_base_entry")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/knowledge_base/entries"

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.create_knowledge_base_entry")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.create_knowledge_base_entry", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantCreateKnowledgeBaseEntryResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}
		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()

		if httpResp.StatusCode < 299 {
			resp.Body = bodyBytes
			return resp, nil
		} else {
			// For all non-success responses
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
