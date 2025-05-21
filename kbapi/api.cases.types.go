package kbapi

import (
	"encoding/json"
)

type CasesObjectRequest struct {
	// Assignees An array containing users that are assigned to the case.
	// Not more than 10 elements.
	Assignees *[]CasesAssignee `json:"assignees"`

	// Category A word or phrase that categorizes the case.
	// Maximum length is 50.
	Category *string `json:"category,omitempty"`

	Connector json.RawMessage `json:"connector"`

	// CustomFields Custom field values for a case. Any optional custom fields that are not specified in the request are set to null.
	CustomFields *CasesCustomField `json:"customFields,omitempty"`

	// Description The description for the case.
	Description string `json:"description"`

	// ID The identifier for the case.
	// Only used to update a case not during creation
	ID string `json:"id,omitempty"`

	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner string `json:"owner"`

	// Settings An object that contains the case settings.
	Settings CasesSettings `json:"settings"`

	// Severity The severity of the case.
	// Values are critical, high, low, or medium. Default value is low.
	Severity *string `json:"severity,omitempty"`

	// Status The status of the case.
	// Values are closed, in-progress, or open.
	// Only used to update a case not during creation.
	Status string `json:"status,omitempty"`

	// Tags The words and phrases that help categorize cases. It can be an empty array.
	// Not more than 200 elements. Maximum length of each is 256.
	Tags []string `json:"tags"`

	// Title A title for the case.
	// Maximum length is 160.
	Title string `json:"title"`
}

// CasesObjectResponse represents a case object.
// The Connector field can have different types based on the connector type.
// Use GetConnectorType() to determine the type, and then use Get<TYPE>Connector() to get the appropriately typed connector.
// Alternatively, use GetConnector() to automatically get the properly typed connector.
//
// The Comments field contains an array of comments with different types.
// Use GetCommentTypes() to get the types of all comments, or use GetAllTypedComments()
// to get all comments with their proper types.
type CasesObjectResponse struct {
	// Assignees An array containing users that are assigned to the case.
	// Not more than 10 elements.
	Assignees *[]CasesAssignee `json:"assignees"`

	// Category A word or phrase that categorizes the case.
	// Maximum length is 50.
	Category *string `json:"category,omitempty"`

	Connector json.RawMessage `json:"connector"`

	ClosedAt  *string           `json:"closed_at"`
	ClosedBy  *UserObject       `json:"closed_by"`
	Comments  []json.RawMessage `json:"comments"`
	CreatedAt string            `json:"created_at"`
	CreatedBy UserObject        `json:"created_by"`
	// CustomFields Custom field values for a case. Any optional custom fields that are not specified in the request are set to null.
	CustomFields []CasesCustomField `json:"customFields,omitempty"`

	// Description The description for the case.
	Description string `json:"description"`
	// Duration The elapsed time from the creation of the case to its closure (in seconds).
	// If the case has not been closed, the duration is set to null. If the case was closed after less than half a second, the duration is rounded down to zero.
	// If the case was closed after less than half a second, the duration is rounded down to zero.
	Duration        *int                  `json:"duration"`
	ExternalService *CasesExternalService `json:"external_service"`
	ID              string                `json:"id"`

	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner string `json:"owner"`

	// Settings An object that contains the case settings.
	Settings CasesSettings `json:"settings"`

	// Severity The severity of the case.
	// Values are critical, high, low, or medium. Default value is low.
	Severity *string `json:"severity,omitempty"`

	// Status The status of the case.
	Status string `json:"status"`

	// Tags The words and phrases that help categorize cases. It can be an empty array.
	// Not more than 200 elements. Maximum length of each is 256.
	Tags []string `json:"tags"`

	// Title A title for the case.
	// Maximum length is 160.
	Title        string      `json:"title"`
	TotalAlerts  int         `json:"totalAlerts"`
	TotalComment int         `json:"totalComment"`
	UpdatedAt    *string     `json:"updated_at,omitempty"`
	UpdatedBy    *UserObject `json:"updated_by,omitempty"`
	Version      string      `json:"version"`
}

