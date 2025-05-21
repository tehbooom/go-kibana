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
// FleetPackagePoliciesBulkDeleteResponse wraps the response from a <todo> call
type FleetPackagePoliciesBulkDeleteResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesBulkDeleteResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesBulkDeleteResponseBody struct {
	Body struct {
		Message string `json:"message"`
	} `json:"body"`

	ID string `json:"id"`

	// Name Package policy name (should be unique)
	Name string `json:"name"`

	OutputID *string `json:"output_id"`

	Package *struct {
		ExperimentalDataStreamFeatures *[]struct {
			DataStream string `json:"data_stream"`
			Features   struct {
				DocValueOnlyNumeric *bool `json:"doc_value_only_numeric,omitempty"`
				DocValueOnlyOther   *bool `json:"doc_value_only_other,omitempty"`
				SyntheticSource     *bool `json:"synthetic_source,omitempty"`
				TSDB                *bool `json:"tsdb,omitempty"`
			} `json:"features"`
		} `json:"experimental_data_stream_features,omitempty"`

		// Name Package name
		Name         string  `json:"name"`
		RequiresRoot *bool   `json:"requires_root,omitempty"`
		Title        *string `json:"title,omitempty"`

		// Version Package version
		Version string `json:"version"`
	} `json:"package,omitempty"`

	PolicyID  *string   `json:"policy_id,omitempty"`
	PolicyIDs *[]string `json:"policy_ids,omitempty"`

	StatusCode *float32 `json:"statusCode"`
	Success    *bool    `json:"success"`
}

type FleetPackagePoliciesBulkDeleteRequest struct {
	Body FleetPackagePoliciesBulkDeleteRequestBody
}

type FleetPackagePoliciesBulkDeleteRequestBody struct {
	// Ids list of package policy ids
	IDs   []string `json:"ids"`
	Force *bool    `json:"force,omitempty"`
}

// newFleetPackagePoliciesBulkDelete returns a function that performs POST /api/fleet/package_policies/delete API requests
func (api *API) newFleetPackagePoliciesBulkDelete() func(context.Context, *FleetPackagePoliciesBulkDeleteRequest, ...RequestOption) (*FleetPackagePoliciesBulkDeleteResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesBulkDeleteRequest, opts ...RequestOption) (*FleetPackagePoliciesBulkDeleteResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.bulk.delete")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/package_policies/delete"

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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.bulk.delete")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.bulk.delete", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesBulkDeleteResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesBulkDeleteResponseBody

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
