package kbapi

import (
	"encoding/json"
	"fmt"
)

func UnmarshalRule(data []byte) (SecurityDetectionsRule, error) {
	var typeContainer struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeContainer); err != nil {
		return nil, fmt.Errorf("error determining rule type: %w", err)
	}

	var rule SecurityDetectionsRule
	switch typeContainer.Type {
	case "esql":
		var r SecurityDetectionsESQLRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling ESQL rule: %w", err)
		}
		rule = &r
	case "new_terms":
		var r SecurityDetectionsNewTermsRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling New Terms rule: %w", err)
		}
		rule = &r
	case "machine_learning":
		var r SecurityDetectionsMachineLearningRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling Machine Learning rule: %w", err)
		}
		rule = &r
	case "threat_match":
		var r SecurityDetectionsThreatMatchRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling Threat Match rule: %w", err)
		}
		rule = &r
	case "threshold":
		var r SecurityDetectionsThresholdRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling Threshold rule: %w", err)
		}
		rule = &r
	case "saved_query":
		var r SecurityDetectionsSavedQueryRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling Saved Query rule: %w", err)
		}
		rule = &r
	case "query":
		var r SecurityDetectionsQueryRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling Query rule: %w", err)
		}
		rule = &r
	case "eql":
		var r SecurityDetectionsEQLRule
		if err := json.Unmarshal(data, &r); err != nil {
			return nil, fmt.Errorf("error unmarshalling EQL rule: %w", err)
		}
		rule = &r
	default:
		return nil, fmt.Errorf("unknown rule type: %s", typeContainer.Type)
	}

	return rule, nil
}

// GetStringValue extracts a string value from the ECS mapping field.
func (m SecurityDetectionsECSMapping) GetStringValue(ecsField string) (string, bool) {
	if mapping, ok := m[ecsField]; ok && mapping.Value != nil {
		var str string
		err := json.Unmarshal(*mapping.Value, &str)
		if err == nil {
			return str, true
		}
	}
	return "", false
}

// GetStringSliceValue extracts a string slice value from the ECS mapping field.
func (m SecurityDetectionsECSMapping) GetStringSliceValue(ecsField string) ([]string, bool) {
	if mapping, ok := m[ecsField]; ok && mapping.Value != nil {
		var strSlice []string
		err := json.Unmarshal(*mapping.Value, &strSlice)
		if err == nil {
			return strSlice, true
		}
	}
	return nil, false
}

// SetStringValue sets a string value for the specified ECS field.
// Creates the mapping entry if it doesn't exist.
func (m SecurityDetectionsECSMapping) SetStringValue(ecsField, value string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal string value: %w", err)
	}

	raw := json.RawMessage(data)

	// Initialize the map entry if it doesn't exist
	if _, ok := m[ecsField]; !ok {
		m[ecsField] = struct {
			Field *string          `json:"field,omitempty"`
			Value *json.RawMessage `json:"value,omitempty"`
		}{
			Value: &raw,
		}
	} else {
		// Update existing entry
		entry := m[ecsField]
		entry.Value = &raw
		m[ecsField] = entry
	}

	return nil
}

// SetStringSliceValue sets a string slice value for the specified ECS field.
// Creates the mapping entry if it doesn't exist.
func (m SecurityDetectionsECSMapping) SetStringSliceValue(ecsField string, values []string) error {
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal string slice value: %w", err)
	}

	raw := json.RawMessage(data)

	// Initialize the map entry if it doesn't exist
	if _, ok := m[ecsField]; !ok {
		m[ecsField] = struct {
			Field *string          `json:"field,omitempty"`
			Value *json.RawMessage `json:"value,omitempty"`
		}{
			Value: &raw,
		}
	} else {
		// Update existing entry
		entry := m[ecsField]
		entry.Value = &raw
		m[ecsField] = entry
	}

	return nil
}

// NewECSMapping creates a new empty SecurityDetectionsECSMapping.
func NewECSMapping() SecurityDetectionsECSMapping {
	return make(SecurityDetectionsECSMapping)
}

