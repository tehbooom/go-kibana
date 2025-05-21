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
// APMSourcemapsGetResponse wraps the response from a <todo> call
type APMSourcemapsGetResponse struct {
	StatusCode int
	Body       *APMSourcemapsGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMSourcemapsGetResponseBody struct {
	Artifacts *[]struct {
		Body *struct {
			BundleFilepath *string `json:"bundleFilepath,omitempty"`
			ServiceName    *string `json:"serviceName,omitempty"`
			ServiceVersion *string `json:"serviceVersion,omitempty"`
			SourceMap      *struct {
				File           *string   `json:"file,omitempty"`
				Mappings       *string   `json:"mappings,omitempty"`
				SourceRoot     *string   `json:"sourceRoot,omitempty"`
				Sources        *[]string `json:"sources,omitempty"`
				SourcesContent *[]string `json:"sourcesContent,omitempty"`
				Version        *float32  `json:"version,omitempty"`
			} `json:"sourceMap,omitempty"`
		} `json:"body,omitempty"`
		CompressionAlgorithm *string  `json:"compressionAlgorithm,omitempty"`
		Created              *string  `json:"created,omitempty"`
		DecodedSha256        *string  `json:"decodedSha256,omitempty"`
		DecodedSize          *float32 `json:"decodedSize,omitempty"`
		EncodedSha256        *string  `json:"encodedSha256,omitempty"`
		EncodedSize          *float32 `json:"encodedSize,omitempty"`
		EncryptionAlgorithm  *string  `json:"encryptionAlgorithm,omitempty"`
		ID                   *string  `json:"id,omitempty"`
		Identifier           *string  `json:"identifier,omitempty"`
		PackageName          *string  `json:"packageName,omitempty"`
		RelativeUrl          *string  `json:"relative_url,omitempty"`
		Type                 *string  `json:"type,omitempty"`
	} `json:"artifacts,omitempty"`
}

type APMSourcemapsGetRequest struct {
	Params APMSourcemapsGetRequestParams
}

type APMSourcemapsGetRequestParams struct {
	Page    *int `form:"page,omitempty" json:"page,omitempty"`
	PerPage *int `form:"perPage,omitempty" json:"perPage,omitempty"`
}

// newAPMSourcemapsGet returns a function that performs GET /api/apm/sourcemaps API requests
func (api *API) newAPMSourcemapsGet() func(context.Context, *APMSourcemapsGetRequest, ...RequestOption) (*APMSourcemapsGetResponse, error) {
	return func(ctx context.Context, req *APMSourcemapsGetRequest, opts ...RequestOption) (*APMSourcemapsGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.sourcemaps.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/apm/sourcemaps"

		// Build query parameters
		params := make(map[string]string)
		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PerPage != nil {
			params["per_page"] = strconv.Itoa(*req.Params.PerPage)
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
			instrument.BeforeRequest(httpReq, "apm.sourcemaps.get")
			if reader := instrument.RecordRequestBody(ctx, "apm.sourcemaps.get", httpReq.Body); reader != nil {
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
		resp := &APMSourcemapsGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMSourcemapsGetResponseBody

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
