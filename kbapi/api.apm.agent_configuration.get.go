package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// APMAgentConfigurationGetResponse wraps the response from a <todo> call
type APMAgentConfigurationGetResponse struct {
	StatusCode int
	Body       *APMAgentConfigurationGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAgentConfigurationGetResponseBody struct {
	ID             string           `json:"id"`
	AgentName      *string          `json:"agent_name,omitempty"`
	AppliedByAgent *bool            `json:"applied_by_agent,omitempty"`
	AtTimestamp    float32          `json:"@Timestamp"`
	Etag           string           `json:"etag"`
	Service        APMServiceObject `json:"service"`
	// Settings Agent configuration settings
	Settings map[string]string `json:"settings"`
}

type APMAgentConfigurationGetRequest struct {
	Params APMAgentConfigurationGetRequestParams
}

type APMAgentConfigurationGetRequestParams struct {
	Name        *string `form:"name,omitempty" json:"name,omitempty"`
	Environment *string `form:"environment,omitempty" json:"environment,omitempty"`
}

// newAPMAgentConfigurationGet returns a function that performs GET /api/ API requests
func (api *API) newAPMAgentConfigurationGet() func(context.Context, *APMAgentConfigurationGetRequest, ...RequestOption) (*APMAgentConfigurationGetResponse, error) {
	return func(ctx context.Context, req *APMAgentConfigurationGetRequest, opts ...RequestOption) (*APMAgentConfigurationGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.agent_configuration.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Name != nil {
			params["name"] = *req.Params.Name
		}

		if req.Params.Environment != nil {
			params["environment"] = *req.Params.Environment
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
			instrument.BeforeRequest(httpReq, "apm.agent_configuration.get")
			if reader := instrument.RecordRequestBody(ctx, "apm.agent_configuration.get", httpReq.Body); reader != nil {
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
		resp := &APMAgentConfigurationGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAgentConfigurationGetResponseBody

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
