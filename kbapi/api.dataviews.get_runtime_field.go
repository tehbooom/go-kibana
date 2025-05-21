package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// DataViewsGetRuntimeFieldResponse wraps the response from a <todo> call
type DataViewsGetRuntimeFieldResponse struct {
	StatusCode int
	Body       *DataViewsGetRuntimeFieldResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsGetRuntimeFieldResponseBody struct {
	Fields   []DataViewsFields `json:"fields"`
	DataView DataViewsObject   `json:"data_view"`
}

type DataViewsGetRuntimeFieldRequest struct {
	// ID an identifier for the data view.
	ID string
	// FieldName the name of the runtime field.
	FieldName string
}

// newDataViewsGetRuntimeField returns a function that performs GET /api/data_views/data_view/{viewId}/runtime_field/{fieldName} API requests
func (api *API) newDataViewsGetRuntimeField() func(context.Context, *DataViewsGetRuntimeFieldRequest, ...RequestOption) (*DataViewsGetRuntimeFieldResponse, error) {
	return func(ctx context.Context, req *DataViewsGetRuntimeFieldRequest, opts ...RequestOption) (*DataViewsGetRuntimeFieldResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.get_runtime_field")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/data_views/data_view/%s/runtime_field/%s", req.ID, req.FieldName)

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
			instrument.BeforeRequest(httpReq, "dataviews.get_runtime_field")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.get_runtime_field", httpReq.Body); reader != nil {
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
		resp := &DataViewsGetRuntimeFieldResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result DataViewsGetRuntimeFieldResponseBody

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
