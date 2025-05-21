package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// LogstashGetPipelineResponse wraps the response from a <todo> call
type LogstashGetPipelineResponse struct {
	StatusCode int
	Body       *LogstashGetPipelineResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type LogstashGetPipelineResponseBody struct {
	ID          string                   `json:"id"`
	Description string                   `json:"description"`
	Username    string                   `json:"username"`
	Pipeline    string                   `json:"pipeline"`
	Settings    LogstashPipelineSettings `json:"settings"`
}

type LogstashPipelineSettings struct {
	PipelineWorkers          *int    `json:"pipeline.workers,omitempty"`
	PipelineBatchSize        *int    `json:"pipeline.batch.size,omitempty"`
	PipelineBatchDelay       *int    `json:"pipeline.batch.delay,omitempty"`
	PipelineECSCompatibility *string `json:"pipeline.ecs_compatibility,omitempty"`
	PipelineOrdered          *bool   `json:"pipeline.ordered,omitempty"`
	QueueType                string  `json:"queue.type,omitempty"`
	QueueMaxBytes            *int    `json:"queue.max_bytes,omitempty"`
	QueueCheckpointWrites    *int    `json:"queue.checkpoint.writes,omitempty"`
}

type LogstashGetPipelineRequest struct {
	ID string
}

// newLogstashGetPipeline returns a function that performs GET /api/logstash/pipeline/{id} API requests
func (api *API) newLogstashGetPipeline() func(context.Context, *LogstashGetPipelineRequest, ...RequestOption) (*LogstashGetPipelineResponse, error) {
	return func(ctx context.Context, req *LogstashGetPipelineRequest, opts ...RequestOption) (*LogstashGetPipelineResponse, error) {
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
			newCtx = instrument.Start(ctx, "logstash.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/logstash/pipeline/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "logstash.get")
			if reader := instrument.RecordRequestBody(ctx, "logstash.get", httpReq.Body); reader != nil {
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
		resp := &LogstashGetPipelineResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result LogstashGetPipelineResponseBody

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
