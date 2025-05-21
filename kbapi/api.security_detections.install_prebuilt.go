package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsInstallPrebuiltResponse wraps the response from a InstallPrebuilt call
type SecurityDetectionsInstallPrebuiltResponse struct {
	StatusCode int
	Body       *SecurityDetectionsInstallPrebuiltResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsInstallPrebuiltResponseBody struct {
	// RulesInstalled The number of rules installed
	RulesInstalled int `json:"rules_installed"`
	// RulesUpdated The number of rules updated
	RulesUpdated int `json:"rules_updated"`
	// TimelinesInstalled The number of timelines installed
	TimelinesInstalled int `json:"timelines_installed"`
	// TimelinesUpdated The number of timelines updated
	TimelinesUpdated int `json:"timelines_updated"`
}

// newSecurityDetectionsInstallPrebuilt returns a function that performs PUT /api/detection_engine/rules/prepackaged API requests
func (api *API) newSecurityDetectionsInstallPrebuilt() func(context.Context, ...RequestOption) (*SecurityDetectionsInstallPrebuiltResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsInstallPrebuiltResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_detections.install_prebuilt")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/prepackaged"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "security_detections.install_prebuilt")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.install_prebuilt", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsInstallPrebuiltResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsInstallPrebuiltResponseBody

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
