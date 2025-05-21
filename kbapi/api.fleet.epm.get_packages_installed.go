package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// FleetEPMGetInstalledPackagesResponse wraps the response from a FleetEPMGetInstalledPackages  call
type FleetEPMGetInstalledPackagesResponse struct {
	StatusCode int
	Body       *FleetEPMGetInstalledPackagesResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetEPMGetInstalledPackagesResponseBody struct {
	Items []struct {
		DataStreams []struct {
			Name  string `json:"name"`
			Title string `json:"title"`
		} `json:"dataStreams"`
		Description *string `json:"description,omitempty"`
		Icons       *[]struct {
			DarkMode *bool   `json:"dark_mode,omitempty"`
			Path     *string `json:"path,omitempty"`
			Size     *string `json:"size,omitempty"`
			Src      string  `json:"src"`
			Title    *string `json:"title,omitempty"`
			Type     *string `json:"type,omitempty"`
		} `json:"icons,omitempty"`
		Name    string  `json:"name"`
		Status  string  `json:"status"`
		Title   *string `json:"title,omitempty"`
		Version string  `json:"version"`
	} `json:"items"`
	SearchAfter []interface{} `json:"searchAfter,omitempty"`
	Total       float32       `json:"total"`
}

// FleetEPMGetInstalledPackagesRequest  is the request for newFleetBulkGetAgentPolicies
type FleetEPMGetInstalledPackagesRequest struct {
	Params FleetEPMGetInstalledPackagesRequestParams
}

type FleetEPMGetInstalledPackagesRequestParams struct {
	// Values are logs, metrics, traces, synthetics, or profiling.
	DataStreamType            *string       `form:"dataStreamType,omitempty" json:"dataStreamType,omitempty"`
	ShowOnlyActiveDataStreams *bool         `form:"showOnlyActiveDataStreams,omitempty" json:"showOnlyActiveDataStreams,omitempty"`
	NameQuery                 *string       `form:"nameQuery,omitempty" json:"nameQuery,omitempty"`
	SearchAfter               []interface{} `form:"searchAfter,omitempty" json:"searchAfter,omitempty"`
	// Default value is 15.
	PerPage *float64 `form:"perPage,omitempty" json:"perPage,omitempty"`
	// Values are asc or desc. Default value is asc.
	SortOrder *string `form:"sortOrder,omitempty" json:"sortOrder,omitempty"`
}

// newFleetEPMGetInstalledPackages returns a function that performs GET /api/fleet/epm/packages/installed API requests
func (api *API) newFleetEPMGetInstalledPackages() func(context.Context, *FleetEPMGetInstalledPackagesRequest, ...RequestOption) (*FleetEPMGetInstalledPackagesResponse, error) {
	return func(ctx context.Context, req *FleetEPMGetInstalledPackagesRequest, opts ...RequestOption) (*FleetEPMGetInstalledPackagesResponse, error) {
		if req == nil {
			req = &FleetEPMGetInstalledPackagesRequest{}
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.epm.get_packages_installed")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		params := make(map[string]string)

		if req.Params.DataStreamType != nil {
			params["dataStreamType"] = *req.Params.DataStreamType
		}
		if req.Params.ShowOnlyActiveDataStreams != nil {
			params["showOnlyActiveDataStreams"] = strconv.FormatBool(*req.Params.ShowOnlyActiveDataStreams)
		}
		if req.Params.NameQuery != nil {
			params["nameQuery"] = *req.Params.NameQuery
		}
		if req.Params.SearchAfter != nil {
			var searchAfterVal []string
			for _, value := range req.Params.SearchAfter {
				switch v := value.(type) {
				case string:
					searchAfterVal = append(searchAfterVal, v)
				case float32:
					searchAfterVal = append(searchAfterVal, strconv.FormatFloat(float64(v), 'f', -1, 32))
				case float64:
					searchAfterVal = append(searchAfterVal, strconv.FormatFloat(v, 'f', -1, 64))
				case int:
					searchAfterVal = append(searchAfterVal, strconv.Itoa(v))
				case bool:
					searchAfterVal = append(searchAfterVal, strconv.FormatBool(v))
				default:
					// For complex types, use JSON
					bytes, err := json.Marshal(v)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal searchAfter value: %v", err)
					}
					searchAfterVal = append(searchAfterVal, string(bytes))
				}
			}
			params["searchAfter"] = strings.Join(searchAfterVal, ",")
		}

		if req.Params.PerPage != nil {
			params["perPage"] = strconv.FormatFloat(float64(*req.Params.PerPage), 'f', -1, 32)
		}
		if req.Params.SortOrder != nil {
			params["sortOrder"] = *req.Params.SortOrder
		}

		path := "/api/fleet/epm/packages/installed"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				return nil, err
			}
		}

		// Add query parameters
		if len(params) > 0 {
			q := httpReq.URL.Query()
			for k, v := range params {
				q.Set(k, v)
			}
			httpReq.URL.RawQuery = q.Encode()
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.epm.get_packages_installed")
			if reader := instrument.RecordRequestBody(ctx, "fleet.epm.get_packages_installed", httpReq.Body); reader != nil {
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
		resp := &FleetEPMGetInstalledPackagesResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetEPMGetInstalledPackagesResponseBody

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