// CasesAssignee An array containing users that are assigned to the case.
type CasesAssignee struct {
	// UID A unique identifier for the user profile. These identifiers can be found by using the suggest user profile API.
	UID string `json:"uid"`
}

type CasesCustomField struct {
	// Key A unique key for the custom field. Must be lower case and composed only of a-z, 0-9, '_', and '-' characters.
	// It is used in API calls to refer to a specific custom field.
	// Minimum length is 1, maximum length is 36.
	Key string `json:"key"`
	// Type the type of the custom field.
	// Values are text or toggle.	Type string `json:"type"`
	Type string `json:"type"`
	// DefaultValue a default value for the custom field. If the type is text, the default value must be a string.
	// If the type is toggle, the default value must be boolean.
	// TODO: String or boolean here
	DefaultValue json.RawMessage `json:"defaultValue"`
	// Label the custom field label that is displayed in the case.
	// Minimum length is 1, maximum length is 50.
	Label string `json:"label"`
	// Required indicates whether the field is required.
	// If false, the custom field can be set to null or omitted when a case is created or updated.
	Required bool `json:"required"`
}

type CasesExternalService struct {
	ConnectorID   *string     `json:"connector_id,omitempty"`
	ConnectorName *string     `json:"connector_name,omitempty"`
	ExternalID    *string     `json:"external_id,omitempty"`
	ExternalTitle *string     `json:"external_title,omitempty"`
	ExternalURL   *string     `json:"external_url,omitempty"`
	PushedAt      *string     `json:"pushed_at,omitempty"`
	PushedBy      *UserObject `json:"pushed_by"`
}

// CasesSettings An object that contains the case settings.
type CasesSettings struct {
	// SyncAlerts Turns alert syncing on or off.
	SyncAlerts bool `json:"syncAlerts"`
}

type UserObject struct {
	Email      *string `json:"email"`
	FullName   *string `json:"full_name"`
	ProfileUID *string `json:"profile_uid,omitempty"`
	Username   *string `json:"username"`
}

// Connector types
type JiraConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.jira`
	Type   string     `json:"type"`
	Fields JiraFields `json:"fields"`
}

type JiraFields struct {
	IssueType string `json:"issueType"`
	Parent    string `json:"parent"`
	Priority  string `json:"priority"`
}

type ResilientConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.resilient`
	Type   string          `json:"type"`
	Fields ResilientFields `json:"fields"`
}

type ResilientFields struct {
	IssueTypes   string `json:"issueTypes"`
	SeverityCode string `json:"severityCode"`
}

type ServiceNowConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.servicenow`
	Type   string           `json:"type"`
	Fields ServiceNowFields `json:"fields"`
}

type ServiceNowFields struct {
	Category    string `json:"category"`
	Impact      string `json:"impact"`
	Severity    string `json:"severity"`
	Subcategory string `json:"subcategory"`
	Urgency     string `json:"urgency"`
}

type ServiceNowSIRConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.servicenow-sir`
	Type   string              `json:"type"`
	Fields ServiceNowSIRFields `json:"fields"`
}

type ServiceNowSIRFields struct {
	Category    string `json:"category"`
	DestIP      bool   `json:"destIp"`
	MalwareHash bool   `json:"malwareHash"`
	MalwareURL  bool   `json:"malwareUrl"`
	Priority    string `json:"priority"`
	SourceIP    bool   `json:"sourceIp"`
	Subcategory string `json:"subcategory"`
}

type SwimlaneConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.swimlane`
	Type   string         `json:"type"`
	Fields SwimlaneFields `json:"fields"`
}

type SwimlaneFields struct {
	CaseID string `json:"caseId"`
}

type WebhookConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.cases-webhook`
	Type   string `json:"type"`
	Fields string `json:"fields"`
}

type NoneConnector struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Type value is `.none`
	Type   string  `json:"type"`
	Fields *string `json:"fields"`
}

