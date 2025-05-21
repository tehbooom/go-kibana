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
// FleetOutputsCreateResponse wraps the response from a <todo> call
type FleetOutputsCreateResponse struct {
	StatusCode int
	Body       *FleetOutputsCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetOutputsCreateResponseBody struct {
	Item FleetOutputsResponseBodyItem `json:"item"`
}

type FleetOutputsCreateRequest struct {
	Body json.RawMessage
}

// NewElasticsearchOutputRequest creates a request body for an Elasticsearch output
func NewElasticsearchOutputRequest(output *ElasticsearchOutput) (*FleetOutputsCreateRequest, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	return &FleetOutputsCreateRequest{Body: data}, nil
}

// NewLogstashOutputRequest creates a request body for a Logstash output
func NewLogstashOutputRequest(output *LogstashOutput) (*FleetOutputsCreateRequest, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	return &FleetOutputsCreateRequest{Body: data}, nil
}

// NewKafkaOutputRequest creates a request body for a Kafka output
func NewKafkaOutputRequest(output *KafkaOutput) (*FleetOutputsCreateRequest, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	return &FleetOutputsCreateRequest{Body: data}, nil
}

// NewRemoteElasticsearchOutputRequest creates a request body for a remote Elasticsearch output
func NewRemoteElasticsearchOutputRequest(output *RemoteElasticsearchOutput) (*FleetOutputsCreateRequest, error) {
	data, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	return &FleetOutputsCreateRequest{Body: data}, nil
}

// newFleetOutputsCreate returns a function that performs POST /api/fleet/outputs API requests
func (api *API) newFleetOutputsCreate() func(context.Context, *FleetOutputsCreateRequest, ...RequestOption) (*FleetOutputsCreateResponse, error) {
	return func(ctx context.Context, req *FleetOutputsCreateRequest, opts ...RequestOption) (*FleetOutputsCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.outputs.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/outputs"

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
			instrument.BeforeRequest(httpReq, "fleet.outputs.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.outputs.create", httpReq.Body); reader != nil {
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
		resp := &FleetOutputsCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetOutputsCreateResponseBody

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
