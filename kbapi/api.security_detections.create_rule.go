package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsCreateRuleResponse wraps the response from a CreateRule call
// To properly handle the response, use the UnmarshalRule function to convert the raw JSON into
// the appropriate rule type struct:
//
//	// Get update response
//	resp, _ := client.CreateRule(ctx, ruleID, updateReq)
//
//	// Unmarshal the response into the appropriate rule type
//	rule, _ := UnmarshalRule(resp.Body)
//
//	// Now you can use the rule with type assertions to access specific fields
//	switch r := rule.(type) {
//	case *SecurityDetectionsEQLRuleResponse:
//	    fmt.Println("Updated EQL rule:", r.Name)
//	case *SecurityDetectionsQueryRuleResponse:
//	    fmt.Println("Updated Query rule:", r.Name)
//	// Handle other rule types...
//	}
//	// Or access common fields directly
//	fmt.Println("Rule ID:", rule.GetCommonFields().ID)
//	fmt.Println("Updated at:", rule.GetCommonFields().UpdatedAt)
type SecurityDetectionsCreateRuleResponse struct {
	StatusCode int
	Body       json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

// SecurityDetectionsCreateRuleRequest is used to create a new detection rule.
//
// For the Body field, provide a JSON-serialized instance of one of the following rule type structs
// that corresponds to the type of rule you're updating:
//
// - SecurityDetectionsESQLRuleResponse: For ESQL rules
// - SecurityDetectionsNewTermsRuleResponse: For New Terms rules
// - SecurityDetectionsMachineLearningRuleResponse: For Machine Learning rules
// - SecurityDetectionsThreatMatchRuleResponse: For Threat Match rules
// - SecurityDetectionsThresholdRuleResponse: For Threshold rules
// - SecurityDetectionsSavedQueryRuleResponse: For Saved Query rules
// - SecurityDetectionsQueryRuleResponse: For Query rules
// - SecurityDetectionsEQLRuleResponse: For EQL rules
// Alternatively, you can use the SetBody method to set the Body field directly from a rule struct:
//
//	req := SecurityDetectionsCreateRuleRequest{}
//	req.SetBody(eqlRule)
type SecurityDetectionsCreateRuleRequest struct {
	Body json.RawMessage
}

// SetBody sets the Body field using a rule struct. The provided rule must be one of the supported
// rule types (e.g., SecurityDetectionsEQLRuleResponse, SecurityDetectionsQueryRuleResponse, etc.)
func (r *SecurityDetectionsCreateRuleRequest) SetBody(rule interface{}) error {
	data, err := json.Marshal(rule)
	if err != nil {
		return fmt.Errorf("failed to marshal rule: %w", err)
	}
	r.Body = data
	return nil
}

// newSecurityDetectionsCreateRule returns a function that performs POST /api/detection_engine/rules API requests
func (api *API) newSecurityDetectionsCreateRule() func(context.Context, *SecurityDetectionsCreateRuleRequest, ...RequestOption) (*SecurityDetectionsCreateRuleResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsCreateRuleRequest, opts ...RequestOption) (*SecurityDetectionsCreateRuleResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.create_rule")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "security_detections.create_rule")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.create_rule", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsCreateRuleResponse{
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
