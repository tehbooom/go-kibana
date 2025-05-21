package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// SecurityDetectionsPreviewAlertsResponse wraps the response from a PreviewAlerts call
type SecurityDetectionsPreviewAlertsResponse struct {
	StatusCode int
	Body       *SecurityDetectionsPreviewAlertsResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type SecurityDetectionsPreviewAlertsResponseBody struct {
	IsAborted bool                                          `json:"isAborted"`
	Logs      []SecurityDetectionsPreviewAlertsResponseLogs `json:"logs,omitempty"`
	PreviewID string                                        `json:"previewId"`
}

type SecurityDetectionsPreviewAlertsResponseLogs struct {
	Duration  int                                                   `json:"duration"`
	Errors    []string                                              `json:"errors"`
	Requests  []SecurityDetectionsPreviewAlertsResponseLogsRequests `json:"requests"`
	StartedAt string                                                `json:"startedAt"`
	Warning   []string                                              `json:"warnings"`
}

type SecurityDetectionsPreviewAlertsResponseLogsRequests struct {
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Request     string `json:"request"`
	RequestType string `json:"request_type"`
}

type SecurityDetectionsPreviewAlertsRequest struct {
	Params SecurityDetectionsPreviewAlertsRequestParams
	Body   SecurityDetectionsPreviewAlertsRequestBody
}

type SecurityDetectionsPreviewAlertsRequestParams struct {
	EnableLoggedRequests *bool
}

type SecurityDetectionsPreviewAlertsRequestBody struct {
	// Actions Array defining the automated actions (notifications) taken when alerts are generated.
	Actions []SecurityDetectionsRuleAction `json:"actions"`
	// AlertSuppression Defines alert suppression configuration.
	AlertSuppression *SecurityDetectionsAlertSuppression `json:"alert_suppression,omitempty"`
	// AliasPurpose Values are savedObjectConversion or savedObjectImport.
	AliasPurpose  *string `json:"alias_purpose,omitempty"`
	AliasTargetID *string `json:"alias_target_id,omitempty"`
	// Author The rule's author.
	Author []string `json:"author"`
	// BuildingBlockType Determines if the rule acts as a building block. If yes, the value must be `default`.
	// By default, building-block alerts are not displayed in the UI. These rules are used as a foundation for other rules that do generate alerts.
	// For more information, refer to [About building block rules](https://www.elastic.co/guide/en/security/current/building-block-rule.html).
	BuildingBlockType *string `json:"building_block_type,omitempty"`
	DataViewID        *string `json:"data_view_id,omitempty"`
	// Description The rule's description.
	Description string `json:"description"`
	// Enabled Determines whether the rule is enabled. Defaults to true.
	Enabled        bool                                  `json:"enabled"`
	ExceptionsList []SecurityDetectionsRuleExceptionList `json:"exceptions_list"`
	// FalsePositives String array used to describe common reasons why the rule may issue false-positive alerts. Defaults to an empty array.
	FalsePositives []string `json:"false_positives"`
	// From Time from which data is analyzed each time the rule runs, using a date math range. For example, now-4200s means the rule analyzes data from 70 minutes before its start time. Defaults to now-6m (analyzes data from 6 minutes before the start time).
	From string `json:"from"`
	// Index Indices on which the rule functions. Defaults to the Security Solution indices defined on the Kibana Advanced Settings page (Kibana → Stack Management → Advanced Settings → `securitySolution:defaultIndex`).
	Index *[]string `json:"index,omitempty"`
	// Interval Frequency of rule execution, using a date math range. For example, "1h" means the rule runs every hour. Defaults to 5m (5 minutes).
	Interval string `json:"interval"`
	// InvestigationFields Schema for fields relating to investigation fields. These are user defined fields we use to highlight
	// in various features in the UI such as alert details flyout and exceptions auto-population from alert.
	InvestigationFields *SecurityDetectionsInvestigationFields `json:"investigation_fields,omitempty"`
	InvocationCount     int                                    `json:"invocationCount"`
	// Language Values are kuery or lucene.
	Language string `json:"language"`
	// License The rule's license.
	License *string `json:"license,omitempty"`
	// MaxSignals Maximum number of alerts the rule can create during a single run (the rule's Max alerts per run [advanced setting](https://www.elastic.co/guide/en/security/current/rules-ui-create.html#rule-ui-advanced-params) value).
	// > info
	// > This setting can be superseded by the [Kibana configuration setting](https://www.elastic.co/guide/en/kibana/current/alert-action-settings-kb.html#alert-settings) `xpack.alerting.rules.run.alerts.max`, which determines the maximum alerts generated by any rule in the Kibana alerting framework. For example, if `xpack.alerting.rules.run.alerts.max` is set to 1000, the rule can generate no more than 1000 alerts even if `max_signals` is set higher.
	MaxSignals int `json:"max_signals"`
	// Meta Placeholder for metadata about the rule.
	// > info
	// > This field is overwritten when you save changes to the rule's settings.
	Meta *map[string]any `json:"meta,omitempty"`
	// Name A human-readable name for the rule.
	Name string `json:"name"`
	// Namespace Has no effect.
	Namespace *string `json:"namespace,omitempty"`
	// Note Notes to help investigate alerts produced by the rule.
	Note *string `json:"note,omitempty"`
	// OutputIndex (deprecated) Has no effect.
	// Deprecated:
	OutputIndex *string `json:"output_index,omitempty"`
	// Outcome Values are exactMatch, aliasMatch, or conflict.
	Outcome *string `json:"outcome,omitempty"`
	// Query [Query](https://www.elastic.co/guide/en/kibana/8.17/search.html) used by the rule to create alerts.
	// - For indicator match rules, only the query’s results are used to determine whether an alert is generated.
	Query string `json:"query"`
	// References Array containing notes about or references to relevant information about the rule. Defaults to an empty array.
	References          []string                               `json:"references"`
	RelatedIntegrations []SecurityDetectionsRelatedIntegration `json:"related_integrations"`
	RequiredFields      []SecurityDetectionsRequiredField      `json:"required_fields"`
	// ResponseActions defines the automated response actions to be taken when alerts are generated.
	// Use the following methods to work with response actions:
	// - AddOsqueryResponseAction(): Add an osquery response action
	// - AddEndpointResponseActionWithDefaultParams(): Add an endpoint action with default params
	// - AddEndpointResponseActionWithProcessesParams(): Add an endpoint action with processes params
	// - GetResponseActions(): Get all response actions
	// - GetOsqueryResponseActions(): Get all osquery actions
	// - GetEndpointResponseActions(): Get all endpoint actions
	// - ClearResponseActions(): Remove all response actions
	// - RemoveResponseActionByIndex(): Remove a specific response action
	ResponseActions *[]json.RawMessage `json:"response_actions,omitempty"`
	// RiskScore A numerical representation of the alert's severity from 0 to 100, where:
	// * `0` - `21` represents low severity
	// * `22` - `47` represents medium severity
	// * `48` - `73` represents high severity
	// * `74` - `100` represents critical severity
	RiskScore int `json:"risk_score"`
	// RiskScoreMapping Overrides generated alerts' risk_score with a value from the source event
	RiskScoreMapping []SecurityDetectionsRiskScoreMapping `json:"risk_score_mapping"`
	// RuleID A stable unique identifier for the rule object. It can be assigned during rule creation. It can be any string, but often is a UUID. It should be unique not only within a given Kibana space, but also across spaces and Elastic environments. The same prebuilt Elastic rule, when installed in two different Kibana spaces or two different Elastic environments, will have the same `rule_id`s.
	RuleID string `json:"rule_id"`
	// RuleNameOverride Sets which field in the source event is used to populate the alert's `signal.rule.name` value (in the UI, this value is displayed on the Rules page in the Rule column). When unspecified, the rule's `name` value is used. The source field must be a string data type.
	RuleNameOverride *string `json:"rule_name_override,omitempty"`
	// SavedID Kibana [saved search](https://www.elastic.co/guide/en/kibana/current/save-open-search.html) used by the rule to create alerts.
	SavedID *string `json:"saved_id,omitempty"`
	// Setup Populates the rule's setup guide with instructions on rule prerequisites such as required integrations, configuration steps, and anything else needed for the rule to work correctly.
	Setup string `json:"setup"`
	// Severity Severity level of alerts produced by the rule, which must be one of the following:
	// * `low`: Alerts that are of interest but generally not considered to be security incidents
	// * `medium`: Alerts that require investigation
	// * `high`: Alerts that require immediate investigation
	// * `critical`: Alerts that indicate it is highly likely a security incident has occurred
	Severity string `json:"severity"`
	// SeverityMapping Overrides generated alerts' severity with values from the source event
	SeverityMapping []SecurityDetectionsSeverityMapping `json:"severity_mapping"`
	// Tags String array containing words and phrases to help categorize, filter, and search rules. Defaults to an empty array.
	Tags   []string                   `json:"tags"`
	Threat []SecurityDetectionsThreat `json:"threat"`
	// Throttle Defines how often rule actions are taken.
	Throttle     *string `json:"throttle,omitempty"`
	TimeframeEnd string  `json:"timeframEnd"`
	// TimelineID Timeline template ID
	TimelineID *string `json:"timeline_id,omitempty"`
	// TimelineTitle Timeline template title
	TimelineTitle *string `json:"timeline_title,omitempty"`
	// TimestampOverride Sets the time field used to query indices. When unspecified, rules query the `@timestamp` field. The source field must be an Elasticsearch date data type.
	TimestampOverride *string `json:"timestamp_override,omitempty"`
	// TimestampOverrideFallbackDisabled Disables the fallback to the event's @timestamp field
	TimestampOverrideFallbackDisabled *bool  `json:"timestamp_override_fallback_disabled,omitempty"`
	To                                string `json:"to"`
	// Type Rule type
	Type string `json:"type"`
	// Version The rule's version number.
	// - For prebuilt rules it represents the version of the rule's content in the source [detection-rules](https://github.com/elastic/detection-rules) repository (and the corresponding `security_detection_engine` Fleet package that is used for distributing prebuilt rules).
	// - For custom rules it is set to `1` when the rule is created.
	// - It is not incremented on each update. Compare this to the `revision` field.
	Version int `json:"version"`
}

// newSecurityDetectionsPreviewAlerts returns a function that performs POST /api/detection_engine/rules/preview API requests
func (api *API) newSecurityDetectionsPreviewAlerts() func(context.Context, *SecurityDetectionsPreviewAlertsRequest, ...RequestOption) (*SecurityDetectionsPreviewAlertsResponse, error) {
	return func(ctx context.Context, req *SecurityDetectionsPreviewAlertsRequest, opts ...RequestOption) (*SecurityDetectionsPreviewAlertsResponse, error) {
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
			newCtx = instrument.Start(ctx, "security_detections.preview_alerts")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := "/api/detection_engine/rules/preview"

		// Build query parameters
		params := make(map[string]string)

		if req.Params.EnableLoggedRequests != nil {
			params["enable_logged_requests"] = strconv.FormatBool(*req.Params.EnableLoggedRequests)
		}

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
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
			instrument.BeforeRequest(httpReq, "security_detections.preview_alerts")
			if reader := instrument.RecordRequestBody(ctx, "security_detections.preview_alerts", httpReq.Body); reader != nil {
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
		resp := &SecurityDetectionsPreviewAlertsResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result SecurityDetectionsPreviewAlertsResponseBody

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
