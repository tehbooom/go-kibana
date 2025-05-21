package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// TODO: Update the call
// CasesAttachFileResponse wraps the response from a <todo> call
type CasesAttachFileResponse struct {
	StatusCode int
	Body       *CasesObjectResponse
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesAttachFileRequest struct {
	ID   string
	Body CasesAttachFileRequestBody
}

type CasesAttachFileRequestBody struct {
	// The file being attached to the case.
	File []byte `json:"file,omitempty"`
	// The desired name of the file being attached to the case,
	// it can be different than the name of the file in the filesystem. This should not include the file extension.
	Filename string `json:"filename"`
}

// newCasesAttachFile returns a function that performs POST /api/cases/{caseId}/files API requests
func (api *API) newCasesAttachFile() func(context.Context, *CasesAttachFileRequest, ...RequestOption) (*CasesAttachFileResponse, error) {
	return func(ctx context.Context, req *CasesAttachFileRequest, opts ...RequestOption) (*CasesAttachFileResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.attach_file")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/cases/%s/files", req.ID)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", "file.txt")
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

		// Add the filename field to the form
		if req.Body.Filename == "" {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, fmt.Errorf("filename cannot be empty")
		}

		if err := writer.WriteField("filename", req.Body.Filename); err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, fmt.Errorf("failed to add filename field: %w", err)
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
			instrument.BeforeRequest(httpReq, "cases.attach_file")
			if reader := instrument.RecordRequestBody(ctx, "cases.attach_file", httpReq.Body); reader != nil {
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
		resp := &CasesAttachFileResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesObjectResponse

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
