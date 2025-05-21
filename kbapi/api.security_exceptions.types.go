package kbapi

import "encoding/json"

type SecurityExceptionsList struct {
	// UnderscoreVersion The version id, normally returned by the API when the item was retrieved.
	// Use it ensure updates are done against the latest version.
	// Readonly: for CreateList
	UnderscoreVersion *string `json:"_version,omitempty"`
	// CreatedBy time item as created
	// Readonly:
	CreatedAt *string `json:"created_at,omitempty"`
	// CreatedBy user who created the item
	// Readonly:
	CreatedBy *string `json:"created_by,omitempty"`
	// Description describes the excpetion list
	Description string `json:"description"`
	// ID
	ID string `json:"id,omitempty"`
	// Immutable
	// Readonly:
	Immutable bool `json:"immutable"`
	// ListID
	// Readonly: for CreateList
	ListID *string `json:"list_id,omitempty"`
	Meta   any     `json:"meta,omitempty"`
	Name   string  `json:"name"`
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string `json:"namespace_type,omitempty"`
	// OSTypes Use this field to specify the operating system.
	// Values are linux, macos, or windows. Default value is [] (empty).
	OSTypes *[]string `json:"os_types,omitempty"`
	// Tags String array containing words and phrases to help categorize exception items.
	Tags *[]string `json:"tags,omitempty"`
	// TieBreakerID Field used in search to ensure all containers are sorted and returned correctly.
	// Readonly:
	TieBreakerID string `json:"tie_breaker_id,omitempty"`
	// Type The type of exception list to be created. Different list types may denote where they can be utilized.
	// Values are detection, rule_default, endpoint, endpoint_trusted_apps, endpoint_events,
	// endpoint_host_isolation_exceptions, or endpoint_blocklists.
	Type string `json:"type"`
	// UpdatedAt
	// Readonly:
	UpdatedAt string `json:"updated_at,omitempty"`
	// UpdatedBy
	// Readonly:
	UpdatedBy string `json:"updated_by,omitempty"`
	// Version The document version, automatically increasd on updates.
	Version *string `json:"version,omitempty"`
}

type SecurityExceptionsItem struct {
	Comments []SecurityExceptionsComment `json:"comments,omitempty"`
	// CreatedBy time item as created
	// Readonly:
	CreatedAt *string `json:"created_at,omitempty"`
	// CreatedBy user who created the item
	// Readonly:
	CreatedBy *string `json:"created_by,omitempty"`
	// Description describes the excpetion list
	Description string `json:"description"`
	// Entries any of
	// - SecurityExceptionsItemEntryMatch
	// - SecurityExceptionsItemEntryMatchAny
	// - SecurityExceptionsItemEntryList
	// - SecurityExceptionsItemEntryExists
	// - SecurityExceptionsItemEntryNested
	// - SecurityExceptionsItemEntryWildCard
	Entries    []json.RawMessage `json:"entries"`
	ExpireTime *string           `json:"expire_time,omitempty"`
	// ID
	// Readonly:
	ID string `json:"id,omitempty"`
	// ItemID Human readable string identifier, e.g. trusted-linux-processes
	ItemID string `json:"item_id"`
	// ListID
	// Readonly: for create_items
	ListID *string `json:"list_id,omitempty"`
	Meta   any     `json:"meta,omitempty"`
	Name   string  `json:"name"`
	// NamespaceType Determines whether the exception container is available in all Kibana spaces or just the space in which it is created, where:
	// - single: Only available in the Kibana space in which it is created.
	// - agnostic: Available in all Kibana spaces.
	NamespaceType *string `json:"namespace_type,omitempty"`
	// OSTypes Use this field to specify the operating system.
	// Values are linux, macos, or windows. Default value is [] (empty).
	OSTypes *[]string `json:"os_types,omitempty"`
	// Tags String array containing words and phrases to help categorize exception items.
	Tags *[]string `json:"tags,omitempty"`
	// TieBreakerID Field used in search to ensure all containers are sorted and returned correctly.
	// Readonly:
	TieBreakerID string `json:"tie_breaker_id,omitempty"`
	// Type Value is simple.
	Type string `json:"type"`
	// UpdatedAt
	// Readonly:
	UpdatedAt string `json:"updated_at,omitempty"`
	// UpdatedBy
	// Readonly:
	UpdatedBy string `json:"updated_by,omitempty"`
	// Version
	// Readonly:
	Version *string `json:"_version,omitempty"`
}

type SecurityExceptionsComment struct {
	Comment string `json:"comment"`
}

type SecurityExceptionsItemEntryMatch struct {
	// Field A string that does not contain only whitespace characters
	Field string `json:"field"`
	// Operator Values are excluded or included.
	Operator string `json:"operator"`
	// Type Value is match.
	Type string `json:"type"`
	// Value A string that does not contain only whitespace characters
	Value string `json:"value"`
}

type SecurityExceptionsItemEntryMatchAny struct {
	// Field A string that does not contain only whitespace characters
	Field string `json:"field"`
	// Operator Values are excluded or included.
	Operator string `json:"operator"`
	// Type Value is match_any.
	Type string `json:"type"`
	// Value A string that does not contain only whitespace characters
	Value []string `json:"value"`
}

type SecurityExceptionsItemEntryList struct {
	// Field A string that does not contain only whitespace characters
	Field string                                `json:"field"`
	List  SecurityExceptionsItemEntryListObject `json:"list"`
	// Operator Values are excluded or included.
	Operator string `json:"operator"`
	// Type Value is list.
	Type string `json:"type"`
}

type SecurityExceptionsItemEntryListObject struct {
	// ID Value list's identifier.
	ID string `json:"id"`
	// Type Specifies the Elasticsearch data type of excludes the list container holds.
	// Values are binary, boolean, byte, date, date_nanos, date_range, double, double_range, float,
	// float_range, geo_point, geo_shape, half_float, integer, integer_range, ip, ip_range, keyword,
	// long, long_range, shape, short, or text.
	Type string `json:"type"`
}

type SecurityExceptionsItemEntryExists struct {
	// Field A string that does not contain only whitespace characters
	Field string `json:"field"`
	// Operator Values are excluded or included.
	Operator string `json:"operator"`
	// Type Value is exists.
	Type string `json:"type"`
}

type SecurityExceptionsItemEntryNested struct {
	// Entries any of
	// - SecurityExceptionsItemEntryMatch
	// - SecurityExceptionsItemEntryMatchAny
	// - SecurityExceptionsItemEntryExists
	Entries []json.RawMessage `json:"entries"`
	// Field A string that does not contain only whitespace characters
	Field string `json:"field"`
	// Type Value is nested.
	Type string `json:"type"`
}

type SecurityExceptionsItemEntryWildCard struct {
	// Field A string that does not contain only whitespace characters
	Field string `json:"field"`
	// Operator Values are excluded or included.
	Operator string `json:"operator"`
	// Type Value is wildcard.
	Type string `json:"type"`
	// Value A string that does not contain only whitespace characters
	Value string `json:"value"`
}

type SecurityExceptionsErrorDetail struct {
	Error  SecurityExceptionsErrorInfo `json:"error"`
	ItemID string                      `json:"item_id,omitempty"`
	ListID string                      `json:"list_id"`
}

type SecurityExceptionsErrorInfo struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
