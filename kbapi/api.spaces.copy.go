package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SpacesCopyObjectsResponse  wraps the response from a SavedObjectImport call
type SpacesCopyObjectsResponse struct {
	StatusCode int
	Body       *SpacesCopyObjectsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SpacesCopyObjectsResponseBody struct {
	Spaces map[string]CopyResult `json:",inline"`
}

// CopyResult represents the import results for a specific namespace
type CopyResult struct {
	// Success indicates whether the import was successful
	Success bool `json:"success"`

	// SuccessCount indicates the number of successfully imported objects
	SuccessCount int `json:"successCount,omitempty"`

	// SuccessResults contains the details of successfully imported objects
	// Using the exact type you specified
	SuccessResults *[]map[string]interface{} `json:"successResults,omitempty"`
}

// SavedObjectExportRequest is the request for newFleetBulkGetAgentPolicies
type SpacesCopyObjectsRequest struct {
	Body SpacesCopyObjectsRequestBody
}

// SavedObjectImportRequestBody defines the body for ExportSavedObjectsDefault.
type SpacesCopyObjectsRequestBody struct {
	// CompatibilityMode Apply various adjustments to the saved objects that are being copied to maintain compatibility between different Kibana versions. Use this option only if you encounter issues with copied saved objects. This option cannot be used with the `createNewCopies` option.
	CompatibilityMode *bool `json:"compatibilityMode,omitempty"`

	// CreateNewCopies Create new copies of saved objects, regenerate each object identifier, and reset the origin. When used, potential conflict errors are avoided.  This option cannot be used with the `overwrite` and `compatibilityMode` options.
	CreateNewCopies *bool `json:"createNewCopies,omitempty"`

	// IncludeReferences When set to true, all saved objects related to the specified saved objects will also be copied into the target spaces.
	IncludeReferences *bool    `json:"includeReferences,omitempty"`
	Objects           []Object `json:"objects"`

	// Overwrite When set to true, all conflicts are automatically overridden. When a saved object with a matching type and identifier exists in the target space, that version is replaced with the version from the source space. This option cannot be used with the `createNewCopies` option.
	Overwrite *bool    `json:"overwrite,omitempty"`
	Spaces    []string `json:"spaces"`
}

// newSavedObjectImport returns a function that performs POST /api/saved_objects/_import API requests
func (api *API) newSpacesCopyObjects() func(context.Context, *SpacesCopyObjectsRequest, ...RequestOption) (*SpacesCopyObjectsResponse, error) {
	return func(ctx context.Context, req *SpacesCopyObjectsRequest, opts ...RequestOption) (*SpacesCopyObjectsResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Objects is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "spaces.copy")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/_copy_saved_objects"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "spaces.copy")
			if reader := instrument.RecordRequestBody(ctx, "spaces.copy", httpReq.Body); reader != nil {
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
		resp := &SpacesCopyObjectsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesCopyObjectsResponseBody

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
