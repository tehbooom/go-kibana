package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// StatusResponse wraps the response from a StatusGet call
type StatusResponse struct {
	StatusCode int
	Body       *KibanaStatusResponse
	Error      interface{}
	RawBody    io.ReadCloser
}

// StatusRedactedResponse wraps the response from a StatusGetRedacted call
type StatusRedactedResponse struct {
	StatusCode int
	Body       *KibanaStatusRedactedResponse
	Error      interface{}
	RawBody    io.ReadCloser
}

// newStatusFunc returns a function that performs status API requests
func (api *API) newStatusFunc() func(context.Context, *GetStatusRequest, ...RequestOption) (*StatusResponse, error) {
	return func(ctx context.Context, req *GetStatusRequest, opts ...RequestOption) (*StatusResponse, error) {
		if req == nil {
			req = &GetStatusRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/status"

		// Build query parameters
		params := make(map[string]string)

		if req.V7format != nil {
			params["v7format"] = strconv.FormatBool(*req.V7format)
		}
		if req.V8format != nil {
			params["v8format"] = strconv.FormatBool(*req.V8format)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
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
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "status")
			if reader := instrument.RecordRequestBody(ctx, "status", httpReq.Body); reader != nil {
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
		resp := &StatusResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		// Parse body based on request type
		var result KibanaStatusResponse

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

// newStatusFunc returns a function that performs status API requests
func (api *API) newStatusRedactedFunc() func(context.Context, *GetStatusRequest, ...RequestOption) (*StatusRedactedResponse, error) {
	return func(ctx context.Context, req *GetStatusRequest, opts ...RequestOption) (*StatusRedactedResponse, error) {
		if req == nil {
			req = &GetStatusRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "status")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/status"

		// Build query parameters
		params := make(map[string]string)
		if req.V7format != nil {
			params["v7format"] = strconv.FormatBool(*req.V7format)
		}
		if req.V8format != nil {
			params["v8format"] = strconv.FormatBool(*req.V8format)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
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
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "status")
			if reader := instrument.RecordRequestBody(ctx, "status", httpReq.Body); reader != nil {
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
		resp := &StatusRedactedResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result KibanaStatusRedactedResponse

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

// KibanaHTTPAPIsCoreStatusRedactedResponse A minimal representation of Kibana's operational status.
type KibanaStatusRedactedResponse struct {
	Status struct {
		Overall struct {
			// Level Service status levels as human and machine readable values.
			Level string `json:"level"`
		} `json:"overall"`
	} `json:"status"`
}

// KibanaHTTPAPIsCoreStatusResponse Kibana's operational status as well as a detailed breakdown of plugin statuses indication of various loads (like event loop utilization and network traffic) at time of request.
type KibanaStatusResponse struct {
	// Metrics Metric groups collected by Kibana.
	Metrics struct {
		// CollectionIntervalInMillis The interval at which metrics should be collected.
		CollectionIntervalInMillis float32 `json:"collection_interval_in_millis"`

		// ElasticsearchClient Current network metrics of Kibana's Elasticsearch client.
		ElasticsearchClient struct {
			// TotalActiveSockets Count of network sockets currently in use.
			TotalActiveSockets float32 `json:"totalActiveSockets"`

			// TotalIdleSockets Count of network sockets currently idle.
			TotalIdleSockets float32 `json:"totalIdleSockets"`

			// TotalQueuedRequests Count of requests not yet assigned to sockets.
			TotalQueuedRequests float32 `json:"totalQueuedRequests"`
		} `json:"elasticsearch_client"`

		// LastUpdated The time metrics were collected.
		LastUpdated string `json:"last_updated"`
	} `json:"metrics"`

	// Name Kibana instance name.
	Name   string `json:"name"`
	Status struct {
		// Core Statuses of core Kibana services.
		Core struct {
			Elasticsearch struct {
				// Detail Human readable detail of the service status.
				Detail *string `json:"detail,omitempty"`

				// DocumentationUrl A URL to further documentation regarding this service.
				DocumentationUrl *string `json:"documentationUrl,omitempty"`

				// Level Service status levels as human and machine readable values.
				Level string `json:"level"`

				// Meta An unstructured set of extra metadata about this service.
				Meta map[string]interface{} `json:"meta"`

				// Summary A human readable summary of the service status.
				Summary string `json:"summary"`
			} `json:"elasticsearch"`
			SavedObjects struct {
				// Detail Human readable detail of the service status.
				Detail *string `json:"detail,omitempty"`

				// DocumentationUrl A URL to further documentation regarding this service.
				DocumentationUrl *string `json:"documentationUrl,omitempty"`

				// Level Service status levels as human and machine readable values.
				Level string `json:"level"`

				// Meta An unstructured set of extra metadata about this service.
				Meta map[string]interface{} `json:"meta"`

				// Summary A human readable summary of the service status.
				Summary string `json:"summary"`
			} `json:"savedObjects"`
		} `json:"core"`
		Overall struct {
			// Detail Human readable detail of the service status.
			Detail *string `json:"detail,omitempty"`

			// DocumentationUrl A URL to further documentation regarding this service.
			DocumentationUrl *string `json:"documentationUrl,omitempty"`

			// Level Service status levels as human and machine readable values.
			Level string `json:"level"`

			// Meta An unstructured set of extra metadata about this service.
			Meta map[string]interface{} `json:"meta"`

			// Summary A human readable summary of the service status.
			Summary string `json:"summary"`
		} `json:"overall"`

		// Plugins A dynamic mapping of plugin ID to plugin status.
		Plugins map[string]struct {
			// Detail Human readable detail of the service status.
			Detail *string `json:"detail,omitempty"`

			// DocumentationUrl A URL to further documentation regarding this service.
			DocumentationUrl *string `json:"documentationUrl,omitempty"`

			// Level Service status levels as human and machine readable values.
			Level string `json:"level"`

			// Meta An unstructured set of extra metadata about this service.
			Meta map[string]interface{} `json:"meta"`

			// Summary A human readable summary of the service status.
			Summary string `json:"summary"`
		} `json:"plugins"`
	} `json:"status"`

	// Uuid Unique, generated Kibana instance UUID. This UUID should persist even if the Kibana process restarts.
	Uuid    string `json:"uuid"`
	Version struct {
		// BuildDate The date and time of this build.
		BuildDate string `json:"build_date"`

		// BuildFlavor The build flavour determines configuration and behavior of Kibana. On premise users will almost always run the "traditional" flavour, while other flavours are reserved for Elastic-specific use cases.
		BuildFlavor string `json:"build_flavor"`

		// BuildHash A unique hash value representing the git commit of this Kibana build.
		BuildHash string `json:"build_hash"`

		// BuildNumber A monotonically increasing number, each subsequent build will have a higher number.
		BuildNumber float32 `json:"build_number"`

		// BuildSnapshot Whether this build is a snapshot build.
		BuildSnapshot bool `json:"build_snapshot"`

		// Number A semantic version number.
		Number string `json:"number"`
	} `json:"version"`
}

// GetStatusParams defines parameters for GetStatus.
type GetStatusRequest struct {
	// V7format Set to "true" to get the response in v7 format.
	V7format *bool `form:"v7format,omitempty" json:"v7format,omitempty"`

	// V8format Set to "true" to get the response in v8 format.
	V8format *bool `form:"v8format,omitempty" json:"v8format,omitempty"`
}
