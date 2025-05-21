package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// APMAgentConfigurationGetNameResponse wraps the response from a <todo> call
type APMAgentConfigurationGetNameResponse struct {
	StatusCode int
	Body       *APMAgentConfigurationGetNameResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAgentConfigurationGetNameResponseBody struct {
	AgentName string `json:"agentName"`
}

type APMAgentConfigurationGetNameRequest struct {
	Params APMAgentConfigurationGetNameRequestParams
}

type APMAgentConfigurationGetNameRequestParams struct {
	ServiceName string `form:"serviceName" json:"serviceName"`
}

// newAPMAgentConfigurationGetName returns a function that performs GET /api/apm/settings/agent-configuration/agent_name API requests
func (api *API) newAPMAgentConfigurationGetName() func(context.Context, *APMAgentConfigurationGetNameRequest, ...RequestOption) (*APMAgentConfigurationGetNameResponse, error) {
	return func(ctx context.Context, req *APMAgentConfigurationGetNameRequest, opts ...RequestOption) (*APMAgentConfigurationGetNameResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.agent_configuration.get_name")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/apm/settings/agent-configuration/agent_name"

		// Build query parameters
		params := make(map[string]string)

		params["serviceName"] = req.Params.ServiceName

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
			instrument.BeforeRequest(httpReq, "apm.agent_configuration.get_name")
			if reader := instrument.RecordRequestBody(ctx, "apm.agent_configuration.get_name", httpReq.Body); reader != nil {
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
		resp := &APMAgentConfigurationGetNameResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAgentConfigurationGetNameResponseBody

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