// GetStringField extracts a string value from the Field field.
func (t *SecurityDetectionsThreshold) GetStringField() (string, bool) {
	if t.Field == nil {
		return "", false
	}

	var str string
	err := json.Unmarshal(t.Field, &str)
	if err == nil {
		return str, true
	}
	return "", false
}

// GetStringSliceField extracts a string slice value from the Field field.
func (t *SecurityDetectionsThreshold) GetStringSliceField() ([]string, bool) {
	if t.Field == nil {
		return nil, false
	}

	var strSlice []string
	err := json.Unmarshal(t.Field, &strSlice)
	if err == nil {
		return strSlice, true
	}
	return nil, false
}

// SetStringField sets a string value for the Field field.
func (t *SecurityDetectionsThreshold) SetStringField(value string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal string field: %w", err)
	}

	t.Field = json.RawMessage(data)
	return nil
}

// SetStringSliceField sets a string slice value for the Field field.
// Use an empty slice ([]string{}) to specify an empty array.
func (t *SecurityDetectionsThreshold) SetStringSliceField(values []string) error {
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal string slice field: %w", err)
	}

	t.Field = json.RawMessage(data)
	return nil
}

// IsFieldEmpty checks if the Field is an empty array ([]).
// Returns true if it's an empty array, false otherwise.
func (t *SecurityDetectionsThreshold) IsFieldEmpty() bool {
	if t.Field == nil {
		return true
	}

	// Check for empty array "[]"
	if string(t.Field) == "[]" {
		return true
	}

	// Also check for empty slice
	fields, ok := t.GetStringSliceField()
	return ok && len(fields) == 0
}

// NewThreshold creates a new SecurityDetectionsThreshold with initialized fields.
func NewThreshold(value int) *SecurityDetectionsThreshold {
	return &SecurityDetectionsThreshold{
		Value: value,
	}
}

// GetMachineLearningJobIDString extracts a string value from the MachineLearningJobID field.
// Returns the string value and true if the field is a string.
// Returns empty string and false if the field is not a string or is nil.
func (r *SecurityDetectionsMachineLearningRule) GetMachineLearningJobIDString() (string, bool) {
	if r.MachineLearningJobID == nil {
		return "", false
	}

	var str string
	err := json.Unmarshal(r.MachineLearningJobID, &str)
	if err == nil {
		return str, true
	}
	return "", false
}

// GetMachineLearningJobIDSlice extracts a string slice value from the MachineLearningJobID field.
func (r *SecurityDetectionsMachineLearningRule) GetMachineLearningJobIDSlice() ([]string, bool) {
	if r.MachineLearningJobID == nil {
		return nil, false
	}

	var strSlice []string
	err := json.Unmarshal(r.MachineLearningJobID, &strSlice)
	if err == nil {
		return strSlice, true
	}
	return nil, false
}

// SetMachineLearningJobIDString sets a string value for the MachineLearningJobID field.
func (r *SecurityDetectionsMachineLearningRule) SetMachineLearningJobIDString(value string) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal string MachineLearningJobID: %w", err)
	}

	r.MachineLearningJobID = json.RawMessage(data)
	return nil
}

// SetMachineLearningJobIDSlice sets a string slice value for the MachineLearningJobID field.
func (r *SecurityDetectionsMachineLearningRule) SetMachineLearningJobIDSlice(values []string) error {
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal string slice MachineLearningJobID: %w", err)
	}

	r.MachineLearningJobID = json.RawMessage(data)
	return nil
}

// AddOsqueryResponseAction adds an osquery response action to the rule.
func (r *SecurityDetectionsCommonRule) AddOsqueryResponseAction(params SecurityDetectionsOSqueryParams) error {
	action := SecurityDetectionsOSqueryResponseAction{
		ActionTypeID: ActionTypeOsquery,
		Params:       params,
	}

	data, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal osquery response action: %w", err)
	}

	if r.ResponseActions == nil {
		r.ResponseActions = &[]json.RawMessage{}
	}

	*r.ResponseActions = append(*r.ResponseActions, json.RawMessage(data))
	return nil
}

