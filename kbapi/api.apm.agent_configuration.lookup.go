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
// APMAgentConfigurationLookupResponse wraps the response from a <todo> call
type APMAgentConfigurationLookupResponse struct {
	StatusCode int
	Body       *APMAgentConfigurationLookupResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type APMAgentConfigurationLookupResponseBody struct {
	ID    *string  `json:"_id,omitempty"`
	Index *string  `json:"_index,omitempty"`
	Score *float32 `json:"_score,omitempty"`
	// Source Agent configuration
	Source *APMUIAgentConfigurationObject `json:"_source,omitempty"`
}

type APMAgentConfigurationLookupRequest struct {
	Body APMAgentConfigurationLookupRequestBody
}

type APMAgentConfigurationLookupRequestBody struct {
	// Etag If etags match then `applied_by_agent` field will be set to `true`
	Etag *string `json:"etag,omitempty"`
	// MarkAsAppliedByAgent `markAsAppliedByAgent=true` means "force setting it to true regardless of etag".
	// This is needed for Jaeger agent that doesn't have etags
	MarkAsAppliedByAgent *bool            `json:"mark_as_applied_by_agent,omitempty"`
	Service              APMServiceObject `json:"service"`
}

// newAPMAgentConfigurationLookup returns a function that performs POST /api/apm/settings/agent-configuration/search API requests
func (api *API) newAPMAgentConfigurationLookup() func(context.Context, *APMAgentConfigurationLookupRequest, ...RequestOption) (*APMAgentConfigurationLookupResponse, error) {
	return func(ctx context.Context, req *APMAgentConfigurationLookupRequest, opts ...RequestOption) (*APMAgentConfigurationLookupResponse, error) {
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
			newCtx = instrument.Start(ctx, "apm.agent_configuration.lookup")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/apm/settings/agent-configuration/search"

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
			instrument.BeforeRequest(httpReq, "apm.agent_configuration.lookup")
			if reader := instrument.RecordRequestBody(ctx, "apm.agent_configuration.lookup", httpReq.Body); reader != nil {
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
		resp := &APMAgentConfigurationLookupResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result APMAgentConfigurationLookupResponseBody

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
