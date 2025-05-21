package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// TODO: Update the call
// FleetEPMListDataStreamsResponse wraps the response from a <todo> call
type FleetEPMListDataStreamsResponse struct {
	StatusCode int
	Body       *FleetEPMListDataStreamsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMListDataStreamsResponseBody struct {
	Items []struct {
		Name string `json:"name"`
	} `json:"items"`
}

type FleetEPMListDataStreamsRequest struct {
	Params FleetEPMListDataStreamsRequestParams
}

type FleetEPMListDataStreamsRequestParams struct {
	// Values are logs, metrics, traces, synthetics, or profiling.
	Type         *string `form:"type,omitempty" json:"type,omitempty"`
	DatasetQuery *string `form:"datasetQuery,omitempty" json:"datasetQuery,omitempty"`
	// Values are asc or desc. Default value is asc.
	SortOrder         *string `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
	UncategorisedOnly *bool   `form:"uncategorisedOnly,omitempty" json:"uncategorisedOnly,omitempty"`
}

// newFleetEPMListDataStreams returns a function that performs GET /api/fleet/epm/data_streams API requests
func (api *API) newFleetEPMListDataStreams() func(context.Context, *FleetEPMListDataStreamsRequest, ...RequestOption) (*FleetEPMListDataStreamsResponse, error) {
	return func(ctx context.Context, req *FleetEPMListDataStreamsRequest, opts ...RequestOption) (*FleetEPMListDataStreamsResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.epm.list_datastreams")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/epm/data_streams"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Type != nil {
			params["type"] = *req.Params.Type
		}
		if req.Params.DatasetQuery != nil {
			params["datasetQuery"] = *req.Params.DatasetQuery
		}
		if req.Params.SortOrder != nil {
			params["sortOrder"] = *req.Params.SortOrder
		}
		if req.Params.UncategorisedOnly != nil {
			params["uncategorisedOnly"] = strconv.FormatBool(*req.Params.UncategorisedOnly)
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
			instrument.BeforeRequest(httpReq, "fleet.epm.list_datastreams")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.list_datastreams", httpReq.Body); reader != nil {
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
		resp := &FleetEPMListDataStreamsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMListDataStreamsResponseBody

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
