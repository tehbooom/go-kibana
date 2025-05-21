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
// ShortURLCreateResponse wraps the response from a <todo> call
type ShortURLCreateResponse struct {
	StatusCode int
	Body       *ShortURLCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type ShortURLCreateResponseBody struct {
	AccessCount int             `json:"accessCount"`
	AccessDate  string          `json:"accessDate"`
	CreateDate  string          `json:"createDate"`
	ID          string          `json:"id"`
	Locator     ShortURLLocator `json:"locator"`
	Slug        string          `json:"slug"`
}

type ShortURLLocator struct {
	ID      string                 `json:"id"`
	State   map[string]interface{} `json:"state"`
	Version string                 `json:"version"`
}

type ShortURLCreateRequest struct {
	Body ShortURLCreateRequestBody
}

type ShortURLCreateRequestBody struct {
	// When the slug parameter is omitted, the API will generate a random human-readable slug if humanReadableSlug is set to true.
	HumanReadableSlug bool `json:"humanReadableSlug"`
	// The identifier for the locator.
	LocatorID string `json:"locatorId"`
	// An object which contains all necessary parameters for the given locator to resolve to a Kibana location.
	Params map[string]interface{} `json:"params"`
	// A custom short URL slug. The slug is the part of the short URL that identifies it.
	// You can provide a custom slug which consists of latin alphabet letters, numbers, and -._ characters.
	// The slug must be at least 3 characters long, but no longer than 255 characters.
	Slug *string `json:"slug,omitempty"`
}

// newShortURLCreate returns a function that performs POST /api/short_url API requests
func (api *API) newShortURLCreate() func(context.Context, *ShortURLCreateRequest, ...RequestOption) (*ShortURLCreateResponse, error) {
	return func(ctx context.Context, req *ShortURLCreateRequest, opts ...RequestOption) (*ShortURLCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "short_url.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/short_url"

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
			instrument.BeforeRequest(httpReq, "short_url.create")
			if reader := instrument.RecordRequestBody(ctx, "short_url.create", httpReq.Body); reader != nil {
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
		resp := &ShortURLCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result ShortURLCreateResponseBody

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
