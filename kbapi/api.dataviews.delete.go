package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// DataViewsDeleteResponse wraps the response from a <todo> call
type DataViewsDeleteResponse struct {
	StatusCode int
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsDeleteRequest struct {
	ID string
}

// newDataViewsDelete returns a function that performs DELETE /api/data_views/data_view/{viewId} API requests
func (api *API) newDataViewsDelete() func(context.Context, *DataViewsDeleteRequest, ...RequestOption) (*DataViewsDeleteResponse, error) {
	return func(ctx context.Context, req *DataViewsDeleteRequest, opts ...RequestOption) (*DataViewsDeleteResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.delete")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/data_views/data_view/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
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
			instrument.BeforeRequest(httpReq, "dataviews.delete")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.delete", httpReq.Body); reader != nil {
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
		resp := &DataViewsDeleteResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		if httpResp.StatusCode < 299 {
			httpResp.Body.Close()
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
