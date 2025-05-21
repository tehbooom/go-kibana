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
// FleetPackagePoliciesListResponse wraps the response from a <todo> call
type FleetPackagePoliciesListResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesListResponseBody struct {
	Items   []PackagePolicy `json:"items"`
	Page    float32         `json:"page"`
	PerPage float32         `json:"perPage"`
	Total   float32         `json:"total"`
}

type FleetPackagePoliciesListRequest struct {
	Params FleetPackagePoliciesListRequestParams
}

type FleetPackagePoliciesListRequestParams struct {
	Page      *float32 `form:"page,omitempty" json:"page,omitempty"`
	PerPage   *float32 `form:"perPage,omitempty" json:"perPage,omitempty"`
	SortField *string  `form:"sortField,omitempty" json:"sortField,omitempty"`
	// Values are desc or asc.
	SortOrder       *string `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
	ShowUpgradeable *bool   `form:"showUpgradeable,omitempty" json:"showUpgradeable,omitempty"`
	Kuery           *string `form:"kuery,omitempty" json:"kuery,omitempty"`
	// Values are simplified or legacy.
	Format         *string `form:"format,omitempty" json:"format,omitempty"`
	WithAgentCount *bool   `form:"withAgentCount,omitempty" json:"withAgentCount,omitempty"`
}

// newFleetPackagePoliciesList returns a function that performs GET /api/fleet/package_policies API requests
func (api *API) newFleetPackagePoliciesList() func(context.Context, *FleetPackagePoliciesListRequest, ...RequestOption) (*FleetPackagePoliciesListResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesListRequest, opts ...RequestOption) (*FleetPackagePoliciesListResponse, error) {
		if req == nil {
			req = &FleetPackagePoliciesListRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.package_policies.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/package_policies"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.FormatFloat(float64(*req.Params.Page), 'f', -1, 32)
		}
		if req.Params.PerPage != nil {
			params["perPage"] = strconv.FormatFloat(float64(*req.Params.PerPage), 'f', -1, 32)
		}
		if req.Params.Kuery != nil {
			params["kuery"] = *req.Params.Kuery
		}
		if req.Params.SortOrder != nil {
			params["sortOrder"] = *req.Params.SortOrder
		}
		if req.Params.SortField != nil {
			params["sortField"] = *req.Params.SortField
		}
		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
		}
		if req.Params.ShowUpgradeable != nil {
			params["showUpgradeable"] = strconv.FormatBool(*req.Params.ShowUpgradeable)
		}

		if req.Params.WithAgentCount != nil {
			params["withAgentCount"] = strconv.FormatBool(*req.Params.WithAgentCount)
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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.list", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesListResponseBody

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