func (req *CasesObjectRequest) SetJiraConnector(connector JiraConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetServiceNowConnector(connector ServiceNowConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetResilientConnector(connector ResilientConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetServiceNowSIRConnector(connector ServiceNowSIRConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetSwimlaneConnector(connector SwimlaneConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetWebhookConnector(connector WebhookConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

func (req *CasesObjectRequest) SetNoneConnector(connector NoneConnector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}
	req.Connector = data
	return nil
}

type CasesAlertCommentResult struct {
	BaseComment          BaseComment
	UserCommentResponse  *UserCommentResponse
	AlertCommentResponse *AlertCommentResponse
}

// Comment types
type BaseComment struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	CreatedAt string     `json:"created_at"`
	CreatedBy UserObject `json:"created_by"`
	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner     string     `json:"owner"`
	PushedAt  string     `json:"pushed_at"`
	PushedBy  UserObject `json:"pushed_by"`
	UpdatedAt string     `json:"updated_at"`
	UpdateBy  UserObject `json:"updated_by"`
	Version   string     `json:"version"`
}

type UserCommentResponse struct {
	BaseComment
	Comment string `json:"comment"`
}

type AlertCommentResponse struct {
	BaseComment
	AlertID []string               `json:"alert_id"`
	Index   []string               `json:"index"`
	Rule    AlertCommentRuleObject `json:"rule"`
}

type AlertCommentRuleObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ConnectorType string

const (
	ConnectorTypeJira          ConnectorType = ".jira"
	ConnectorTypeNone          ConnectorType = ".none"
	ConnectorTypeResilient     ConnectorType = ".resilient"
	ConnectorTypeServiceNow    ConnectorType = ".servicenow"
	ConnectorTypeServiceNowSIR ConnectorType = ".servicenow-sir"
	ConnectorTypeSwimlane      ConnectorType = ".swimlane"
	ConnectorTypeWebhook       ConnectorType = ".cases-webhook"
	ConnectorTypeUnknown       ConnectorType = "unknown"
)

type CommentType string

const (
	CommentTypeUser    CommentType = "user"
	CommentTypeAlert   CommentType = "alert"
	CommentTypeUnknown CommentType = "unknown"
)

type AlertCommentRequest struct {
	AlertID []string `json:"alert_id"`
	Index   []string `json:"index"`
	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner string                 `json:"owner"`
	Rule  AlertCommentRuleObject `json:"rule"`
	Type  string                 `json:"type"`
}

type UserCommentRequest struct {
	Comment string `json:"comment"`
	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner string `json:"owner"`
	Type  string `json:"type"`
}

type CasesUserActionObject struct {
	// Values are add, create, delete, push_to_service, or update.
	Action    string     `json:"action"`
	CommentID *string    `json:"comment_id"`
	CreatedAt string     `json:"created_at"`
	CreatedBy UserObject `json:"created_by"`
	ID        string     `json:"id"`

	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	Owner   string          `json:"owner"`
	Payload json.RawMessage `json:"payload"`

	// Type The type of action.
	// Values are assignees, create_case, comment, connector, description, pushed,
	// tags, title, status, settings, or severity.
	Type    string `json:"type"`
	Version string `json:"version"`
}

type CasesPayloadAlert struct {
	Comment struct {
		AlertID []string `json:"alert_id"`
		Index   []string `json:"index"`
		// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
		// Values are cases, observability, or securitySolution.
		Owner string                 `json:"owner"`
		Rule  AlertCommentRuleObject `json:"rule"`
		Type  string                 `json:"type"`
	} `json:"comment"`
}

type CasesPayloadComment struct {
	Comment struct {
		Comment string `json:"comment"`
		// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
		// Values are cases, observability, or securitySolution.
		Owner string `json:"owner"`
		Type  string `json:"type"`
	} `json:"comment"`
}

type CasesPayloadAssigness struct {
	// Assigness an array containing users that are assigned to the case.
	// Not more than 10 elements.
	Assignees *[]CasesAssignee `json:"assigness"`
}

// CasesPayloadConnector defines model for Cases_payload_connector.
type CasesPayloadConnector struct {
	Connector *struct {
		// Fields An object containing the connector fields. To create a case without a connector, specify null. If you want to omit any individual field, specify null as its value.
		Fields *struct {
			// CaseID The case identifier for Swimlane connectors.
			CaseID *string `json:"caseId,omitempty"`

			// Category The category of the incident for ServiceNow ITSM and ServiceNow SecOps connectors.
			Category *string `json:"category,omitempty"`

			// DestIP Indicates whether cases will send a comma-separated list of destination IPs for ServiceNow SecOps connectors.
			DestIP *bool `json:"destIp"`

			// Impact The effect an incident had on business for ServiceNow ITSM connectors.
			Impact *string `json:"impact,omitempty"`

			// IssueType The type of issue for Jira connectors.
			IssueType *string `json:"issueType,omitempty"`

			// IssueTypes The type of incident for IBM Resilient connectors.
			IssueTypes *[]string `json:"issueTypes,omitempty"`

			// MalwareHash Indicates whether cases will send a comma-separated list of malware hashes for ServiceNow SecOps connectors.
			MalwareHash *bool `json:"malwareHash"`

			// MalwareURL Indicates whether cases will send a comma-separated list of malware URLs for ServiceNow SecOps connectors.
			MalwareURL *bool `json:"malwareUrl"`

			// Parent The key of the parent issue, when the issue type is sub-task for Jira connectors.
			Parent *string `json:"parent,omitempty"`

			// Priority The priority of the issue for Jira and ServiceNow SecOps connectors.
			Priority *string `json:"priority,omitempty"`

			// Severity The severity of the incident for ServiceNow ITSM connectors.
			Severity *string `json:"severity,omitempty"`

			// SeverityCode The severity code of the incident for IBM Resilient connectors.
			SeverityCode *string `json:"severityCode,omitempty"`

			// SourceIP Indicates whether cases will send a comma-separated list of source IPs for ServiceNow SecOps connectors.
			SourceIP *bool `json:"sourceIp"`

			// Subcategory The subcategory of the incident for ServiceNow ITSM connectors.
			Subcategory *string `json:"subcategory,omitempty"`

			// Urgency The extent to which the incident resolution can be delayed for ServiceNow ITSM connectors.
			Urgency *string `json:"urgency,omitempty"`
		} `json:"fields"`

		// ID The identifier for the connector. To create a case without a connector, use `none`.
		ID *string `json:"id,omitempty"`

		// Name The name of the connector. To create a case without a connector, use `none`.
		Name *string `json:"name,omitempty"`

		// Type The type of connector.
		Type *string `json:"type,omitempty"`
	} `json:"connector,omitempty"`
}

// CasesPayloadCreateCase defines model for Cases_payload_create_case.
type CasesPayloadCreateCase struct {
	// Assignees An array containing users that are assigned to the case.
	// Not more than 10 elements.
	Assignees *[]CasesAssignee `json:"assignees"`
	Connector *struct {
		// Fields An object containing the connector fields. To create a case without a connector, specify null. If you want to omit any individual field, specify null as its value.
		Fields *struct {
			// CaseID The case identifier for Swimlane connectors.
			CaseID *string `json:"caseId,omitempty"`

			// Category The category of the incident for ServiceNow ITSM and ServiceNow SecOps connectors.
			Category *string `json:"category,omitempty"`

			// DestIP Indicates whether cases will send a comma-separated list of destination IPs for ServiceNow SecOps connectors.
			DestIP *bool `json:"destIp"`

			// Impact The effect an incident had on business for ServiceNow ITSM connectors.
			Impact *string `json:"impact,omitempty"`

			// IssueType The type of issue for Jira connectors.
			IssueType *string `json:"issueType,omitempty"`

			// IssueTypes The type of incident for IBM Resilient connectors.
			IssueTypes *[]string `json:"issueTypes,omitempty"`

			// MalwareHash Indicates whether cases will send a comma-separated list of malware hashes for ServiceNow SecOps connectors.
			MalwareHash *bool `json:"malwareHash"`

			// MalwareURL Indicates whether cases will send a comma-separated list of malware URLs for ServiceNow SecOps connectors.
			MalwareURL *bool `json:"malwareUrl"`

			// Parent The key of the parent issue, when the issue type is sub-task for Jira connectors.
			Parent *string `json:"parent,omitempty"`

			// Priority The priority of the issue for Jira and ServiceNow SecOps connectors.
			Priority *string `json:"priority,omitempty"`

			// Severity The severity of the incident for ServiceNow ITSM connectors.
			Severity *string `json:"severity,omitempty"`

			// SeverityCode The severity code of the incident for IBM Resilient connectors.
			SeverityCode *string `json:"severityCode,omitempty"`

			// SourceIP Indicates whether cases will send a comma-separated list of source IPs for ServiceNow SecOps connectors.
			SourceIP *bool `json:"sourceIp"`

			// Subcategory The subcategory of the incident for ServiceNow ITSM connectors.
			Subcategory *string `json:"subcategory,omitempty"`

			// Urgency The extent to which the incident resolution can be delayed for ServiceNow ITSM connectors.
			Urgency *string `json:"urgency,omitempty"`
		} `json:"fields"`

		// ID The identifier for the connector. To create a case without a connector, use `none`.
		ID *string `json:"id,omitempty"`

		// Name The name of the connector. To create a case without a connector, use `none`.
		Name *string `json:"name,omitempty"`

		// Type The type of connector.
		Type *string `json:"type,omitempty"`
	} `json:"connector,omitempty"`
	Description *string `json:"description,omitempty"`

	// Owner The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner *string `json:"owner,omitempty"`

	// Settings An object that contains the case settings.
	Settings *CasesSettings `json:"settings,omitempty"`

	// Severity The severity of the case.
	// Values are critical, high, low, or medium. Default value is low.
	Severity *string `json:"severity,omitempty"`

	// Status The status of the case.
	// Values are closed, in-progress, or open.
	Status *string   `json:"status,omitempty"`
	Tags   *[]string `json:"tags,omitempty"`
	Title  *string   `json:"title,omitempty"`
}

// If the action is delete and the type is delete_case, the payload is nullable.
type CasesPayloadDelete struct{}

type CasesPayloadDescription struct {
	Description *string `json:"description,omitempty"`
}

type CasesPayloadPushed struct {
	ExternalService *CasesExternalService `json:"externalService"`
}

type CasesPayloadSettings struct {
	// Settings an object that contains the case settings.
	Settings *CasesSettings `json:"settings,omitempty"`
}

type CasesPayloadSeverity struct {
	// Severity the severity of the case.
	// Values are critical, high, low, or medium. Default value is low.
	Severity *string `json:"severity,omitempty"`
}

type CasesPayloadStatus struct {
	// Status the status of the case.
	// Values are closed, in-progress, or open.
	Status *string `json:"status,omitempty"`
}

type CasesPayloadTags struct {
	Tags *[]string `json:"tags,omitempty"`
}

type CasesPayloadTitle struct {
	Title *string `json:"title,omitempty"`
}

type UserActionType string

const (
	UserActionTypeAssignees   UserActionType = "assignees"
	UserActionTypeCreateCase  UserActionType = "create_case"
	UserActionTypeComment     UserActionType = "comment"
	UserActionTypeConnector   UserActionType = "connector"
	UserActionTypeDescription UserActionType = "description"
	UserActionTypePushed      UserActionType = "pushed"
	UserActionTypeTags        UserActionType = "tags"
	UserActionTypeTitle       UserActionType = "title"
	UserActionTypeStatus      UserActionType = "status"
	UserActionTypeSettings    UserActionType = "settings"
	UserActionTypeSeverity    UserActionType = "severity"
)

type CasesSettingsRequest struct {
	// ClosureType indicates whether a case is automatically closed when it is pushed
	// to external systems (close-by-pushing) or not automatically closed (close-by-user).
	// Values are close-by-pushing or close-by-user.
	ClosureType string `json:"closure_type"`
	// Connector an object that contains the connector configuration.
	Connector CaseSettingsConnector `json:"connector"`
	// CustomFields custom field values for a case.
	// Any optional custom fields that are not specified in the request are set to null.
	// At least 0 but not more than 10 elements.
	CustomFields *[]CasesCustomField `json:"customFields,omitempty"`
	// The application that owns the cases: Stack Management, Observability, or Elastic Security.
	// Values are cases, observability, or securitySolution.
	Owner     string          `json:"owner"`
	Templates []CasesTemplate `json:"templates"`
}

type CasesSettingsResponse struct {
	// ClosureType indicates whether a case is automatically closed when it is pushed
	// to external systems (close-by-pushing) or not automatically closed (close-by-user).
	// Values are close-by-pushing or close-by-user.
	ClosureType string                `json:"closure_type"`
	Connector   CaseSettingsConnector `json:"connector"`
	CreatedAt   string                `json:"created_at"`
	CreatedBy   UserObject            `json:"created_by"`
	// CustomFields Custom field values for a case. Any optional custom fields that are not specified in the request are set to null.
	CustomFields *[]CasesCustomField `json:"customFields,omitempty"`
	Error        *string             `json:"error"`
	ID           string              `json:"id,omitempty"`
	Mappings     *[]CaseMapping      `json:"mappings,omitempty"`
	Owner        string              `json:"owner"`
	Templates    []CasesTemplate     `json:"templates"`
	UpdatedAt    *string             `json:"updated_at,omitempty"`
	UpdatedBy    *UserObject         `json:"updated_by,omitempty"`
	Version      string              `json:"version"`
}

type CaseMapping struct {
	ActionType *string `json:"action_type,omitempty"`
	Source     *string `json:"source,omitempty"`
	Target     *string `json:"target,omitempty"`
}

type CaseSettingsConnector struct {
	// Fields the fields specified in the case configuration are not used and are not propagated to individual cases,
	// therefore it is recommended to set it to `null`.
	Fields *map[string]interface{} `json:"fields"`
	// ID The identifier for the connector. If you do not want a default connector, use `none`.
	// To retrieve connector IDs, use the find connectors API.
	ID *string `json:"id,omitempty"`
	// Name the name of the connector. If you do not want a default connector, use `none`.
	// To retrieve connector names, use the find connectors API.
	Name *string `json:"name,omitempty"`
	// Type the type of connector.
	Type *string `json:"type,omitempty"`
}

type CasesTemplate struct {
	CaseFields *struct {
		// Assignees An array containing users that are assigned to the case.
		Assignees *[]CasesAssignee `json:"assignees"`

		// Category A word or phrase that categorizes the case.
		Category  *string                `json:"category,omitempty"`
		Connector *CaseSettingsConnector `json:"connector,omitempty"`

		// CustomFields Custom field values in the template.
		CustomFields *[]CasesCustomField `json:"customFields,omitempty"`

		// Description The description for the case.
		Description *string `json:"description,omitempty"`

		// Settings An object that contains the case settings.
		Settings *CasesSettings `json:"settings,omitempty"`

		// Severity The severity of the case.
		Severity *string `json:"severity,omitempty"`

		// Tags The words and phrases that help categorize cases. It can be an empty array.
		Tags *[]string `json:"tags,omitempty"`

		// Title A title for the case.
		Title *string `json:"title,omitempty"`
	} `json:"caseFields,omitempty"`

	// Description A description for the template.
	Description *string `json:"description,omitempty"`

	// Key A unique key for the template. Must be lower case and composed only of a-z, 0-9, '_', and '-' characters.
	// It is used in API calls to refer to a specific template.
	Key *string `json:"key,omitempty"`

	// Name The name of the template.
	Name *string `json:"name,omitempty"`

	// Tags The words and phrases that help categorize templates. It can be an empty array.
	Tags *[]string `json:"tags,omitempty"`
}
