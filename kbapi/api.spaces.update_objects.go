package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SpacesUpdateObjectsResponse wraps the response from a SpacesUpdateObjects call
type SpacesUpdateObjectsResponse struct {
	StatusCode int
	Body       *SpacesUpdateObjectsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SpacesUpdateObjectsResponseBody struct {
	Objects []SpaceObjects `json:"objects"`
}

type SpaceObjects struct {
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Spaces []string `json:"spaces"`
}

// SpacesUpdateObjectsRequest is the request for newSpacesUpdateObjects
type SpacesUpdateObjectsRequest struct {
	Body SpacesUpdateObjectsRequestBody
}

type SpacesUpdateObjectsRequestBody struct {
	Objects        []Object `json:"objects"`
	SpacesToAdd    []string `json:"spacesToAdd"`
	SpacesToRemove []string `json:"spacesToRemove"`
}

// newSpacesUpdateObjects returns a function that performs POST /api/spaces/_update_objects_spaces API requests
func (api *API) newSpacesUpdateObjects() func(context.Context, *SpacesUpdateObjectsRequest, ...RequestOption) (*SpacesUpdateObjectsResponse, error) {
	return func(ctx context.Context, req *SpacesUpdateObjectsRequest, opts ...RequestOption) (*SpacesUpdateObjectsResponse, error) {
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
			newCtx = instrument.Start(ctx, "spaces.update_objects")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/_update_objects_spaces"

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
			instrument.BeforeRequest(httpReq, "spaces.update_objects")
			if reader := instrument.RecordRequestBody(ctx, "spaces.update_objects", httpReq.Body); reader != nil {
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
		resp := &SpacesUpdateObjectsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesUpdateObjectsResponseBody

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
