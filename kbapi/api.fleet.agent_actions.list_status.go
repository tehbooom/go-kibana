package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// FleetAgentActionsListStatusResponse wraps the response from a <todo> call
type FleetAgentActionsListStatusResponse struct {
	StatusCode int
	Body       *FleetAgentActionsListStatusResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetAgentActionsListStatusResponseBody struct {
	Items []struct {
		ActionID         string  `json:"actionId"`
		CancellationTime *string `json:"cancellationTime,omitempty"`
		CompletionTime   *string `json:"completionTime,omitempty"`

		// CreationTime creation time of action
		CreationTime     string  `json:"creationTime"`
		Expiration       *string `json:"expiration,omitempty"`
		HasRolloutPeriod *bool   `json:"hasRolloutPeriod,omitempty"`
		LatestErrors     *[]struct {
			AgentID   string  `json:"agentId"`
			Error     string  `json:"error"`
			Hostname  *string `json:"hostname,omitempty"`
			Timestamp string  `json:"timestamp"`
		} `json:"latestErrors,omitempty"`

		// NBAgentsAck number of agents that acknowledged the action
		NBAgentsAck float32 `json:"nbAgentsAck"`

		// NBAgentsActionCreated number of agents included in action from kibana
		NBAgentsActionCreated float32 `json:"nbAgentsActionCreated"`

		// NBAgentsActioned number of agents actioned
		NBAgentsActioned float32 `json:"nbAgentsActioned"`

		// NBAgentsFailed number of agents that failed to execute the action
		NBAgentsFailed float32 `json:"nbAgentsFailed"`

		// NewPolicyID new policy id (POLICY_REASSIGN action)
		NewPolicyID *string `json:"newPolicyId,omitempty"`

		// PolicyID policy id (POLICY_CHANGE action)
		PolicyID *string `json:"policyId,omitempty"`

		// Revision new policy revision (POLICY_CHANGE action)
		Revision *float32 `json:"revision,omitempty"`

		// StartTime start time of action (scheduled actions)
		StartTime *string `json:"startTime,omitempty"`
		Status    string  `json:"status"`
		Type      string  `json:"type"`

		// Version agent version number (UPGRADE action)
		Version *string `json:"version,omitempty"`
	} `json:"items"`
}

type FleetAgentActionsListStatusRequest struct {
	Params FleetAgentActionsListStatusRequestParams
}

type FleetAgentActionsListStatusRequestParams struct {
	Page      *float32 `form:"page,omitempty" json:"page,omitempty"`
	PerPage   *float32 `form:"perPage,omitempty" json:"perPage,omitempty"`
	Date      *string  `form:"date,omitempty" json:"date,omitempty"`
	Latest    *int     `form:"latest,omitempty" json:"latest,omitempty"`
	ErrorSize *int     `form:"errorSize,omitempty" json:"errorSize,omitempty"`
}

// newFleetAgentActionsListStatus returns a function that performs GET /api/fleet/agents/action_status API requests
func (api *API) newFleetAgentActionsListStatus() func(context.Context, *FleetAgentActionsListStatusRequest, ...RequestOption) (*FleetAgentActionsListStatusResponse, error) {
	return func(ctx context.Context, req *FleetAgentActionsListStatusRequest, opts ...RequestOption) (*FleetAgentActionsListStatusResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_actions.list_status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agents/action_status"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.FormatFloat(float64(*req.Params.Page), 'f', -1, 32)
		}
		if req.Params.PerPage != nil {
			params["perPage"] = strconv.FormatFloat(float64(*req.Params.PerPage), 'f', -1, 32)
		}
		if req.Params.Date != nil {
			params["date"] = *req.Params.Date
		}
		if req.Params.Latest != nil {
			params["latest"] = strconv.Itoa(*req.Params.Latest)
		}
		if req.Params.ErrorSize != nil {
			params["errorSize"] = strconv.Itoa(*req.Params.ErrorSize)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
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
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.agent_actions.list_status")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_actions.list_status", httpReq.Body); reader != nil {
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
		resp := &FleetAgentActionsListStatusResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetAgentActionsListStatusResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
