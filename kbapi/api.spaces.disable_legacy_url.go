package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SpacesDisableLegacyURLResponse wraps the response from a SpacesDisableLegacyURL call
type SpacesDisableLegacyURLResponse struct {
	StatusCode int
	Body       *SpacesDisableLegacyURLResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SpacesDisableLegacyURLResponseBody struct {
}

// SpacesDisableLegacyURLRequest is the request for newSpacesDisableLegacyURL
type SpacesDisableLegacyURLRequest struct {
	// The space identifier
	ID   string `json:"id"`
	Body SpacesDisableLegacyURLRequestBody
}

// SpacesDisableLegacyURLRequestBody  defines the body for SpacesCreateRequest.
type SpacesDisableLegacyURLRequestBody struct {
	Aliases []SpacesAlias `json:"aliases"`
}

type SpacesAlias struct {
	// SourceId The alias source object identifier. This is the legacy object identifier.
	SourceId string `json:"sourceId"`

	// TargetSpace The space where the alias target object exists.
	TargetSpace string `json:"targetSpace"`

	// TargetType The type of alias target object.
	TargetType string `json:"targetType"`
}

// newSpacesDisableLegacyURL returns a function that performs POST /api/spaces/_disable_legacy_url_aliases API requests
func (api *API) newSpacesDisableLegacyURL() func(context.Context, *SpacesDisableLegacyURLRequest, ...RequestOption) (*SpacesDisableLegacyURLResponse, error) {
	return func(ctx context.Context, req *SpacesDisableLegacyURLRequest, opts ...RequestOption) (*SpacesDisableLegacyURLResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("request cannot be nil")
		}

		if req.ID == "" {
			return nil, fmt.Errorf("ID not specified")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "spaces.disable_legacy_url")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/_disable_legacy_url_aliases"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "spaces.disable_legacy_url")
			if reader := instrument.RecordRequestBody(ctx, "spaces.disable_legacy_url", httpReq.Body); reader != nil {
				httpReq.Body = reader
			}
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

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
		resp := &SpacesDisableLegacyURLResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesDisableLegacyURLResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
