package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetBulkGetAgentPolicies wraps the response from a FleetBulkGetAgentPolicies call
type FleetListAgentsResponse struct {
	StatusCode int
	Body       *FleetListAgentsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetListAgentsRequest struct {
	Params FleetListAgentsRequestParams
}

// GetFleetAgentsParams defines parameters for GetFleetAgents.
type FleetListAgentsRequestParams struct {
	Page             *float32 `form:"page,omitempty" json:"page,omitempty"`
	PerPage          *float32 `form:"perPage,omitempty" json:"perPage,omitempty"`
	Kuery            *string  `form:"kuery,omitempty" json:"kuery,omitempty"`
	ShowInactive     *bool    `form:"showInactive,omitempty" json:"showInactive,omitempty"`
	WithMetrics      *bool    `form:"withMetrics,omitempty" json:"withMetrics,omitempty"`
	ShowUpgradeable  *bool    `form:"showUpgradeable,omitempty" json:"showUpgradeable,omitempty"`
	GetStatusSummary *bool    `form:"getStatusSummary,omitempty" json:"getStatusSummary,omitempty"`
	SortOrder        *string  `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// newFleetListAgents returns a function that performs GET /api/fleet/agents API requests
func (api *API) newFleetListAgents() func(context.Context, *FleetListAgentsRequest, ...RequestOption) (*FleetListAgentsResponse, error) {
	return func(ctx context.Context, req *FleetListAgentsRequest, opts ...RequestOption) (*FleetListAgentsResponse, error) {
		if req == nil {
			req = &FleetListAgentsRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents"

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
		if req.Params.ShowInactive != nil {
			params["showInactive"] = strconv.FormatBool(*req.Params.ShowInactive)
		}
		if req.Params.WithMetrics != nil {
			params["withMetrics"] = strconv.FormatBool(*req.Params.WithMetrics)
		}
		if req.Params.ShowUpgradeable != nil {
			params["showUpgradeable"] = strconv.FormatBool(*req.Params.ShowUpgradeable)
		}
		if req.Params.GetStatusSummary != nil {
			params["getStatusSummary"] = strconv.FormatBool(*req.Params.GetStatusSummary)
		}
		if req.Params.SortOrder != nil {
			params["sortOrder"] = *req.Params.SortOrder
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
			instrument.BeforeRequest(httpReq, "fleet.agents.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.list", httpReq.Body); reader != nil {
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
		resp := &FleetListAgentsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetListAgentsResponseBody

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
