package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsGetStatusPrebuiltResponse wraps the response from a GetStatusPrebuilt call
type SecurityDetectionsGetStatusPrebuiltResponse struct {
	StatusCode int
	Body       *SecurityDetectionsGetStatusPrebuiltResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsGetStatusPrebuiltResponseBody struct {
	// RulesCustomInstalled The total number of custom rules
	RulesCustomInstalled int `json:"rules_custom_installed"`
	// RulesInstalled The total number of installed prebuilt rules
	RulesInstalled int `json:"rules_installed"`
	// RulesNotInstalled The total number of available prebuilt rules that are not installed
	RulesNotInstalled int `json:"rules_not_installed"`
	// RulesNotUpdated The total number of outdated prebuilt rules
	RulesNotUpdated int `json:"rules_not_updated"`
	// TimelinesInstalled The total number of installed prebuilt timelines
	TimelinesInstalled int `json:"timelines_installed"`
	// TimelinesNotInstalled The total number of available prebuilt timelines that are not installed
	TimelinesNotInstalled int `json:"timelines_not_installed"`
	// TimelinesNotUpdated The total number of outdated prebuilt timelines
	TimelinesNotUpdated int `json:"timelines_not_updated"`
}

// newSecurityDetectionsGetStatusPrebuilt returns a function that performs GET /api/detection_engine/rules/prepackaged/_status API requests
func (api *API) newSecurityDetectionsGetStatusPrebuilt() func(context.Context, ...RequestOption) (*SecurityDetectionsGetStatusPrebuiltResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsGetStatusPrebuiltResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_detections.get_status_prebuilt")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/prepackaged/_status"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_detections.get_status_prebuilt")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.get_status_prebuilt", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsGetStatusPrebuiltResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsGetStatusPrebuiltResponseBody

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
