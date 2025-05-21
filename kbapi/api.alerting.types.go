package kbapi

type AlertingGetTypesResponseItem struct {
	ID                     string              `json:"id"`
	Name                   string              `json:"name"`
	Alerts                 *Alerts             `json:"alerts,omitempty"`
	Category               string              `json:"category"`
	Producer               string              `json:"producer"`
	ActionGroups           []ActionGroup       `json:"action_groups"`
	IsExportable           bool                `json:"is_exportable"`
	ActionVariables        ActionVariables     `json:"action_variables"`
	RuleTaskTimeout        string              `json:"rule_task_timeout"`
	EnabledInLicense       bool                `json:"enabled_in_license"`
	HasAlertsMappings      bool                `json:"has_alerts_mappings"`
	AuthorizedConsumers    AuthorizedConsumers `json:"authorized_consumers"`
	HasFieldsForAAD        bool                `json:"has_fields_for_a_a_d"`
	RecoveryActionGroup    ActionGroup         `json:"recovery_action_group"`
	DefaultActionGroupID   string              `json:"default_action_group_id"`
	MinimumLicenseRequired string              `json:"minimum_license_required"`
	DoesSetRecoveryContext bool                `json:"does_set_recovery_context"`
}

// Alerts contains alert configuration
type Alerts struct {
	Context     string   `json:"context"`
	Mappings    Mappings `json:"mappings"`
	ShouldWrite bool     `json:"shouldWrite"`
}

// Mappings defines field mappings for alerts
type Mappings struct {
	FieldMap map[string]FieldDefinition `json:"fieldMap"`
}

