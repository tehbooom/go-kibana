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
// FleetPackagePoliciesUpgradeDryRunResponse wraps the response from a <todo> call
type FleetPackagePoliciesUpgradeDryRunResponse struct {
	StatusCode int
	Body       *FleetPackagePoliciesUpgradeDryRunResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetPackagePoliciesUpgradeDryRunResponseBody []struct {
	AgentDiff [][]struct {
		DataStream struct {
			Namespace string `json:"namespace"`
		} `json:"data_stream"`
		ID   string `json:"id"`
		Meta struct {
			Package struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"package"`
		} `json:"meta"`
		Name            string `json:"name"`
		PackagePolicyID string `json:"package_policy_id"`
		Processors      []struct {
			AddFields struct {
				Fields map[string]interface{} `json:"fields"`
				Target string                 `json:"target"`
			} `json:"add_fields"`
		} `json:"processors"`
		Revision float64 `json:"revision"`
		Streams  []struct {
			DataStream struct {
				Dataset string `json:"dataset"`
				Type    string `json:"type"`
			} `json:"data_stream"`
			ID string `json:"id"`
		} `json:"streams"`
		Type      string `json:"type"`
		UseOutput string `json:"use_output"`
	} `json:"agent_diff"`
	Body struct {
		Message string `json:"message"`
	} `json:"body"`
	Diff       []PackagePolicy `json:"diff"`
	HasErrors  bool            `json:"hasErrors"`
	Name       string          `json:"name"`
	StatusCode float64         `json:"statusCode"`
}

type FleetPackagePoliciesUpgradeDryRunRequest struct {
	Body FleetPackagePoliciesUpgradeDryRunRequestBody
}

type FleetPackagePoliciesUpgradeDryRunRequestBody struct {
	PackagePolicyIDs []string `json:"packagePolicyIds"`
	PackageVersion   string   `json:"packageVersion"`
}

// newFleetPackagePoliciesUpgradeDryRun returns a function that performs POST /api/fleet/package_policies/upgrade/dryrun API requests
func (api *API) newFleetPackagePoliciesUpgradeDryRun() func(context.Context, *FleetPackagePoliciesUpgradeDryRunRequest, ...RequestOption) (*FleetPackagePoliciesUpgradeDryRunResponse, error) {
	return func(ctx context.Context, req *FleetPackagePoliciesUpgradeDryRunRequest, opts ...RequestOption) (*FleetPackagePoliciesUpgradeDryRunResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.package_policies.upgrade_dry_run")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/package_policies/upgrade/dryrun"

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
			instrument.BeforeRequest(httpReq, "fleet.package_policies.upgrade_dry_run")
			if reader := instrument.RecordRequestBody(ctx, "fleet.package_policies.upgrade_dry_run", httpReq.Body); reader != nil {
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
		resp := &FleetPackagePoliciesUpgradeDryRunResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetPackagePoliciesUpgradeDryRunResponseBody

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
