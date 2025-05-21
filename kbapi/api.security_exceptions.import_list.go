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

// SecurityExceptionsImportListResponse wraps the response from a ImportList call
type SecurityExceptionsImportListResponse struct {
	StatusCode int
	Body       *SecurityExceptionsImportListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityExceptionsImportListResponseBody struct {
	Errors                         []SecurityExceptionsErrorDetail `json:"errors"`
	Success                        string                          `json:"success"`
	SuccessCount                   string                          `json:"success_count"`
	SuccessExceptionLists          string                          `json:"success_exception_lists"`
	SuccessExceptionListItems      string                          `json:"success_exception_list_items"`
	SuccessCountExceptionLists     string                          `json:"success_count_exception_lists"`
	SuccessCountExceptionListItems int                             `json:"success_count_exception_list_items"`
}

type SecurityExceptionsImportListRequest struct {
	Params SecurityExceptionsImportListRequestParams
	Body   SecurityExceptionsImportListRequestBody
}

type SecurityExceptionsImportListRequestParams struct {
	// AsNewList Determines whether the list being imported will have a new list_id generated.
	// Additional item_id's are generated for each exception item. Both the exception list and its items are overwritten.
	// Default value is false.
	AsNewList *bool
	// Overwrite Determines whether existing exception lists with the same list_id are overwritten.
	// If any exception items have the same item_id, those are also overwritten.
	// Default value is false.
	Overwrite *bool
}

type SecurityExceptionsImportListRequestBody struct {
	File []byte
}

// newSecurityExceptionsImportList returns a function that performs POST /api/exception_lists/_import API requests
func (api *API) newSecurityExceptionsImportList() func(context.Context, *SecurityExceptionsImportListRequest, ...RequestOption) (*SecurityExceptionsImportListResponse, error) {
	return func(ctx context.Context, req *SecurityExceptionsImportListRequest, opts ...RequestOption) (*SecurityExceptionsImportListResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_exceptions.import_list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/exception_lists/_import"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.AsNewList != nil {
			params["as_new_list"] = strconv.FormatBool(*req.Params.AsNewList)
		}
		if req.Params.Overwrite != nil {
			params["overwrite"] = strconv.FormatBool(*req.Params.Overwrite)
		}

		// Set up multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "export.ndjson")
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := part.Write(req.Body.File); err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, fmt.Errorf("failed to write data to form: %w", err)
		}

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

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
			instrument.BeforeRequest(httpReq, "security_exceptions.import_list")
			if reader := instrument.RecordRequestBody(ctx, "security_exceptions.import_list", httpReq.Body); reader != nil {
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
		resp := &SecurityExceptionsImportListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityExceptionsImportListResponseBody

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
