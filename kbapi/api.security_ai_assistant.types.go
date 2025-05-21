package kbapi

import "encoding/json"

type SecurityAIAssistantAnonymizationFieldObject struct {
	Allowed    *bool   `json:"allowed,omitempty"`
	Anonymized *bool   `json:"anonymized,omitempty"`
	CreatedAt  *string `json:"createdAt,omitempty"`
	CreatedBy  *string `json:"createdBy,omitempty"`
	Field      string  `json:"field"`
	ID         string  `json:"id"`
	Namespace  *string `json:"namespace,omitempty"`
	Timestamp  *string `json:"timestamp,omitempty"`
	UpdatedAt  *string `json:"updatedAt,omitempty"`
	UpdatedBy  *string `json:"updatedBy,omitempty"`
}

type SecurityAIAssistantAnonymizationCreateUpdateObject struct {
	Allowed    *bool  `json:"allowed,omitempty"`
	Anonymized *bool  `json:"anonymized,omitempty"`
	Field      string `json:"field"`
}

type SecurityAIAssistantDeleteObject struct {
	IDs   []string `json:"ids"`
	Query *string  `json:"query,omitempty"`
}

type SecurityAIAssistantBulkActionSummary struct {
	Failed    int `json:"failed"`
	Skipped   int `json:"skipped"`
	Succeeded int `json:"succeeded"`
	Total     int `json:"total"`
}

type SecurityAIAssistantBulkActionAnonymizationResults struct {
	Created []SecurityAIAssistantAnonymizationFieldObject `json:"created"`
	Deleted []string                                      `json:"deleted"`
	Skipped []SecurityAIAssistantBulkActionSkipResult     `json:"skipped"`
	Updated []SecurityAIAssistantAnonymizationFieldObject `json:"updated"`
}

type SecurityAIAssistantBulkActionSkipResult struct {
	ID         string  `json:"id"`
	Name       *string `json:"name,omitempty"`
	SkipReason string  `json:"skip_reason"`
}

type SecurityAIAssistantBulkActionAnonymizationErrors struct {
	AnonymizationFields []SecurityAIAssistantBulkActionErrorDetailsInError `json:"anonymization_fields"`
	ErrCode             *string                                            `json:"err_code,omitempty"`
	Message             string                                             `json:"message"`
	StatusCode          int                                                `json:"status_code"`
}

