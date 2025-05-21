package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// SpacesGetAllResponse wraps the response from a SpacesGetAll call
type SpacesGetAllResponse struct {
	StatusCode int
	Body       *SpacesGetAllResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type Space struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	ImageURL         string   `json:"imageUrl,omitempty"`
	Reserved         bool     `json:"_reserved,omitempty"`
	Description      string   `json:"description,omitempty"`
	DisabledFeatures []string `json:"disabledFeatures,omitempty"`
	Color            *string  `json:"color,omitempty"`
	Initials         string   `json:"initials,omitempty"`
	Solution         string   `json:"solution,omitempty"`
}

type SpacesGetAllResponseBody []Space

// SpacesGetAllRequest is the request for newFleetBulkGetAgentPolicies
type SpacesGetAllRequest struct {
	Params SpacesGetAllRequestParams
}

// SpacesGetAllRequestParams defines the body for SpacesGetAllRequest.
type SpacesGetAllRequestParams struct {
	// Purpose Specifies which authorization checks are applied to the API call. The default value is `any`.
	Purpose *string `form:"purpose,omitempty" json:"purpose,omitempty"`

	// IncludeAuthorizedPurposes When enabled, the API returns any spaces that the user is authorized to access in any capacity and each space will contain the purposes for which the user is authorized. This can be useful to determine which spaces a user can read but not take a specific action in. If the security plugin is not enabled, this parameter has no effect, since no authorization checks take place. This parameter cannot be used in with the `purpose` parameter.
	IncludeAuthorizedPurposes *bool `form:"include_authorized_purposes,omitempty" json:"include_authorized_purposes,omitempty"`
}

// newSpacesGetAll returns a function that performs GET /api/spaces/space API requests
func (api *API) newSpacesGetAll() func(context.Context, *SpacesGetAllRequest, ...RequestOption) (*SpacesGetAllResponse, error) {
	return func(ctx context.Context, req *SpacesGetAllRequest, opts ...RequestOption) (*SpacesGetAllResponse, error) {
		if req == nil {
			req = &SpacesGetAllRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "spaces.get_all")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/spaces/space"
		//
		// Build query parameters
		params := make(map[string]string)

		if req.Params.Purpose != nil {
			params["purpose"] = *req.Params.Purpose
		}
		if req.Params.IncludeAuthorizedPurposes != nil {
			params["include_authorized_purposes"] = strconv.FormatBool(*req.Params.IncludeAuthorizedPurposes)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "spaces.get_all")
			if reader := instrument.RecordRequestBody(ctx, "spaces.get_all", httpReq.Body); reader != nil {
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
		resp := &SpacesGetAllResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SpacesGetAllResponseBody

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
