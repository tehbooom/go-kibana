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
// SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse wraps the response from a <todo> call
type SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse struct {
	StatusCode int
	Body       *SecurityAIAssistantBulkActionKnowledgeBaseEntryResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityAIAssistantBulkActionKnowledgeBaseEntryResponseBody struct {
	KnowledgeBaseEntriesCount *int `json:"knowledgeBaseEntriesCount,omitempty"`
	Attributes                struct {
		Errors  *[]SecurityAIAssistantBulkActionKBErrors `json:"errors,omitempty"`
		Results SecurityAIAssistantBulkActionKBResults   `json:"results"`
		Summary SecurityAIAssistantBulkActionSummary     `json:"summary"`
	} `json:"attributes"`
	Message    *string `json:"message,omitempty"`
	StatusCode *int    `json:"status_code,omitempty"`
	Success    *bool   `json:"success,omitempty"`
}

func (resp *SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse) GetUpdatedDocumentResults() (*[]SecurityAIAssistantDocumentEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var docResponse []SecurityAIAssistantDocumentEntryResponse
	updatedResults := resp.Body.Attributes.Results.Updated
	for _, result := range updatedResults {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(result, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "document":
			var documentResult SecurityAIAssistantDocumentEntryResponse
			err := json.Unmarshal(result, &documentResult)
			if err != nil {
				return nil, err
			}
			docResponse = append(docResponse, documentResult)
		}
	}
	return &docResponse, nil
}

func (resp *SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse) GetUpdatedIndexResults() (*[]SecurityAIAssistantIndexEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var indexResponse []SecurityAIAssistantIndexEntryResponse
	updatedResults := resp.Body.Attributes.Results.Updated
	for _, result := range updatedResults {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(result, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "index":
			var indexResult SecurityAIAssistantIndexEntryResponse
			err := json.Unmarshal(result, &indexResult)
			if err != nil {
				return nil, err
			}
			indexResponse = append(indexResponse, indexResult)
		}
	}
	return &indexResponse, nil
}

func (resp *SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse) GetCreatedDocumentResults() (*[]SecurityAIAssistantDocumentEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var docResponse []SecurityAIAssistantDocumentEntryResponse
	createdResults := resp.Body.Attributes.Results.Created
	for _, result := range createdResults {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(result, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "document":
			var documentResult SecurityAIAssistantDocumentEntryResponse
			err := json.Unmarshal(result, &documentResult)
			if err != nil {
				return nil, err
			}
			docResponse = append(docResponse, documentResult)
		}
	}
	return &docResponse, nil
}

func (resp *SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse) GetCreatedIndexResults() (*[]SecurityAIAssistantIndexEntryResponse, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	var indexResponse []SecurityAIAssistantIndexEntryResponse
	createdResults := resp.Body.Attributes.Results.Created
	for _, result := range createdResults {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(result, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "index":
			var indexResult SecurityAIAssistantIndexEntryResponse
			err := json.Unmarshal(result, &indexResult)
			if err != nil {
				return nil, err
			}
			indexResponse = append(indexResponse, indexResult)
		}
	}
	return &indexResponse, nil
}

type SecurityAIAssistantBulkActionKnowledgeBaseEntryRequest struct {
	Body SecurityAIAssistantBulkActionKnowledgeBaseEntryRequestBody
}

type SecurityAIAssistantBulkActionKnowledgeBaseEntryRequestBody struct {
	Create []json.RawMessage               `json:"create,omitempty"`
	Delete SecurityAIAssistantDeleteObject `json:"delete,omitempty"`
	Update []json.RawMessage               `json:"update,omitempty"`
}

func (body *SecurityAIAssistantBulkActionKnowledgeBaseEntryRequestBody) SetCreateDocumentEntry(entry SecurityAIAssistantDocumentEntryRequest) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	body.Create = append(body.Create, data)
	return nil
}

func (body *SecurityAIAssistantBulkActionKnowledgeBaseEntryRequestBody) SetCreateIndexEntry(entry SecurityAIAssistantIndexEntryResponse) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	body.Update = append(body.Create, data)
	return nil
}

// newSecurityAIAssistantBulkActionKnowledgeBaseEntry returns a function that performs POST /api/security_ai_assistant/knowledge_base/entries/_bulk_action API requests
func (api *API) newSecurityAIAssistantBulkActionKnowledgeBaseEntry() func(context.Context, *SecurityAIAssistantBulkActionKnowledgeBaseEntryRequest, ...RequestOption) (*SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse, error) {
	return func(ctx context.Context, req *SecurityAIAssistantBulkActionKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_ai_assistant.bulk_action_knowledge_base_entry")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/security_ai_assistant/knowledge_base/entries/_bulk_action"

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
			instrument.BeforeRequest(httpReq, "security_ai_assistant.bulk_action_knowledge_base_entry")
			if reader := instrument.RecordRequestBody(ctx, "security_ai_assistant.bulk_action_knowledge_base_entry", httpReq.Body); reader != nil {
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
		resp := &SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityAIAssistantBulkActionKnowledgeBaseEntryResponseBody

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
