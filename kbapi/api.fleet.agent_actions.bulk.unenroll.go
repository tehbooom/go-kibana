package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetBulkUnenrollAgentsResponse  wraps the response from a FleetUnenrollAgent call
type FleetBulkUnenrollAgentsResponse struct {
	StatusCode int
	Body       *FleetBulkUnenrollAgentsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetBulkUnenrollAgentsResponseBody struct {
	ActionId string `json:"actionId"`
}

type FleetBulkUnenrollAgentsRequest struct {
	Body FleetBulkUnenrollAgentsRequestBody
}

type FleetBulkUnenrollAgentsRequestBody struct {
	Agents          json.RawMessage `json:"agents"`
	Force           *bool           `json:"force,omitempty"`
	IncludeInactive *bool           `json:"includeInactive,omitempty"`
	Revoke          *bool           `json:"revoke,omitempty"`
}

// SetAgentsQuery sets the Agents field as a KQL query string, leave empty to action all agents
func (body *FleetBulkUnenrollAgentsRequestBody) SetAgentsQuery(query string) error {
	data, err := json.Marshal(query)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// SetAgentsList sets the Agents field as a list of agent IDs
func (body *FleetBulkUnenrollAgentsRequestBody) SetAgentsList(agents []string) error {
	data, err := json.Marshal(agents)
	if err != nil {
		return err
	}
	body.Agents = data
	return nil
}

// newFleetBulkUnenrollAgents returns a function that performs POST /api/fleet/agent/bulk_unenroll API requests
func (api *API) newFleetBulkUnenrollAgents() func(context.Context, *FleetBulkUnenrollAgentsRequest, ...RequestOption) (*FleetBulkUnenrollAgentsResponse, error) {
	return func(ctx context.Context, req *FleetBulkUnenrollAgentsRequest, opts ...RequestOption) (*FleetBulkUnenrollAgentsResponse, error) {
		if req.Body.Agents == nil {
			return nil, fmt.Errorf("Agent IDs is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agents.bulk.unenroll")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/bulk_unenroll"

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
			instrument.BeforeRequest(httpReq, "fleet.agents.bulk.unenroll")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agents.bulk.unenroll", httpReq.Body); reader != nil {
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
		resp := &FleetBulkUnenrollAgentsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBulkUnenrollAgentsResponseBody

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
