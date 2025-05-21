package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// SecurityDetectionsImportRulesResponse wraps the response from a ImportRules call
type SecurityDetectionsImportRulesResponse struct {
	StatusCode int
	Body       *SecurityDetectionsImportRulesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsImportRulesResponseBody struct {
}

type SecurityDetectionsImportRulesRequest struct {
	Params SecurityDetectionsImportRulesRequestParams
	Body   SecurityDetectionsImportRulesRequestBody
}

type SecurityDetectionsImportRulesRequestParams struct {
	// Overwrite Determines whether existing rules with the same `rule_id` are overwritten.
	Overwrite *bool `form:"overwrite,omitempty" json:"overwrite,omitempty"`
	// OverwriteExceptions Determines whether existing exception lists with the same `list_id` are overwritten. Both the exception list container and its items are overwritten.
	OverwriteExceptions *bool `form:"overwrite_exceptions,omitempty" json:"overwrite_exceptions,omitempty"`
	// OverwriteActionConnectors Determines whether existing actions with the same `kibana.alert.rule.actions.id` are overwritten.
	OverwriteActionConnectors *bool `form:"overwrite_action_connectors,omitempty" json:"overwrite_action_connectors,omitempty"`
	// AsNewList Generates a new list ID for each imported exception list.
	AsNewList *bool `form:"as_new_list,omitempty" json:"as_new_list,omitempty"`
}

type SecurityDetectionsImportRulesRequestBody struct {
	File []byte
}

// newSecurityDetectionsImportRules returns a function that performs POST /api/detection_engine/rules/_import API requests
func (api *API) newSecurityDetectionsImportRules() func(context.Context, *SecurityDetectionsImportRulesRequest, ...RequestOption) (*SecurityDetectionsImportRulesResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsImportRulesRequest, opts ...RequestOption) (*SecurityDetectionsImportRulesResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.import_rules")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Overwrite != nil {
			params["overwrite"] = strconv.FormatBool(*req.Params.Overwrite)
		}
		if req.Params.OverwriteExceptions != nil {
			params["overwrite_exceptions"] = strconv.FormatBool(*req.Params.OverwriteExceptions)
		}
		if req.Params.OverwriteActionConnectors != nil {
			params["overwrite_action_connectors"] = strconv.FormatBool(*req.Params.OverwriteActionConnectors)
		}
		if req.Params.AsNewList != nil {
			params["as_new_list"] = strconv.FormatBool(*req.Params.AsNewList)
		}

		// Set up multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "export.ndjson")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := part.Write(req.Body.File); err != nil {
			return nil, fmt.Errorf("failed to write data to form: %w", err)
		}

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

		path := "/api/detection_engine/rules/_import"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		// Set the content type for multipart form data
		httpReq.Header.Set("Content-Type", writer.FormDataContentType())

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
			instrument.BeforeRequest(httpReq, "security_detections.import_rules")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.import_rules", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsImportRulesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsImportRulesResponseBody

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