// AddEndpointResponseActionWithDefaultParams adds an endpoint response action with default parameters.
func (r *SecurityDetectionsCommonRule) AddEndpointResponseActionWithDefaultParams(command string, comment *string) error {
	params := SecurityDetectionsAPIEndpointResponseActionDefaultParams{
		Command: command,
		Comment: comment,
	}

	paramsData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoint default params: %w", err)
	}

	action := SecurityDetectionsEndpointResponseAction{
		ActionTypeId: ActionTypeEndpoint,
		Params:       json.RawMessage(paramsData),
	}

	actionData, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoint response action: %w", err)
	}

	if r.ResponseActions == nil {
		r.ResponseActions = &[]json.RawMessage{}
	}

	*r.ResponseActions = append(*r.ResponseActions, json.RawMessage(actionData))
	return nil
}

// AddEndpointResponseActionWithProcessesParams adds an endpoint response action with processes parameters.
func (r *SecurityDetectionsCommonRule) AddEndpointResponseActionWithProcessesParams(
	command string,
	comment *string,
	field string,
	overwrite *bool,
) error {
	params := SecurityDetectionsAPIEndpointResponseActionProcessesParams{
		Command: command,
		Comment: comment,
		Config: struct {
			Field     string `json:"field"`
			Overwrite *bool  `json:"overwrite,omitempty"`
		}{
			Field:     field,
			Overwrite: overwrite,
		},
	}

	paramsData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoint processes params: %w", err)
	}

	action := SecurityDetectionsEndpointResponseAction{
		ActionTypeId: ActionTypeEndpoint,
		Params:       json.RawMessage(paramsData),
	}

	actionData, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("failed to marshal endpoint response action: %w", err)
	}

	if r.ResponseActions == nil {
		r.ResponseActions = &[]json.RawMessage{}
	}

	*r.ResponseActions = append(*r.ResponseActions, json.RawMessage(actionData))
	return nil
}

// GetResponseActions returns all response actions.
func (r *SecurityDetectionsCommonRule) GetResponseActions() []json.RawMessage {
	if r.ResponseActions == nil {
		return []json.RawMessage{}
	}
	return *r.ResponseActions
}

// GetOsqueryResponseActions returns all osquery response actions.
func (r *SecurityDetectionsCommonRule) GetOsqueryResponseActions() ([]SecurityDetectionsOSqueryResponseAction, error) {
	if r.ResponseActions == nil {
		return []SecurityDetectionsOSqueryResponseAction{}, nil
	}

	var result []SecurityDetectionsOSqueryResponseAction

	for _, actionData := range *r.ResponseActions {
		var typeContainer struct {
			ActionTypeID string `json:"action_type_id"`
		}

		if err := json.Unmarshal(actionData, &typeContainer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response action type: %w", err)
		}

		if typeContainer.ActionTypeID == ActionTypeOsquery {
			var action SecurityDetectionsOSqueryResponseAction
			if err := json.Unmarshal(actionData, &action); err != nil {
				return nil, fmt.Errorf("failed to unmarshal osquery response action: %w", err)
			}
			result = append(result, action)
		}
	}

	return result, nil
}

// GetEndpointResponseActions returns all endpoint response actions.
func (r *SecurityDetectionsCommonRule) GetEndpointResponseActions() ([]SecurityDetectionsEndpointResponseAction, error) {
	if r.ResponseActions == nil {
		return []SecurityDetectionsEndpointResponseAction{}, nil
	}

	var result []SecurityDetectionsEndpointResponseAction

	for _, actionData := range *r.ResponseActions {
		var typeContainer struct {
			ActionTypeID string `json:"action_type_id"`
		}

		if err := json.Unmarshal(actionData, &typeContainer); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response action type: %w", err)
		}

		if typeContainer.ActionTypeID == ActionTypeEndpoint {
			var action SecurityDetectionsEndpointResponseAction
			if err := json.Unmarshal(actionData, &action); err != nil {
				return nil, fmt.Errorf("failed to unmarshal endpoint response action: %w", err)
			}
			result = append(result, action)
		}
	}

	return result, nil
}

