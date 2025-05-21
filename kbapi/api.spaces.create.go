package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SpacesCreateResponse wraps the response from a SpacesCreate call
type SpacesCreateResponse struct {
	StatusCode int
	Body       *SpacesCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SpacesCreateResponseBody struct {
}

// SpacesCreateRequest  is the request for newFleetBulkGetAgentPolicies
type SpacesCreateRequest struct {
	Body SpacesCreateRequestBody
}

// SpacesCreateRequestBody defines the body for SpacesCreateRequest.
type SpacesCreateRequestBody struct {
	Reserved *bool `json:"_reserved,omitempty"`

	// Color The hexadecimal color code used in the space avatar. By default, the color is automatically generated from the space name.
	Color *string `json:"color,omitempty"`

	// Description A description for the space.
	Description *string `json:"description,omitempty"`

	//The list of features that are turned off in the space.
	DisabledFeatures *[]string `json:"disabledFeatures,omitempty"`

	// Id The space ID that is part of the Kibana URL when inside the space. Space IDs are limited to lowercase alphanumeric, underscore, and hyphen characters (a-z, 0-9, _, and -). You are cannot change the ID with the update operation.
	ID string `json:"id"`

	// ImageUrl The data-URL encoded image to display in the space avatar. If specified, initials will not be displayed and the color will be visible as the background color for transparent images. For best results, your image should be 64x64. Images will not be optimized by this API call, so care should be taken when using custom images.
	ImageUrl *string `json:"imageUrl,omitempty"`

	// Initials One or two characters that are shown in the space avatar. By default, the initials are automatically generated from the space name.
	Initials *string `json:"initials,omitempty"`

	// Name The display name for the space.
	Name string `json:"name"`

	// Values are security, oblt, es, or classic.
	Solution *string `json:"solution,omitempty"`
}

// newSpacesCreate returns a function that performs POST /api/spaces/space API requests
func (api *API) newSpacesCreate() func(context.Context, *SpacesCreateRequest, ...RequestOption) (*SpacesCreateResponse, error) {
	return func(ctx context.Context, req *SpacesCreateRequest, opts ...RequestOption) (*SpacesCreateResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Name or ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "spaces.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/space"

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

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "spaces.create")
			if reader := instrument.RecordRequestBody(ctx, "spaces.create", httpReq.Body); reader != nil {
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
		resp := &SpacesCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesCreateResponseBody

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
