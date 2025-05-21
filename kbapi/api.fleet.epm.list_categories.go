package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// FleetEPMListPkgCategoriesResponse wraps the response from a FleetEPMListPkgCategories  call
type FleetEPMListPkgCategoriesResponse struct {
	StatusCode int
	Body       *FleetEPMListPkgCategoriesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMListPkgCategoriesResponseBody struct {
	Items    []FleetEPMPkgCategoryItem `json:"items"`
	Response []FleetEPMPkgCategoryItem `json:"response"`
}

type FleetEPMPkgCategoryItem struct {
	Count       float32 `json:"count"`
	Id          string  `json:"id"`
	ParentId    *string `json:"parent_id,omitempty"`
	ParentTitle *string `json:"parent_title,omitempty"`
	Title       string  `json:"title"`
}

// FleetEPMListPkgCategoriesRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMListPkgCategoriesRequest struct {
	Params FleetEPMListPkgCategoriesRequestParams
}

type FleetEPMListPkgCategoriesRequestParams struct {
	// Whether to include prerelease packages in categories count (e.g. beta, rc, preview)
	Prerelease             *bool `form:"prerelease,omitempty" json:"prerelease,omitempty"`
	IncludePolicyTemplates *bool `form:"include_policy_templates,omitempty" json:"include_policy_templates,omitempty"`
}

// newFleetEPMListPkgCategories returns a function that performs GET /api/fleet/epm/categories API requests
func (api *API) newFleetEPMListPkgCategories() func(context.Context, *FleetEPMListPkgCategoriesRequest, ...RequestOption) (*FleetEPMListPkgCategoriesResponse, error) {
	return func(ctx context.Context, req *FleetEPMListPkgCategoriesRequest, opts ...RequestOption) (*FleetEPMListPkgCategoriesResponse, error) {
		if req == nil {
			req = &FleetEPMListPkgCategoriesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.list_categories")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Prerelease != nil {
			params["prerelease"] = strconv.FormatBool(*req.Params.Prerelease)
		}
		if req.Params.IncludePolicyTemplates != nil {
			params["include_policy_templates"] = strconv.FormatBool(*req.Params.IncludePolicyTemplates)
		}

		path := "/api/fleet/epm/categories"

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
			instrument.BeforeRequest(httpReq, "fleet.epm.list_categories")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.list_categories", httpReq.Body); reader != nil {
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
		resp := &FleetEPMListPkgCategoriesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMListPkgCategoriesResponseBody

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
