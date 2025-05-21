package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// SavedObjectRotateKeyResponse wraps the response from a <todo> call
type SavedObjectRotateKeyResponse struct {
	StatusCode int
	Body       *SavedObjectRotateKeyResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SavedObjectRotateKeyResponseBody struct {
	Total      int `json:"total"`
	Failed     int `json:"failed"`
	Successful int `json:"successful"`
}

type SavedObjectRotateKeyRequest struct {
	Params SavedObjectRotateKeyRequestParams
}

type SavedObjectRotateKeyRequestParams struct {
	// BatchSize Specifies a maximum number of saved objects that Kibana can process in a single batch.
	// Bulk key rotation is an iterative process since Kibana may not be able to fetch and process all required
	// saved objects in one go and splits processing into consequent batches.
	// By default, the batch size is 10000, which is also a maximum allowed value.
	BatchSize *int `form:"batch_size,omitempty" json:"batch_size,omitempty"`

	// Type Limits encryption key rotation only to the saved objects with the specified type.
	// By default, Kibana tries to rotate the encryption key for all saved object types that may contain encrypted attributes.
	Type *string `form:"type,omitempty" json:"type,omitempty"`
}

// newSavedObjectRotateKey returns a function that performs POST /api/encrypted_saved_objects/_rotate_key API requests
func (api *API) newSavedObjectRotateKey() func(context.Context, *SavedObjectRotateKeyRequest, ...RequestOption) (*SavedObjectRotateKeyResponse, error) {
	return func(ctx context.Context, req *SavedObjectRotateKeyRequest, opts ...RequestOption) (*SavedObjectRotateKeyResponse, error) {
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
			newCtx = instrument.Start(ctx, "saved_objects.rotate_key")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/encrypted_saved_objects/_rotate_key"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.BatchSize != nil {
			params["batch_size"] = strconv.Itoa(*req.Params.BatchSize)
		}
		if req.Params.Type != nil {
			params["type"] = *req.Params.Type
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
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
			instrument.BeforeRequest(httpReq, "saved_objects.rotate_key")
			if reader := instrument.RecordRequestBody(ctx, "saved_objects.rotate_key", httpReq.Body); reader != nil {
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
		resp := &SavedObjectRotateKeyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SavedObjectRotateKeyResponseBody

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
