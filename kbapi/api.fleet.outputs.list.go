package kbapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// FleetOutputsListResponse wraps the response from a <todo> call
type FleetOutputsListResponse struct {
	StatusCode int
	Body       *FleetOutputsListResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type FleetOutputsListResponseBody struct {
	Items   []FleetOutputsResponseBodyItem `json:"items"`
	Page    int                            `json:"page"`
	PerPage int                            `json:"perPage"`
	Total   int                            `json:"total"`
}

// GetOutputByID finds an output by its ID
func (body *FleetOutputsListResponseBody) GetOutputByID(id string) *FleetOutputsResponseBodyItem {
	for i := range body.Items {
		if body.Items[i].ID == id {
			return &body.Items[i]
		}
	}
	return nil
}

// GetOutputsByType filters outputs by type
func (body *FleetOutputsListResponseBody) GetOutputsByType(outputType string) []FleetOutputsResponseBodyItem {
	var results []FleetOutputsResponseBodyItem
	for _, item := range body.Items {
		if item.Type == outputType {
			results = append(results, item)
		}
	}
	return results
}

// newFleetOutputsList returns a function that performs GET /api/fleet/outputs API requests
func (api *API) newFleetOutputsList() func(context.Context, ...RequestOption) (*FleetOutputsListResponse, error) {
	return func(ctx context.Context, opts ...RequestOption) (*FleetOutputsListResponse, error) {
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "fleet.outputs.list")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/fleet/outputs"

		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		if instrument != nil {
			instrument.BeforeRequest(httpReq, "fleet.outputs.list")
			if reader := instrument.RecordRequestBody(ctx, "fleet.outputs.list", httpReq.Body); reader != nil {
				httpReq.Body = reader
			}
		}

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

		resp := &FleetOutputsListResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result FleetOutputsListResponseBody

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
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
