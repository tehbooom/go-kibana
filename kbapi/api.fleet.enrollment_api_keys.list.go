package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEnrollmentAPIKeysListResponse  wraps the response from a FleetBulkGetAgentPolicies call
type FleetEnrollmentAPIKeysListResponse struct {
	StatusCode int
	Body       *FleetEnrollmentAPIKeysListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEnrollmentAPIKeysListResponseBody struct {
	Action  string             `json:"action"`
	Items   []EnrollmentApiKey `json:"items"`
	List    []EnrollmentApiKey `json:"list"`
	Page    float32            `json:"page"`
	PerPage float32            `json:"perPage"`
	Total   float32            `json:"total"`
}

type EnrollmentApiKey struct {
	// Active When false, the enrollment API key is revoked and cannot be used for enrolling Elastic Agents.
	Active bool `json:"active"`

	// ApiKey The enrollment API key (token) used for enrolling Elastic Agents.
	APIKey string `json:"api_key"`

	// ApiKeyId The ID of the API key in the Security API.
	ApiKeyID  string `json:"api_key_id"`
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`

	// Name The name of the enrollment API key.
	Name *string `json:"name,omitempty"`

	// PolicyId The ID of the agent policy the Elastic Agent will be enrolled in.
	PolicyID *string `json:"policy_id,omitempty"`
}

// FleetEnrollmentAPIKeysListRequest  defines parameters for GetFleetAgents.
type FleetEnrollmentAPIKeysListRequest struct {
	Params FleetEnrollmentAPIKeysListRequestParams
}

type FleetEnrollmentAPIKeysListRequestParams struct {
	Page    *float32 `form:"page,omitempty" json:"page,omitempty"`
	PerPage *float32 `form:"perPage,omitempty" json:"perPage,omitempty"`
	Kuery   *string  `form:"kuery,omitempty" json:"kuery,omitempty"`
}

// newFleetEnrollmentAPIKeysList returns a function that performs GET /api/fleet/enrollment_api_keys API requests
func (api *API) newFleetEnrollmentAPIKeysList() func(context.Context, *FleetEnrollmentAPIKeysListRequest, ...RequestOption) (*FleetEnrollmentAPIKeysListResponse, error) {
	return func(ctx context.Context, req *FleetEnrollmentAPIKeysListRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysListResponse, error) {
		if req == nil {
			req = &FleetEnrollmentAPIKeysListRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.enrollment_api_keys.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/enrollment_api_keys"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.FormatFloat(float64(*req.Params.Page), 'f', -1, 32)
		}
		if req.Params.PerPage != nil {
			params["perPage"] = strconv.FormatFloat(float64(*req.Params.PerPage), 'f', -1, 32)
		}
		if req.Params.Kuery != nil {
			params["kuery"] = *req.Params.Kuery
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
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
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.enrollment_api_keys.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.enrollment_api_keys.list", httpReq.Body); reader != nil {
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
		resp := &FleetEnrollmentAPIKeysListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEnrollmentAPIKeysListResponseBody

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
