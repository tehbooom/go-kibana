package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// APMAgentKeyCreateResponse wraps the response from a <todo> call
type APMAgentKeyCreateResponse struct {
	StatusCode int
	Body       *APMAgentKeyCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAgentKeyCreateResponseBody struct {
	AgentKey *struct {
		APIKey     string `json:"api_key"`
		Encoded    string `json:"encoded"`
		Expiration *int64 `json:"expiration,omitempty"`
		ID         string `json:"id"`
		Name       string `json:"name"`
	} `json:"agentKey,omitempty"`
}

type APMAgentKeyCreateRequest struct {
	Body APMAgentKeyCreateRequestBody
}

type APMAgentKeyCreateRequestBody struct {
	// Name The name of the APM agent key.
	Name string `json:"name"`
	// Privileges The APM agent key privileges. It can take one or more of the following values:
	//
	// `event:write`, which is required for ingesting APM agent events.
	// `config_agent:read`, which is required for APM agents to read agent configuration remotely.
	Privileges []string `json:"privileges"`
}

// newAPMAgentKeyCreate returns a function that performs POST /api/apm/agent_keys API requests
func (api *API) newAPMAgentKeyCreate() func(context.Context, *APMAgentKeyCreateRequest, ...RequestOption) (*APMAgentKeyCreateResponse, error) {
	return func(ctx context.Context, req *APMAgentKeyCreateRequest, opts ...RequestOption) (*APMAgentKeyCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.agent_key.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/apm/agent_keys"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

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
			instrument.BeforeRequest(httpReq, "apm.agent_key.create")
			if reader := instrument.RecordRequestBody(ctx, "apm.agent_key.create", httpReq.Body); reader != nil {
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
		resp := &APMAgentKeyCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAgentKeyCreateResponseBody

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
