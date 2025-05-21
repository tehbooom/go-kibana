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
// FleetPackagePoliciesBulkGetResponse wraps the response from a <todo> call
type FleetPackagePoliciesBulkGetResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesBulkGetResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesBulkGetResponseBody struct {
	Items []PackagePolicy `json:"items"`
}

type FleetPackagePoliciesBulkGetRequest struct {
	Params FleetPackagePoliciesBulkGetRequestParams
	Body   FleetPackagePoliciesBulkGetRequestBody
}

type FleetPackagePoliciesBulkGetRequestParams struct {
	// Values are simplified or legacy
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

type FleetPackagePoliciesBulkGetRequestBody struct {
	// Ids list of package policy ids
	IDs           []string `json:"ids"`
	IgnoreMissing *bool    `json:"ignoreMissing,omitempty"`
}

// newFleetPackagePoliciesBulkGet returns a function that performs POST /api/fleet/package_policies/_bulk_get API requests
func (api *API) newFleetPackagePoliciesBulkGet() func(context.Context, *FleetPackagePoliciesBulkGetRequest, ...RequestOption) (*FleetPackagePoliciesBulkGetResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesBulkGetRequest, opts ...RequestOption) (*FleetPackagePoliciesBulkGetResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.bulk.get")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/package_policies/_bulk_get"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.bulk.get")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.bulk.get", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesBulkGetResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesBulkGetResponseBody

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
