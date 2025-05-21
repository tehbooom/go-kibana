package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// TaskManagerHealthResponse wraps the response from a <todo> call
type TaskManagerHealthResponse struct {
	StatusCode int
	Body       *TaskManagerHealthResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type TaskManagerHealthResponseBody struct {
	ID         *string `json:"id,omitempty"`
	LastUpdate *string `json:"last_update,omitempty"`
	Stats      *struct {
		// CapacityEstimation This object provides a rough estimate about the sufficiency of its capacity. These are estimates based on historical data and should not be used as predictions.
		CapacityEstimation *map[string]interface{} `json:"capacity_estimation,omitempty"`

		// Configuration This object summarizes the current configuration of Task Manager. This includes dynamic configurations that change over time, such as `poll_interval` and `max_workers`, which can adjust in reaction to changing load on the system.
		Configuration *map[string]interface{} `json:"configuration,omitempty"`

		// Runtime This object tracks runtime performance of Task Manager, tracking task drift, worker load, and stats broken down by type, including duration and run results.
		Runtime *map[string]interface{} `json:"runtime,omitempty"`

		// Workload This object summarizes the work load across the cluster, including the tasks in the system, their types, and current status.
		Workload *map[string]interface{} `json:"workload,omitempty"`
	} `json:"stats,omitempty"`
	Status    *string `json:"status,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}

// newTaskManagerHealth returns a function that performs GET /api/task_manager/_health API requests
func (api *API) newTaskManagerHealth() func(context.Context, ...RequestOption) (*TaskManagerHealthResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*TaskManagerHealthResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "task_manager.health")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/task_manager/_health"

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
			instrument.BeforeRequest(httpReq, "task_manager.health")
			if reader := instrument.RecordRequestBody(ctx, "task_manager.health", httpReq.Body); reader != nil {
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
		resp := &TaskManagerHealthResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result TaskManagerHealthResponseBody

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
