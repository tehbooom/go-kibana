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

// SavedObjectImportResponse wraps the response from a SavedObjectImport call
type SavedObjectImportResponse struct {
	StatusCode int
	Body       *SavedObjectImportResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SavedObjectImportResponseBody struct {
	// Errors Indicates the import was unsuccessful and specifies the objects that failed to import.
	//
	// NOTE: One object may result in multiple errors, which requires separate steps to resolve. For instance, a `missing_references` error and conflict error.
	Errors *[]map[string]interface{} `json:"errors,omitempty"`
	// Success Indicates when the import was successfully completed. When set to false, some objects may not have been created. For additional information, refer to the `errors` and `successResults` properties.
	Success *bool `json:"success,omitempty"`

	// SuccessCount Indicates the number of successfully imported records.
	SuccessCount *int `json:"successCount,omitempty"`

	// SuccessResults Indicates the objects that are successfully imported, with any metadata if applicable.
	//
	// NOTE: Objects are created only when all resolvable errors are addressed, including conflicts and missing references. If objects are created as new copies, each entry in the `successResults` array includes a `destinationId` attribute.
	SuccessResults *[]map[string]interface{} `json:"successResults,omitempty"`
}

// SavedObjectExportRequest   is the request for newFleetBulkGetAgentPolicies
type SavedObjectImportRequest struct {
	Params PostSavedObjectImportResponseParams
	Body   SavedObjectImportRequestBody
}

type PostSavedObjectImportResponseParams struct {
	// CreateNewCopies Creates copies of saved objects, regenerates each object ID, and resets the origin. When used, potential conflict errors are avoided. NOTE: This option cannot be used with the `overwrite` and `compatibilityMode` options.
	CreateNewCopies *bool `form:"createNewCopies,omitempty" json:"createNewCopies,omitempty"`

	// Overwrite Overwrites saved objects when they already exist. When used, potential conflict errors are automatically resolved by overwriting the destination object. NOTE: This option cannot be used with the `createNewCopies` option.
	Overwrite *bool `form:"overwrite,omitempty" json:"overwrite,omitempty"`

	// CompatibilityMode Applies various adjustments to the saved objects that are being imported to maintain compatibility between different Kibana versions. Use this option only if you encounter issues with imported saved objects. NOTE: This option cannot be used with the `createNewCopies` option.
	CompatibilityMode *bool `form:"compatibilityMode,omitempty" json:"compatibilityMode,omitempty"`
}

// SavedObjectImportRequestBody defines the body for ExportSavedObjectsDefault.
type SavedObjectImportRequestBody struct {
	// File A file exported using the export API.
	// NOTE: The `savedObjects.maxImportExportSize` configuration setting limits the number of saved objects which may be included in this file. Similarly, the `savedObjects.maxImportPayloadBytes` setting limits the overall size of the file that can be imported.
	File []byte `json:"file,omitempty"`
}

// newSavedObjectImport returns a function that performs POST /api/saved_objects/_import API requests
func (api *API) newSavedObjectImport() func(context.Context, *SavedObjectImportRequest, ...RequestOption) (*SavedObjectImportResponse, error) {
	return func(ctx context.Context, req *SavedObjectImportRequest, opts ...RequestOption) (*SavedObjectImportResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("File is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "saved_objects.import")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.CreateNewCopies != nil {
			params["createNewCopies"] = strconv.FormatBool(*req.Params.CreateNewCopies)
		}
		if req.Params.Overwrite != nil {
			params["overwrite"] = strconv.FormatBool(*req.Params.Overwrite)
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

		// Close the multipart writer
		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close writer: %w", err)
		}

		path := "/api/saved_objects/_import"

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
				return nil, err
			}
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "saved_objects.import")
			if reader := instrument.RecordRequestBody(ctx, "saved_objects.import", httpReq.Body); reader != nil {
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
		resp := &SavedObjectImportResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SavedObjectImportResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
