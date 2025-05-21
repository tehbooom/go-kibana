package kbapi

type SecurityEndpointManagementAction struct {
	ID            string   `json:"id"`
	Agents        []string `json:"agents"`
	Command       string   `json:"command"`
	AgentType     string   `json:"agentType"`
	CreatedBy     string   `json:"createdBy"`
	IsExpired     bool     `json:"isExpired"`
	StartedAt     string   `json:"startedAt"`
	CompletedAt   string   `json:"completedAt"`
	IsCompleted   bool     `json:"isCompleted"`
	WasSuccessful bool     `json:"wasSuccessful"`
}

type SecurityEndpointManagementGetActionStatusRequestQueryParam struct {
	// AgentIDs At least 1 but not more than 50 elements. Minimum length is 1.
	AgentIDs []string `json:"agent_ids"`
}
