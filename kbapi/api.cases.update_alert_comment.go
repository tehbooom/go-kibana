package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// CasesUpdateCommentAlertResponse wraps the response from a <todo> call
type CasesUpdateCommentAlertResponse struct {
	StatusCode int
	Body       *CasesUpdateCommentAlertResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesUpdateCommentAlertResponseBody struct{}

type CasesUpdateCommentAlertRequest struct {
	ID   string
	Body json.RawMessage
}

func (req *CasesUpdateCommentAlertRequest) SetAlertBody(body AlertCommentRequest) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req.Body = data
	return nil
}

func (req *CasesUpdateCommentAlertRequest) SetUserBody(body UserCommentRequest) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req.Body = data
	return nil
}

// newCasesUpdateCommentAlert returns a function that performs PATCH /api/cases/{caseId}/comments API requests
func (api *API) newCasesUpdateCommentAlert() func(context.Context, *CasesUpdateCommentAlertRequest, ...RequestOption) (*CasesUpdateCommentAlertResponse, error) {
	return func(ctx context.Context, req *CasesUpdateCommentAlertRequest, opts ...RequestOption) (*CasesUpdateCommentAlertResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.update_alert_comment")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/cases/%s/comments", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPatch, path, nil)
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
			instrument.BeforeRequest(httpReq, "cases.update_alert_comment")
			if reader := instrument.RecordRequestBody(ctx, "cases.update_alert_comment", httpReq.Body); reader != nil {
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
		resp := &CasesUpdateCommentAlertResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesUpdateCommentAlertResponseBody

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