type SecurityAIAssistantBulkActionErrorDetailsInError struct {
	ID   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

type SecurityAIAssistantChatMessage struct {
	Content           *string         `json:"content,omitempty"`
	Data              *map[string]any `json:"data,omitempty"`
	FieldsToAnonymize *[]string       `json:"fields_to_anonymize,omitempty"`
	// Role message role.
	// Values are system, user, or assistant.
	Role string `json:"role"`
}

type SecurityAIAssistantAPIConfig struct {
	ActionTypeID          string  `json:"actionTypeId"`
	ConnectorID           string  `json:"connectorId"`
	DefaultSystemPromptID *string `json:"defaultSystemPromptId,omitempty"`
	Model                 *string `json:"model,omitempty"`
	Provider              *string `json:"provider,omitempty"`
}

type SecurityAIAssistantMessage struct {
	Content   string                              `json:"content"`
	IsError   *bool                               `json:"isError,omitempty"`
	Metadata  *SecurityAIAssistantMessageMetadata `json:"metadata,omitempty"`
	Reader    *map[string]any                     `json:"reader,omitempty"`
	Role      string                              `json:"role"`
	Timestamp string                              `json:"timestamp"`
	TraceData *SecurityAIAssistantTraceData       `json:"traceData,omitempty"`
}

type SecurityAIAssistantMessageMetadata struct {
	ContentReferences *map[string]any `json:"contentReferences,omitempty"`
}

type SecurityAIAssistantTraceData struct {
	TraceID       *string `json:"traceId,omitempty"`
	TransactionID *string `json:"transactionId,omitempty"`
}
type SecurityAIAssistantConversationResponse struct {
	ApiConfig                          *SecurityAIAssistantAPIConfig          `json:"apiConfig,omitempty"`
	Category                           *string                                `json:"category,omitempty"`
	CreatedAt                          string                                 `json:"createdAt"`
	ExcludeFromLastConversationStorage *bool                                  `json:"excludeFromLastConversationStorage,omitempty"`
	ID                                 *string                                `json:"id,omitempty"`
	IsDefault                          bool                                   `json:"isDefault"`
	Messages                           *[]SecurityAIAssistantMessage          `json:"messages,omitempty"`
	Namespace                          string                                 `json:"namespace"`
	Replacements                       *map[string]string                     `json:"replacements,omitempty"`
	Summary                            SecurityAIAssistantConversationSummary `json:"summary"`
	Timestamp                          string                                 `json:"timestamp"`
	Title                              string                                 `json:"title"`
	UpdatedAt                          string                                 `json:"updatedAt"`
	Users                              []SecurityAIAssistantUsers             `json:"users"`
}

type SecurityAIAssistantConversationSummary struct {
	Confidence string `json:"confidence"`
	Content    string `json:"content"`
	Public     bool   `json:"public"`
	Timestamp  string `json:"timestamp"`
}

type SecurityAIAssistantUsers struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SecurityAIAssistantIndexEntryRequest struct {
	// Global Whether this Knowledge Base Entry is global, defaults to false
	Global bool
	// Name Name of the Knowledge Base Entry
	Name string
	// Namespace Kibana Space, defaults to 'default' space
	Namespace string
	// Type Entry type
	// Value is document or index.
	Type string `json:"type"`
	// Users Users who have access to the Knowledge Base Entry, defaults to current user.
	// Empty array provides access to all users.
	Users *[]SecurityAIAssistantUsers `json:"users,omitempty"`
	// Description Description for when this index or data stream should be queried for Knowledge Base content.
	// Passed to the LLM as a tool description
	Description string `json:"description"`
	// Field Field to query for Knowledge Base content
	Field string `json:"field"`
	// Index Index or Data Stream to query for Knowledge Base content
	Index string `json:"index"`
	// InputSchema Array of objects defining the input schema, allowing the LLM to extract
	// structured data to be used in retrieval
	InputSchema *[]SecurityAIAssistantInputSchema `json:"inputSchema,omitempty"`
	// OutputFields Fields to extract from the query result, defaults to all fields if not provided or empty
	OutputFields *[]string `json:"outputFields,omitempty"`
	// QueryDescription Description of query field used to fetch Knowledge Base content.
	// Passed to the LLM as part of the tool input schema
	QueryDescription string `json:"queryDescription"`
}

type SecurityAIAssistantInputSchema struct {
	// Description Description of the field
	Description string `json:"description"`
	// FieldName Name of the field
	FieldName string `json:"fieldName"`
	// FieldType Type of the field
	FieldType string `json:"fieldType"`
}

type SecurityAIAssistantDocumentEntryRequest struct {
	// Global Whether this Knowledge Base Entry is global, defaults to false
	Global bool
	// Name Name of the Knowledge Base Entry
	Name string
	// Namespace Kibana Space, defaults to 'default' space
	Namespace string
	// Type Entry type
	// Value is document or index.
	Type string `json:"type"`
	// Users Users who have access to the Knowledge Base Entry, defaults to current user.
	// Empty array provides access to all users.
	Users *[]SecurityAIAssistantUsers `json:"users,omitempty"`
	// KBResource Knowledge Base resource name for grouping entries, e.g. 'security_labs', 'user', etc
	// Values are security_labs or user.
	KBResource string `json:"kbResource"`
	// Required Whether this resource should always be included, defaults to false
	Required *bool `json:"required,omitempty"`
	// Source Source document name or filepath
	Source string `json:"source"`
	// Text Knowledge Base Entry content
	Text string `json:"text"`
	// Vector Object containing Knowledge Base Entry text embeddings and modelId used to create the embeddings
	Vector *SecurityAIAssistantVector `json:"vector,omitempty"`
}

type SecurityAIAssistantVector struct {
	// ModelID ID of the model used to create the embeddings
	ModelID string `json:"modelId"`
	// Tokens Tokens with their corresponding values
	Tokens map[string]float32 `json:"tokens"`
}

type SecurityAIAssistantDocumentEntryResponse struct {
	// Global Whether this Knowledge Base Entry is global, defaults to false
	Global bool
	// ID A string that does not contain only whitespace characters
	ID string `json:"id"`
	// Name Name of the Knowledge Base Entry
	Name string
	// Namespace Kibana Space, defaults to 'default' space
	Namespace string
	// CreatedAt Time the Knowledge Base Entry was created
	CreatedAt string `json:"createdAt"`
	// CreatedBy User who created the Knowledge Base Entry
	CreatedBy string `json:"createdBy"`
	// Type Entry type
	// Value is document or index.
	Type string `json:"type"`
	// UpdatedAt Time the Knowledge Base Entry was last updated
	UpdatedAt string `json:"updatedAt"`
	// UpdatedBy User who last updated the Knowledge Base Entry
	UpdatedBy string `json:"updatedBy"`
	// Users Users who have access to the Knowledge Base Entry, defaults to current user.
	// Empty array provides access to all users.
	Users *[]SecurityAIAssistantUsers `json:"users,omitempty"`
	// KBResource Knowledge Base resource name for grouping entries, e.g. 'security_labs', 'user', etc
	// Values are security_labs or user.
	KBResource string `json:"kbResource"`
	// Required Whether this resource should always be included, defaults to false
	Required *bool `json:"required,omitempty"`
	// Source Source document name or filepath
	Source string `json:"source"`
	// Text Knowledge Base Entry content
	Text string `json:"text"`
	// Vector Object containing Knowledge Base Entry text embeddings and modelId used to create the embeddings
	Vector *SecurityAIAssistantVector `json:"vector,omitempty"`
}

type SecurityAIAssistantIndexEntryResponse struct {
	// Global Whether this Knowledge Base Entry is global, defaults to false
	Global bool
	// ID A string that does not contain only whitespace characters
	ID string `json:"id"`
	// Name Name of the Knowledge Base Entry
	Name string
	// Namespace Kibana Space, defaults to 'default' space
	Namespace string
	// CreatedAt Time the Knowledge Base Entry was created
	CreatedAt string `json:"createdAt"`
	// CreatedBy User who created the Knowledge Base Entry
	CreatedBy string `json:"createdBy"`
	// Type Entry type
	// Value is document or index.
	Type string `json:"type"`
	// UpdatedAt Time the Knowledge Base Entry was last updated
	UpdatedAt string `json:"updatedAt"`
	// UpdatedBy User who last updated the Knowledge Base Entry
	UpdatedBy string `json:"updatedBy"`
	// Users Users who have access to the Knowledge Base Entry, defaults to current user.
	// Empty array provides access to all users.
	Users *[]SecurityAIAssistantUsers `json:"users,omitempty"`
	// Description Description for when this index or data stream should be queried for Knowledge Base content.
	// Passed to the LLM as a tool description
	Description string `json:"description"`
	// Field Field to query for Knowledge Base content
	Field string `json:"field"`
	// Index Index or Data Stream to query for Knowledge Base content
	Index string `json:"index"`
	// InputSchema Array of objects defining the input schema, allowing the LLM to extract
	// structured data to be used in retrieval
	InputSchema *[]SecurityAIAssistantInputSchema `json:"inputSchema,omitempty"`
	// OutputFields Fields to extract from the query result, defaults to all fields if not provided or empty
	OutputFields *[]string `json:"outputFields,omitempty"`
	// QueryDescription Description of query field used to fetch Knowledge Base content.
	// Passed to the LLM as part of the tool input schema
	QueryDescription string `json:"queryDescription"`
}

type SecurityAIAssistantBulkActionKBErrors struct {
	KnowledgeBaseEntries []SecurityAIAssistantBulkActionErrorDetailsInError `json:"knowledgeBaseEntries"`
	ErrCode              *string                                            `json:"err_code,omitempty"`
	Message              string                                             `json:"message"`
	StatusCode           int                                                `json:"status_code"`
}

type SecurityAIAssistantBulkActionKBResults struct {
	Created []json.RawMessage                         `json:"created"`
	Deleted []string                                  `json:"deleted"`
	Skipped []SecurityAIAssistantBulkActionSkipResult `json:"skipped"`
	Updated []json.RawMessage                         `json:"updated"`
}

type SecurityAIAssistantPromptResponse struct {
	Categories               *[]string `json:"categories,omitempty"`
	Color                    *string   `json:"color,omitempty"`
	Consumer                 *string   `json:"consumer,omitempty"`
	Content                  string    `json:"content"`
	CreatedAt                *string   `json:"createdAt,omitempty"`
	CreatedBy                *string   `json:"createdBy,omitempty"`
	ID                       string    `json:"id"`
	IsDefault                *bool     `json:"isDefault,omitempty"`
	IsNewConversationDefault *bool     `json:"isNewConversationDefault,omitempty"`
	Name                     string    `json:"name"`
	Namespace                *string   `json:"namespace,omitempty"`
	// PromptType
	// Values are system or quick.
	PromptType string                      `json:"promptType"`
	Timestamp  *string                     `json:"timestamp,omitempty"`
	UpdatedAt  *string                     `json:"updatedAt,omitempty"`
	UpdatedBy  *string                     `json:"updatedBy,omitempty"`
	Users      *[]SecurityAIAssistantUsers `json:"users,omitempty"`
}

type SecurityAIAssistantPromptRequest struct {
	Categories *[]string `json:"categories,omitempty"`
	Color      *string   `json:"color,omitempty"`
	Consumer   *string   `json:"consumer,omitempty"`
	Content    string    `json:"content"`
	// ID can only be used for the Update object in a bulk action request.
	ID                       *string `json:"id,omitempty"`
	IsDefault                *bool   `json:"isDefault,omitempty"`
	IsNewConversationDefault *bool   `json:"isNewConversationDefault,omitempty"`
	Name                     string  `json:"name"`
	// PromptType
	// Values are system or quick.
	PromptType string `json:"promptType"`
}

type SecurityAIAssistantBulkActionPromptsErrors struct {
	Prompts    []SecurityAIAssistantBulkActionErrorDetailsInError `json:"knowledgeBaseEntries"`
	ErrCode    *string                                            `json:"err_code,omitempty"`
	Message    string                                             `json:"message"`
	StatusCode int                                                `json:"status_code"`
}

type SecurityAIAssistantBulkActionPromptsResults struct {
	Created []SecurityAIAssistantPromptResponse       `json:"created"`
	Deleted []string                                  `json:"deleted"`
	Skipped []SecurityAIAssistantBulkActionSkipResult `json:"skipped"`
	Updated []SecurityAIAssistantPromptResponse       `json:"updated"`
}
