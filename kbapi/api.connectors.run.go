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
// ConnectorsRunResponse wraps the response from a <todo> call
type ConnectorsRunResponse struct {
	StatusCode int
	Body       *ConnectorsRunResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type ConnectorsRunResponseBody struct {
	Data        map[string]interface{} `json:"data"`
	Status      string                 `json:"status"`
	ConnectorID string                 `json:"connector_id"`
}

type ConnectorsRunRequest struct {
	ID   string
	Body ConnectorsRunRequestBody
}

type ConnectorsRunRequestBody struct {
	Params json.RawMessage `json:"params"`
}

// SetRunAcknowledgeResolvePagerduty sets the params for RunAcknowledgeResolvePagerduty subaction for Pager Duty connectors
func (body *ConnectorsRunRequestBody) SetRunAcknowledgeResolvePagerduty(params RunAcknowledgeResolvePagerduty) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data

	return nil
}

// SetRunAddevent sets the params for AddEvent subaction for ServiceNow ITOM connectors
func (body *ConnectorsRunRequestBody) SetRunAddevent(params RunAddevent) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunClosealert sets the params for CloseAlert subaction for Opsgenie connectors
func (body *ConnectorsRunRequestBody) SetRunClosealert(params RunClosealert) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunCloseincident sets the params for CloseIncident subaction for ServiceNow ITSM connectors
func (body *ConnectorsRunRequestBody) SetRunCloseincident(params RunCloseincident) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunCreatealert sets the params for CreateAlert subaction for Opsgenie and TheHive connectors
func (body *ConnectorsRunRequestBody) SetRunCreatealert(params RunCreatealert) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunDocuments sets the params for indexing documents into Elasticsearch
func (body *ConnectorsRunRequestBody) SetRunDocuments(params RunDocuments) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunFieldsbyissuetype sets the params for FieldsByIssueType subaction for Jira connectors
func (body *ConnectorsRunRequestBody) SetRunFieldsbyissuetype(params RunFieldsbyissuetype) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunGetchoices sets the params for GetChoices subaction for ServiceNow connectors
func (body *ConnectorsRunRequestBody) SetRunGetchoices(params RunGetchoices) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunGetfields sets the params for GetFields subaction for various connectors
func (body *ConnectorsRunRequestBody) SetRunGetfields(params RunGetfields) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunGetincident sets the params for GetIncident subaction for various connectors
func (body *ConnectorsRunRequestBody) SetRunGetincident(params RunGetincident) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunIssue sets the params for Issue subaction for Jira connectors
func (body *ConnectorsRunRequestBody) SetRunIssue(params RunIssue) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunIssues sets the params for Issues subaction for Jira connectors
func (body *ConnectorsRunRequestBody) SetRunIssues(params RunIssues) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunIssuetypes sets the params for IssueTypes subaction for Jira connectors
func (body *ConnectorsRunRequestBody) SetRunIssuetypes(params RunIssuetypes) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunMessageEmail sets the params for sending an email message
func (body *ConnectorsRunRequestBody) SetRunMessageEmail(params RunMessageEmail) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunMessageServerlog sets the params for writing to the Kibana server log
func (body *ConnectorsRunRequestBody) SetRunMessageServerlog(params RunMessageServerlog) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunMessageSlack sets the params for sending a message to Slack
func (body *ConnectorsRunRequestBody) SetRunMessageSlack(params RunMessageSlack) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunPostmessage sets the params for PostMessage subaction for Slack API
func (body *ConnectorsRunRequestBody) SetRunPostmessage(params RunPostmessage) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunPushtoservice sets the params for PushToService subaction for various connectors
func (body *ConnectorsRunRequestBody) SetRunPushtoservice(params RunPushtoservice) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunTriggerPagerduty sets the params for triggering a PagerDuty alert
func (body *ConnectorsRunRequestBody) SetRunTriggerPagerduty(params RunTriggerPagerduty) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// SetRunValidchannelid sets the params for ValidChannelId subaction for Slack API
func (body *ConnectorsRunRequestBody) SetRunValidchannelid(params RunValidchannelid) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	body.Params = data
	return nil
}

// newConnectorsRun returns a function that performs POST /api/actions/connector/{id}/_execute API requests
func (api *API) newConnectorsRun() func(context.Context, *ConnectorsRunRequest, ...RequestOption) (*ConnectorsRunResponse, error) {
	return func(ctx context.Context, req *ConnectorsRunRequest, opts ...RequestOption) (*ConnectorsRunResponse, error) {
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
			newCtx = instrument.Start(ctx, "connectors.run")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/actions/connector/%s/_execute", req.ID)

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
			instrument.BeforeRequest(httpReq, "connectors.run")
			if reader := instrument.RecordRequestBody(ctx, "connectors.run", httpReq.Body); reader != nil {
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
		resp := &ConnectorsRunResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result ConnectorsRunResponseBody

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
