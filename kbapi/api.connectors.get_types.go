package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// ConnectorsGetTypesResponse wraps the response from a <todo> call
type ConnectorsGetTypesResponse struct {
	StatusCode int
	Body       *ConnectorsGetTypesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type ConnectorsGetTypesResponseBody []ConnectorsGetTypesResponseBodyItem

type ConnectorsGetTypesResponseBodyItem struct {
	ID                     string   `json:"id"`
	Name                   string   `json:"name"`
	Enabled                bool     `json:"enabled"`
	EnabledInConfig        bool     `json:"enabled_in_config"`
	EnabledInLicense       bool     `json:"enabled_in_license"`
	IsSystemActionType     bool     `json:"is_system_action_type"`
	SupportedFeatureIDs    []string `json:"supported_feature_ids"`
	MinimumLicenseRequired string   `json:"minimum_license_required"`
}

type ConnectorsGetTypesRequest struct {
	Params ConnectorsGetTypesRequestParams
}

type ConnectorsGetTypesRequestParams struct {
	// FeatureId A filter to limit the retrieved connector types to those that support a specific feature (such as alerting or cases).
	FeatureId *string `form:"feature_id,omitempty" json:"feature_id,omitempty"`
}

// newConnectorsGetTypes returns a function that performs GET /api/actions/connector_types API requests
func (api *API) newConnectorsGetTypes() func(context.Context, *ConnectorsGetTypesRequest, ...RequestOption) (*ConnectorsGetTypesResponse, error) {
	return func(ctx context.Context, req *ConnectorsGetTypesRequest, opts ...RequestOption) (*ConnectorsGetTypesResponse, error) {
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
			newCtx = instrument.Start(ctx, "connectors.get_types")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/actions/connector_types"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.FeatureId != nil {
			params["feature_id"] = *req.Params.FeatureId
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
			instrument.BeforeRequest(httpReq, "connectors.get_types")
			if reader := instrument.RecordRequestBody(ctx, "connectors.get_types", httpReq.Body); reader != nil {
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
		resp := &ConnectorsGetTypesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result ConnectorsGetTypesResponseBody

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
