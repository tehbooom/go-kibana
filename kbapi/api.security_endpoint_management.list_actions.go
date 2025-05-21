package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// SecurityEndpointManagementListActionsResponse wraps the response from a List call
type SecurityEndpointManagementListActionsResponse struct {
	StatusCode int
	Body       *SecurityEndpointManagementListActionsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityEndpointManagementListActionsResponseBody struct {
	Data            []SecurityEndpointManagementAction `json:"data"`
	Page            int                                `json:"page"`
	PageSize        int                                `json:"pageSize"`
	StartDate       string                             `json:"startDate,omitempty"`
	EndDate         string                             `json:"endDate,omitempty"`
	ElasticAgentIDs []string                           `json:"elasticAgentIds,omitempty"`
	Total           int                                `json:"total"`
}

type SecurityEndpointManagementListActionsRequest struct {
	Params SecurityEndpointManagementListActionsRequestParams
}

type SecurityEndpointManagementListActionsRequestParams struct {
	Page     *int
	PageSize *int
	// Commands A list of response action command names. Minimum length of each is 1.
	// Values are isolate, unisolate, kill-process, suspend-process, running-processes,
	// get-file, execute, upload, or scan.
	Commands *[]string
	// AgentIDs A list of agent IDs. Max of 50.
	AgentIDs *[]string
	UserIDs  *[]string
	// StartDate A start date in ISO 8601 format or Date Math format.
	StartDate *string
	// EndDate An end date in ISO format or Date Math format.
	EndDate *string
	// WithOutputs A list of action IDs that should include the complete output of the action.
	WithOutputs *[]string
	// Types List of types of response actions
	// Values are automated or manual.
	Types *[]string
	// AgentTypes List of agent types to retrieve. Defaults to endpoint.
	// Values are endpoint, sentinel_one, crowdstrike, or microsoft_defender_endpoint.
	AgentTypes *string
}

// newSecurityEndpointManagementListActions returns a function that performs GET /api/endpoint/action API requests
func (api *API) newSecurityEndpointManagementListActions() func(context.Context, *SecurityEndpointManagementListActionsRequest, ...RequestOption) (*SecurityEndpointManagementListActionsResponse, error) {
	return func(ctx context.Context, req *SecurityEndpointManagementListActionsRequest, opts ...RequestOption) (*SecurityEndpointManagementListActionsResponse, error) {
		if req == nil {
			req = &SecurityEndpointManagementListActionsRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_endpoint_management.list_actions")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/endpoint/action"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Commands != nil {
			params["commands"] = strings.Join(*req.Params.Commands, ",")
		}
		if req.Params.AgentIDs != nil {
			params["agentIds"] = strings.Join(*req.Params.AgentIDs, ",")
		}
		if req.Params.UserIDs != nil {
			params["userIds"] = strings.Join(*req.Params.UserIDs, ",")
		}
		if req.Params.StartDate != nil {
			params["startDate"] = *req.Params.StartDate
		}
		if req.Params.EndDate != nil {
			params["endDate"] = *req.Params.EndDate
		}
		if req.Params.AgentTypes != nil {
			params["agentTypes"] = *req.Params.AgentTypes
		}
		if req.Params.WithOutputs != nil {
			params["withOutputs"] = strings.Join(*req.Params.WithOutputs, ",")
		}
		if req.Params.Types != nil {
			params["types"] = strings.Join(*req.Params.Types, ",")
		}
		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PageSize != nil {
			params["pageSize"] = strconv.Itoa(*req.Params.PageSize)
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
			instrument.BeforeRequest(httpReq, "security_endpoint_management.list_actions")
			if reader := instrument.RecordRequestBody(ctx, "security_endpoint_management.list_actions", httpReq.Body); reader != nil {
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
		resp := &SecurityEndpointManagementListActionsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityEndpointManagementListActionsResponseBody

		if httpResp.StatusCode < 299 {
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
			// For all non-success responses
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
