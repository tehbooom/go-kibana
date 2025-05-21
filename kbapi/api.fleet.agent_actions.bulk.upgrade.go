package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkUpgradeAgentsResponse wraps the response from a FleetUpgradeAgent call
type FleetBulkUpgradeAgentsResponse struct {
	StatusCode int
	Body       *FleetBulkUpgradeAgentsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetBulkUpgradeAgentsResponseBody struct {
	Force              *bool   `json:"force"`
	SkipRateLimitCheck *bool   `json:"skipRateLimitCheck"`
	SourceUri          *string `json:"source_uri"`
	Version            string  `json:"version"`
}

type FleetBulkUpgradeAgentsRequest struct {
	Body FleetBulkUpgradeAgentsRequestBody
}

type FleetBulkUpgradeAgentsRequestBody struct {
	Agents json.RawMessage `json:"agents"`
	// Force upgrade, skipping validation (should be used with caution)
	Force *bool `json:"force,omitempty"`
	// rolling upgrade window duration in seconds
	RolloutDurationSeconds *float32 `json:"rollout_duration_seconds,omitempty"`
	// Skip rate limit check for upgrade
	SkipRateLimitCheck *bool `json:"skipRateLimitCheck,omitempty"`
	// alternative upgrade binary download url
	SourceUri *string `json:"source_uri,omitempty"`
	// start time of upgrade in ISO 8601 format
	StartTime *string `json:"start_time,omitempty"`
	// version to upgrade to
	Version string `json:"version"`
}

// SetAgentsQuery sets the Agents field as a KQL query string, leave empty to action all agents
func (body *FleetBulkUpgradeAgentsRequestBody) SetAgentsQuery(query string) error {
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// SetAgentsList sets the Agents field as a list of agent IDs
func (body *FleetBulkUpgradeAgentsRequestBody) SetAgentsList(agents []string) error {
	data, err := json.Marshal(agents)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// newFleetBulkUpgradeAgents returns a function that performs POST /api/fleet/agent/bulk_upgrade API requests
func (api *API) newFleetBulkUpgradeAgents() func(context.Context, *FleetBulkUpgradeAgentsRequest, ...RequestOption) (*FleetBulkUpgradeAgentsResponse, error) {
	return func(ctx context.Context, req *FleetBulkUpgradeAgentsRequest, opts ...RequestOption) (*FleetBulkUpgradeAgentsResponse, error) {
		if req.Body.Agents == nil {
			return nil, fmt.Errorf("Agent ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.bulk.upgrade")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/bulk_upgrade"

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
			instrument.BeforeRequest(httpReq, "fleet.agents.bulk.upgrade")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.bulk.upgrade", httpReq.Body); reader != nil {
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
		resp := &FleetBulkUpgradeAgentsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBulkUpgradeAgentsResponseBody

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
