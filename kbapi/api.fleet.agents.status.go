package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetAgentStatusResponse wraps the response from a FleetAgentStatus call
type FleetAgentStatusResponse struct {
	StatusCode int
	Body       *FleetAgentStatusResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetAgentStatusResponseBody struct {
	Results struct {
		Active      float32  `json:"active"`
		All         float32  `json:"all"`
		Error       float32  `json:"error"`
		Events      float32  `json:"events"`
		Inactive    float32  `json:"inactive"`
		Offline     float32  `json:"offline"`
		Online      float32  `json:"online"`
		Orphaned    *float32 `json:"orphaned,omitempty"`
		Other       float32  `json:"other"`
		Unenrolled  float32  `json:"unenrolled"`
		Uninstalled *float32 `json:"uninstalled,omitempty"`
		Updating    float32  `json:"updating"`
	} `json:"results"`
}
type FleetAgentStatusRequest struct {
	Params FleetAgentStatusRequestParams
}

type FleetAgentStatusRequestParams struct {
	PolicyId *string `form:"policyId,omitempty" json:"policyId,omitempty"`
	Kuery    *string `form:"kuery,omitempty" json:"kuery,omitempty"`
}

// newFleetAgentStatusFunc returns a function that performs GET /api/fleet/agent_status API requests
func (api *API) newFleetAgentStatusFunc() func(context.Context, *FleetAgentStatusRequest, ...RequestOption) (*FleetAgentStatusResponse, error) {
	return func(ctx context.Context, req *FleetAgentStatusRequest, opts ...RequestOption) (*FleetAgentStatusResponse, error) {
		if req == nil {
			req = &FleetAgentStatusRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_status"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.PolicyId != nil {
			params["policyId"] = *req.Params.PolicyId
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
			instrument.BeforeRequest(httpReq, "fleet.agents.status")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.status", httpReq.Body); reader != nil {
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
		resp := &FleetAgentStatusResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetAgentStatusResponseBody

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
