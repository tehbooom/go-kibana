package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// APMAnnotationSearchResponse wraps the response from a <todo> call
type APMAnnotationSearchResponse struct {
	StatusCode int
	Body       *APMAnnotationSearchResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAnnotationSearchResponseBody struct {
	// Annotations Annotations
	Annotations *[]struct {
		Timestamp *float32 `json:"@timestamp,omitempty"`
		ID        *string  `json:"id,omitempty"`
		Text      *string  `json:"text,omitempty"`
		Type      *string  `json:"type,omitempty"`
	} `json:"annotations,omitempty"`
}

type APMAnnotationSearchRequest struct {
	ServiceName string
	Params      APMAnnotationSearchRequestParams
}

type APMAnnotationSearchRequestParams struct {
	// Environment The environment to filter annotations by
	Environment *string `form:"environment,omitempty" json:"environment,omitempty"`
	// Start The start date for the search
	Start *string `form:"start,omitempty" json:"start,omitempty"`
	// End The end date for the search
	End *string `form:"end,omitempty" json:"end,omitempty"`
}

// newAPMAnnotationSearch returns a function that performs GET /api/apm/services/{serviceName}/annotation/search API requests
func (api *API) newAPMAnnotationSearch() func(context.Context, *APMAnnotationSearchRequest, ...RequestOption) (*APMAnnotationSearchResponse, error) {
	return func(ctx context.Context, req *APMAnnotationSearchRequest, opts ...RequestOption) (*APMAnnotationSearchResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.annotation.search")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/apm/services/%s/annotation/search", req.ServiceName)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Environment != nil {
			params["environment"] = *req.Params.Environment
		}
		if req.Params.Start != nil {
			params["start"] = *req.Params.Start
		}
		if req.Params.End != nil {
			params["end"] = *req.Params.End
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
			instrument.BeforeRequest(httpReq, "apm.annotation.search")
			if reader := instrument.RecordRequestBody(ctx, "apm.annotation.search", httpReq.Body); reader != nil {
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
		resp := &APMAnnotationSearchResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAnnotationSearchResponseBody

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
