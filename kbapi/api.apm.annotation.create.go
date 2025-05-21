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
// APMAnnotationCreateResponse wraps the response from a <todo> call
type APMAnnotationCreateResponse struct {
	StatusCode int
	Body       *APMAnnotationCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAnnotationCreateResponseBody struct {
	ID     *string `json:"_id,omitempty"`
	Index  *string `json:"_index,omitempty"`
	Source *struct {
		Timestamp  *string `json:"@timestamp,omitempty"`
		Annotation *struct {
			Title *string `json:"title,omitempty"`
			Type  *string `json:"type,omitempty"`
		} `json:"annotation,omitempty"`
		Event *struct {
			Created *string `json:"created,omitempty"`
		} `json:"event,omitempty"`
		Message *string `json:"message,omitempty"`
		Service *struct {
			Environment *string `json:"environment,omitempty"`
			Name        *string `json:"name,omitempty"`
			Version     *string `json:"version,omitempty"`
		} `json:"service,omitempty"`
		Tags *[]string `json:"tags,omitempty"`
	} `json:"_source,omitempty"`
}

type APMAnnotationCreateRequest struct {
	ServiceName string
	Body        APMAnnotationCreateRequestBody
}

type APMAnnotationCreateRequestBody struct {
	// AtTimestamp The date and time of the annotation. It must be in ISO 8601 format.
	AtTimestamp string `json:"@timestamp"`
	// Message The message displayed in the annotation. It defaults to service.version.
	Message *string `json:"message,omitempty"`
	// Service The service that identifies the configuration to create or update.
	Service struct {
		Environment *string `json:"environment,omitempty"`
		Version     string  `json:"version"`
	} `json:"service"`
	// Tags tags are used by the Applications UI to distinguish APM annotations from other annotations.
	// Tags may have additional functionality in future releases. It defaults to [apm].
	// While you can add additional tags, you cannot remove the apm tag.
	Tags *[]string `json:"tags,omitempty"`
}

// newAPMAnnotationCreate returns a function that performs POST /api/apm/services/{serviceName}/annotation API requests
func (api *API) newAPMAnnotationCreate() func(context.Context, *APMAnnotationCreateRequest, ...RequestOption) (*APMAnnotationCreateResponse, error) {
	return func(ctx context.Context, req *APMAnnotationCreateRequest, opts ...RequestOption) (*APMAnnotationCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.annotation.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/apm/services/%s/annotation", req.ServiceName)

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
			instrument.BeforeRequest(httpReq, "apm.annotation.create")
			if reader := instrument.RecordRequestBody(ctx, "apm.annotation.create", httpReq.Body); reader != nil {
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
		resp := &APMAnnotationCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAnnotationCreateResponseBody

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
