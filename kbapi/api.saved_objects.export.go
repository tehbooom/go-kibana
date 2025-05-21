package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// SavedObjectExport  wraps the response from a FleetEPMBulkGetAssets  call
type SavedObjectExportResponse struct {
	StatusCode int
	Body       []json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

// SavedObjectExportRequest   is the request for newFleetBulkGetAgentPolicies
type SavedObjectExportRequest struct {
	Body SavedObjectExportRequestBody
}

// Object represents an individual object to be exported
type Object struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type SavedObjectExportResponseBody struct {
	ID                   string                 `json:"id"`
	Type                 string                 `json:"type"`
	Managed              bool                   `json:"managed"`
	Version              string                 `json:"version"`
	Attributes           map[string]interface{} `json:"attributes"`
	CreatedAt            string                 `json:"created_at"`
	References           []SavedObjectReference `json:"references"`
	UpdatedAt            string                 `json:"updated_at"`
	CoreMigrationVersion string                 `json:"coreMigrationVersion"`
	TypeMigrationVersion string                 `json:"typeMigrationVersion"`
}

// SavedObjectReference represents a reference to another saved object
type SavedObjectReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// SavedObjectExportRequestBody  defines parameters for ExportSavedObjectsDefault.
type SavedObjectExportRequestBody struct {
	// ExcludeExportDetails Do not add export details entry at the end of the stream.
	ExcludeExportDetails *bool `json:"excludeExportDetails,omitempty"`
	// IncludeReferencesDeep Includes all of the referenced objects in the exported objects.
	IncludeReferencesDeep *bool `json:"includeReferencesDeep,omitempty"`
	// Objects A list of objects to export.
	Objects []Object `json:"objects"`
	// Type The saved object types to include in the export. Use `*` to export all the types.
	Type json.RawMessage `json:"type,omitempty"`
}

// WithTypeString sets the Type field to a string value
func (body *SavedObjectExportRequestBody) WriteTypeString(value string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	body.Type = data
	return nil
}

// WithTypeList sets the Type field to a list of strings
func (body *SavedObjectExportRequestBody) WriteTypeList(values []string) error {
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	body.Type = data
	return nil
}

// newSavedObjectExport returns a function that performs POST /api/saved_objects/_export API requests
func (api *API) newSavedObjectExport() func(context.Context, *SavedObjectExportRequest, ...RequestOption) (*SavedObjectExportResponse, error) {
	return func(ctx context.Context, req *SavedObjectExportRequest, opts ...RequestOption) (*SavedObjectExportResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request body is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "saved_objects.export")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/saved_objects/_export"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "saved_objects.export")
			if reader := instrument.RecordRequestBody(ctx, "connector.post", httpReq.Body); reader != nil {
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
		resp := &SavedObjectExportResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		resp.RawBody = io.NopCloser(bytes.NewReader(bodyBytes))

		if httpResp.StatusCode == http.StatusOK {
			var objects []json.RawMessage
			lines := bytes.Split(bodyBytes, []byte("\n"))

			for _, line := range lines {
				if len(bytes.TrimSpace(line)) == 0 {
					continue
				}

				if bytes.HasPrefix(line, []byte(",")) {
					line = line[1:]
				}

				objects = append(objects, json.RawMessage(line))
			}

			resp.Body = objects
			return resp, nil
		} else {
			// For all non-200 responses
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

// WriteToFile writes the response body to the specified path in NDJSON format.
// See https://github.com/ndjson/ndjson-spec
func (d *SavedObjectExportResponse) WriteToFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	for _, obj := range d.Body {
		if _, err := file.Write(obj); err != nil {
			return fmt.Errorf("failed to write object to file: %w", err)
		}
		if _, err := file.Write([]byte("\n")); err != nil {
			return fmt.Errorf("failed to write newline to file: %w", err)
		}
	}

	return nil
}
