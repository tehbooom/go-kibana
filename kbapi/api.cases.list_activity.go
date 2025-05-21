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

// TODO: Update the call
// CasesListActivityResponse wraps the response from a <todo> call
type CasesListActivityResponse struct {
	StatusCode int
	Body       *CasesListActivityResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type CasesListActivityResponseBody struct {
	Page        int                     `json:"page"`
	PerPage     int                     `json:"per_page"`
	Total       int                     `json:"total"`
	UserActions []CasesUserActionObject `json:"userActions"`
}

func (resp *CasesListActivityResponseBody) GetUserActionTypes() ([]UserActionType, error) {
	if len(resp.UserActions) == 0 {
		return []UserActionType{}, nil
	}

	types := make([]UserActionType, 0, len(resp.UserActions))
	for _, userAction := range resp.UserActions {
		switch userAction.Type {
		case "assignees":
			types = append(types, UserActionTypeAssignees)
		case "create_case":
			types = append(types, UserActionTypeCreateCase)
		case "comment":
			types = append(types, UserActionTypeComment)
		case "connector":
			types = append(types, UserActionTypeConnector)
		case "description":
			types = append(types, UserActionTypeDescription)
		case "pushed":
			types = append(types, UserActionTypePushed)
		case "tags":
			types = append(types, UserActionTypeTags)
		case "title":
			types = append(types, UserActionTypeTitle)
		case "status":
			types = append(types, UserActionTypeStatus)
		case "settings":
			types = append(types, UserActionTypeSettings)
		case "severity":
			types = append(types, UserActionTypeSeverity)
		}
	}

	return types, nil
}

func (resp *CasesListActivityResponseBody) GetAllTypedUserActions() ([]interface{}, error) {
	if len(resp.UserActions) == 0 {
		return []interface{}{}, nil
	}

	userActions := make([]interface{}, 0, len(resp.UserActions))
	for _, userAction := range resp.UserActions {
		switch userAction.Type {
		case "assignees":
			var action CasesPayloadAssigness
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "create_case":
			var action CasesPayloadCreateCase
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "comment":
			// create a temporary struct to get the type of comment (alert or comment)
			var typeContainer struct {
				Type string `json:"type"`
			}
			err := json.Unmarshal(userAction.Payload, &typeContainer)
			if err != nil {
				return nil, err
			}

			if typeContainer.Type == "alert" {
				var action CasesPayloadAlert
				err = json.Unmarshal(userAction.Payload, &action)
				if err != nil {
					return nil, err
				}
				userActions = append(userActions, action)
			} else {
				var action CasesPayloadComment
				err = json.Unmarshal(userAction.Payload, &action)
				if err != nil {
					return nil, err
				}
				userActions = append(userActions, action)
			}
		case "connector":
			var action CasesPayloadConnector
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "description":
			var action CasesPayloadDescription
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "pushed":
			var action CasesPayloadPushed
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "tags":
			var action CasesPayloadTags
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "title":
			var action CasesPayloadTitle
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "status":
			var action CasesPayloadStatus
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "settings":
			var action CasesPayloadSettings
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		case "severity":
			var action CasesPayloadSeverity
			err := json.Unmarshal(userAction.Payload, &action)
			if err != nil {
				return nil, err
			}
			userActions = append(userActions, action)
		}
	}

	return userActions, nil
}

type CasesListActivityRequest struct {
	ID     string
	Params CasesListActivityRequestParams
}

type CasesListActivityRequestParams struct {
	// PerPage The number of rules to return per page.
	// Maximum value is 100. Default value is 20.
	PerPage *int `form:"per_page,omitempty" json:"per_page,omitempty"`
	// Page The page number to return.
	// Default value is 1
	Page *int `form:"page,omitempty" json:"page,omitempty"`
	// SortOrder Determines the sort order.
	// Values are asc or desc. Default value is desc.
	SortOrder *string `form:"sort_order,omitempty" json:"sort_order,omitempty"`
	// Types determines the types of user actions to return.
	// Values are action, alert, assignees, attachment, comment, connector, create_case,
	// description, pushed, settings, severity, status, tags, title, or user.
	Types *[]string `form:"types,omitempty" json:"types,omitempty"`
}

// newCasesListActivity returns a function that performs GET /api/cases/{caseId}/user_actions/_find API requests
func (api *API) newCasesListActivity() func(context.Context, *CasesListActivityRequest, ...RequestOption) (*CasesListActivityResponse, error) {
	return func(ctx context.Context, req *CasesListActivityRequest, opts ...RequestOption) (*CasesListActivityResponse, error) {
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
			newCtx = instrument.Start(ctx, "cases.list_activity")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/cases/%s/user_actions/_find", req.ID)

		// Build query parameters
		params := make(map[string]string)

		if req.Params.Page != nil {
			params["page"] = strconv.Itoa(*req.Params.Page)
		}
		if req.Params.PerPage != nil {
			params["per_page"] = strconv.Itoa(*req.Params.PerPage)
		}
		if req.Params.SortOrder != nil {
			params["sort_order"] = *req.Params.SortOrder
		}
		if req.Params.Types != nil {
			params["types"] = strings.Join(*req.Params.Types, ",")
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
			instrument.BeforeRequest(httpReq, "cases.list_activity")
			if reader := instrument.RecordRequestBody(ctx, "cases.list_activity", httpReq.Body); reader != nil {
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
		resp := &CasesListActivityResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result CasesListActivityResponseBody

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
