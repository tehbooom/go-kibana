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
// DataViewsUpdateRuntimeFieldResponse wraps the response from a <todo> call
type DataViewsUpdateRuntimeFieldResponse struct {
	StatusCode int
	Body       *DataViewsUpdateRuntimeFieldResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsUpdateRuntimeFieldResponseBody struct{}

type DataViewsUpdateRuntimeFieldRequest struct {
	// ID an identifier for the data view.
	ID string
	// FieldName the name of the runtime field.
	FieldName string
	Body      DataViewsUpdateRuntimeFieldRequestBody
}

type DataViewsUpdateRuntimeFieldRequestBody struct {
	RuntimeField DataViewsRuntimeFieldMap `json:"runtimeField"`
}

// newDataViewsUpdateRuntimeField returns a function that performs POST /api/data_views/data_view/{viewId}/runtime_field/{fieldName} API requests
func (api *API) newDataViewsUpdateRuntimeField() func(context.Context, *DataViewsUpdateRuntimeFieldRequest, ...RequestOption) (*DataViewsUpdateRuntimeFieldResponse, error) {
	return func(ctx context.Context, req *DataViewsUpdateRuntimeFieldRequest, opts ...RequestOption) (*DataViewsUpdateRuntimeFieldResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.update_runtime_field")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/data_views/data_view/%s/runtime_field/%s", req.ID, req.FieldName)

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
			instrument.BeforeRequest(httpReq, "dataviews.update_runtime_field")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.update_runtime_field", httpReq.Body); reader != nil {
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
		resp := &DataViewsUpdateRuntimeFieldResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result DataViewsUpdateRuntimeFieldResponseBody

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
