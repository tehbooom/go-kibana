package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkGetDiagnosticsAgentResponse  wraps the response from a FleetBulkGetAgentPolicies call
type FleetBulkGetDiagnosticsAgentResponse struct {
	StatusCode int
	Body       *FleetBulkGetDiagnosticsAgentResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetBulkGetDiagnosticsAgentResponseBody struct {
	ActionId string `json:"actionId"`
}

type FleetBulkGetDiagnosticsAgentRequest struct {
	Body FleetBulkGetDiagnosticsAgentRequestBody
}

type FleetBulkGetDiagnosticsAgentRequestBody struct {
	Agents            json.RawMessage `json:"agents"`
	AdditionalMetrics *[]string       `json:"additional_metrics,omitempty"`
}

// SetAgentsQuery sets the Agents field as a KQL query string, leave empty to action all agents
func (body *FleetBulkGetDiagnosticsAgentRequestBody) SetAgentsQuery(query string) error {
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// SetAgentsList sets the Agents field as a list of agent IDs
func (body *FleetBulkGetDiagnosticsAgentRequestBody) SetAgentsList(agents []string) error {
	data, err := json.Marshal(agents)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// newFleetBulkGetDiagnosticsAgents returns a function that performs POST /api/fleet/agent/bulk_request_diagnostics API requests
func (api *API) newFleetBulkGetDiagnosticsAgents() func(context.Context, *FleetBulkGetDiagnosticsAgentRequest, ...RequestOption) (*FleetBulkGetDiagnosticsAgentResponse, error) {
	return func(ctx context.Context, req *FleetBulkGetDiagnosticsAgentRequest, opts ...RequestOption) (*FleetBulkGetDiagnosticsAgentResponse, error) {
		if req.Body.Agents == nil {
			return nil, fmt.Errorf("Agent IDs are not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.bulk.diagnostics")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/request_diagnostics"

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
			instrument.BeforeRequest(httpReq, "fleet.agents.bulk.diagnostics")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.bulk.diagnostics", httpReq.Body); reader != nil {
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
		resp := &FleetBulkGetDiagnosticsAgentResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBulkGetDiagnosticsAgentResponseBody

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
