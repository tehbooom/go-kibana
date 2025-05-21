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
// DataViewsPreviewSavedObjectSwapResponse wraps the response from a <todo> call
type DataViewsPreviewSavedObjectSwapResponse struct {
	StatusCode int
	Body       *DataViewsPreviewSavedObjectSwapResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsPreviewSavedObjectSwapResponseBody struct {
	Result []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"result"`
}

type DataViewsPreviewSavedObjectSwapRequest struct {
	Body DataViewsPreviewSavedObjectSwapRequestBody
}

type DataViewsPreviewSavedObjectSwapRequestBody struct {
	// Delete deletes referenced saved object if all references are removed.
	Delete *bool `json:"delete,omitempty"`

	// ForID limit the affected saved objects to one or more by identifier.
	ForID *[]string `json:"forId,omitempty"`

	// ForType limit the affected saved objects by type.
	ForType *string `json:"forType,omitempty"`

	// FromID the saved object reference to change.
	FromID string `json:"fromId"`

	// FromType specify the type of the saved object reference to alter. The default value is `index-pattern` for data views.
	FromType *string `json:"fromType,omitempty"`

	// ToID new saved object reference value to replace the old value.
	ToID string `json:"toId"`
}

// newDataViewsPreviewSavedObjectSwap returns a function that performs POST /api/data_views/swap_references/_preview API requests
func (api *API) newDataViewsPreviewSavedObjectSwap() func(context.Context, *DataViewsPreviewSavedObjectSwapRequest, ...RequestOption) (*DataViewsPreviewSavedObjectSwapResponse, error) {
	return func(ctx context.Context, req *DataViewsPreviewSavedObjectSwapRequest, opts ...RequestOption) (*DataViewsPreviewSavedObjectSwapResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.preview_saved_object_swap")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/data_views/swap_references/_preview"

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
			instrument.BeforeRequest(httpReq, "dataviews.preview_saved_object_swap")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.preview_saved_object_swap", httpReq.Body); reader != nil {
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
		resp := &DataViewsPreviewSavedObjectSwapResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result DataViewsPreviewSavedObjectSwapResponseBody

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