// GetDefaultParams returns the default parameters if the endpoint action has default params.
func (a *SecurityDetectionsEndpointResponseAction) GetDefaultParams() (*SecurityDetectionsAPIEndpointResponseActionDefaultParams, bool) {
	if a.Params == nil {
		return nil, false
	}

	var params SecurityDetectionsAPIEndpointResponseActionDefaultParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	var testProcesses struct {
		Config *struct{} `json:"config"`
	}
	err = json.Unmarshal(a.Params, &testProcesses)
	if err == nil && testProcesses.Config != nil {
		return nil, false
	}

	return &params, true
}

// GetProcessesParams returns the processes parameters if the endpoint action has process params.
func (a *SecurityDetectionsEndpointResponseAction) GetProcessesParams() (*SecurityDetectionsAPIEndpointResponseActionProcessesParams, bool) {
	if a.Params == nil {
		return nil, false
	}

	var params SecurityDetectionsAPIEndpointResponseActionProcessesParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	if params.Config.Field == "" {
		return nil, false
	}

	return &params, true
}

// GetParams returns the params as an interface{} which will be either DefaultParams or ProcessesParams.
func (a *SecurityDetectionsEndpointResponseAction) GetParams() (interface{}, error) {
	if a.Params == nil {
		return nil, fmt.Errorf("params is nil")
	}

	if params, ok := a.GetDefaultParams(); ok {
		return params, nil
	}

	if params, ok := a.GetProcessesParams(); ok {
		return params, nil
	}

	return nil, fmt.Errorf("unknown params format")
}

// SetDefaultParams sets default parameters for the endpoint action.
func (a *SecurityDetectionsEndpointResponseAction) SetDefaultParams(params SecurityDetectionsAPIEndpointResponseActionDefaultParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal default params: %w", err)
	}

	a.Params = json.RawMessage(data)
	return nil
}

// SetProcessesParams sets processes parameters for the endpoint action.
func (a *SecurityDetectionsEndpointResponseAction) SetProcessesParams(params SecurityDetectionsAPIEndpointResponseActionProcessesParams) error {
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal processes params: %w", err)
	}

	a.Params = json.RawMessage(data)
	return nil
}

// ClearResponseActions removes all response actions from the rule.
func (r *SecurityDetectionsCommonRule) ClearResponseActions() {
	r.ResponseActions = nil
}

// RemoveResponseActionByIndex removes a response action by its index.
func (r *SecurityDetectionsCommonRule) RemoveResponseActionByIndex(index int) error {
	if r.ResponseActions == nil || index < 0 || index >= len(*r.ResponseActions) {
		return fmt.Errorf("invalid index: %d", index)
	}

	actions := *r.ResponseActions
	*r.ResponseActions = append(actions[:index], actions[index+1:]...)

	if len(*r.ResponseActions) == 0 {
		r.ResponseActions = nil
	}

	return nil
}

// GetRuleSourceType determines the type of rule source without fully unmarshaling the data.
func (r *SecurityDetectionsCommonRule) GetRuleSourceType() (string, error) {
	if r.RuleSource == nil {
		return "", nil
	}

	var typeContainer struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(r.RuleSource, &typeContainer); err != nil {
		return "", fmt.Errorf("failed to unmarshal rule source type: %w", err)
	}

	return typeContainer.Type, nil
}

// GetInternalRuleSource extracts the internal rule source if this is an internal rule.
func (r *SecurityDetectionsCommonRule) GetInternalRuleSource() (*SecurityDetectionsInternalRuleSource, bool) {
	if r.RuleSource == nil {
		return nil, false
	}

	var source SecurityDetectionsInternalRuleSource
	if err := json.Unmarshal(r.RuleSource, &source); err != nil {
		return nil, false
	}

	if source.Type != RuleSourceTypeInternal {
		return nil, false
	}

	return &source, true
}

// GetExternalRuleSource extracts the external rule source if this is an external rule.
func (r *SecurityDetectionsCommonRule) GetExternalRuleSource() (*SecurityDetectionsExternalRuleSource, bool) {
	if r.RuleSource == nil {
		return nil, false
	}

	var source SecurityDetectionsExternalRuleSource
	if err := json.Unmarshal(r.RuleSource, &source); err != nil {
		return nil, false
	}

	if source.Type != RuleSourceTypeExternal {
		return nil, false
	}

	return &source, true
}

