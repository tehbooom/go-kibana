package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// FleetDownloadAgentPolicyResponse wraps the response from a FleetDownloadAgentPolicy call
type FleetDownloadAgentPolicyResponse struct {
	StatusCode int
	Body       *string
	Error      interface{}
	RawBody    io.ReadCloser
}

// FleetDownloadAgentPolicyRequest is the request for newFleetDownloadAgentPolicy
type FleetDownloadAgentPolicyRequest struct {
	// ID of agent policy
	ID     string
	Params FleetDownloadAgentPolicyRequestParams
}

type FleetDownloadAgentPolicyRequestParams struct {
	Download   *bool `form:"download,omitempty" json:"download,omitempty"`
	Standalone *bool `form:"standalone,omitempty" json:"standalone,omitempty"`
	Kubernetes *bool `form:"kubernetes,omitempty" json:"kubernetes,omitempty"`
}

// newFleetDownloadAgentPolicy returns a function that performs GET /api/fleet/agent_policies/{agentPolicyId}/download API requests
func (api *API) newFleetDownloadAgentPolicy() func(context.Context, *FleetDownloadAgentPolicyRequest, ...RequestOption) (*FleetDownloadAgentPolicyResponse, error) {
	return func(ctx context.Context, req *FleetDownloadAgentPolicyRequest, opts ...RequestOption) (*FleetDownloadAgentPolicyResponse, error) {
		if req.ID == "" {
			return nil, fmt.Errorf("Required Agent Policy ID is not defined")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.agent_policies.download")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/fleet/agent_policies/%s/download", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Download != nil {
			params["download"] = strconv.FormatBool(*req.Params.Download)
		}
		if req.Params.Standalone != nil {
			params["standalone"] = strconv.FormatBool(*req.Params.Standalone)
		}
		if req.Params.Kubernetes != nil {
			params["kubernetes"] = strconv.FormatBool(*req.Params.Kubernetes)
		}

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
			instrument.BeforeRequest(httpReq, "fleet.agent_policies.download")
			if reader := instrument.RecordRequestBody(ctx, "fleet.agent_policies.download", httpReq.Body); reader != nil {
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
		resp := &FleetDownloadAgentPolicyResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		httpResp.Body.Close()

		if httpResp.StatusCode == 200 {
			resp.Body = StrPtr(string(bodyBytes))
			return resp, nil
		} else {
			// For all non-200 responses

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

func (d *FleetDownloadAgentPolicyResponse) WriteToFile(filepath string) error {
	if d.Body == nil {
		return fmt.Errorf("response body is nil, nothing to write")
	}

	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(*d.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
