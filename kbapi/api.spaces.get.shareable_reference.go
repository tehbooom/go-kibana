package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SpacesShareableReferencesResponse wraps the response from a SpacesShareableReferences call
type SpacesShareableReferencesResponse struct {
	StatusCode int
	Body       *SpacesShareableReferencesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SpacesShareableReferencesResponseBody struct {
	Objects []SpaceObjectsReferences `json:"objects"`
}

type SpaceObjectsReferences struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	Spaces           []string `json:"spaces"`
	InboudReferences []string `json:"inboundReferences"`
}

// SpacesShareableReferencesRequest is the request for newSpacesShareableReferences
type SpacesShareableReferencesRequest struct {
	Body SpacesShareableReferencesRequestBody
}

// SpacesShareableReferencesRequestBody  defines the body for SpacesCreateRequest.
type SpacesShareableReferencesRequestBody struct {
	Objects []Object `json:"objects"`
}

// newSpacesShareableReferences returns a function that performs POST /api/spaces/_get_shareable_references API requests
func (api *API) newSpacesShareableReferences() func(context.Context, *SpacesShareableReferencesRequest, ...RequestOption) (*SpacesShareableReferencesResponse, error) {
	return func(ctx context.Context, req *SpacesShareableReferencesRequest, opts ...RequestOption) (*SpacesShareableReferencesResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "spaces.shareable_references")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/_get_shareable_references"

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
			instrument.BeforeRequest(httpReq, "spaces.shareable_references")
			if reader := instrument.RecordRequestBody(ctx, "spaces.shareable_references", httpReq.Body); reader != nil {
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
		resp := &SpacesShareableReferencesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesShareableReferencesResponseBody

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
