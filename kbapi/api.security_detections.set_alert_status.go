package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsSetAlertStatusResponse wraps the response from a SetAlertStatus call
type SecurityDetectionsSetAlertStatusResponse struct {
	StatusCode int
	Body       *SecurityDetectionsSetAlertStatusResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsSetAlertStatusResponseBody struct {
	Took                 string         `json:"took"`
	Noops                string         `json:"noops"`
	Total                string         `json:"total"`
	Batches              string         `json:"batches"`
	Deleted              string         `json:"deleted"`
	Retries              map[string]any `json:"retries"`
	Updated              string         `json:"updated"`
	Failures             []any          `json:"failures"`
	TimedOut             string         `json:"timed_out"`
	ThrottledMillis      string         `json:"throttled_millis"`
	VersionConflicts     string         `json:"version_conflicts"`
	RequestsPerSecond    string         `json:"requests_per_second"`
	ThrottledUntilMillis string         `json:"throttled_until_millis"`
}

// SecurityDetectionsSetAlertStatusRequest to set the body provied
// For the Body field, provide a JSON-serialized instance of one of the following  structs
// - SecurityDetectionsSetAlertStatusRequestQuery
// - SecurityDetectionsSetAlertStatusRequestIDs
// Alternatively, you can use the SetBody method to set the Body field directly from a rule struct:
// req := SecurityDetectionsSetAlertStatusRequest{}
// req.SetBody(eqlRule)
type SecurityDetectionsSetAlertStatusRequest struct {
	Body json.RawMessage
}

// SetBody sets the Body field
func (r *SecurityDetectionsSetAlertStatusRequest) SetBody(body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal rule: %w", err)
	}
	r.Body = data
	return nil
}

// newSecurityDetectionsSetAlertStatus returns a function that performs POST /api/detection_engine/signals/status API requests
func (api *API) newSecurityDetectionsSetAlertStatus() func(context.Context, *SecurityDetectionsSetAlertStatusRequest, ...RequestOption) (*SecurityDetectionsSetAlertStatusResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsSetAlertStatusRequest, opts ...RequestOption) (*SecurityDetectionsSetAlertStatusResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.set_alert_status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/signals/status"

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
			instrument.BeforeRequest(httpReq, "security_detections.set_alert_status")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.set_alert_status", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsSetAlertStatusResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsSetAlertStatusResponseBody

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
