package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// SecurityDetectionsGetPrivilegesSpaceResponse wraps the response from a <todo> call
type SecurityDetectionsGetPrivilegesSpaceResponse struct {
	StatusCode int
	Body       *SecurityDetectionsGetPrivilegesSpaceResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsGetPrivilegesSpaceResponseBody struct {
	Index struct {
		AlertsSecurityAlertsDefault struct {
			All               bool `json:"all"`
			Read              bool `json:"read"`
			Index             bool `json:"index"`
			Write             bool `json:"write"`
			Create            bool `json:"create"`
			Delete            bool `json:"delete"`
			Manage            bool `json:"manage"`
			Monitor           bool `json:"monitor"`
			CreateDoc         bool `json:"create_doc"`
			Maintenance       bool `json:"maintenance"`
			CreateIndex       bool `json:"create_index"`
			DeleteIndex       bool `json:"delete_index"`
			ViewIndexMetadata bool `json:"view_index_metadata"`
		} `json:".alerts-security.alerts-default"`
	} `json:"index,omitempty"`
	Cluster struct {
		All                  bool `json:"all"`
		Manage               bool `json:"manage"`
		Monitor              bool `json:"monitor"`
		ManageML             bool `json:"manage_ml"`
		MonitorML            bool `json:"monitor_ml"`
		ManageAPIKey         bool `json:"manage_api_key"`
		ManagePipeline       bool `json:"manage_pipeline"`
		ManageSecurity       bool `json:"manage_security"`
		ManageTransform      bool `json:"manage_transform"`
		MonitorTransform     bool `json:"monitor_transform"`
		ManageOwnAPIKey      bool `json:"manage_own_api_key"`
		ManageIndexTemplates bool `json:"manage_index_templates"`
	} `json:"cluster,omitempty"`
	Username         string                 `json:"username,omitempty"`
	Application      map[string]interface{} `json:"application,omitempty"`
	IsAuthenticated  bool                   `json:"is_authenticated"`
	HasAllRequested  bool                   `json:"has_all_requested,omitempty"`
	HasEncryptionKey bool                   `json:"has_encryption_key"`
}

// newSecurityDetectionsGetPrivilegesSpace returns a function that performs GET /api/detection_engine/privileges API requests
func (api *API) newSecurityDetectionsGetPrivilegesSpace() func(context.Context, ...RequestOption) (*SecurityDetectionsGetPrivilegesSpaceResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsGetPrivilegesSpaceResponse, error) {

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "security_detections.get_privileges_space")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/privileges"

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
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
			instrument.BeforeRequest(httpReq, "security_detections.get_privileges_space")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.get_privileges_space", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsGetPrivilegesSpaceResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsGetPrivilegesSpaceResponseBody

		if httpResp.StatusCode < 299 {
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
			// For all non-success responses
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
