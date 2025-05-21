package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// EndpointExceptionsGetResponse wraps the response from a <todo> call
type EndpointExceptionsGetResponse struct {
	StatusCode int
	Body       *EndpointExceptionsGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type EndpointExceptionsGetResponseBody = []EndpointExceptionsListItem

type EndpointExceptionsGetRequest struct {
	Params EndpointExceptionsGetRequestParams
}

type EndpointExceptionsGetRequestParams struct {
	// ID Either `id` or `item_id` must be specified
	ID *string `form:"id,omitempty" json:"id,omitempty"`

	// ItemID Either `id` or `item_id` must be specified
	ItemID *string `form:"item_id,omitempty" json:"item_id,omitempty"`
}

// newEndpointExceptionsGet returns a function that performs GET /api/endpoint_list/items API requests
func (api *API) newEndpointExceptionsGet() func(context.Context, *EndpointExceptionsGetRequest, ...RequestOption) (*EndpointExceptionsGetResponse, error) {
	return func(ctx context.Context, req *EndpointExceptionsGetRequest, opts ...RequestOption) (*EndpointExceptionsGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "endpoint.exceptions.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/endpoint_list/items"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.ID != nil {
			params["id"] = *req.Params.ID
		}
		if req.Params.ItemID != nil {
			params["item_id"] = *req.Params.ItemID
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "endpoint.exceptions.get")
			if reader := instrument.RecordRequestBody(ctx, "endpoint.exceptions.get", httpReq.Body); reader != nil {
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
		resp := &EndpointExceptionsGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result EndpointExceptionsGetResponseBody

		if httpResp.StatusCode == 200 {
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
			// For all non-200 responses
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
