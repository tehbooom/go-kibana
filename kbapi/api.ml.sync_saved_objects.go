package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// MLSyncSavedObjects  wraps the response from a FleetEPMBulkGetAssets  call
type MLSyncSavedObjectsResponse struct {
	StatusCode int
	Body       *MLSyncSavedObjectsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type MLSyncSavedObjectsResponseBody struct {
	// DatafeedsAdded If a saved object for an anomaly detection job is missing a datafeed identifier, it is added when you run the sync machine learning saved objects API.
	DatafeedsAdded *map[string]MachineLearningAPIsMlSyncResponseDatafeeds `json:"datafeedsAdded,omitempty"`

	// DatafeedsRemoved If a saved object for an anomaly detection job references a datafeed that no longer exists, it is deleted when you run the sync machine learning saved objects API.
	DatafeedsRemoved *map[string]MachineLearningAPIsMlSyncResponseDatafeeds `json:"datafeedsRemoved,omitempty"`

	// SavedObjectsCreated If saved objects are missing for machine learning jobs or trained models, they are created when you run the sync machine learning saved objects API.
	SavedObjectsCreated *MachineLearningAPIsMlSyncResponseSavedObjectsCreated `json:"savedObjectsCreated,omitempty"`

	// SavedObjectsDeleted If saved objects exist for machine learning jobs or trained models that no longer exist, they are deleted when you run the sync machine learning saved objects API.
	SavedObjectsDeleted *MachineLearningAPIsMlSyncResponseSavedObjectsDeleted `json:"savedObjectsDeleted,omitempty"`
}

// MachineLearningAPIsMlSyncResponseSavedObjectsCreated If saved objects are missing for machine learning jobs or trained models, they are created when you run the sync machine learning saved objects API.
type MachineLearningAPIsMlSyncResponseSavedObjectsCreated struct {
	// AnomalyDetector If saved objects are missing for anomaly detection jobs, they are created.
	AnomalyDetector *map[string]MachineLearningAPIsMlSyncResponseAnomalyDetectors `json:"anomaly-detector,omitempty"`

	// DataFrameAnalytics If saved objects are missing for data frame analytics jobs, they are created.
	DataFrameAnalytics *map[string]MachineLearningAPIsMlSyncResponseDataFrameAnalytics `json:"data-frame-analytics,omitempty"`

	// TrainedModel If saved objects are missing for trained models, they are created.
	TrainedModel *map[string]MachineLearningAPIsMlSyncResponseTrainedModels `json:"trained-model,omitempty"`
}

// MachineLearningAPIsMlSyncResponseDatafeeds The sync machine learning saved objects API response contains this object when there are datafeeds affected by the synchronization. There is an object for each relevant datafeed, which contains the synchronization status.
type MachineLearningAPIsMlSyncResponseDatafeeds struct {
	// Success The success or failure of the synchronization.
	Success *bool `json:"success,omitempty"`
}

// MachineLearningAPIsMlSyncResponseAnomalyDetectors The sync machine learning saved objects API response contains this object when there are anomaly detection jobs affected by the synchronization. There is an object for each relevant job, which contains the synchronization status.
type MachineLearningAPIsMlSyncResponseAnomalyDetectors struct {
	// Success The success or failure of the synchronization.
	Success *bool `json:"success,omitempty"`
}

type MLSyncSavedObjectsRequestParams struct {
	// Simulate When true, simulates the synchronization by returning only the list of actions that would be performed.
	Simulate *bool `form:"simulate,omitempty" json:"simulate,omitempty"`
}

// MachineLearningAPIsMlSyncResponseSavedObjectsDeleted If saved objects exist for machine learning jobs or trained models that no longer exist, they are deleted when you run the sync machine learning saved objects API.
type MachineLearningAPIsMlSyncResponseSavedObjectsDeleted struct {
	// AnomalyDetector If there are saved objects exist for nonexistent anomaly detection jobs, they are deleted.
	AnomalyDetector *map[string]MachineLearningAPIsMlSyncResponseAnomalyDetectors `json:"anomaly-detector,omitempty"`

	// DataFrameAnalytics If there are saved objects exist for nonexistent data frame analytics jobs, they are deleted.
	DataFrameAnalytics *map[string]MachineLearningAPIsMlSyncResponseDataFrameAnalytics `json:"data-frame-analytics,omitempty"`

	// TrainedModel If there are saved objects exist for nonexistent trained models, they are deleted.
	TrainedModel *map[string]MachineLearningAPIsMlSyncResponseTrainedModels `json:"trained-model,omitempty"`
}

// MachineLearningAPIsMlSyncResponseTrainedModels The sync machine learning saved objects API response contains this object when there are trained models affected by the synchronization. There is an object for each relevant trained model, which contains the synchronization status.
type MachineLearningAPIsMlSyncResponseTrainedModels struct {
	// Success The success or failure of the synchronization.
	Success *bool `json:"success,omitempty"`
}

// MachineLearningAPIsMlSyncResponseDataFrameAnalytics The sync machine learning saved objects API response contains this object when there are data frame analytics jobs affected by the synchronization. There is an object for each relevant job, which contains the synchronization status.
type MachineLearningAPIsMlSyncResponseDataFrameAnalytics struct {
	// Success The success or failure of the synchronization.
	Success *bool `json:"success,omitempty"`
}

// MLSyncSavedObjectsRequest   is the request for newFleetBulkGetAgentPolicies
type MLSyncSavedObjectsRequest struct {
	Params MLSyncSavedObjectsRequestParams
}

// newMLSyncSavedObjects returns a function that performs POST /api/saved_objects/_export API requests
func (api *API) newMLSyncSavedObjects() func(context.Context, *MLSyncSavedObjectsRequest, ...RequestOption) (*MLSyncSavedObjectsResponse, error) {
	return func(ctx context.Context, req *MLSyncSavedObjectsRequest, opts ...RequestOption) (*MLSyncSavedObjectsResponse, error) {
		if req == nil {
			req = &MLSyncSavedObjectsRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "ml.sync_saved_objects")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Simulate != nil {
			params["simulate"] = strconv.FormatBool(*req.Params.Simulate)
		}

		path := "/api/saved_objects/_export"

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
			instrument.BeforeRequest(httpReq, "ml.sync_saved_objects")
			if reader := instrument.RecordRequestBody(ctx, "ml.sync_saved_objects", httpReq.Body); reader != nil {
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
		resp := &MLSyncSavedObjectsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result MLSyncSavedObjectsResponseBody

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
