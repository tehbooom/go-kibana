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
// DataViewsCreateUpdateRuntimeFieldResponse wraps the response from a <todo> call
type DataViewsCreateUpdateRuntimeFieldResponse struct {
	StatusCode int
	Body       *DataViewsCreateUpdateRuntimeFieldResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsCreateUpdateRuntimeFieldResponseBody struct {
	Fields   []DataViewsFields `json:"fields"`
	DataView DataViewsObject   `json:"data_view"`
}

type DataViewsCreateUpdateRuntimeFieldRequest struct {
	ID   string
	Body DataViewsCreateUpdateRuntimeFieldRequestBody
}

type DataViewsCreateUpdateRuntimeFieldRequestBody struct {
	Name         string                   `json:"name"`
	RuntimeField DataViewsRuntimeFieldMap `json:"runtimeField"`
}

// newDataViewsCreateUpdateRuntimeField returns a function that performs PUT /api/data_views/data_view/{viewId}/runtime_field API requests
func (api *API) newDataViewsCreateUpdateRuntimeField() func(context.Context, *DataViewsCreateUpdateRuntimeFieldRequest, ...RequestOption) (*DataViewsCreateUpdateRuntimeFieldResponse, error) {
	return func(ctx context.Context, req *DataViewsCreateUpdateRuntimeFieldRequest, opts ...RequestOption) (*DataViewsCreateUpdateRuntimeFieldResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.create_update_runtime_field")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/data_views/data_view/%s/runtime_field", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "dataviews.create_update_runtime_field")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.create_update_runtime_field", httpReq.Body); reader != nil {
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
		resp := &DataViewsCreateUpdateRuntimeFieldResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result DataViewsCreateUpdateRuntimeFieldResponseBody

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
