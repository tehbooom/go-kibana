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

// TODO: Update the call
// SavedObjectResolveImportsResponse wraps the response from a <todo> call
type SavedObjectResolveImportsResponse struct {
	StatusCode int
	Body       *SavedObjectResolveImportsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SavedObjectResolveImportsResponseBody struct {
	// Errors Specifies the objects that failed to resolve.
	// NOTE: One object can result in multiple errors, which requires separate steps to resolve.
	// For instance, a missing_references error and a conflict error.
	Errors []any `json:"errors,omitempty"`
	// Success Indicates a successful import. When set to false, some objects may not have been created.
	// For additional information, refer to the errors and successResults properties.
	Success bool `json:"success"`
	// SuccessCount Indicates the number of successfully resolved records.
	SuccessCount int `json:"successCount"`
	// SuccessResults Indicates the objects that are successfully imported, with any metadata if applicable.
	// NOTE: Objects are only created when all resolvable errors are addressed, including conflict and missing references.
	SuccessResults []struct {
		ID   string `json:"id,omitempty"`
		Meta struct {
			Icon  string `json:"icon,omitempty"`
			Title string `json:"title,omitempty"`
		} `json:"meta,omitempty"`
		Type string `json:"type,omitempty"`
	} `json:"successResults,omitempty"`
}

type SavedObjectResolveImportsRequest struct {
	Params SavedObjectResolveImportsRequestParams
	Body   SavedObjectResolveImportsRequestBody
}

type SavedObjectResolveImportsRequestParams struct {
	// CompatibilityMode Applies various adjustments to the saved objects that are being imported to maintain compatibility between different
	// Kibana versions. When enabled during the initial import, also enable when resolving import errors.
	// This option cannot be used with the `createNewCopies` option.
	CompatibilityMode *bool `form:"compatibilityMode,omitempty" json:"compatibilityMode,omitempty"`

	// CreateNewCopies Creates copies of the saved objects, regenerates each object ID, and resets the origin.
	// When enabled during the initial import, also enable when resolving import errors.
	CreateNewCopies *bool `form:"createNewCopies,omitempty" json:"createNewCopies,omitempty"`
}

type SavedObjectResolveImportsRequestBody struct {
	// File the same file given to the import API.
	File []byte `json:"file,omitempty"`
	// Retries The retry operations, which can specify how to resolve different types of errors.
	Retries []SavedObjectRetryOperationItem `json:"retries"`
}

type SavedObjectRetryOperationItem struct {
	// DestinationID Specifies the destination ID that the imported object should have, if different from the current ID.
	DestinationID *string `json:"destinationId,omitempty"`

	// ID The saved object ID.
	ID string `json:"id"`

	// IgnoreMissingReferences When set to `true`, ignores missing reference errors. When set to `false`, does nothing.
	IgnoreMissingReferences *bool `json:"ignoreMissingReferences,omitempty"`

	// Overwrite When set to `true`, the source object overwrites the conflicting destination object. When set to `false`, does nothing.
	Overwrite *bool `json:"overwrite,omitempty"`

	// ReplaceReferences A list of `type`, `from`, and `to` used to change the object references.
	ReplaceReferences *[]struct {
		From *string `json:"from,omitempty"`
		To   *string `json:"to,omitempty"`
		Type *string `json:"type,omitempty"`
	} `json:"replaceReferences,omitempty"`

	// Type The saved object type.
	Type string `json:"type"`
}

// newSavedObjectResolveImports returns a function that performs POST /api/saved_objects/_resolve_import_errors API requests
func (api *API) newSavedObjectResolveImports() func(context.Context, *SavedObjectResolveImportsRequest, ...RequestOption) (*SavedObjectResolveImportsResponse, error) {
	return func(ctx context.Context, req *SavedObjectResolveImportsRequest, opts ...RequestOption) (*SavedObjectResolveImportsResponse, error) {
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
			newCtx = instrument.Start(ctx, "saved_objects.resolve_imports")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.CreateNewCopies != nil {
			params["createNewCopies"] = strconv.FormatBool(*req.Params.CreateNewCopies)
		}
		if req.Params.CompatibilityMode != nil {
			params["compatibilityMode"] = strconv.FormatBool(*req.Params.CompatibilityMode)
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

		for i, retry := range req.Body.Retries {
			retryJSON, err := json.Marshal(retry)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal retry data: %w", err)
			}

			fieldName := fmt.Sprintf("retries[%d]", i)
			retryField, err := writer.CreateFormField(fieldName)
			if err != nil {
				return nil, fmt.Errorf("failed to create retries field: %w", err)
			}

			if _, err := retryField.Write(retryJSON); err != nil {
				return nil, fmt.Errorf("failed to write retry data to form: %w", err)
			}
		}

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

		path := "/api/saved_objects/_resolve_import_errors"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, body)
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
			instrument.BeforeRequest(httpReq, "saved_objects.resolve_imports")
			if reader := instrument.RecordRequestBody(ctx, "saved_objects.resolve_imports", httpReq.Body); reader != nil {
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
		resp := &SavedObjectResolveImportsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SavedObjectResolveImportsResponseBody

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