// FieldDefinition defines properties of a field
type FieldDefinition struct {
	Type       string              `json:"type"`
	Array      bool                `json:"array"`
	Required   bool                `json:"required"`
	Dynamic    *bool               `json:"dynamic,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
}

// Property defines a property of a complex field
type Property struct {
	Type string `json:"type"`
}

// ActionGroup defines an action group
type ActionGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ActionVariables defines variables available for actions
type ActionVariables struct {
	State   []interface{}     `json:"state"`
	Params  []interface{}     `json:"params"`
	Context []ContextVariable `json:"context"`
}

// ContextVariable defines a context variable
type ContextVariable struct {
	Name                           string `json:"name"`
	Description                    string `json:"description"`
	UseWithTripleBracesInTemplates bool   `json:"useWithTripleBracesInTemplates,omitempty"`
}

// AuthorizedConsumers defines which consumers can use this rule type
type AuthorizedConsumers struct {
	ML             ConsumerAccess `json:"ml"`
	APM            ConsumerAccess `json:"apm"`
	SLO            ConsumerAccess `json:"slo"`
	Logs           ConsumerAccess `json:"logs"`
	SIEM           ConsumerAccess `json:"siem"`
	Alerts         ConsumerAccess `json:"alerts"`
	Uptime         ConsumerAccess `json:"uptime"`
	Discover       ConsumerAccess `json:"discover"`
	Monitoring     ConsumerAccess `json:"monitoring"`
	StackAlerts    ConsumerAccess `json:"stackAlerts"`
	Infrastructure ConsumerAccess `json:"infrastructure"`
}

// ConsumerAccess defines access levels for a consumer
type ConsumerAccess struct {
	All  bool `json:"all"`
	Read bool `json:"read"`
}

// TopRecord represents the structure for top_records array items
type TopRecord struct {
	Actual              float64 `json:"actual"`
	JobID               string  `json:"job_id"`
	Typical             float64 `json:"typical"`
	Function            string  `json:"function"`
	Timestamp           string  `json:"timestamp"`
	FieldName           string  `json:"field_name"`
	IsInterim           bool    `json:"is_interim"`
	RecordScore         float64 `json:"record_score"`
	ByFieldName         string  `json:"by_field_name"`
	ByFieldValue        string  `json:"by_field_value"`
	DetectorIndex       int     `json:"detector_index"`
	OverFieldName       string  `json:"over_field_name"`
	OverFieldValue      string  `json:"over_field_value"`
	InitialRecordScore  float64 `json:"initial_record_score"`
	PartitionFieldName  string  `json:"partition_field_name"`
	PartitionFieldValue string  `json:"partition_field_value"`
}

// TopInfluencer represents the structure for top_influencers array items
type TopInfluencer struct {
	JobID                  string  `json:"job_id"`
	Timestamp              string  `json:"timestamp"`
	IsInterim              bool    `json:"is_interim"`
	InfluencerScore        float64 `json:"influencer_score"`
	InfluencerFieldName    string  `json:"influencer_field_name"`
	InfluencerFieldValue   string  `json:"influencer_field_value"`
	InitialInfluencerScore float64 `json:"initial_influencer_score"`
}

// ActionResponse extends Action with response-specific fields
type ActionResponse struct {
	AlertsFilter            *AlertsFilter  `json:"alerts_filter,omitempty"`
	ConnectorTypeID         string         `json:"connector_type_id"`
	Frequency               *Frequency     `json:"frequency,omitempty"`
	Group                   *string        `json:"group,omitempty"`
	ID                      string         `json:"id"`
	Params                  map[string]any `json:"params"`
	UseAlertDataForTemplate *bool          `json:"use_alert_data_for_template,omitempty"`
	UUID                    *string        `json:"uuid,omitempty"`
}

// ExecutionStatus represents the current execution state of a rule
type ExecutionStatus struct {
	Error             *Error   `json:"error,omitempty"`
	LastDuration      *float32 `json:"last_duration,omitempty"`
	LastExecutionDate string   `json:"last_execution_date"`
	Status            string   `json:"status"`
	Warning           *Warning `json:"warning,omitempty"`
}

// Error provides details about an error
type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// Warning provides details about a warning
type Warning struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// LastRun contains information about the last execution of the rule
type LastRun struct {
	AlertsCount  AlertsCount `json:"alerts_count"`
	Outcome      string      `json:"outcome"`
	OutcomeMsg   []string    `json:"outcome_msg,omitempty"`
	OutcomeOrder *float32    `json:"outcome_order,omitempty"`
	Warning      *string     `json:"warning,omitempty"`
}

// AlertsCount provides counts of different alert states
type AlertsCount struct {
	Active    *float32 `json:"active"`
	Ignored   *float32 `json:"ignored"`
	New       *float32 `json:"new"`
	Recovered *float32 `json:"recovered"`
}

// Monitoring contains performance and execution metrics for a rule
type Monitoring struct {
	Run MonitoringRun `json:"run"`
}

// MonitoringRun contains details about rule execution
type MonitoringRun struct {
	CalculatedMetrics CalculatedMetrics `json:"calculated_metrics"`
	History           []HistoryItem     `json:"history"`
	LastRun           MonitoringLastRun `json:"last_run"`
}

// CalculatedMetrics contains performance metrics
type CalculatedMetrics struct {
	P50          *float32 `json:"p50,omitempty"`
	P95          *float32 `json:"p95,omitempty"`
	P99          *float32 `json:"p99,omitempty"`
	SuccessRatio float32  `json:"success_ratio"`
}

// HistoryItem represents a single historical execution
type HistoryItem struct {
	Duration  *float32 `json:"duration,omitempty"`
	Outcome   *string  `json:"outcome,omitempty"`
	Success   bool     `json:"success"`
	Timestamp float32  `json:"timestamp"`
}

// MonitoringLastRun contains details about the most recent execution
type MonitoringLastRun struct {
	Metrics   RunMetrics `json:"metrics"`
	Timestamp string     `json:"timestamp"`
}

// RunMetrics contains performance measurements
type RunMetrics struct {
	Duration            *float32  `json:"duration,omitempty"`
	GapDurationS        *float32  `json:"gap_duration_s"`
	GapRange            *GapRange `json:"gap_range"`
	TotalAlertsCreated  *float32  `json:"total_alerts_created"`
	TotalAlertsDetected *float32  `json:"total_alerts_detected"`
	TotalIndexingTimeMS *float32  `json:"total_indexing_duration_ms"`
	TotalSearchTimeMS   *float32  `json:"total_search_duration_ms"`
}

// GapRange represents a time range
type GapRange struct {
	Gte string `json:"gte"`
	Lte string `json:"lte"`
}

// SnoozeSchedule defines when alerts are snoozed
type SnoozeSchedule struct {
	Duration        float32  `json:"duration"`
	ID              *string  `json:"id,omitempty"`
	RRule           RRule    `json:"rRule"`
	SkipRecurrences []string `json:"skipRecurrences,omitempty"`
}

// RRule defines recurrence rules for snooze schedules
type RRule struct {
	ByHour     []float32 `json:"byhour,omitempty"`
	ByMinute   []float32 `json:"byminute,omitempty"`
	ByMonth    []float32 `json:"bymonth,omitempty"`
	ByMonthDay []float32 `json:"bymonthday,omitempty"`
	BySecond   []float32 `json:"bysecond,omitempty"`
	BySetPos   []float32 `json:"bysetpos,omitempty"`
	ByWeekday  []string  `json:"byweekday,omitempty"`
	ByWeekNo   []float32 `json:"byweekno,omitempty"`
	ByYearDay  []float32 `json:"byyearday,omitempty"`
	Count      *float32  `json:"count,omitempty"`
	DTStart    string    `json:"dtstart"`
	Freq       *int      `json:"freq,omitempty"`
	Interval   *float32  `json:"interval,omitempty"`
	TzID       string    `json:"tzid"`
	Until      *string   `json:"until,omitempty"`
	WkSt       *string   `json:"wkst,omitempty"`
}

// Action defines the actions to be taken when alert conditions are met
type Action struct {
	AlertsFilter            *AlertsFilter   `json:"alerts_filter,omitempty"`
	Frequency               *Frequency      `json:"frequency,omitempty"`
	Group                   *string         `json:"group,omitempty"`
	ID                      string          `json:"id"`
	Params                  *map[string]any `json:"params,omitempty"`
	UseAlertDataForTemplate *bool           `json:"use_alert_data_for_template,omitempty"`
	UUID                    *string         `json:"uuid,omitempty"`
}

// AlertsFilter defines filters for when alerts should trigger actions
type AlertsFilter struct {
	Query     *Query     `json:"query,omitempty"`
	Timeframe *Timeframe `json:"timeframe,omitempty"`
}

// Query defines the search criteria for alerts
type Query struct {
	// DSL TODO: This cannot be used in EQL query
	DSL *string `json:"dsl,omitempty"`
	// Filters is an array of filter objects as defined in kbn-es-query package
	Filters []Filter `json:"filters"`
	// KQL is a Kibana Query Language string
	KQL string `json:"kql"`
}

// Filter A filter written in Elasticsearch Query Domain Specific Language (DSL) as defined in the `kbn-es-query` package.
type Filter struct {
	State *map[string]interface{} `json:"$state,omitempty"`
	Meta  *struct {
		Alias        *string                 `json:"alias"`
		ControlledBy *string                 `json:"controlledBy,omitempty"`
		Disabled     *bool                   `json:"disabled,omitempty"`
		Field        *string                 `json:"field,omitempty"`
		Group        *string                 `json:"group,omitempty"`
		Index        *string                 `json:"index,omitempty"`
		IsMultiIndex *bool                   `json:"isMultiIndex,omitempty"`
		Key          *string                 `json:"key,omitempty"`
		Negate       *bool                   `json:"negate,omitempty"`
		Params       *map[string]interface{} `json:"params,omitempty"`
		Type         *string                 `json:"type,omitempty"`
		Value        *string                 `json:"value,omitempty"`
	} `json:"meta,omitempty"`
	Query *map[string]interface{} `json:"query,omitempty"`
}

// FilterState represents the state of a filter
type FilterState struct {
	Store string `json:"store"`
}

// Timeframe defines the period during which actions can run
type Timeframe struct {
	Days     []int          `json:"days"`
	Hours    TimeframeHours `json:"hours"`
	Timezone string         `json:"timezone"`
}

// TimeframeHours defines the daily time window for actions
type TimeframeHours struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

// Frequency defines how often alerts generate actions
type Frequency struct {
	NotifyWhen string  `json:"notify_when"`
	Summary    bool    `json:"summary"`
	Throttle   *string `json:"throttle"`
}

// AlertDelay specifies conditions for delaying alert notifications
type AlertDelay struct {
	Active float32 `json:"active"`
}

// Flapping defines parameters for detecting rapidly changing alert states
type Flapping struct {
	LookBackWindow        float32 `json:"look_back_window"`
	StatusChangeThreshold float32 `json:"status_change_threshold"`
}

// Schedule defines how frequently the rule is evaluated
type Schedule struct {
	Interval string `json:"interval"`
}

// AlertingResponseBase contains common fields for Alerting responses
type AlertingResponseBase struct {
	ID                  string           `json:"id"`
	Name                string           `json:"name"`
	Tags                []string         `json:"tags"`
	Params              map[string]any   `json:"params"`
	Actions             []ActionResponse `json:"actions"`
	Enabled             bool             `json:"enabled"`
	Running             bool             `json:"running"`
	Consumer            string           `json:"consumer"`
	LastRun             LastRun          `json:"last_run"`
	MuteAll             bool             `json:"mute_all"`
	NextRun             string           `json:"next_run"`
	Revision            int              `json:"revision"`
	Schedule            Schedule         `json:"schedule"`
	Throttle            *string          `json:"throttle"`
	CreatedAt           string           `json:"created_at"`
	CreatedBy           string           `json:"created_by"`
	UpdatedAt           string           `json:"updated_at"`
	UpdatedBy           string           `json:"updated_by"`
	RuleTypeID          string           `json:"rule_type_id"`
	APIKeyOwner         string           `json:"api_key_owner"`
	MutedAlertIDs       []string         `json:"muted_alert_ids"`
	ExecutionStatus     ExecutionStatus  `json:"execution_status"`
	ScheduledTaskID     string           `json:"scheduled_task_id"`
	APIKeyCreatedByUser bool             `json:"api_key_created_by_user"`
}

// ActionCreate defines an action configuration when creating a rule
type ActionCreate struct {
	ID        string         `json:"id"`
	Group     string         `json:"group"`
	Params    map[string]any `json:"params"`
	Frequency Frequency      `json:"frequency"`
}
