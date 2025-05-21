package kbapi

// APMUIAgentConfigurationObject Agent configuration
type APMUIAgentConfigurationObject struct {
	AgentName      *string          `json:"agent_name,omitempty"`
	AppliedByAgent *bool            `json:"applied_by_agent,omitempty"`
	AtTimestamp    float32          `json:"@Timestamp"`
	Etag           string           `json:"etag"`
	Service        APMServiceObject `json:"service"`
	// Settings Agent configuration settings
	Settings map[string]string `json:"settings"`
}

type APMServiceObject struct {
	Environment *string `json:"environment,omitempty"`
	Name        *string `json:"name,omitempty"`
}

type APMEnvironmentObject struct {
	AlreadyConfigured *bool   `json:"alreadyConfigured,omitempty"`
	Name              *string `json:"name,omitempty"`
}
