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

// SecurityExceptionsExportListResponse wraps the response from a ExportList call
type SecurityExceptionsExportListResponse struct {
	StatusCode int
	Body       []json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityExceptionsExportListRequest struct {
	Params SecurityExceptionsExportListRequestParams
}

type SecurityExceptionsExportListRequestParams struct {
	// ID Exception list's identifier.
	// Either id or list_id must be specified.
	ID *string
	// ListID Human readable exception list string identifier, e.g. trusted-linux-processes.
	// Either id or list_id must be specified.
	ListID *string
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string
	// IncludeExpiredExceptions Determines whether to include expired exceptions in the duplicated list.
	// Expiration date defined by expire_time.
	// Values are true or false. Default value is true.
	IncludeExpiredExceptions *string
}

// newSecurityExceptionsExportList returns a function that performs POST /api/ API requests
func (api *API) newSecurityExceptionsExportList() func(context.Context, *SecurityExceptionsExportListRequest, ...RequestOption) (*SecurityExceptionsExportListResponse, error) {
	return func(ctx context.Context, req *SecurityExceptionsExportListRequest, opts ...RequestOption) (*SecurityExceptionsExportListResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_exceptions.export_list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.IncludeExpiredExceptions != nil {
			params["include_expired_exceptions"] = *req.Params.IncludeExpiredExceptions
		}
		if req.Params.ListID != nil {
			params["list_id"] = *req.Params.ListID
		}
		if req.Params.NamespaceType != nil {
			params["namespace_type"] = *req.Params.NamespaceType
		}
		if req.Params.ID != nil {
			params["id"] = *req.Params.ID
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_exceptions.export_list")
			if reader := instrument.RecordRequestBody(ctx, "security_exceptions.export_list", httpReq.Body); reader != nil {
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
		resp := &SecurityExceptionsExportListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}

		resp.RawBody = io.NopCloser(bytes.NewReader(bodyBytes))

		if httpResp.StatusCode < 299 {
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
			// For all non-success responses
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

// WriteToFile writes the response body to the specified path in NDJSON format.
// See https://github.com/ndjson/ndjson-spec
func (d *SecurityExceptionsExportListResponse) WriteToFile(filepath string) error {
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
