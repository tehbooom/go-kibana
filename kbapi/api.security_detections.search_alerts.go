package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"io"
	"net/http"
)

// SecurityDetectionsSearchAlertsResponse wraps the response from a SearchAlerts call
type SecurityDetectionsSearchAlertsResponse struct {
	StatusCode int
	Body       *SecurityDetectionsSearchAlertsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsSearchAlertsResponseBody struct {
	Aggregations map[string]types.Aggregate `json:"aggregations,omitempty"`
	Clusters_    *types.ClusterStatistics   `json:"_clusters,omitempty"`
	Fields       map[string]json.RawMessage `json:"fields,omitempty"`
	// Hits The returned documents and metadata.
	Hits            types.HitsMetadata `json:"hits"`
	MaxScore        *types.Float64     `json:"max_score,omitempty"`
	NumReducePhases *int64             `json:"num_reduce_phases,omitempty"`
	PitId           *string            `json:"pit_id,omitempty"`
	Profile         *types.Profile     `json:"profile,omitempty"`
	// ScrollId_ The identifier for the search and its search context.
	// You can use this scroll ID with the scroll API to retrieve the next batch of
	// search results for the request.
	// This property is returned only if the `scroll` query parameter is specified
	// in the request.
	ScrollId_ *string `json:"_scroll_id,omitempty"`
	// Shards_ A count of shards used for the request.
	Shards_         types.ShardStatistics      `json:"_shards"`
	Suggest         map[string][]types.Suggest `json:"suggest,omitempty"`
	TerminatedEarly *bool                      `json:"terminated_early,omitempty"`
	// TimedOut If `true`, the request timed out before completion; returned results may be
	// partial or empty.
	TimedOut bool `json:"timed_out"`
	// Took The number of milliseconds it took Elasticsearch to run the request.
	// This value is calculated by measuring the time elapsed between receipt of a
	// request on the coordinating node and the time at which the coordinating node
	// is ready to send the response.
	// It includes:
	//
	// * Communication time between the coordinating node and data nodes
	// * Time the request spends in the search thread pool, queued for execution
	// * Actual run time
	//
	// It does not include:
	//
	// * Time needed to send the request to Elasticsearch
	// * Time needed to serialize the JSON response
	// * Time needed to send the response to a client
	Took int64 `json:"took"`
}

type SecurityDetectionsSearchAlertsRequest struct {
	Body SecurityDetectionsSearchAlertsRequestBody
}

type SecurityDetectionsSearchAlertsRequestBody struct {
	// Aggregations Defines the aggregations that are run as part of the search request.
	Aggregations map[string]types.Aggregations `json:"aggregations,omitempty"`
	// Fields An array of wildcard (`*`) field patterns.
	// The request returns values for field names matching these patterns in the
	// `hits.fields` property of the response.
	Fields []string `json:"fields,omitempty"`
	// Query The search definition using the Query DSL.
	Query *types.Query `json:"query,omitempty"`
	// RuntimeMappings One or more runtime fields in the search request.
	// These fields take precedence over mapped fields with the same name.
	RuntimeMappings types.RuntimeFields `json:"runtime_mappings,omitempty"`
	// Size The number of hits to return, which must not be negative.
	// By default, you cannot page through more than 10,000 hits using the `from`
	// and `size` parameters.
	// To page through more hits, use the `search_after` property.
	Size *int `json:"size,omitempty"`
	// Sort A comma-separated list of <field>:<direction> pairs.
	Sort []types.SortCombinations `json:"sort,omitempty"`
	// TrackTotalHits Number of hits matching the query to count accurately.
	// If `true`, the exact number of hits is returned at the cost of some
	// performance.
	// If `false`, the  response does not include the total number of hits matching
	// the query.
	TrackTotalHits bool `json:"track_total_hits,omitempty"`
}

// newSecurityDetectionsSearchAlerts returns a function that performs POST /api/detection_engine/signals/search API requests
func (api *API) newSecurityDetectionsSearchAlerts() func(context.Context, *SecurityDetectionsSearchAlertsRequest, ...RequestOption) (*SecurityDetectionsSearchAlertsResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsSearchAlertsRequest, opts ...RequestOption) (*SecurityDetectionsSearchAlertsResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.search_alerts")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/signals/search"

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
			instrument.BeforeRequest(httpReq, "security_detections.search_alerts")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.search_alerts", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsSearchAlertsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsSearchAlertsResponseBody

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
