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
// FleetBinaryDownloadCreateResponse wraps the response from a <todo> call
type FleetBinaryDownloadCreateResponse struct {
	StatusCode int
	Body       *FleetBinaryDownloadCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetBinaryDownloadCreateResponseBody struct {
	Item FleetBinaryDownloadResponseItem `json:"item"`
}

type FleetBinaryDownloadCreateRequest struct {
	Body FleetBinaryDownloadCreateRequestBody
}

type FleetBinaryDownloadCreateRequestBody struct {
	Host      string `json:"host"`
	ID        string `json:"id"`
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name"`
	ProxyID   string `json:"proxy_id"`
}

// newFleetBinaryDownloadCreate returns a function that performs POST /api/fleet/agent_download_sources API requests
func (api *API) newFleetBinaryDownloadCreate() func(context.Context, *FleetBinaryDownloadCreateRequest, ...RequestOption) (*FleetBinaryDownloadCreateResponse, error) {
	return func(ctx context.Context, req *FleetBinaryDownloadCreateRequest, opts ...RequestOption) (*FleetBinaryDownloadCreateResponse, error) {
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
			newCtx = instrument.Start(ctx, "fleet.binary_download.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/agent_download_sources"

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
			instrument.BeforeRequest(httpReq, "fleet.binary_download.create")
			if reader := instrument.RecordRequestBody(ctx, "fleet.binary_download.create", httpReq.Body); reader != nil {
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
		resp := &FleetBinaryDownloadCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetBinaryDownloadCreateResponseBody

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
