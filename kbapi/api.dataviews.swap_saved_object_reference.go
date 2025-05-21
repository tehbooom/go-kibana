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
// DataViewsSwapSavedObjectReferenceResponse wraps the response from a <todo> call
type DataViewsSwapSavedObjectReferenceResponse struct {
	StatusCode int
	Body       *DataViewsSwapSavedObjectReferenceResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type DataViewsSwapSavedObjectReferenceResponseBody struct {
	DeleteStatus *struct {
		DeletePerformed *bool `json:"deletePerformed,omitempty"`
		RemainingRefs   *int  `json:"remainingRefs,omitempty"`
	} `json:"deleteStatus,omitempty"`
	Result *[]struct {
		// ID A saved object identifier.
		ID *string `json:"id,omitempty"`
		// Type The saved object type.
		Type *string `json:"type,omitempty"`
	} `json:"result,omitempty"`
}

type DataViewsSwapSavedObjectReferenceRequest struct {
	Body DataViewsSwapSavedObjectReferenceRequestBody
}

type DataViewsSwapSavedObjectReferenceRequestBody struct {
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

// newDataViewsSwapSavedObjectReference returns a function that performs POST /api/data_views/swap_references API requests
func (api *API) newDataViewsSwapSavedObjectReference() func(context.Context, *DataViewsSwapSavedObjectReferenceRequest, ...RequestOption) (*DataViewsSwapSavedObjectReferenceResponse, error) {
	return func(ctx context.Context, req *DataViewsSwapSavedObjectReferenceRequest, opts ...RequestOption) (*DataViewsSwapSavedObjectReferenceResponse, error) {
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
			newCtx = instrument.Start(ctx, "dataviews.swap_saved_object_reference")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/data_views/swap_references"

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
			instrument.BeforeRequest(httpReq, "dataviews.swap_saved_object_reference")
			if reader := instrument.RecordRequestBody(ctx, "dataviews.swap_saved_object_reference", httpReq.Body); reader != nil {
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
		resp := &DataViewsSwapSavedObjectReferenceResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result DataViewsSwapSavedObjectReferenceResponseBody

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