// GetRuleSource returns the rule source as an interface{} which will be either
func (r *SecurityDetectionsCommonRule) GetRuleSource() (interface{}, error) {
	if r.RuleSource == nil {
		return nil, nil
	}

	if source, ok := r.GetInternalRuleSource(); ok {
		return source, nil
	}

	if source, ok := r.GetExternalRuleSource(); ok {
		return source, nil
	}

	return nil, fmt.Errorf("unknown rule source format")
}

// IsInternalRule checks if the rule is internal.
func (r *SecurityDetectionsCommonRule) IsInternalRule() bool {
	_, ok := r.GetInternalRuleSource()
	return ok
}

// IsExternalRule checks if the rule is external.
func (r *SecurityDetectionsCommonRule) IsExternalRule() bool {
	_, ok := r.GetExternalRuleSource()
	return ok
}

// IsCustomizedExternalRule checks if the rule is an external rule that has been customized.
func (r *SecurityDetectionsCommonRule) IsCustomizedExternalRule() bool {
	source, ok := r.GetExternalRuleSource()
	return ok && source.IsCustomized
}

// GetSlackParams returns the Slack connector parameters if this is a Slack action.
func (a *SecurityDetectionsRuleAction) GetSlackParams() (*SecurityDetectionsSlackConnectorParams, bool) {
	if a.Params == nil || (a.ActionTypeID != ActionTypeSlack && a.ActionTypeID != ActionTypeSlackAPI) {
		return nil, false
	}

	var params SecurityDetectionsSlackConnectorParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	return &params, true
}

// SetSlackParams sets the Slack connector parameters for this action.
func (a *SecurityDetectionsRuleAction) SetSlackParams(message string) error {
	params := SecurityDetectionsSlackConnectorParams{
		Message: message,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack params: %w", err)
	}

	a.Params = json.RawMessage(data)

	// Set action type if not already set
	if a.ActionTypeID == "" {
		a.ActionTypeID = ActionTypeSlack
	}

	return nil
}

// GetEmailParams returns the Email connector parameters if this is an Email action.
func (a *SecurityDetectionsRuleAction) GetEmailParams() (*SecurityDetectionsEmailConnectorParams, bool) {
	if a.Params == nil || a.ActionTypeID != ActionTypeEmail {
		return nil, false
	}

	var params SecurityDetectionsEmailConnectorParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	return &params, true
}

// SetEmailParams sets the Email connector parameters for this action.
func (a *SecurityDetectionsRuleAction) SetEmailParams(to, cc, bcc, subject, message string) error {
	if to == "" && cc == "" && bcc == "" {
		return fmt.Errorf("at least one of to, cc, or bcc must be specified for Email actions")
	}

	params := SecurityDetectionsEmailConnectorParams{
		To:      to,
		Cc:      cc,
		Bcc:     bcc,
		Subject: subject,
		Message: message,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Email params: %w", err)
	}

	a.Params = json.RawMessage(data)

	// Set action type if not already set
	if a.ActionTypeID == "" {
		a.ActionTypeID = ActionTypeEmail
	}

	return nil
}

// GetWebhookParams returns the Webhook connector parameters if this is a Webhook action.
func (a *SecurityDetectionsRuleAction) GetWebhookParams() (*SecurityDetectionsWebhookConnectorParams, bool) {
	if a.Params == nil || a.ActionTypeID != ActionTypeWebhook {
		return nil, false
	}

	var params SecurityDetectionsWebhookConnectorParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	return &params, true
}

// SetWebhookParams sets the Webhook connector parameters for this action.
func (a *SecurityDetectionsRuleAction) SetWebhookParams(body string) error {
	params := SecurityDetectionsWebhookConnectorParams{
		Body: body,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Webhook params: %w", err)
	}

	a.Params = json.RawMessage(data)

	if a.ActionTypeID == "" {
		a.ActionTypeID = ActionTypeWebhook
	}

	return nil
}

