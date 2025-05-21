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
// FleetInternalUpdateSettingsResponse wraps the response from a fleet.internal.update_settings call
type FleetInternalUpdateSettingsResponse struct {
	StatusCode int
	Body       *FleetInternalUpdateSettingsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetInternalUpdateSettingsResponseBody struct {
	Item struct {
		DeleteUnenrolledAgents *struct {
			Enabled         bool `json:"enabled"`
			IsPreconfigured bool `json:"is_preconfigured"`
		} `json:"delete_unenrolled_agents,omitempty"`
		HasSeenAddDataNotice                *bool     `json:"has_seen_add_data_notice,omitempty"`
		Id                                  string    `json:"id"`
		OutputSecretStorageRequirementsMet  *bool     `json:"output_secret_storage_requirements_met,omitempty"`
		PreconfiguredFields                 *[]string `json:"preconfigured_fields,omitempty"`
		PrereleaseIntegrationsEnabled       *bool     `json:"prerelease_integrations_enabled,omitempty"`
		SecretStorageRequirementsMet        *bool     `json:"secret_storage_requirements_met,omitempty"`
		UseSpaceAwarenessMigrationStartedAt *string   `json:"use_space_awareness_migration_started_at"`
		UseSpaceAwarenessMigrationStatus    *string   `json:"use_space_awareness_migration_status,omitempty"`
		Version                             *string   `json:"version,omitempty"`
	} `json:"item"`
}

type FleetInternalUpdateSettingsRequest struct {
	Body FleetInternalUpdateSettingsRequestBody
}

type FleetInternalUpdateSettingsRequestBody struct {
	AdditionalYamlConfig   *string `json:"additional_yaml_config,omitempty"`
	DeleteUnenrolledAgents *struct {
		Enabled         bool `json:"enabled"`
		IsPreconfigured bool `json:"is_preconfigured"`
	} `json:"delete_unenrolled_agents,omitempty"`
	HasSeenAddDataNotice          *bool     `json:"has_seen_add_data_notice,omitempty"`
	KibanaCaSha256                *string   `json:"kibana_ca_sha256,omitempty"`
	KibanaUrls                    *[]string `json:"kibana_urls,omitempty"`
	PrereleaseIntegrationsEnabled *bool     `json:"prerelease_integrations_enabled,omitempty"`
}

// newFleetInternalUpdateSettings returns a function that performs PUT /api/fleet/settings API requests
func (api *API) newFleetInternalUpdateSettings() func(context.Context, *FleetInternalUpdateSettingsRequest, ...RequestOption) (*FleetInternalUpdateSettingsResponse, error) {
	return func(ctx context.Context, req *FleetInternalUpdateSettingsRequest, opts ...RequestOption) (*FleetInternalUpdateSettingsResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.internal.update_settings")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/settings"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, path, nil)
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
			instrument.BeforeRequest(httpReq, "fleet.internal.update_settings")
			if reader := instrument.RecordRequestBody(ctx, "fleet.internal.update_settings", httpReq.Body); reader != nil {
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
		resp := &FleetInternalUpdateSettingsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetInternalUpdateSettingsResponseBody

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
