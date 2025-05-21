package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// CasesGetAlertCommentResponse wraps the response from a <todo> call
type CasesGetAlertCommentResponse struct {
	StatusCode int
	Body       *CasesGetAlertCommentResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesGetAlertCommentResponseBody struct {
	Comment json.RawMessage
}

func (resp *CasesGetAlertCommentResponseBody) GetTypedComment() (any, error) {
	var typeContainer struct {
		Type string `json:"type"`
	}

	err := json.Unmarshal(resp.Comment, &typeContainer)
	if err != nil {
		return nil, err
	}

	switch typeContainer.Type {
	case "user":
		var comment UserCommentResponse
		err := json.Unmarshal(resp.Comment, &comment)
		if err != nil {
			return nil, err
		}
		return comment, nil
	case "alert":
		var comment AlertCommentResponse
		err := json.Unmarshal(resp.Comment, &comment)
		if err != nil {
			return nil, err
		}
		return comment, nil
	default:
		var comment BaseComment
		err := json.Unmarshal(resp.Comment, &comment)
		if err != nil {
			return nil, err
		}
		return comment, nil
	}
}

type CasesGetAlertCommentRequest struct {
	CaseID    string
	CommentID string
}

// newCasesGetAlertComment returns a function that performs GET /api/cases/{caseId}/comments/{commentId} API requests
func (api *API) newCasesGetAlertComment() func(context.Context, *CasesGetAlertCommentRequest, ...RequestOption) (*CasesGetAlertCommentResponse, error) {
	return func(ctx context.Context, req *CasesGetAlertCommentRequest, opts ...RequestOption) (*CasesGetAlertCommentResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.get_alert_comment")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/cases/%s/comments/%s", req.CaseID, req.CommentID)

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
			instrument.BeforeRequest(httpReq, "cases.get_alert_comment")
			if reader := instrument.RecordRequestBody(ctx, "cases.get_alert_comment", httpReq.Body); reader != nil {
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
		resp := &CasesGetAlertCommentResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
			Body:       &CasesGetAlertCommentResponseBody{},
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()

		if httpResp.StatusCode < 299 {
			resp.Body.Comment = bodyBytes
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
