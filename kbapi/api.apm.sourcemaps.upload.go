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
// APMSourcemapsUploadResponse wraps the response from a <todo> call
type APMSourcemapsUploadResponse struct {
	StatusCode int
	Body       *APMSourcemapsUploadResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMSourcemapsUploadResponseBody struct {
	Body                 *string  `json:"body,omitempty"`
	CompressionAlgorithm *string  `json:"compressionAlgorithm,omitempty"`
	Created              *string  `json:"created,omitempty"`
	DecodedSHA256        *string  `json:"decodedSha256,omitempty"`
	DecodedSize          *float32 `json:"decodedSize,omitempty"`
	EncodedSHA256        *string  `json:"encodedSha256,omitempty"`
	EncodedSize          *float32 `json:"encodedSize,omitempty"`
	EncryptionAlgorithm  *string  `json:"encryptionAlgorithm,omitempty"`
	ID                   *string  `json:"id,omitempty"`
	Identifier           *string  `json:"identifier,omitempty"`
	PackageName          *string  `json:"packageName,omitempty"`
	RelativeURL          *string  `json:"relative_url,omitempty"`
	Type                 *string  `json:"type,omitempty"`
}

type APMSourcemapsUploadRequest struct {
	Body APMSourcemapsUploadRequestBody
}

type APMSourcemapsUploadRequestBody struct {
	// BundleFilepath The absolute path of the final bundle as used in the web application.
	BundleFilepath string `json:"bundle_filepath"`
	// ServiceName The name of the service that the service map should apply to.
	ServiceName string `json:"service_name"`
	// ServiceVersion The version of the service that the service map should apply to.
	ServiceVersion string `json:"service_version"`
	// Sourcemap The source map. It can be a string or file upload. It must follow the
	// [source map format specification](https://tc39.es/ecma426/).
	Sourcemap []byte `json:"sourcemap"`
}

// newAPMSourcemapsUpload returns a function that performs POST /api/apm/sourcemaps API requests
func (api *API) newAPMSourcemapsUpload() func(context.Context, *APMSourcemapsUploadRequest, ...RequestOption) (*APMSourcemapsUploadResponse, error) {
	return func(ctx context.Context, req *APMSourcemapsUploadRequest, opts ...RequestOption) (*APMSourcemapsUploadResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.sourcemaps.upload")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Set up multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("sourcemap", "sourcemap.json")
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := part.Write(req.Body.Sourcemap); err != nil {
			return nil, fmt.Errorf("failed to write data to form: %w", err)
		}

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

		path := "/api/apm/sourcemaps"

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
			instrument.BeforeRequest(httpReq, "apm.sourcemaps.upload")
			if reader := instrument.RecordRequestBody(ctx, "apm.sourcemaps.upload", httpReq.Body); reader != nil {
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
		resp := &APMSourcemapsUploadResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMSourcemapsUploadResponseBody

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
