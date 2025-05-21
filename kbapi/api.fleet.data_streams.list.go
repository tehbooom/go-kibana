package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetDataStreamsListResponse wraps the response from a FleetDataStreamsListResponse call
type FleetDataStreamsListResponse struct {
	StatusCode int
	Body       *FleetDataStreamsListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetDataStreamsListResponseBody struct {
	DataStreams []Datastream `json:"data_streams"`
}

type Datastream struct {
	Dashboards           []Dashboard `json:"dashboards"`
	Dataset              string      `json:"dataset"`
	Index                string      `json:"index"`
	LastActivityMs       float32     `json:"last_activity_ms"`
	Namespace            string      `json:"namespace"`
	Package              string      `json:"package"`
	PackageVersion       string      `json:"package_version"`
	SizeInBytes          float32     `json:"size_in_bytes"`
	SizeInBytesFormatted string      `json:"size_in_bytes_formatted"`
	Type                 string      `json:"type"`
}

type Dashboard struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// newFleetDataStreamsList returns a function that performs GET /api/fleet/data_streams API requests
func (api *API) newFleetDataStreamsList() func(context.Context, ...RequestOption) (*FleetDataStreamsListResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*FleetDataStreamsListResponse, error) {
		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.data_streams.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/data_streams"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.data_streams.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.data_streams.list", httpReq.Body); reader != nil {
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
		resp := &FleetDataStreamsListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetDataStreamsListResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
