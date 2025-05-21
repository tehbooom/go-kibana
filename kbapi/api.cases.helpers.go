package kbapi

import (
	"encoding/json"
	"fmt"
)

func (resp *CasesObjectResponse) GetConnectorType() (ConnectorType, error) {
	if len(resp.Connector) == 0 {
		return ConnectorTypeUnknown, nil
	}

	var typeContainer struct {
		Type string `json:"type"`
	}

	err := json.Unmarshal(resp.Connector, &typeContainer)
	if err != nil {
		return ConnectorTypeUnknown, err
	}

	switch typeContainer.Type {
	case ".jira":
		return ConnectorTypeJira, nil
	case ".servicenow":
		return ConnectorTypeServiceNow, nil
	case ".servicenow-sir":
		return ConnectorTypeServiceNowSIR, nil
	case ".none":
		return ConnectorTypeNone, nil
	case ".resilient":
		return ConnectorTypeResilient, nil
	case ".swimlane":
		return ConnectorTypeSwimlane, nil
	case ".cases-webhook":
		return ConnectorTypeWebhook, nil
	default:
		return ConnectorTypeUnknown, nil
	}
}

// Generic method to get the connector of the appropriate type
func (resp *CasesObjectResponse) GetConnector() (interface{}, error) {
	connectorType, err := resp.GetConnectorType()
	if err != nil {
		return nil, err
	}

	switch connectorType {
	case ConnectorTypeJira:
		return resp.GetJiraConnector()
	case ConnectorTypeServiceNow:
		return resp.GetServiceNowConnector()
	case ConnectorTypeServiceNowSIR:
		return resp.GetServiceNowSIRConnector()
	case ConnectorTypeNone:
		return resp.GetNoneConnector()
	case ConnectorTypeSwimlane:
		return resp.GetSwimlaneConnector()
	case ConnectorTypeResilient:
		return resp.GetResilientConnector()
	case ConnectorTypeWebhook:
		return resp.GetWebhookConnector()
	default:
		return nil, fmt.Errorf("Unknown connector")
	}
}

func (resp *CasesObjectResponse) GetJiraConnector() (*JiraConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector JiraConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".jira" {
		return nil, fmt.Errorf("connector is not a .jira connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetServiceNowConnector() (*ServiceNowConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector ServiceNowConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".servicenow" {
		return nil, fmt.Errorf("connector is not a .servicenow connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetServiceNowSIRConnector() (*ServiceNowSIRConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector ServiceNowSIRConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".servicenow-sir" {
		return nil, fmt.Errorf("connector is not a .servicenow-sir connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetSwimlaneConnector() (*SwimlaneConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector SwimlaneConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".swimlane" {
		return nil, fmt.Errorf("connector is not a .swimlane connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetResilientConnector() (*ResilientConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector ResilientConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != "jira" {
		return nil, fmt.Errorf("connector is not a Jira connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetWebhookConnector() (*WebhookConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector WebhookConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".cases-webhook" {
		return nil, fmt.Errorf("connector is not a .cases-webhook connector")
	}

	return &connector, nil
}

func (resp *CasesObjectResponse) GetNoneConnector() (*NoneConnector, error) {
	if len(resp.Connector) == 0 {
		return nil, nil
	}

	var connector NoneConnector
	err := json.Unmarshal(resp.Connector, &connector)
	if err != nil {
		return nil, err
	}

	if connector.Type != ".none" {
		return nil, fmt.Errorf("connector is not a .none connector")
	}

	return &connector, nil
}

// Method to get the types of all comments
func (resp *CasesObjectResponse) GetCommentTypes() ([]CommentType, error) {
	if len(resp.Comments) == 0 {
		return []CommentType{}, nil
	}

	types := make([]CommentType, 0, len(resp.Comments))
	for _, rawComment := range resp.Comments {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(rawComment, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "user":
			types = append(types, CommentTypeUser)
		case "alert":
			types = append(types, CommentTypeAlert)
		default:
			types = append(types, CommentTypeUnknown)
		}
	}

	return types, nil
}

// Method to get all comments with their appropriate types
func (resp *CasesObjectResponse) GetAllTypedComments() ([]interface{}, error) {
	if len(resp.Comments) == 0 {
		return []interface{}{}, nil
	}

	comments := make([]interface{}, 0, len(resp.Comments))
	for _, rawComment := range resp.Comments {
		var typeContainer struct {
			Type string `json:"type"`
		}

		err := json.Unmarshal(rawComment, &typeContainer)
		if err != nil {
			return nil, err
		}

		switch typeContainer.Type {
		case "user":
			var comment UserCommentResponse
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		case "alert":
			var comment AlertCommentResponse
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		default:
			var comment BaseComment
			err := json.Unmarshal(rawComment, &comment)
			if err != nil {
				return nil, err
			}
			comments = append(comments, comment)
		}
	}

	return comments, nil
}
