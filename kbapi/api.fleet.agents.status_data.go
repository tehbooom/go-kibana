package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FleetAgentStatusDataResponse wraps the response from a FleetAgentStatusData call
type FleetAgentStatusDataResponse struct {
	StatusCode int
	Body       *FleetAgentStatusDataResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetAgentStatusDataResponseBody struct {
	DataPreview []interface{} `json:"dataPreview"`
	Items       []map[string]struct {
		Data bool `json:"data"`
	} `json:"items"`
}

type FleetAgentStatusDataRequest struct {
	Params FleetAgentStatusDataRequestParams
}
type FleetAgentStatusDataRequestParams struct {
	AgentsIds   []string `form:"agentsIds" json:"agentsIds"`
	PkgName     *string  `form:"pkgName,omitempty" json:"pkgName,omitempty"`
	PkgVersion  *string  `form:"pkgVersion,omitempty" json:"pkgVersion,omitempty"`
	PreviewData *bool    `form:"previewData,omitempty" json:"previewData,omitempty"`
}

// newFleetAgentStatusDataFunc returns a function that performs get /api/fleet/agent_status/data API requests
func (api *API) newFleetAgentStatusDataFunc() func(context.Context, *FleetAgentStatusDataRequest, ...RequestOption) (*FleetAgentStatusDataResponse, error) {
	return func(ctx context.Context, req *FleetAgentStatusDataRequest, opts ...RequestOption) (*FleetAgentStatusDataResponse, error) {
		if req == nil {
			req = &FleetAgentStatusDataRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.status_data")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_status/data"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.AgentsIds != nil {
			params["agentsIds"] = strings.Join(req.Params.AgentsIds, ",")
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
			instrument.BeforeRequest(httpReq, "fleet.agents.status_data")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.status_data", httpReq.Body); reader != nil {
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
		resp := &FleetAgentStatusDataResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetAgentStatusDataResponseBody

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
