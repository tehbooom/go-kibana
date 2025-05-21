package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkReassignAgentResponse wraps the response from a FleetReassignAgent
type FleetBulkReassignAgentResponse struct {
	StatusCode int
	Body       *FleetBulkReassignAgentResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetBulkReassignAgentResponseBody struct {
	ActionId string `json:"actionId"`
}

type FleetBulkReassignAgentRequest struct {
	Body *FleetBulkReassignAgentRequestBody
}

// PostFleetBulkReassignAgentRequestBody  defines parameters for PostFleetAgentsBulkReassign.
type FleetBulkReassignAgentRequestBody struct {
	Agents          json.RawMessage `json:"agents"`
	BatchSize       *float32        `json:"batchSize,omitempty"`
	IncludeInactive *bool           `json:"includeInactive,omitempty"`
	PolicyId        string          `json:"policy_id"`
}

// SetAgentsQuery sets the Agents field as a KQL query string, leave empty to action all agents
func (body *FleetBulkReassignAgentRequestBody) SetAgentsQuery(query string) error {
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// SetAgentsList sets the Agents field as a list of agent IDs
func (body *FleetBulkReassignAgentRequestBody) SetAgentsList(agents []string) error {
	data, err := json.Marshal(agents)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// newFleetBulkReassignAgent returns a function that performs POST /api/fleet/agents/bulk_reassign API requests
func (api *API) newFleetBulkReassignAgents() func(context.Context, *FleetBulkReassignAgentRequest, ...RequestOption) (*FleetBulkReassignAgentResponse, error) {
	return func(ctx context.Context, req *FleetBulkReassignAgentRequest, opts ...RequestOption) (*FleetBulkReassignAgentResponse, error) {
		if req.Body.PolicyId == "" {
			return nil, fmt.Errorf("Policy ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.bulk.reassign")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/bulk_reassign"

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
			instrument.BeforeRequest(httpReq, "fleet.agents.bulk.reassign")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.bulk.reassign", httpReq.Body); reader != nil {
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
		resp := &FleetBulkReassignAgentResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBulkReassignAgentResponseBody

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
