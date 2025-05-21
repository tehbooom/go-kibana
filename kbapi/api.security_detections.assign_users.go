package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecurityDetectionsAssignUsersResponse wraps the response from a AssignUsers call
type SecurityDetectionsAssignUsersResponse struct {
	StatusCode int
	Body       *SecurityDetectionsAssignUsersResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsAssignUsersResponseBody struct {
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

type SecurityDetectionsAssignUsersRequest struct {
	Body SecurityDetectionsAssignUsersRequestBody
}

type SecurityDetectionsAssignUsersRequestBody struct {
	Assignees SecurityDetectionsAssignUsersRequestAssignees `json:"assignees"`
	IDs       []string                                      `json:"ids"`
}

type SecurityDetectionsAssignUsersRequestAssignees struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

// newSecurityDetectionsAssignUsers returns a function that performs POST /api/detection_engine/signals/assignees API requests
func (api *API) newSecurityDetectionsAssignUsers() func(context.Context, *SecurityDetectionsAssignUsersRequest, ...RequestOption) (*SecurityDetectionsAssignUsersResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsAssignUsersRequest, opts ...RequestOption) (*SecurityDetectionsAssignUsersResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.assign_users")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/signals/assignees"

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
			instrument.BeforeRequest(httpReq, "security_detections.assign_users")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.assign_users", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsAssignUsersResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsAssignUsersResponseBody

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
