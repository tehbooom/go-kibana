package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// APMAgentConfigurationCreateUpdateResponse wraps the response from a <todo> call
type APMAgentConfigurationCreateUpdateResponse struct {
	StatusCode int
	Body       *APMAgentConfigurationCreateUpdateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAgentConfigurationCreateUpdateResponseBody struct{}

type APMAgentConfigurationCreateUpdateRequest struct {
	Params APMAgentConfigurationCreateUpdateRequestParams
	Body   APMAgentConfigurationCreateUpdateRequestBody
}

type APMAgentConfigurationCreateUpdateRequestParams struct {
	// If the config exists ?overwrite=true is required
	Overwrite *bool
}

type APMAgentConfigurationCreateUpdateRequestBody struct {
	AgentName *string           `json:"agent_name,omitempty"`
	Service   APMServiceObject  `json:"service"`
	Settings  map[string]string `json:"settings"`
}

// newAPMAgentConfigurationCreateUpdate returns a function that performs PUT /api/ API requests
func (api *API) newAPMAgentConfigurationCreateUpdate() func(context.Context, *APMAgentConfigurationCreateUpdateRequest, ...RequestOption) (*APMAgentConfigurationCreateUpdateResponse, error) {
	return func(ctx context.Context, req *APMAgentConfigurationCreateUpdateRequest, opts ...RequestOption) (*APMAgentConfigurationCreateUpdateResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.agent_configuration.create_update")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Overwrite != nil {
			params["overwrite"] = strconv.FormatBool(*req.Params.Overwrite)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "apm.agent_configuration.create_update")
			if reader := instrument.RecordRequestBody(ctx, "apm.agent_configuration.create_update", httpReq.Body); reader != nil {
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
		resp := &APMAgentConfigurationCreateUpdateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAgentConfigurationCreateUpdateResponseBody

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