// GetPagerDutyParams returns the PagerDuty connector parameters if this is a PagerDuty action.
func (a *SecurityDetectionsRuleAction) GetPagerDutyParams() (*SecurityDetectionsPagerDutyConnectorParams, bool) {
	if a.Params == nil || a.ActionTypeID != ActionTypePagerDuty {
		return nil, false
	}

	var params SecurityDetectionsPagerDutyConnectorParams
	err := json.Unmarshal(a.Params, &params)
	if err != nil {
		return nil, false
	}

	return &params, true
}

// SetPagerDutyParams sets the PagerDuty connector parameters for this action.
func (a *SecurityDetectionsRuleAction) SetPagerDutyParams(
	severity string,
	eventAction string,
	dedupKey string,
	timestamp string,
	component string,
	group string,
	source string,
	summary string,
	class string,
) error {
	// Validate severity
	validSeverities := map[string]bool{
		PagerDutySeverityCritical: true,
		PagerDutySeverityError:    true,
		PagerDutySeverityWarning:  true,
		PagerDutySeverityInfo:     true,
	}
	if !validSeverities[severity] {
		return fmt.Errorf("invalid PagerDuty severity: %s", severity)
	}

	// Validate event action
	validEventActions := map[string]bool{
		PagerDutyEventActionTrigger:     true,
		PagerDutyEventActionResolve:     true,
		PagerDutyEventActionAcknowledge: true,
	}
	if !validEventActions[eventAction] {
		return fmt.Errorf("invalid PagerDuty event action: %s", eventAction)
	}

	params := SecurityDetectionsPagerDutyConnectorParams{
		Severity:    severity,
		EventAction: eventAction,
		DedupKey:    dedupKey,
		Timestamp:   timestamp,
		Component:   component,
		Group:       group,
		Source:      source,
		Summary:     summary,
		Class:       class,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal PagerDuty params: %w", err)
	}

	a.Params = json.RawMessage(data)

	if a.ActionTypeID == "" {
		a.ActionTypeID = ActionTypePagerDuty
	}

	return nil
}

// GetParams returns the action parameters as an interface{} which will be one of the connector param types.
func (a *SecurityDetectionsRuleAction) GetParams() (interface{}, error) {
	if a.Params == nil {
		return nil, nil
	}

	switch a.ActionTypeID {
	case ActionTypeSlack, ActionTypeSlackAPI:
		if params, ok := a.GetSlackParams(); ok {
			return params, nil
		}
	case ActionTypeEmail:
		if params, ok := a.GetEmailParams(); ok {
			return params, nil
		}
	case ActionTypeWebhook:
		if params, ok := a.GetWebhookParams(); ok {
			return params, nil
		}
	case ActionTypePagerDuty:
		if params, ok := a.GetPagerDutyParams(); ok {
			return params, nil
		}
	}

	return nil, fmt.Errorf("could not determine params type for action type %s", a.ActionTypeID)
}

// NewSlackAction creates a new Slack notification action.
func NewSlackAction(connectorID, message string) *SecurityDetectionsRuleAction {
	action := &SecurityDetectionsRuleAction{
		ActionTypeID: ActionTypeSlack,
		ID:           connectorID,
	}
	action.SetSlackParams(message)
	return action
}

// NewEmailAction creates a new Email notification action.
func NewEmailAction(connectorID, to, cc, bcc, subject, message string) (*SecurityDetectionsRuleAction, error) {
	action := &SecurityDetectionsRuleAction{
		ActionTypeID: ActionTypeEmail,
		ID:           connectorID,
	}
	if err := action.SetEmailParams(to, cc, bcc, subject, message); err != nil {
		return nil, err
	}
	return action, nil
}

// NewWebhookAction creates a new Webhook notification action.
func NewWebhookAction(connectorID, body string) *SecurityDetectionsRuleAction {
	action := &SecurityDetectionsRuleAction{
		ActionTypeID: ActionTypeWebhook,
		ID:           connectorID,
	}
	action.SetWebhookParams(body)
	return action
}

