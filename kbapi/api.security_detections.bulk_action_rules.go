package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// SecurityDetectionsBulkActionRulesResponse wraps the response from a BulkActionRules call
// To properly handle the response, use the UnmarshalBulkAction function to convert the raw JSON
// to the SecurityDetectionsBulkActionEditResponse struct. If you exported the rules convert the
// Body to a string
type SecurityDetectionsBulkActionRulesResponse struct {
	StatusCode int
	Body       json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

// SecurityDetectionsBulkActionRulesRequest is used to execute bulk actions on rules.
// For the Body field, provide a JSON-serialized instance of one of the following action structs:
// - SecurityDetectionsBulkActionRulesDelete
// - SecurityDetectionsBulkActionRulesDisable
// - SecurityDetectionsBulkActionRulesEnable
// - SecurityDetectionsBulkActionRulesExport
// - SecurityDetectionsBulkActionRulesDuplicate
// - SecurityDetectionsBulkActionRulesRun
// - SecurityDetectionsBulkActionRulesEdit
// Alternatively, you can use the SetBody method to set the Body field directly from an action struct:
//
//	req := SecurityDetectionsCreateRuleRequest{}
//	req.SetBody(editAction)
type SecurityDetectionsBulkActionRulesRequest struct {
	Params SecurityDetectionsBulkActionRulesRequestParams
	Body   json.RawMessage
}

func (r *SecurityDetectionsBulkActionRulesRequest) SetBody(action interface{}) error {
	data, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal rule: %w", err)
	}
	r.Body = data
	return nil
}

type SecurityDetectionsBulkActionRulesRequestParams struct {
	// DryRun Enables dry run mode for the request call.
	// Enable dry run mode to verify that bulk actions can be applied to specified rules.
	// Certain rules, such as prebuilt Elastic rules on a Basic subscription, can’t be edited and will return errors in the request response.
	// Error details will contain an explanation, the rule name and/or ID, and additional troubleshooting information.
	DryRun *bool
}

// newSecurityDetectionsBulkActionRules returns a function that performs POST /api/detection_engine/rules/_bulk_action API requests
func (api *API) newSecurityDetectionsBulkActionRules() func(context.Context, *SecurityDetectionsBulkActionRulesRequest, ...RequestOption) (*SecurityDetectionsBulkActionRulesResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsBulkActionRulesRequest, opts ...RequestOption) (*SecurityDetectionsBulkActionRulesResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.bulk_action_rules")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/_bulk_action"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.DryRun != nil {
			params["dry_run"] = strconv.FormatBool(*req.Params.DryRun)
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
			instrument.BeforeRequest(httpReq, "security_detections.bulk_action_rules")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.bulk_action_rules", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsBulkActionRulesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()

		if httpResp.StatusCode < 299 {
			resp.Body = bodyBytes
			return resp, nil
		} else {
			// For all non-success responses
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
