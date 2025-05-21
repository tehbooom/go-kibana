package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FleetEPMBulkGetAssetsResponse wraps the response from a FleetEPMBulkGetAssets  call
type FleetEPMBulkGetAssetsResponse struct {
	StatusCode int
	Body       *PostFleetEPMBulkAssetsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type PostFleetEPMBulkAssetsResponseBody struct {
	Items []struct {
		AppLink    *string `json:"appLink,omitempty"`
		Attributes struct {
			Description *string `json:"description,omitempty"`
			Service     *string `json:"service,omitempty"`
			Title       *string `json:"title,omitempty"`
		} `json:"attributes"`
		Id        string  `json:"id"`
		Type      string  `json:"type"`
		UpdatedAt *string `json:"updatedAt,omitempty"`
	} `json:"items"`
}

// FleetEPMBulkGetAssetsRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMBulkGetAssetsRequest struct {
	Body FleetEPMBulkGetAssetsRequestBody
}

type FleetEPMBulkGetAssetsRequestBody struct {
	AssetIDs []string `json:"assetIds"`
}

// newFleetEPMBulkGetAssets returns a function that performs POST /api/fleet/epm/bulk_assets API requests
func (api *API) newFleetEPMBulkGetAssets() func(context.Context, *FleetEPMBulkGetAssetsRequest, ...RequestOption) (*FleetEPMBulkGetAssetsResponse, error) {
	return func(ctx context.Context, req *FleetEPMBulkGetAssetsRequest, opts ...RequestOption) (*FleetEPMBulkGetAssetsResponse, error) {
		if req.Body.AssetIDs == nil {
			return nil, fmt.Errorf("Required Asset IDs is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.bulk.get_assets")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/epm/bulk_assets"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.bulk.get_assets")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.bulk.get_assets", httpReq.Body); reader != nil {
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
		resp := &FleetEPMBulkGetAssetsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result PostFleetEPMBulkAssetsResponseBody

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