// NewPagerDutyAction creates a new PagerDuty notification action.
func NewPagerDutyAction(
	connectorID,
	severity,
	eventAction,
	dedupKey,
	timestamp,
	component,
	group,
	source,
	summary,
	class string,

) (*SecurityDetectionsRuleAction, error) {
	action := &SecurityDetectionsRuleAction{
		ActionTypeID: ActionTypePagerDuty,
		ID:           connectorID,
	}
	if err := action.SetPagerDutyParams(
		severity,
		eventAction,
		dedupKey,
		timestamp,
		component,
		group,
		source,
		summary,
		class,
	); err != nil {
		return nil, err
	}
	return action, nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddTags(operation SecurityDetectionsBulkActionRulesEditTags) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddIndexPatterns(operation SecurityDetectionsBulkActionRulesEditIndexPatterns) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddInvestigationFields(operation SecurityDetectionsBulkActionRulesEditInvestigationFields) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddTimeline(operation SecurityDetectionsBulkActionRulesEditTimeline) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddSchedule(operation SecurityDetectionsBulkActionRulesEditSchedule) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

func (e *SecurityDetectionsBulkActionRulesEdit) AddActions(operation SecurityDetectionsBulkActionRulesEditActions) error {
	data, err := json.Marshal(operation)
	if err != nil {
		return fmt.Errorf("failed to marshal operation: %w", err)
	}

	e.Edit = append(e.Edit, data)

	return nil
}

// SetSlackParams sets the Slack connector parameters for this action.
func (a *SecurityDetectionsBulkActionRulesEditActionsItem) SetSlackParams(message string) error {
	params := SecurityDetectionsSlackConnectorParams{
		Message: message,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack params: %w", err)
	}

	a.Params = json.RawMessage(data)

	return nil
}

// SetEmailParams sets the Email connector parameters for this action.
func (a *SecurityDetectionsBulkActionRulesEditActionsItem) SetEmailParams(to, cc, bcc, subject, message string) error {
	if to == "" && cc == "" && bcc == "" {
		return fmt.Errorf("at least one of to, cc, or bcc must be specified for Email actions")
	}

	params := SecurityDetectionsEmailConnectorParams{
		To:      to,
		Cc:      cc,
		Bcc:     bcc,
		Subject: subject,
		Message: message,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Email params: %w", err)
	}

	a.Params = json.RawMessage(data)

	return nil
}

// SetWebhookParams sets the Webhook connector parameters for this action.
func (a *SecurityDetectionsBulkActionRulesEditActionsItem) SetWebhookParams(body string) error {
	params := SecurityDetectionsWebhookConnectorParams{
		Body: body,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal Webhook params: %w", err)
	}

	a.Params = json.RawMessage(data)

	return nil
}

// SetPagerDutyParams sets the PagerDuty connector parameters for this action.
func (a *SecurityDetectionsBulkActionRulesEditActionsItem) SetPagerDutyParams(
	severity string,
	eventAction string,
	dedupKey string,
	timestamp string,
	component string,
	group string,
	source string,
	summary string,
	class string,
) error {
	// Validate severity
	validSeverities := map[string]bool{
		PagerDutySeverityCritical: true,
		PagerDutySeverityError:    true,
		PagerDutySeverityWarning:  true,
		PagerDutySeverityInfo:     true,
	}
	if !validSeverities[severity] {
		return fmt.Errorf("invalid PagerDuty severity: %s", severity)
	}

	// Validate event action
	validEventActions := map[string]bool{
		PagerDutyEventActionTrigger:     true,
		PagerDutyEventActionResolve:     true,
		PagerDutyEventActionAcknowledge: true,
	}
	if !validEventActions[eventAction] {
		return fmt.Errorf("invalid PagerDuty event action: %s", eventAction)
	}

	params := SecurityDetectionsPagerDutyConnectorParams{
		Severity:    severity,
		EventAction: eventAction,
		DedupKey:    dedupKey,
		Timestamp:   timestamp,
		Component:   component,
		Group:       group,
		Source:      source,
		Summary:     summary,
		Class:       class,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal PagerDuty params: %w", err)
	}

	a.Params = json.RawMessage(data)

	return nil
}

func (r *SecurityDetectionsBulkActionRulesResponse) UnmarshalBulkAction() (SecurityDetectionsBulkActionEditResponse, error) {
	var resp SecurityDetectionsBulkActionEditResponse
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		return resp, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return resp, nil
}
