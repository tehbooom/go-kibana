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
// FleetPackagePoliciesUpdateResponse wraps the response from a <todo> call
type FleetPackagePoliciesUpdateResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesUpdateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesUpdateResponseBody struct {
	Item PackagePolicy `json:"item"`
}

type FleetPackagePoliciesUpdateRequest struct {
	PackagePolicyId string
	Params          FleetPackagePoliciesUpdateRequestParams
	Body            FleetPackagePoliciesUpdateRequestBody
}

type FleetPackagePoliciesUpdateRequestParams struct {
	// Values are simplified or legacy
	Format *string `form:"format,omitempty" json:"format,omitempty"`
}

type FleetPackagePoliciesUpdateRequestBody struct {
	// Description Package policy description
	Description *string `json:"description,omitempty"`
	Enabled     bool    `json:"enabled"`
	ID          string  `json:"id"`
	// Force package policy creation even if package is not verified, or if the agent policy is managed.
	Force *bool

	// Inputs Package policy inputs (see integration documentation to know what inputs are available)
	Inputs    []PackagePolicyInput `json:"inputs"`
	IsManaged *bool                `json:"is_managed,omitempty"`

	// Name Package policy name (should be unique)
	Name string `json:"name"`

	// Namespace The package policy namespace. Leave blank to inherit the agent policy's namespace.
	Namespace *string `json:"namespace,omitempty"`
	OutputID  *string `json:"output_id"`

	// Overrides Override settings that are defined in the package policy. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides *struct {
		Inputs *map[string]interface{} `json:"inputs,omitempty"`
	} `json:"overrides"`
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

	PolicyIDs *[]string `json:"policy_ids,omitempty"`
	// SupportsAgentless Indicates whether the package policy belongs to an agentless agent policy.
	SupportsAgentless *bool                   `json:"supports_agentless"`
	Vars              *map[string]interface{} `json:"vars,omitempty"`
}

// newFleetPackagePoliciesUpdate returns a function that performs PUT /api/fleet/package_policies/{packagePolicyId} API requests
func (api *API) newFleetPackagePoliciesUpdate() func(context.Context, *FleetPackagePoliciesUpdateRequest, ...RequestOption) (*FleetPackagePoliciesUpdateResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesUpdateRequest, opts ...RequestOption) (*FleetPackagePoliciesUpdateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.update")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/package_policies/%s", req.PackagePolicyId)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Format != nil {
			params["format"] = *req.Params.Format
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.update")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.update", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesUpdateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesUpdateResponseBody

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
