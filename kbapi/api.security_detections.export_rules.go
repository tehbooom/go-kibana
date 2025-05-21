package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// SecurityDetectionsExportRulesResponse wraps the response from a ExportRules call
type SecurityDetectionsExportRulesResponse struct {
	StatusCode int
	Body       []json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsExportRulesRequest struct {
	Params SecurityDetectionsExportRulesRequestParams
	Body   SecurityDetectionsExportRulesRequestBody
}

type SecurityDetectionsExportRulesRequestParams struct {
	// ExcludeExportDetails Determines whether a summary of the exported rules is returned.
	// Default value is false
	ExcludeExportDetails *bool
}

type SecurityDetectionsExportRulesRequestBody struct {
	Objects []SecurityDetectionsExportRulesRequestBodyItems `json:"objects"`
}

type SecurityDetectionsExportRulesRequestBodyItems struct {
	RuleID string `json:"rule_id"`
}

// newSecurityDetectionsExportRules returns a function that performs POST /api/detection_engine/rules/_export API requests
func (api *API) newSecurityDetectionsExportRules() func(context.Context, *SecurityDetectionsExportRulesRequest, ...RequestOption) (*SecurityDetectionsExportRulesResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsExportRulesRequest, opts ...RequestOption) (*SecurityDetectionsExportRulesResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.export_rules")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/_export"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ExcludeExportDetails != nil {
			params["exclude_export_details"] = strconv.FormatBool(*req.Params.ExcludeExportDetails)
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
			instrument.BeforeRequest(httpReq, "security_detections.export_rules")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.export_rules", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsExportRulesResponse{
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
func (d *SecurityDetectionsExportRulesResponse) WriteToFile(filepath string) error {
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
