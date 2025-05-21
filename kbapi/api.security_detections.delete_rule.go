package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsDeleteRuleResponse wraps the response from a DeleteRule call
// To properly handle the response, use the UnmarshalRule function to convert the raw JSON into
// the appropriate rule type struct:
//
//	// Get update response
//	resp, _ := client.DeleteRule(ctx, ruleID, updateReq)
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
type SecurityDetectionsDeleteRuleResponse struct {
	StatusCode int
	Body       json.RawMessage
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsDeleteRuleRequest struct {
	Params SecurityDetectionsDeleteRuleRequestParams
}

type SecurityDetectionsDeleteRuleRequestParams struct {
	// The difference between the id and rule_id is that the id is a unique rule
	// identifier that is randomly generated when a rule is created and cannot be
	// set, whereas rule_id is a stable rule identifier that can be assigned during rule creation.

	// ID The rule's id value. Cannot be set if you are using RuleID
	ID *string
	// RuleID The rule's rule_id value. Cannot be set if you are using ID
	RuleID *string
}

// newSecurityDetectionsDeleteRule returns a function that performs DELETE /api/detection_engine/rules API requests
func (api *API) newSecurityDetectionsDeleteRule() func(context.Context, *SecurityDetectionsDeleteRuleRequest, ...RequestOption) (*SecurityDetectionsDeleteRuleResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsDeleteRuleRequest, opts ...RequestOption) (*SecurityDetectionsDeleteRuleResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.delete_rule")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ID != nil {
			params["id"] = *req.Params.ID
		}

		if req.Params.RuleID != nil {
			params["rule_id"] = *req.Params.RuleID
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_detections.delete_rule")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.delete_rule", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsDeleteRuleResponse{
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
