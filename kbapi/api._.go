package kbapi

import (
	"context"
)

type API struct {
	transport Transport
	Alerting
	APM
	Cases
	Connectors
	Dataviews
	Fleet
	Endpoint
	Logstash
	ML
	Roles
	SavedObjects
	SecurityAIAssistant
	SecurityDetections
	SecurityEndpointManagement
	SecurityExceptions
	ShortURL
	Spaces
	Status
	TaskManager
	Uptime
}

type Alerting struct {
	// Create creates the specified alerting rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id
	Create func(ctx context.Context, req *AlertingCreateRequest, opts ...RequestOption) (*AlertingCreateResponse, error)
	// Delete deletes the specified alerting rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-alerting-rule-id
	Delete func(ctx context.Context, req *AlertingDeleteRequest, opts ...RequestOption) (*AlertingDeleteResponse, error)
	// Disable disables the specified alerting rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id-disable
	Disable func(ctx context.Context, req *AlertingDisableRequest, opts ...RequestOption) (*AlertingDisableResponse, error)
	// Enable enables the specified alerting rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id-enable
	Enable func(ctx context.Context, req *AlertingEnableRequest, opts ...RequestOption) (*AlertingEnableResponse, error)
	// Get returns the specified rule details. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-alerting-rule-id
	Get func(ctx context.Context, req *AlertingGetRequest, opts ...RequestOption) (*AlertingGetResponse, error)
	// GetTypes returns the rule types. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getruletypes
	GetTypes func(ctx context.Context, opts ...RequestOption) (*AlertingGetTypesResponse, error)
	// Health returns the alerting framework health. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getalertinghealth
	Health func(ctx context.Context, opts ...RequestOption) (*AlertingHealthResponse, error)
	// List returns all alerting rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-alerting-rules-find
	List func(ctx context.Context, req *AlertingListRequest, opts ...RequestOption) (*AlertingListResponse, error)
	// Mute mutes the specified alert for the rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-rule-id-alert-alert-id-mute
	Mute func(ctx context.Context, req *AlertingMuteRequest, opts ...RequestOption) (*AlertingMuteResponse, error)
	// MuteAll mutes all alerts for the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id-mute-all
	MuteAll func(ctx context.Context, req *AlertingMuteAllRequest, opts ...RequestOption) (*AlertingMuteAllResponse, error)
	// Unmute unmutes the specified alert for the rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-rule-id-alert-alert-id-unmute
	Unmute func(ctx context.Context, req *AlertingUnmuteRequest, opts ...RequestOption) (*AlertingUnmuteResponse, error)
	// UnmuteAll unmutes all alerts for the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id-unmute-all
	UnmuteAll func(ctx context.Context, req *AlertingUnmuteAllRequest, opts ...RequestOption) (*AlertingUnmuteAllResponse, error)
	// Update updates the specified alerting rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-alerting-rule-id
	Update func(ctx context.Context, req *AlertingUpdateRequest, opts ...RequestOption) (*AlertingUpdateResponse, error)
	// UpdateAPIkey updates the API key for the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-alerting-rule-id-update-api-key
	UpdateAPIkey func(ctx context.Context, req *AlertingUpdateAPIKeyRequest, opts ...RequestOption) (*AlertingUpdateAPIKeyResponse, error)
}

type APM struct {
	AgentConfiguration
	AgentKey
	Annotation
	ServerSchema
	SourceMaps
}

type AgentConfiguration struct {
	// CreateUpdate creates or updates the specified APM Agent configuration. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createupdateagentconfiguration
	CreateUpdate func(ctx context.Context, req *APMAgentConfigurationCreateUpdateRequest, opts ...RequestOption) (*APMAgentConfigurationCreateUpdateResponse, error)
	// Get returns the specified APM Agent configuration. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getsingleagentconfiguration
	Get func(ctx context.Context, req *APMAgentConfigurationGetRequest, opts ...RequestOption) (*APMAgentConfigurationGetResponse, error)
	// GetEnvironments reutnrs the environments for a service. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getenvironmentsforservice
	GetEnvironments func(ctx context.Context, req *APMAgentConfigurationGetEnvironmentRequest, opts ...RequestOption) (*APMAgentConfigurationGetEnvironmentResponse, error)
	// GetName returns the Agent name for a service. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getagentnameforservice
	GetName func(ctx context.Context, req *APMAgentConfigurationGetNameRequest, opts ...RequestOption) (*APMAgentConfigurationGetNameResponse, error)
	// Delete deletes the specified APM Agent configuration. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteagentconfiguration
	Delete func(ctx context.Context, req *APMAgentConfigurationDeleteRequest, opts ...RequestOption) (*APMAgentConfigurationDeleteResponse, error)
	// List returns all APM Agent configurations. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getagentconfigurations
	List func(ctx context.Context, opts ...RequestOption) (*APMAgentConfigurationListResponse, error)
	// Lookup searches for a single agent configuration and update the 'applied_by_agent' field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-searchsingleconfiguration
	Lookup func(ctx context.Context, req *APMAgentConfigurationLookupRequest, opts ...RequestOption) (*APMAgentConfigurationLookupResponse, error)
}

type AgentKey struct {
	// Create creates a new agent key for APM. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createagentkey
	Create func(ctx context.Context, req *APMAgentKeyCreateRequest, opts ...RequestOption) (*APMAgentKeyCreateResponse, error)
}

type Annotation struct {
	// Create creates a new annotation for a specific service. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createannotation
	Create func(ctx context.Context, req *APMAnnotationCreateRequest, opts ...RequestOption) (*APMAnnotationCreateResponse, error)
	// Search searches for annotations related to a specific service. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getannotation
	Search func(ctx context.Context, req *APMAnnotationSearchRequest, opts ...RequestOption) (*APMAnnotationSearchResponse, error)
}

type ServerSchema struct {
	// Save saves the APM server schema. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-saveapmserverschema
	Save func(ctx context.Context, req *APMServerSchemaSaveRequest, opts ...RequestOption) (*APMServerSchemaSaveResponse, error)
}

type SourceMaps struct {
	// Delete deletes a previously uploaded source map. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletesourcemap
	Delete func(ctx context.Context, req *APMSourcemapsDeleteRequest, opts ...RequestOption) (*APMSourcemapsDeleteResponse, error)
	// Get returns an array of Fleet artifacts, including source map uploads. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getsourcemaps
	Get func(ctx context.Context, req *APMSourcemapsGetRequest, opts ...RequestOption) (*APMSourcemapsGetResponse, error)
	// Upload uploads a source map for a specific service and version. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-uploadsourcemap
	Upload func(ctx context.Context, req *APMSourcemapsUploadRequest, opts ...RequestOption) (*APMSourcemapsUploadResponse, error)
}

type Cases struct {
	// AddCommentAlert adds a comment or alert to the specified case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-addcasecommentdefaultspace
	AddCommentAlert func(ctx context.Context, req *CasesAddCommentAlertRequest, opts ...RequestOption) (*CasesAddCommentAlertResponse, error)
	// AddSettings adds case settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-setcaseconfigurationdefaultspace
	AddSettings func(ctx context.Context, req *CasesAddSettingsRequest, opts ...RequestOption) (*CasesAddSettingsResponse, error)
	// AttachFile attaches a file to the specified case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-addcasefiledefaultspace
	AttachFile func(ctx context.Context, req *CasesAttachFileRequest, opts ...RequestOption) (*CasesAttachFileResponse, error)
	// Create creates a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createcasedefaultspace
	Create func(ctx context.Context, req *CasesCreateRequest, opts ...RequestOption) (*CasesCreateResponse, error)
	// Delete deletes a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletecasedefaultspace
	Delete func(ctx context.Context, req *CasesDeleteRequest, opts ...RequestOption) (*CasesDeleteResponse, error)
	// DeleteAlertComment deletes teh specified alert or comment for the case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletecasecommentdefaultspace
	DeleteAlertComment func(ctx context.Context, req *CasesDeleteAlertCommentRequest, opts ...RequestOption) (*CasesDeleteAlertCommentResponse, error)
	// DeleteAllAlertsComments deletes all comments and alerts from a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletecasecommentsdefaultspace
	DeleteAllAlertsComments func(ctx context.Context, req *CasesDeleteAllAlertsCommentsRequest, opts ...RequestOption) (*CasesDeleteAllAlertsCommentsResponse, error)
	// Get returns the specified case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasedefaultspace
	Get func(ctx context.Context, req *CasesGetRequest, opts ...RequestOption) (*CasesGetResponse, error)
	// GetAlertComment returns the specified comment for a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasecommentdefaultspace
	GetAlertComment func(ctx context.Context, req *CasesGetAlertCommentRequest, opts ...RequestOption) (*CasesGetAlertCommentResponse, error)
	// GetAllAlerts returns all alerts for the specified case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasealertsdefaultspace
	GetAllAlerts func(ctx context.Context, req *CasesGetAllAlertsRequest, opts ...RequestOption) (*CasesGetAllAlertsResponse, error)
	// GetConnectors returns information about connectors that are supported for use in cases. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findcaseconnectorsdefaultspace
	GetConnectors func(ctx context.Context, opts ...RequestOption) (*CasesGetConnectorsResponse, error)
	// GetCreators returns information about the users who opened cases. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasereportersdefaultspace
	GetCreators func(ctx context.Context, req *CasesGetCreatorsRequest, opts ...RequestOption) (*CasesGetCreatorsResponse, error)
	// GetSettings returns setting details. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcaseconfigurationdefaultspace
	GetSettings func(ctx context.Context, req *CasesGetSettingsRequest, opts ...RequestOption) (*CasesGetSettingsResponse, error)
	// GetTags returns a list of case tags. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasetagsdefaultspace
	GetTags func(ctx context.Context, req *CasesGetTagsRequest, opts ...RequestOption) (*CasesGetTagsResponse, error)
	// ListAlertsComments retrieves a paginated list of comments for a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findcasecommentsdefaultspace
	ListAlertsComments func(ctx context.Context, req *CasesListCommentsAlertsRequest, opts ...RequestOption) (*CasesListCommentsAlertsResponse, error)
	// ListActivity retrives a paginated list of user activity for the specified case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findcaseactivitydefaultspace
	ListActivity func(ctx context.Context, req *CasesListActivityRequest, opts ...RequestOption) (*CasesListActivityResponse, error)
	// ListFromAlert returns all cases for the specified alert. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getcasesbyalertdefaultspace
	ListFromAlert func(ctx context.Context, req *CasesListFromAlertRequest, opts ...RequestOption) (*CasesListFromAlertResponse, error)
	// Push pushes a case to an external service. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-pushcasedefaultspace
	Push func(ctx context.Context, req *CasesPushRequest, opts ...RequestOption) (*CasesPushResponse, error)
	// Search searches for cases. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findcasesdefaultspace
	Search func(ctx context.Context, req *CasesSearchRequest, opts ...RequestOption) (*CasesSearchResponse, error)
	// Update updates a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updatecasedefaultspace
	Update func(ctx context.Context, req *CasesUpdateRequest, opts ...RequestOption) (*CasesUpdateResponse, error)
	// UpdateAlertComment updates the specified commnet or alert for a case. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updatecasecommentdefaultspace
	UpdateAlertComment func(ctx context.Context, req *CasesUpdateCommentAlertRequest, opts ...RequestOption) (*CasesUpdateCommentAlertResponse, error)
	// UpdateSettings updates the specified configuration settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updatecaseconfigurationdefaultspace
	UpdateSettings func(ctx context.Context, req *CasesUpdateSettingsRequest, opts ...RequestOption) (*CasesUpdateSettingsResponse, error)
}

type Connectors struct {
	// Create creates a connector. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-actions-connector-id
	Create func(ctx context.Context, req *ConnectorsCreateRequest, opts ...RequestOption) (*ConnectorsCreateResponse, error)
	// Delete deletes the specified connector. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-actions-connector-id
	Delete func(ctx context.Context, req *ConnectorsDeleteRequest, opts ...RequestOption) (*ConnectorsDeleteResponse, error)
	// Get returns the specified connector. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-actions-connector-id
	Get func(ctx context.Context, req *ConnectorsGetRequest, opts ...RequestOption) (*ConnectorsGetResponse, error)
	// GetTypes returns all connector types. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-actions-connector-types
	GetTypes func(ctx context.Context, req *ConnectorsGetTypesRequest, opts ...RequestOption) (*ConnectorsGetTypesResponse, error)
	// List returns all connectors. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-actions-connectors
	List func(ctx context.Context, opts ...RequestOption) (*ConnectorsListResponse, error)
	// Run executes a connector action. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-actions-connector-id-execute
	Run func(ctx context.Context, req *ConnectorsRunRequest, opts ...RequestOption) (*ConnectorsRunResponse, error)
	// Update updates the specified connector. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-actions-connector-id
	Update func(ctx context.Context, req *ConnectorsUpdateRequest, opts ...RequestOption) (*ConnectorsUpdateResponse, error)
}

type Dataviews struct {
	// Create creates a data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createdataviewdefaultw
	Create func(ctx context.Context, req *DataViewsCreateRequest, opts ...RequestOption) (*DataViewsCreateResponse, error)
	// CreateRuntimeField creates a runtime field for the specified data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createruntimefielddefault
	CreateRuntimeField func(ctx context.Context, req *DataViewsCreateRuntimeFieldRequest, opts ...RequestOption) (*DataViewsCreateRuntimeFieldResponse, error)
	// CreateUpdateRuntimeField creates or updates a runtime field for the specified data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createupdateruntimefielddefault
	CreateUpdateRuntimeField func(ctx context.Context, req *DataViewsCreateUpdateRuntimeFieldRequest, opts ...RequestOption) (*DataViewsCreateUpdateRuntimeFieldResponse, error)
	// Delete deletes the specified data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletedataviewdefault
	Delete func(ctx context.Context, req *DataViewsDeleteRequest, opts ...RequestOption) (*DataViewsDeleteResponse, error)
	// DeleteRuntimeField deletes the specified runtime field from a data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteruntimefielddefault
	DeleteRuntimeField func(ctx context.Context, req *DataViewsDeleteRuntimeFieldRequest, opts ...RequestOption) (*DataViewsDeleteRuntimeFieldResponse, error)
	// Get returns the specified data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getdataviewdefault
	Get func(ctx context.Context, req *DataViewsGetRequest, opts ...RequestOption) (*DataViewsGetResponse, error)
	// GetDefault returns the default data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getdefaultdataviewdefault
	GetDefault func(ctx context.Context, opts ...RequestOption) (*DataViewsGetDefaultResponse, error)
	// GetRuntimeField returns the specified runtime field for a data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getruntimefielddefault
	GetRuntimeField func(ctx context.Context, req *DataViewsGetRuntimeFieldRequest, opts ...RequestOption) (*DataViewsGetRuntimeFieldResponse, error)
	// List returns all data views. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-getalldataviewsdefault
	List func(ctx context.Context, opts ...RequestOption) (*DataViewsListResponse, error)
	// PreviewSavedObjectSwap preview the impact of swapping saved object references from one data view identifier to another. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-previewswapdataviewsdefault
	PreviewSavedObjectSwap func(ctx context.Context, req *DataViewsPreviewSavedObjectSwapRequest, opts ...RequestOption) (*DataViewsPreviewSavedObjectSwapResponse, error)
	// SetDefault sets the default data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-setdefaultdatailviewdefault
	SetDefault func(ctx context.Context, req *DataViewsSetDefaultRequest, opts ...RequestOption) (*DataViewsSetDefaultResponse, error)
	// SwapSavedObjectReference Changes saved object references from one data view identifier to another. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-swapdataviewsdefault
	// WARNING: Misuse can break large numbers of saved objects!
	SwapSavedObjectReference func(ctx context.Context, req *DataViewsSwapSavedObjectReferenceRequest, opts ...RequestOption) (*DataViewsSwapSavedObjectReferenceResponse, error)
	// Update updates the specified data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updatedataviewdefault
	Update func(ctx context.Context, req *DataViewsUpdateRequest, opts ...RequestOption) (*DataViewsUpdateResponse, error)
	// UpdateFieldMetadata updates the specified fields metadata. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updatefieldsmetadatadefault
	UpdateFieldMetadata func(ctx context.Context, req *DataViewsUpdateFieldMetadataRequest, opts ...RequestOption) (*DataViewsUpdateFieldMetadataResponse, error)
	// UpdateRuntimeField updates the specified runtime field for a data view. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updateruntimefielddefault
	UpdateRuntimeField func(ctx context.Context, req *DataViewsUpdateRuntimeFieldRequest, opts ...RequestOption) (*DataViewsUpdateRuntimeFieldResponse, error)
}

type Fleet struct {
	AgentActions          AgentActions
	Agents                Agents
	AgentPolicies         AgentPolicies
	BinaryDownloadSources BinaryDownloadSources
	DataStreams           DataStreams
	EPM                   EPM
	EnrollmentAPIKeys     EnrollmentAPIKeys
	Internal              Internal
	MessageSigningService MessageSigningService
	Outputs               Outputs
	PackagePolicies       PackagePolicies
	Proxies               Proxies
	ServerHost            ServerHost
	ServiceToken          ServiceToken
	UninstallTokens       UninstallTokens
}

type Agents struct {
	// Status returns all agents status. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-status
	Status func(ctx context.Context, req *FleetAgentStatusRequest, opts ...RequestOption) (*FleetAgentStatusResponse, error)
	// StatusData returns the specified agents incoming data. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-status-data
	StatusData func(ctx context.Context, req *FleetAgentStatusDataRequest, opts ...RequestOption) (*FleetAgentStatusDataResponse, error)
	// ListAgents returns all agents. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents
	List func(ctx context.Context, req *FleetListAgentsRequest, opts ...RequestOption) (*FleetListAgentsResponse, error)
	// ListByActionID returns a list of agents by action ID. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents
	ListByActionID func(ctx context.Context, req *FleetListAgentsByActionIDRequest, opts ...RequestOption) (*FleetListAgentsByActionIDResponse, error)
	// GetAgentByID returns the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-agentid
	GetAgent func(ctx context.Context, req *FleetGetAgentRequest, opts ...RequestOption) (*FleetGetAgentResponse, error)
	// UpdateAgent updates the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-agents-agentid
	UpdateAgent func(ctx context.Context, req *FleetUpdateAgentRequest, opts ...RequestOption) (*FleetUpdateAgentResponse, error)
	// ListFiles returns a list of files uploaded by the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-agentid-uploads
	ListFiles func(ctx context.Context, req *FleetListAgentUploadsRequest, opts ...RequestOption) (*FleetListAgentUploadsResponse, error)
	// DeleteFile deletes the specified file. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-agents-files-fileid
	DeleteFile func(ctx context.Context, req *FleetDeleteFileRequest, opts ...RequestOption) (*FleetDeleteFileResponse, error)
	// GetFile returns the specified file. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-files-fileid-filename
	GetFile func(ctx context.Context, req *FleetGetAgentFileRequest, opts ...RequestOption) (*FleetGetAgentFileResponse, error)
	// GetSetup returns the current setup state. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-setup
	GetSetup func(ctx context.Context, opts ...RequestOption) (*FleetGetAgentSetupResponse, error)
	// InitiateSetup initiates Fleet. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-setup
	InitiateSetup func(ctx context.Context, req *FleetInitiateSetupRequest, opts ...RequestOption) (*FleetInitiateSetupResponse, error)
	// ListTags returns all agent tags in Fleet. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-tags
	ListTags func(ctx context.Context, opts ...RequestOption) (*FleetListTagsReponse, error)
}

type AgentActions struct {
	// BulkReassign reassigns the specified agents. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-bulk-reassign
	BulkReassign func(ctx context.Context, req *FleetBulkReassignAgentRequest, opts ...RequestOption) (*FleetBulkReassignAgentResponse, error)
	// BulkGetDiagnostics gets diagnostics for the specified agents. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-bulk-request-diagnostics
	BulkGetDiagnostics func(ctx context.Context, req *FleetBulkGetDiagnosticsAgentRequest, opts ...RequestOption) (*FleetBulkGetDiagnosticsAgentResponse, error)
	// BulkUpdateAgentTags updates tags for the specified agents. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-bulk-update-agent-tags
	BulkUpdateAgentTags func(ctx context.Context, req *FleetBulkUpdateAgentTagsRequest, opts ...RequestOption) (*FleetBulkUpdateAgentTagsResponse, error)
	// BulkUpgrade upgrades the specified agents. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-bulk-upgrade
	BulkUpgrade func(ctx context.Context, req *FleetBulkUpgradeAgentsRequest, opts ...RequestOption) (*FleetBulkUpgradeAgentsResponse, error)
	// Cancel cancles the specified agent action. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-actions-actionid-cancel
	Cancel func(ctx context.Context, req *FleetAgentActionsCancelRequest, opts ...RequestOption) (*FleetAgentActionsCancelResponse, error)
	// Create creates an agent action for the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-agentid-actions
	Create func(ctx context.Context, req *FleetAgentActionsCreateRequest, opts ...RequestOption) (*FleetAgentActionsCreateResponse, error)
	// GetDiagnostics gets diagnostics for the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-agentid-request-diagnostics
	GetDiagnostics func(ctx context.Context, req *FleetGetDiagnosticsAgentRequest, opts ...RequestOption) (*FleetGetDiagnosticsAgentResponse, error)
	// ListStatus gets all agent actions status. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agents-action-status
	ListStatus func(ctx context.Context, req *FleetAgentActionsListStatusRequest, opts ...RequestOption) (*FleetAgentActionsListStatusResponse, error)
	// Reassign reassigns the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-agentid-reassign
	Reassign func(ctx context.Context, req *FleetReassignAgentRequest, opts ...RequestOption) (*FleetReassignAgentResponse, error)
	// Unenroll unenrolls the specified agent. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-agentid-unenroll
	Unenroll func(ctx context.Context, req *FleetUnenrollAgentRequest, opts ...RequestOption) (*FleetUnenrollAgentResponse, error)
	// Upgrade upgrades the speicifed agent to the requested version. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agents-agentid-upgrade
	Upgrade func(ctx context.Context, req *FleetUpgradeAgentRequest, opts ...RequestOption) (*FleetUpgradeAgentResponse, error)
}

type AgentPolicies struct {
	// List returns a list of agent policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-policies
	List func(ctx context.Context, req *FleetAgentPoliciesRequest, opts ...RequestOption) (*FleetListAgentPoliciesResponse, error)
	// Get gets specified agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-policies-agentpolicyid
	Get func(ctx context.Context, req *FleetGetAgentPolicyRequest, opts ...RequestOption) (*FleetGetAgentPolicyResponse, error)
	// Create creates a new agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agent-policies
	Create func(ctx context.Context, req *FleetCreateAgentPolicyRequest, opts ...RequestOption) (*FleetCreateAgentPolicyResponse, error)
	// BulkGet gets specified agent policies in bulk. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agent-policies-bulk-get
	BulkGet func(ctx context.Context, req *FleetBulkGetAgentPoliciesRequest, opts ...RequestOption) (*FleetBulkGetAgentPoliciesResponse, error)
	// Update updates the specified agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-agent-policies-agentpolicyid
	Update func(ctx context.Context, req *FleetUpdateAgentPolicyRequest, opts ...RequestOption) (*FleetUpdateAgentPolicyResponse, error)
	// Copy copies the specified agent policy to a new agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-agent-policies-agentpolicyid
	Copy func(ctx context.Context, req *FleetCopyAgentPolicyRequest, opts ...RequestOption) (*FleetCopyAgentPolicyResponse, error)
	// Download downloads the specified agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-policies-agentpolicyid-download
	Download func(ctx context.Context, req *FleetDownloadAgentPolicyRequest, opts ...RequestOption) (*FleetDownloadAgentPolicyResponse, error)
	// GetFull gets the policy for the specified agent policy. There are two possible outputs. If you set kubernetes as true
	// use AsKubernetes method to get the yaml. Otherwise, use AsJSON method to get the json response.
	// See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-policies-agentpolicyid-full
	GetFull func(ctx context.Context, req *FleetGetFullAgentPolicyRequest, opts ...RequestOption) (*FleetGetFullAgentPolicyResponse, error)
	// Delete  deletes the specified agent policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agent-policies-delete
	Delete func(ctx context.Context, req *FleetDeleteAgentPolicyRequest, opts ...RequestOption) (*FleetDeleteAgentPolicyResponse, error)
}

type BinaryDownloadSources struct {
	// Create creates a binary download source. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-agent-download-sources
	Create func(ctx context.Context, req *FleetBinaryDownloadCreateRequest, opts ...RequestOption) (*FleetBinaryDownloadCreateResponse, error)
	// Delete deletes the specified binary download source. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-agent-download-sources-sourceid
	Delete func(ctx context.Context, req *FleetBinaryDownloadDeleteRequest, opts ...RequestOption) (*FleetBinaryDownloadDeleteResponse, error)
	// Get returns the specified binary download source. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-download-sources-sourceid
	Get func(ctx context.Context, req *FleetBinaryDownloadGetRequest, opts ...RequestOption) (*FleetBinaryDownloadGetResponse, error)
	// List returns all download sources. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-agent-download-sources
	List func(ctx context.Context, opts ...RequestOption) (*FleetBinaryDownloadListResponse, error)
	// Update updates the specified binary download source. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-agent-download-sources-sourceid
	Update func(ctx context.Context, req *FleetBinaryDownloadUpdateRequest, opts ...RequestOption) (*FleetBinaryDownloadUpdateResponse, error)
}
type DataStreams struct {
	// ListAll returns all datastreams and their metadata. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-data-streams-list
	ListAll func(ctx context.Context, opts ...RequestOption) (*FleetDataStreamsListResponse, error)
}

type EPM struct {
	// AuthorizeTransforms authroizes the specified transform. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-epm-packages-pkgname-pkgversion-transforms-authorize
	AuthorizeTransforms func(ctx context.Context, req *FleetEPMAuthorizeTransformsRequest, opts ...RequestOption) (*FleetEPMAuthorizeTransformsResponse, error)
	// BulkInstallPackages installs packages in bulk. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-epm-packages-bulk
	BulkInstallPackages func(ctx context.Context, req *FleetEPMBulkInstallPackagesRequest, opts ...RequestOption) (*FleetEPMBulkInstallPackagesResponse, error)
	// CreateCustomIntegration creates a custom integration. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-fleet-epm-custom-integrations
	CreateCustomIntegration func(ctx context.Context, req *FleetEPMCreateCustomIntegrationRequest, opts ...RequestOption) (*FleetEPMCreateCustomIntegrationResponse, error)
	// DeletePackage deletes the specified package. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-epm-packages-pkgname-pkgversion
	DeletePackage func(ctx context.Context, req *FleetEPMDeletePackageRequest, opts ...RequestOption) (*FleetEPMDeletePackageResponse, error)
	// GetInputsTemplate gets the specified packages input templates. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-templates-pkgname-pkgversion-inputs
	GetInputsTemplate func(ctx context.Context, req *FleetEPMGetInputsTemplateRequest, opts ...RequestOption) (*FleetEPMGetInputsTemplateResponse, error)
	// GetPackage gets the specified package. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-packages-pkgname-pkgversion
	GetPackage func(ctx context.Context, req *FleetEPMGetPackageRequest, opts ...RequestOption) (*FleetEPMGetPackageResponse, error)
	// GetPackageFile returns the specified file. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-packages-pkgname-pkgversion-filepath
	GetPackageFile func(ctx context.Context, req *FleetEPMGetPackageFileRequest, opts ...RequestOption) (*FleetEPMGetPackageFileResponse, error)
	// GetPackageSignatureVerificiationID  returns the package signature verification ID. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-verification-key-id
	GetPackageSignatureVerificiationID func(ctx context.Context, opts ...RequestOption) (*FleetEPMGetPackageSignatureVerificationIDResponse, error)
	// GetPackagesInstalled gets all installed packages. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-packages-installed
	GetPackagesInstalled func(ctx context.Context, req *FleetEPMGetInstalledPackagesRequest, opts ...RequestOption) (*FleetEPMGetInstalledPackagesResponse, error)
	// GetPackagesLimited gets a limited package list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-packages-limited
	GetPackagesLimited func(ctx context.Context, opts ...RequestOption) (*FleetEPMGetPackagesLimitedResponse, error)
	// GetPackageStats returns the stats for the specified package. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-packages-pkgname-stats
	GetPackageStats func(ctx context.Context, req *FleetEPMGetPackageStatsRequest, opts ...RequestOption) (*FleetEPMGetPackageStatsResponse, error)
	// InstallPackageRegistry installs a package using the registry. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-epm-packages-pkgname-pkgversion
	InstallPackageRegistry func(ctx context.Context, req *FleetEPMInstallPackageRegistryRequest, opts ...RequestOption) (*FleetEPMInstallPackageRegistryResponse, error)
	// InstallPackageUpload installs a package by upload. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-epm-packages
	InstallPackageUpload func(ctx context.Context, req *FleetEPMInstallPackageUploadRequest, opts ...RequestOption) (*FleetEPMInstallPackageUploadResponse, error)
	// ListCategories lists package categories. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-get-package-categories
	ListCategories func(ctx context.Context, req *FleetEPMListPkgCategoriesRequest, opts ...RequestOption) (*FleetEPMListPkgCategoriesResponse, error)
	// ListDataStreams returns datastreams. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-epm-data-streams
	ListDataStreams func(ctx context.Context, req *FleetEPMListDataStreamsRequest, opts ...RequestOption) (*FleetEPMListDataStreamsResponse, error)
	// ListPackages lists packages. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-get-fleet-epm-packages
	ListPackages func(ctx context.Context, req *FleetEPMListPackagesRequest, opts ...RequestOption) (*FleetEPMListPackagesResponse, error)
	// UpdatePackageSettings updates the specified package settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-epm-packages-pkgname-pkgversion
	UpdatePackageSettings func(ctx context.Context, req *FleetEPMUpdatePackageSettingsRequest, opts ...RequestOption) (*FleetEPMUpdatePackageSettingsResponse, error)
}

type EnrollmentAPIKeys struct {
	// Create creates an API Key for the specified policy ID. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-enrollment-api-keys
	Create func(ctx context.Context, req *FleetEnrollmentAPIKeysCreateRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysCreateResponse, error)
	// Get returns the specified enrollment API key. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-enrollment-api-keys-keyid
	Get func(ctx context.Context, req *FleetEnrollmentAPIKeysGetRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysGetResponse, error)
	// List retuns all enrollment API keys. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-enrollment-api-keys
	List func(ctx context.Context, req *FleetEnrollmentAPIKeysListRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysListResponse, error)
	// Revoke revokes an enrollment API key by ID by marking it as inactive. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-enrollment-api-keys-keyid
	Revoke func(ctx context.Context, req *FleetEnrollmentAPIKeysRevokeRequest, opts ...RequestOption) (*FleetEnrollmentAPIKeysRevokeResponse, error)
}

type Internal struct {
	// CheckFleetServerHealth gets the specified Fleet server's health status. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-health-check
	CheckFleetServerHealth func(ctx context.Context, req *FleetInternalCheckFleetServerHealthRequest, opts ...RequestOption) (*FleetInternalCheckFleetServerHealthResponse, error)
	// CheckPermissions checks permissions. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-check-permissions
	CheckPermissions func(ctx context.Context, req *FleetInternalCheckPermissionsRequest, opts ...RequestOption) (*FleetInternalCheckPermissionsResponse, error)
	// GetSettings gets the current Fleet settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-settings
	GetSettings func(ctx context.Context, opts ...RequestOption) (*FleetInternalGetSettingsResponse, error)
	// UpdateSettings updates Fleet settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-settings
	UpdateSettings func(ctx context.Context, req *FleetInternalUpdateSettingsRequest, opts ...RequestOption) (*FleetInternalUpdateSettingsResponse, error)
	// InitiateFleetSetup initiates the fleet server. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-setup
	InitiateFleetSetup func(ctx context.Context, opts ...RequestOption) (*FleetInternalInitiateFleetSetupResponse, error)
}

type MessageSigningService struct {
	// Rotate rotates a Fleet message signing key pair. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-message-signing-service-rotate-key-pair
	Rotate func(ctx context.Context, req *FleetMessageSigningServiceRotateRequest, opts ...RequestOption) (*FleetMessageSigningServiceRotateResponse, error)
}

type Outputs struct {
	// Create creates a new output. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-outputs
	Create func(ctx context.Context, req *FleetOutputsCreateRequest, opts ...RequestOption) (*FleetOutputsCreateResponse, error)
	// Delete deletes the specified output. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-outputs-outputid
	Delete func(ctx context.Context, req *FleetOutputsDeleteRequest, opts ...RequestOption) (*FleetOutputsDeleteResponse, error)
	// GenerateLogstashKey generates a Logstash API key. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-logstash-api-keys
	GenerateLogstashKey func(ctx context.Context, opts ...RequestOption) (*FleetOutputsGenerateLogstashKeyResponse, error)
	// Get returns the specified output. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-outputs-outputid
	Get func(ctx context.Context, req *FleetOutputsGetRequest, opts ...RequestOption) (*FleetOutputsGetResponse, error)
	// Health returns the health of the specified output. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-outputs-outputid-health
	Health func(ctx context.Context, req *FleetOutputsHealthRequest, opts ...RequestOption) (*FleetOutputsHealthResponse, error)
	// List returns all outputs in Fleet. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-outputs
	List func(ctx context.Context, opts ...RequestOption) (*FleetOutputsListResponse, error)
	// Update updates the specified output. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-outputs-outputid
	Update func(ctx context.Context, req *FleetOutputsUpdateRequest, opts ...RequestOption) (*FleetOutputsUpdateResponse, error)
}

type PackagePolicies struct {
	// BulkDelete deletes the specified package policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-package-policies-delete
	BulkDelete func(ctx context.Context, req *FleetPackagePoliciesBulkDeleteRequest, opts ...RequestOption) (*FleetPackagePoliciesBulkDeleteResponse, error)
	// BulkGet returns the specified package policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-package-policies-bulk-get
	BulkGet func(ctx context.Context, req *FleetPackagePoliciesBulkGetRequest, opts ...RequestOption) (*FleetPackagePoliciesBulkGetResponse, error)
	// Create creates the specified package policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-package-policies
	Create func(ctx context.Context, req *FleetPackagePoliciesCreateRequest, opts ...RequestOption) (*FleetPackagePoliciesCreateResponse, error)
	// Delete deletes the specified package policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-package-policies-packagepolicyid
	Delete func(ctx context.Context, req *FleetPackagePoliciesDeleteRequest, opts ...RequestOption) (*FleetPackagePoliciesDeleteResponse, error)
	// Get returns the specified package policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-package-policies-packagepolicyid
	Get func(ctx context.Context, req *FleetPackagePoliciesGetRequest, opts ...RequestOption) (*FleetPackagePoliciesGetResponse, error)
	// List returns all package policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-package-policies
	List func(ctx context.Context, req *FleetPackagePoliciesListRequest, opts ...RequestOption) (*FleetPackagePoliciesListResponse, error)
	// Update updates the specified package policy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-package-policies-packagepolicyid
	Update func(ctx context.Context, req *FleetPackagePoliciesUpdateRequest, opts ...RequestOption) (*FleetPackagePoliciesUpdateResponse, error)
	// Upgrade upgrades the specified package policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-package-policies-upgrade
	Upgrade func(ctx context.Context, req *FleetPackagePoliciesUpgradeRequest, opts ...RequestOption) (*FleetPackagePoliciesUpgradeResponse, error)
	// UpgradeDryRun runs a dry upgrade for the specified package policies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-package-policies-upgrade-dryrun
	UpgradeDryRun func(ctx context.Context, req *FleetPackagePoliciesUpgradeDryRunRequest, opts ...RequestOption) (*FleetPackagePoliciesUpgradeDryRunResponse, error)
}

type Proxies struct {
	// Create creates the specified proxy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-proxies
	Create func(ctx context.Context, req *FleetProxiesCreateRequest, opts ...RequestOption) (*FleetProxiesCreateResponse, error)
	// Delete deletes the specified proxy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-proxies-itemid
	Delete func(ctx context.Context, req *FleetProxiesDeleteRequest, opts ...RequestOption) (*FleetProxiesDeleteResponse, error)
	// Get returns the specified proxy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-proxies-itemid
	Get func(ctx context.Context, req *FleetProxiesGetRequest, opts ...RequestOption) (*FleetProxiesGetResponse, error)
	// List returns all proxies. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-proxies
	List func(ctx context.Context, opts ...RequestOption) (*FleetProxiesListResponse, error)
	// Update updates the specified proxy. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-proxies-itemid
	Update func(ctx context.Context, req *FleetProxiesUpdateRequest, opts ...RequestOption) (*FleetProxiesUpdateResponse, error)
}

type ServerHost struct {
	// Create creates the specified Fleet server. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-fleet-server-hosts
	Create func(ctx context.Context, req *FleetServerHostCreateRequest, opts ...RequestOption) (*FleetServerHostCreateResponse, error)
	// Delete deletes the specified Fleet server. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-fleet-fleet-server-hosts-itemid
	Delete func(ctx context.Context, req *FleetServerHostDeleteRequest, opts ...RequestOption) (*FleetServerHostDeleteResponse, error)
	// Get returns the specified Fleet server. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-fleet-server-hosts-itemid
	Get func(ctx context.Context, req *FleetServerHostGetRequest, opts ...RequestOption) (*FleetServerHostGetResponse, error)
	// List returns Fleet Server Hosts. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-fleet-server-hosts
	List func(ctx context.Context, opts ...RequestOption) (*FleetServerHostListResponse, error)
	// Update updates the specified Fleet server host. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-fleet-fleet-server-hosts-itemid
	Update func(ctx context.Context, req *FleetServerHostUpdateRequest, opts ...RequestOption) (*FleetServerHostUpdateResponse, error)
}

type ServiceToken struct {
	// Create creates a Fleet service token. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-fleet-service-tokens
	Create func(ctx context.Context, req *FleetServiceTokenCreateRequest, opts ...RequestOption) (*FleetServiceTokenCreateResponse, error)
}

type UninstallTokens struct {
	// GetDecrypted returns the specified token decrypted. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-uninstall-tokens-uninstalltokenid
	GetDecrypted func(ctx context.Context, req *FleetUninstallTokensGetDecryptedRequest, opts ...RequestOption) (*FleetUninstallTokensGetDecryptedResponse, error)
	// GetMetadata gets metadata for latest uninstall tokens. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-fleet-uninstall-tokens
	GetMetadata func(ctx context.Context, opts ...RequestOption) (*FleetUninstallTokensGetMetadataResponse, error)
}

type Endpoint struct {
	Exceptions Exceptions
}

type Exceptions struct {
	// CreateItem creates an endpoint exception list item. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createendpointlistitem
	CreateItem func(ctx context.Context, req *EndpointExceptionsCreateItemRequest, opts ...RequestOption) (*EndpointExceptionsCreateItemResponse, error)
	// CreateList creates an endpoint exception list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createendpointlist
	CreateList func(ctx context.Context, opts ...RequestOption) (*EndpointExceptionsCreateListResponse, error)
	// Delete deletes the specified endpoint exception list item. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteendpointlistitem
	Delete func(ctx context.Context, req *EndpointExceptionsDeleteItemRequest, opts ...RequestOption) (*EndpointExceptionsDeleteItemResponse, error)
	// Get returns the specified endpoint exception list item. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readendpointlistitem
	Get func(ctx context.Context, req *EndpointExceptionsGetRequest, opts ...RequestOption) (*EndpointExceptionsGetResponse, error)
	// ListItems returns a list of all endpoint exception list items. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findendpointlistitems
	ListItems func(ctx context.Context, req *EndpointExceptionsListItemsRequest, opts ...RequestOption) (*EndpointExceptionsListItemsResponse, error)
	// Update updates the specified endpoint exception list item. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updateendpointlistitem
	Update func(ctx context.Context, req *EndpointExceptionsUpdateRequest, opts ...RequestOption) (*EndpointExceptionsUpdateResponse, error)
}

type Logstash struct {
	// Delete deletes the specified Logstash pipeline. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-logstash-pipeline
	Delete func(ctx context.Context, req *LogstashDeletePipelineRequest, opts ...RequestOption) (*LogstashDeletePipelineResponse, error)
	// Get returns the specified Logstash pipeline. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-logstash-pipeline
	Get func(ctx context.Context, req *LogstashGetPipelineRequest, opts ...RequestOption) (*LogstashGetPipelineResponse, error)
	// List returns all Logstash pipelines. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-logstash-pipelines
	List func(ctx context.Context, opts ...RequestOption) (*LogstashListPipelinesResponse, error)
	// Put creates or updates the specified Logstash pipeline. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-logstash-pipeline
	Put func(ctx context.Context, req *LogstashPutPipelineRequest, opts ...RequestOption) (*LogstashPutPipelineResponse, error)
}

type ML struct {
	// SyncSavedObjects synchronizes Kibana saved objects for machine learning jobs and trained models in the default space. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-mlsync
	SyncSavedObjects func(ctx context.Context, req *MLSyncSavedObjectsRequest, opts ...RequestOption) (*MLSyncSavedObjectsResponse, error)
}

type Roles struct {
	// CreateOrUpdateMulti creates or updates multiple roles. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-security-roles
	CreateOrUpdateMulti func(ctx context.Context, req *RolesCreateOrUpdateMultiRequest, opts ...RequestOption) (*RolesCreateOrUpdateMultiResponse, error)
	// CreateOrUpdateSingle creates or updates the specified role. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-security-role-name
	CreateOrUpdateSingle func(ctx context.Context, req *RolesCreateUpdateSingleRoleRequest, opts ...RequestOption) (*RolesCreateUpdateSingleRoleResponse, error)
	// Delete deletes the specified role. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-security-role-name
	Delete func(ctx context.Context, req *RolesDeleteRequest, opts ...RequestOption) (*RolesDeleteResponse, error)
	// Get returns the specified role. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-security-role-name
	Get func(ctx context.Context, req *RolesGetRequest, opts ...RequestOption) (*RolesGetResponse, error)
	// List returns all security roles. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-security-role
	List func(ctx context.Context, req *RolesListRequest, opts ...RequestOption) (*RolesListResponse, error)
}

type SavedObjects struct {
	// Export exports the specified objects. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-exportsavedobjectsdefault
	Export func(ctx context.Context, req *SavedObjectExportRequest, opts ...RequestOption) (*SavedObjectExportResponse, error)
	// Import imports the specified objects. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-importsavedobjectsdefault
	Import func(ctx context.Context, req *SavedObjectImportRequest, opts ...RequestOption) (*SavedObjectImportResponse, error)
	// ResolveImport resolves import errors. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-resolveimporterrors
	ResolveImport func(ctx context.Context, req *SavedObjectResolveImportsRequest, opts ...RequestOption) (*SavedObjectResolveImportsResponse, error)
	// RotateKey rorates a key for encrypted saved objects. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-rotateencryptionkey
	RotateKey func(ctx context.Context, req *SavedObjectRotateKeyRequest, opts ...RequestOption) (*SavedObjectRotateKeyResponse, error)
}

type SecurityAIAssistant struct {
	// BulkActionAnonymizationFields applies a bulk action to multiple anonymization fields. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-performanonymizationfieldsbulkaction
	BulkActionAnonymizationFields func(ctx context.Context, req *SecurityAIAssistantBulkActionAnonymizationRequest, opts ...RequestOption) (*SecurityAIAssistantBulkActionAnonymizationResponse, error)
	// BulkActionKnowledgeBaseEntires applies a bulk action to all Knowledge Base Entries that match the filter or to the list of Knowledge Base Entries by their IDs. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-performknowledgebaseentrybulkaction
	BulkActionKnowledgeBaseEntires func(ctx context.Context, req *SecurityAIAssistantBulkActionKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantBulkActionKnowledgeBaseEntryResponse, error)
	// BulkActionPrompts applies a bulk action to all prompts that match the filter or to the list of prompts by their IDs. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-performpromptsbulkaction
	BulkActionPrompts func(ctx context.Context, req *SecurityAIAssistantBulkActionPromptsRequest, opts ...RequestOption) (*SecurityAIAssistantBulkActionPromptsResponse, error)
	// CreateConversation creates a new Security AI Assistant conversation. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createconversation
	CreateConversation func(ctx context.Context, req *SecurityAIAssistantCreateConversationRequest, opts ...RequestOption) (*SecurityAIAssistantCreateConversationResponse, error)
	// CreateKnowledgeBase creates a Knowledge Base. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createknowledgebase
	CreateKnowledgeBase func(ctx context.Context, req *SecurityAIAssistantCreateKnowledgeBaseRequest, opts ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseResponse, error)
	// CreateKnowledgeBaseEntry creates a Knowledge Base Entry. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createknowledgebaseentry
	CreateKnowledgeBaseEntry func(ctx context.Context, req *SecurityAIAssistantCreateKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantCreateKnowledgeBaseEntryResponse, error)
	// CreateModelResponse creates a model response for the given chat conversation. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-chatcomplete
	CreateModelResponse func(ctx context.Context, req *SecurityAIAssistantCreateModelResponseRequest, opts ...RequestOption) (*SecurityAIAssistantCreateModelResponseResponse, error)
	// DeleteConversation deletes an existing conversation using the conversation ID. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteconversation
	DeleteConversation func(ctx context.Context, req *SecurityAIAssistantDeleteConversationRequest, opts ...RequestOption) (*SecurityAIAssistantDeleteConversationResponse, error)
	// DeleteKnowledgeBaseEntry deletes the specified Knowledge Base Entry. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteknowledgebaseentry
	DeleteKnowledgeBaseEntry func(ctx context.Context, req *SecurityAIAssistantDeleteKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantDeleteKnowledgeBaseEntryResponse, error)
	// GetConversation returns the details of the specified conversation. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readconversation
	GetConversation func(ctx context.Context, req *SecurityAIAssistantGetConversationRequest, opts ...RequestOption) (*SecurityAIAssistantGetConversationResponse, error)
	// GetKnowledgeBase returns the specified Knowledge Base. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readknowledgebase
	GetKnowledgeBase func(ctx context.Context, req *SecurityAIAssistantGetKnowledgeBaseRequest, opts ...RequestOption) (*SecurityAIAssistantGetKnowledgeBaseResponse, error)
	// GetKnowledgeBaseEntry returns the specified Knowledge Base Entry. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readknowledgebaseentry
	GetKnowledgeBaseEntry func(ctx context.Context, req *SecurityAIAssistantGetKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantGetKnowledgeBaseEntryResponse, error)
	// ListAnonymizationFields returns a list of all anonymization fields. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findanonymizationfields
	ListAnonymizationFields func(ctx context.Context, req *SecurityAIAssistantListAnonymizationRequest, opts ...RequestOption) (*SecurityAIAssistantListAnonymizationResponse, error)
	// ListConversations returns all conversations for the current user. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findconversations
	ListConversations func(ctx context.Context, req *SecurityAIAssistantListConversationsRequest, opts ...RequestOption) (*SecurityAIAssistantListConversationsResponse, error)
	// ListKnowledgeBaseEntries returns Knowledge Base entires that match the given query. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findknowledgebaseentries
	ListKnowledgeBaseEntries func(ctx context.Context, req *SecurityAIAssistantListKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantListKnowledgeBaseEntryResponse, error)
	// ListPrompts returns all prompts. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-performpromptsbulkaction
	ListPrompts func(ctx context.Context, req *SecurityAIAssistantListPromptsRequest, opts ...RequestOption) (*SecurityAIAssistantListPromptsResponse, error)
	// UpdateKnowledgeBaseEntry updates the specified Knowledge Base Entry. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updateknowledgebaseentry
	UpdateKnowledgeBaseEntry func(ctx context.Context, req *SecurityAIAssistantUpdateKnowledgeBaseEntryRequest, opts ...RequestOption) (*SecurityAIAssistantUpdateKnowledgeBaseEntryResponse, error)
}

type SecurityDetections struct {
	// AssignUsers assigns users to detection alerts, and unassigns them from alerts. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-setalertassignees
	AssignUsers func(ctx context.Context, req *SecurityDetectionsAssignUsersRequest, opts ...RequestOption) (*SecurityDetectionsAssignUsersResponse, error)
	// BulkActionRules applies a bulk action, such as bulk edit, duplicate, or delete, to multiple detection rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-performrulesbulkaction
	BulkActionRules func(ctx context.Context, req *SecurityDetectionsBulkActionRulesRequest, opts ...RequestOption) (*SecurityDetectionsBulkActionRulesResponse, error)
	// CreateIndex creates an alert index. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createalertsindex
	CreateIndex func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsCreateIndexResponse, error)
	// CreateRule creates a new detection rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createrule
	CreateRule func(ctx context.Context, req *SecurityDetectionsCreateRuleRequest, opts ...RequestOption) (*SecurityDetectionsCreateRuleResponse, error)
	// DeleteIndex deletes an alert index. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletealertsindex
	DeleteIndex func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsDeleteIndexResponse, error)
	// DeleteRule deletes the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleterule
	DeleteRule func(ctx context.Context, req *SecurityDetectionsDeleteRuleRequest, opts ...RequestOption) (*SecurityDetectionsDeleteRuleResponse, error)
	// ExportRules exports the specified rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-exportrules
	ExportRules func(ctx context.Context, req *SecurityDetectionsExportRulesRequest, opts ...RequestOption) (*SecurityDetectionsExportRulesResponse, error)
	// GetIndex returns the alert index name if it exists. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readalertsindex
	GetIndex func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsGetIndexResponse, error)
	// GetRule returns the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readrule
	GetRule func(ctx context.Context, req *SecurityDetectionsGetRuleRequest, opts ...RequestOption) (*SecurityDetectionsGetRuleResponse, error)
	// GetStatusPrebuilt returns the status of all Elastic prebuilt detection rules and Timelines. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readprebuiltrulesandtimelinesstatus
	GetStatusPrebuilt func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsGetStatusPrebuiltResponse, error)
	// GetPrivilegesSpace returns whether or not the user is authenticated, and the user's Kibana space and index privileges. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deletealertsindex
	GetPrivilegesSpace func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsGetPrivilegesSpaceResponse, error)
	// ImportRule imports the specified rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-importrules
	ImportRule func(ctx context.Context, req *SecurityDetectionsImportRulesRequest, opts ...RequestOption) (*SecurityDetectionsImportRulesResponse, error)
	// InstallPrebuilt installs and updates all Elastic prebuilt detection rules and Timelines. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-installprebuiltrulesandtimelines
	InstallPrebuilt func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsInstallPrebuiltResponse, error)
	// ListRules returns a list of rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findrules
	ListRules func(ctx context.Context, req *SecurityDetectionsListRulesRequest, opts ...RequestOption) (*SecurityDetectionsListRulesResponse, error)
	// ListTags returns all unique tags from all detection rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readtags
	ListTags func(ctx context.Context, opts ...RequestOption) (*SecurityDetectionsListTagsResponse, error)
	// PatchRule patches an existing rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-patchrule
	PatchRule func(ctx context.Context, req *SecurityDetectionsPatchRuleRequest, opts ...RequestOption) (*SecurityDetectionsPatchRuleResponse, error)
	// PreviewAlerts previews rule alerts generated on specified time range. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-rulepreview
	PreviewAlerts func(ctx context.Context, req *SecurityDetectionsPreviewAlertsRequest, opts ...RequestOption) (*SecurityDetectionsPreviewAlertsResponse, error)
	// SearchAlerts finds and/or aggregates detection alerts that match the given query. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-searchalerts
	SearchAlerts func(ctx context.Context, req *SecurityDetectionsSearchAlertsRequest, opts ...RequestOption) (*SecurityDetectionsSearchAlertsResponse, error)
	// SetAlertStatus sets the status of one or more detection alerts. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-setalertsstatus
	SetAlertStatus func(ctx context.Context, req *SecurityDetectionsSetAlertStatusRequest, opts ...RequestOption) (*SecurityDetectionsSetAlertStatusResponse, error)
	// UpdateRule updates the specified rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updaterule
	UpdateRule func(ctx context.Context, req *SecurityDetectionsUpdateRuleRequest, opts ...RequestOption) (*SecurityDetectionsUpdateRuleResponse, error)
	// UpdateTags adds tags to detection alerts, and removes them from alerts. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-setalerttags
	UpdateTags func(ctx context.Context, req *SecurityDetectionsUpdateTagsRequest, opts ...RequestOption) (*SecurityDetectionsUpdateTagsResponse, error)
}

type SecurityEndpointManagement struct {
	// GetActionStatus returns the status of response actions for the specified agent IDs. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-endpointgetactionsstatus
	GetActionStatus func(ctx context.Context, req *SecurityEndpointManagementGetActionStatusRequest, opts ...RequestOption) (*SecurityEndpointManagementGetActionStatusResponse, error)
	// ListActions returns a list of all response actions. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-endpointgetactionslist
	ListActions func(ctx context.Context, req *SecurityEndpointManagementListActionsRequest, opts ...RequestOption) (*SecurityEndpointManagementListActionsResponse, error)
}

type SecurityExceptions struct {
	// CreateItem creates an exception item and associate it with the specified exception list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createexceptionlistitem
	CreateItem func(ctx context.Context, req *SecurityExceptionsCreateItemRequest, opts ...RequestOption) (*SecurityExceptionsCreateItemResponse, error)
	// CreateItems creates exception items that apply to a single detection rule. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createruleexceptionlistitems
	CreateItems func(ctx context.Context, req *SecurityExceptionsCreateItemsRequest, opts ...RequestOption) (*SecurityExceptionsCreateItemsResponse, error)
	// CreateList creates an exception list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createexceptionlist
	CreateList func(ctx context.Context, req *SecurityExceptionsCreateListRequest, opts ...RequestOption) (*SecurityExceptionsCreateListResponse, error)
	// CreateSharedList creates a shared exception list can apply to multiple detection rules. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-createsharedexceptionlist
	CreateSharedList func(ctx context.Context, req *SecurityExceptionsCreateSharedListRequest, opts ...RequestOption) (*SecurityExceptionsCreateSharedListResponse, error)
	// DeleteItem deletes an exception list item using the id or item_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteexceptionlistitem
	DeleteItem func(ctx context.Context, req *SecurityExceptionsDeleteItemRequest, opts ...RequestOption) (*SecurityExceptionsDeleteItemResponse, error)
	// DeleteList deletes an exception list using the id or list_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-deleteexceptionlist
	DeleteList func(ctx context.Context, req *SecurityExceptionsDeleteListRequest, opts ...RequestOption) (*SecurityExceptionsDeleteListResponse, error)
	// DuplicateList duplicates an existing exception list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-duplicateexceptionlist
	DuplicateList func(ctx context.Context, req *SecurityExceptionsDuplicateListRequest, opts ...RequestOption) (*SecurityExceptionsDuplicateListResponse, error)
	// ExportList exports an exception list and its associated items to an NDJSON file. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-exportexceptionlist
	ExportList func(ctx context.Context, req *SecurityExceptionsExportListRequest, opts ...RequestOption) (*SecurityExceptionsExportListResponse, error)
	// GetList Get the details of an exception list using the id or list_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readexceptionlist
	GetList func(ctx context.Context, req *SecurityExceptionsGetListRequest, opts ...RequestOption) (*SecurityExceptionsGetListResponse, error)
	// GetItem returns the details of an exception list item using the id or item_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readexceptionlistitem
	GetItem func(ctx context.Context, req *SecurityExceptionsGetItemRequest, opts ...RequestOption) (*SecurityExceptionsGetItemResponse, error)
	// GetSummary returns a summary of the specified exception list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-readexceptionlistsummary
	GetSummary func(ctx context.Context, req *SecurityExceptionsGetSummaryRequest, opts ...RequestOption) (*SecurityExceptionsGetSummaryResponse, error)
	// ImportList imports an exception list and its associated items from an NDJSON file. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-importexceptionlist
	ImportList func(ctx context.Context, req *SecurityExceptionsImportListRequest, opts ...RequestOption) (*SecurityExceptionsImportListResponse, error)
	// ListItems returns a list of all exception list items in the specified list. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findexceptionlistitems
	ListItems func(ctx context.Context, req *SecurityExceptionsListItemsRequest, opts ...RequestOption) (*SecurityExceptionsListItemsResponse, error)
	// ListLists returns a list of all exception list containers. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-findexceptionlists
	ListLists func(ctx context.Context, req *SecurityExceptionsListListsRequest, opts ...RequestOption) (*SecurityExceptionsListListsResponse, error)
	// UpdateItem updates an exception list item using the id or item_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updateexceptionlistitem
	UpdateItem func(ctx context.Context, req *SecurityExceptionsUpdateItemRequest, opts ...RequestOption) (*SecurityExceptionsUpdateItemResponse, error)
	// UpdateList updates an exception list using the id or list_id field. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-updateexceptionlist
	UpdateList func(ctx context.Context, req *SecurityExceptionsUpdateListRequest, opts ...RequestOption) (*SecurityExceptionsUpdateListResponse, error)
}

type ShortURL struct {
	// Create creates a short URL. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-post-url
	Create func(ctx context.Context, req *ShortURLCreateRequest, opts ...RequestOption) (*ShortURLCreateResponse, error)
	// Delete deletes the specified Kibana short URL. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-delete-url
	Delete func(ctx context.Context, req *ShortURLDeleteRequest, opts ...RequestOption) (*ShortURLDeleteResponse, error)
	// Get returns the specified Kibana short URL. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-url
	Get func(ctx context.Context, req *ShortURLGetRequest, opts ...RequestOption) (*ShortURLGetResponse, error)
	// Resolve resolves a Kibana short URL by its slug. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-resolve-url
	Resolve func(ctx context.Context, req *ShortURLResolveRequest, opts ...RequestOption) (*ShortURLResolveResponse, error)
}

type Spaces struct {
	// CopyObjects copies saved objects between spaces. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-spaces-copy-saved-objects
	CopyObjects func(ctx context.Context, req *SpacesCopyObjectsRequest, opts ...RequestOption) (*SpacesCopyObjectsResponse, error)
	// Create creates a space. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-spaces-space
	Create func(ctx context.Context, req *SpacesCreateRequest, opts ...RequestOption) (*SpacesCreateResponse, error)
	// GetAll gets all spaces. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-get-spaces-space
	GetAll func(ctx context.Context, req *SpacesGetAllRequest, opts ...RequestOption) (*SpacesGetAllResponse, error)
	// Get gets the specified space. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-get-spaces-space-id
	Get func(ctx context.Context, req *SpacesGetRequest, opts ...RequestOption) (*SpacesGetResponse, error)
	// Update updates the specified space. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-put-spaces-space-id
	Update func(ctx context.Context, req *SpacesUpdateRequest, opts ...RequestOption) (*SpacesUpdateResponse, error)
	// Delete deletes the specified space. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-delete-spaces-space-id
	Delete func(ctx context.Context, req *SpacesDeleteRequest, opts ...RequestOption) (*SpacesDeleteResponse, error)
	// UpdateObjects  updates one or more saved objects to add or remove them from some spaces. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-spaces-update-objects-spaces
	UpdateObjects func(ctx context.Context, req *SpacesUpdateObjectsRequest, opts ...RequestOption) (*SpacesUpdateObjectsResponse, error)
	// GetShareableReferences  collects references and space contexts for saved objects. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-spaces-get-shareable-references
	GetShareableReferences func(ctx context.Context, req *SpacesShareableReferencesRequest, opts ...RequestOption) (*SpacesShareableReferencesResponse, error)
	// SpacesDisableLegacyURLAliases  leaves the alias intact but the legacy URL for the alias will no longer function. See https://www.elastic.co/docs/api/doc/kibana/operation/operation-post-spaces-disable-legacy-url-aliases
	SpacesDisableLegacyURLAliases func(ctx context.Context, req *SpacesDisableLegacyURLRequest, opts ...RequestOption) (*SpacesDisableLegacyURLResponse, error)
}

type Status struct {
	// Get returns Kibana's operational status as well as a detailed breakdown of plugin statuses. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-status#operation-get-status-200-body-application-json
	Get func(ctx context.Context, req *GetStatusRequest, opts ...RequestOption) (*StatusResponse, error)
	// GetRedacted returns a minimal representation of Kibana's operational status. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-status#operation-get-status-200-body-application-json-kibana_http_apis_core_status_redactedresponse-object
	GetRedacted func(ctx context.Context, req *GetStatusRequest, opts ...RequestOption) (*StatusRedactedResponse, error)
}

type TaskManager struct {
	// Health gets the task manager health. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-task-manager-health
	Health func(ctx context.Context, opts ...RequestOption) (*TaskManagerHealthResponse, error)
}

type Uptime struct {
	// GetSettings returns uptime settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-get-uptime-settings
	GetSettings func(ctx context.Context, opts ...RequestOption) (*UptimeGetSettingsResponse, error)
	// UpdateSettings updates uptime settings. See https://www.elastic.co/docs/api/doc/kibana/v9/operation/operation-put-uptime-settings
	UpdateSettings func(ctx context.Context, req *UptimeUpdateSettingsRequest, opts ...RequestOption) (*UptimeUpdateSettingsResponse, error)
}

// New creates a new API
func New(t Transport) *API {
	api := &API{
		transport: t,
	}

	api.Alerting = Alerting{
		Create:       api.newAlertingCreate(),
		Delete:       api.newAlertingDelete(),
		Disable:      api.newAlertingDisable(),
		Enable:       api.newAlertingEnable(),
		Get:          api.newAlertingGet(),
		GetTypes:     api.newAlertingGetTypes(),
		Health:       api.newAlertingHealth(),
		List:         api.newAlertingList(),
		Mute:         api.newAlertingMute(),
		MuteAll:      api.newAlertingMuteAll(),
		Unmute:       api.newAlertingUnmute(),
		UnmuteAll:    api.newAlertingUnmuteAll(),
		Update:       api.newAlertingUpdate(),
		UpdateAPIkey: api.newAlertingUpdateAPIKey(),
	}

	api.APM = APM{
		AgentConfiguration: AgentConfiguration{
			CreateUpdate:    api.newAPMAgentConfigurationCreateUpdate(),
			Get:             api.newAPMAgentConfigurationGet(),
			GetEnvironments: api.newAPMAgentConfigurationGetEnvironment(),
			GetName:         api.newAPMAgentConfigurationGetName(),
			Delete:          api.newAPMAgentConfigurationDelete(),
			List:            api.newAPMAgentConfigurationList(),
		},
		AgentKey: AgentKey{
			Create: api.newAPMAgentKeyCreate(),
		},
		Annotation: Annotation{
			Create: api.newAPMAnnotationCreate(),
			Search: api.newAPMAnnotationSearch(),
		},
		ServerSchema: ServerSchema{
			Save: api.newAPMServerSchemaSave(),
		},
		SourceMaps: SourceMaps{
			Delete: api.newAPMSourcemapsDelete(),
			Get:    api.newAPMSourcemapsGet(),
			Upload: api.newAPMSourcemapsUpload(),
		},
	}

	api.Cases = Cases{
		AddCommentAlert:         api.newCasesAddCommentAlert(),
		AddSettings:             api.newCasesAddSettings(),
		AttachFile:              api.newCasesAttachFile(),
		Create:                  api.newCasesCreate(),
		Delete:                  api.newCasesDelete(),
		DeleteAlertComment:      api.newCasesDeleteAlertComment(),
		DeleteAllAlertsComments: api.newCasesDeleteAllAlertsComments(),
		Get:                     api.newCasesGet(),
		GetAlertComment:         api.newCasesGetAlertComment(),
		GetAllAlerts:            api.newCasesGetAllAlerts(),
		GetConnectors:           api.newCasesGetConnectors(),
		GetCreators:             api.newCasesGetCreators(),
		GetSettings:             api.newCasesGetSettings(),
		GetTags:                 api.newCasesGetTags(),
		Push:                    api.newCasesPush(),
		ListActivity:            api.newCasesListActivity(),
		ListAlertsComments:      api.newCasesListCommentsAlerts(),
		ListFromAlert:           api.newCasesListFromAlert(),
		Search:                  api.newCasesSearch(),
		Update:                  api.newCasesUpdate(),
		UpdateAlertComment:      api.newCasesUpdateCommentAlert(),
		UpdateSettings:          api.newCasesUpdateSettings(),
	}

	api.Connectors = Connectors{
		Create:   api.newConnectorsCreate(),
		Delete:   api.newConnectorsDelete(),
		Get:      api.newConnectorsGet(),
		GetTypes: api.newConnectorsGetTypes(),
		List:     api.newConnectorsList(),
		Run:      api.newConnectorsRun(),
		Update:   api.newConnectorsUpdate(),
	}

	api.Dataviews = Dataviews{
		Create:                   api.newDataViewsCreate(),
		CreateRuntimeField:       api.newDataViewsCreateRuntimeField(),
		CreateUpdateRuntimeField: api.newDataViewsCreateUpdateRuntimeField(),
		Delete:                   api.newDataViewsDelete(),
		DeleteRuntimeField:       api.newDataViewsDeleteRuntimeField(),
		Get:                      api.newDataviewsGet(),
		GetDefault:               api.newDataViewsGetDefault(),
		GetRuntimeField:          api.newDataViewsGetRuntimeField(),
		List:                     api.newDataViewsList(),
		PreviewSavedObjectSwap:   api.newDataViewsPreviewSavedObjectSwap(),
		SetDefault:               api.newDataViewsSetDefault(),
		SwapSavedObjectReference: api.newDataViewsSwapSavedObjectReference(),
		Update:                   api.newDataViewsUpdate(),
		UpdateFieldMetadata:      api.newDataViewsUpdateFieldMetadata(),
		UpdateRuntimeField:       api.newDataViewsUpdateRuntimeField(),
	}

	api.Fleet = Fleet{
		Agents: Agents{
			DeleteFile:     api.newFleetDeleteFile(),
			GetAgent:       api.newFleetGetAgent(),
			GetFile:        api.newFleetGetAgentFile(),
			GetSetup:       api.newFleetGetAgentSetup(),
			InitiateSetup:  api.newFleetInitiateSetup(),
			Status:         api.newFleetAgentStatusFunc(),
			StatusData:     api.newFleetAgentStatusDataFunc(),
			List:           api.newFleetListAgents(),
			ListByActionID: api.newFleetListAgentsByActionID(),
			ListFiles:      api.newFleetListAgentUploads(),
			ListTags:       api.newFleetListTags(),
			UpdateAgent:    api.newFleetUpdateAgent(),
		},
		AgentPolicies: AgentPolicies{
			BulkGet:  api.newFleetBulkGetAgentPolicies(),
			Copy:     api.newFleetCopyAgentPolicy(),
			Create:   api.newFleetCreateAgentPolicyFunc(),
			Delete:   api.newFleetDeleteAgentPolicy(),
			Download: api.newFleetDownloadAgentPolicy(),
			Get:      api.newFleetGetAgentPolicy(),
			GetFull:  api.newFleetGetFullAgentPolicy(),
			List:     api.newFleetAgentListPoliciesFunc(),
			Update:   api.newFleetUpdateAgentPolicy(),
		},
		AgentActions: AgentActions{
			BulkGetDiagnostics:  api.newFleetBulkGetDiagnosticsAgents(),
			BulkReassign:        api.newFleetBulkReassignAgents(),
			BulkUpdateAgentTags: api.newFleetBulkUpdateAgentTags(),
			BulkUpgrade:         api.newFleetBulkUpgradeAgents(),
			Cancel:              api.newFleetAgentActionsCancel(),
			Create:              api.newFleetAgentActionsCreate(),
			GetDiagnostics:      api.newFleetGetDiagnosticsAgent(),
			ListStatus:          api.newFleetAgentActionsListStatus(),
			Reassign:            api.newFleetReassignAgent(),
			Unenroll:            api.newFleetUnenrollAgent(),
			Upgrade:             api.newFleetUpgradeAgent(),
		},
		BinaryDownloadSources: BinaryDownloadSources{
			Create: api.newFleetBinaryDownloadCreate(),
			Delete: api.newFleetBinaryDownloadDelete(),
			Get:    api.newFleetBinaryDownloadGet(),
			List:   api.newFleetBinaryDownloadList(),
			Update: api.newFleetBinaryDownloadUpdate(),
		},
		DataStreams: DataStreams{
			ListAll: api.newFleetDataStreamsList(),
		},
		EPM: EPM{
			AuthorizeTransforms:                api.newFleetEPMAuthorizeTransforms(),
			BulkInstallPackages:                api.newFleetEPMBulkInstallPackages(),
			CreateCustomIntegration:            api.newFleetEPMCreateCustomIntegration(),
			DeletePackage:                      api.newFleetEPMDeletePackage(),
			GetInputsTemplate:                  api.newFleetEPMGetInputsTemplate(),
			GetPackage:                         api.newFleetEPMGetPackage(),
			GetPackageFile:                     api.newFleetEPMGetPackageFile(),
			GetPackageSignatureVerificiationID: api.newFleetEPMGetPackageSignatureVerificationID(),
			GetPackagesInstalled:               api.newFleetEPMGetInstalledPackages(),
			GetPackagesLimited:                 api.newFleetEPMGetPackagesLimited(),
			GetPackageStats:                    api.newFleetEPMGetPackageStats(),
			InstallPackageRegistry:             api.newFleetEPMInstallPackageRegistry(),
			InstallPackageUpload:               api.newFleetEPMInstallPackageUpload(),
			ListCategories:                     api.newFleetEPMListPkgCategories(),
			ListDataStreams:                    api.newFleetEPMListDataStreams(),
			ListPackages:                       api.newFleetEPMListPackages(),
			UpdatePackageSettings:              api.newFleetEPMUpdatePackageSettings(),
		},
		EnrollmentAPIKeys: EnrollmentAPIKeys{
			Create: api.newFleetEnrollmentAPIKeysCreate(),
			Get:    api.newFleetEnrollmentAPIKeysGet(),
			List:   api.newFleetEnrollmentAPIKeysList(),
			Revoke: api.newFleetEnrollmentAPIKeysRevoke(),
		},
		Internal: Internal{
			CheckFleetServerHealth: api.newFleetInternalCheckFleetServerHealth(),
			CheckPermissions:       api.newFleetInternalCheckPermissions(),
			GetSettings:            api.newFleetInternalGetSettings(),
			InitiateFleetSetup:     api.newFleetInternalInitiateFleetSetup(),
			UpdateSettings:         api.newFleetInternalUpdateSettings(),
		},
		MessageSigningService: MessageSigningService{
			Rotate: api.newFleetMessageSigningServiceRotate(),
		},
		Outputs: Outputs{
			Create:              api.newFleetOutputsCreate(),
			Delete:              api.newFleetOutputsDelete(),
			GenerateLogstashKey: api.newFleetOutputsGenerateLogstashKey(),
			Get:                 api.newFleetOutputsGet(),
			Health:              api.newFleetOutputsHealth(),
			List:                api.newFleetOutputsList(),
			Update:              api.newFleetOutputsUpdate(),
		},
		PackagePolicies: PackagePolicies{
			BulkDelete:    api.newFleetPackagePoliciesBulkDelete(),
			BulkGet:       api.newFleetPackagePoliciesBulkGet(),
			Create:        api.newFleetPackagePoliciesCreate(),
			Delete:        api.newFleetPackagePoliciesDelete(),
			Get:           api.newFleetPackagePoliciesGet(),
			List:          api.newFleetPackagePoliciesList(),
			Update:        api.newFleetPackagePoliciesUpdate(),
			Upgrade:       api.newFleetPackagePoliciesUpgrade(),
			UpgradeDryRun: api.newFleetPackagePoliciesUpgradeDryRun(),
		},
		Proxies: Proxies{
			Create: api.newFleetProxiesCreate(),
			Delete: api.newFleetProxiesDelete(),
			Get:    api.newFleetProxiesGet(),
			List:   api.newFleetProxiesList(),
			Update: api.newFleetProxiesUpdate(),
		},
		ServerHost: ServerHost{
			Create: api.newFleetServerHostCreate(),
			Delete: api.newFleetServerHostDelete(),
			Get:    api.newFleetServerHostGet(),
			List:   api.newFleetServerHostList(),
			Update: api.newFleetServerHostUpdate(),
		},
		ServiceToken: ServiceToken{
			Create: api.newFleetServiceTokenCreate(),
		},
		UninstallTokens: UninstallTokens{
			GetDecrypted: api.newFleetUninstallTokensGetDecrypted(),
			GetMetadata:  api.newFleetUninstallTokensGetMetadata(),
		},
	}

	api.Endpoint = Endpoint{
		Exceptions: Exceptions{
			CreateItem: api.newEndpointExceptionsCreateItem(),
			CreateList: api.newEndpointExceptionsCreateList(),
			Delete:     api.newEndpointExceptionsDeleteItem(),
			Get:        api.newEndpointExceptionsGet(),
			ListItems:  api.newEndpointExceptionsListItems(),
			Update:     api.newEndpointExceptionsUpdate(),
		},
	}

	api.Logstash = Logstash{
		Delete: api.newLogstashDeletePipeline(),
		Get:    api.newLogstashGetPipeline(),
		List:   api.newLogstashListPipelines(),
		Put:    api.newLogstashPutPipeline(),
	}

	api.ML = ML{
		SyncSavedObjects: api.newMLSyncSavedObjects(),
	}

	api.Roles = Roles{
		CreateOrUpdateMulti:  api.newRolesCreateOrUpdateMulti(),
		CreateOrUpdateSingle: api.newRolesCreateUpdateSingleRole(),
		Delete:               api.newRolesDelete(),
		Get:                  api.newRolesGet(),
		List:                 api.newRolesList(),
	}

	api.SavedObjects = SavedObjects{
		Export:        api.newSavedObjectExport(),
		Import:        api.newSavedObjectImport(),
		ResolveImport: api.newSavedObjectResolveImports(),
		RotateKey:     api.newSavedObjectRotateKey(),
	}

	api.SecurityAIAssistant = SecurityAIAssistant{
		BulkActionAnonymizationFields:  api.newSecurityAIAssistantBulkActionAnonymization(),
		BulkActionKnowledgeBaseEntires: api.newSecurityAIAssistantBulkActionKnowledgeBaseEntry(),
		BulkActionPrompts:              api.newSecurityAIAssistantBulkActionPrompts(),
		CreateConversation:             api.newSecurityAIAssistantCreateConversation(),
		CreateKnowledgeBase:            api.newSecurityAIAssistantCreateKnowledgeBase(),
		CreateKnowledgeBaseEntry:       api.newSecurityAIAssistantCreateKnowledgeBaseEntry(),
		CreateModelResponse:            api.newSecurityAIAssistantCreateModelResponse(),
		DeleteConversation:             api.newSecurityAIAssistantDeleteConversation(),
		DeleteKnowledgeBaseEntry:       api.newSecurityAIAssistantDeleteKnowledgeBaseEntry(),
		GetConversation:                api.newSecurityAIAssistantGetConversation(),
		GetKnowledgeBase:               api.newSecurityAIAssistantGetKnowledgeBase(),
		GetKnowledgeBaseEntry:          api.newSecurityAIAssistantGetKnowledgeBaseEntry(),
		ListAnonymizationFields:        api.newSecurityAIAssistantListAnonymization(),
		ListConversations:              api.newSecurityAIAssistantListConversations(),
		ListKnowledgeBaseEntries:       api.newSecurityAIAssistantListKnowledgeBaseEntry(),
		ListPrompts:                    api.newSecurityAIAssistantListPrompts(),
		UpdateKnowledgeBaseEntry:       api.newSecurityAIAssistantUpdateKnowledgeBaseEntry(),
	}

	api.SecurityDetections = SecurityDetections{
		AssignUsers:        api.newSecurityDetectionsAssignUsers(),
		BulkActionRules:    api.newSecurityDetectionsBulkActionRules(),
		CreateIndex:        api.newSecurityDetectionsCreateIndex(),
		CreateRule:         api.newSecurityDetectionsCreateRule(),
		ExportRules:        api.newSecurityDetectionsExportRules(),
		DeleteIndex:        api.newSecurityDetectionsDeleteIndex(),
		DeleteRule:         api.newSecurityDetectionsDeleteRule(),
		GetIndex:           api.newSecurityDetectionsGetIndex(),
		GetRule:            api.newSecurityDetectionsGetRule(),
		GetStatusPrebuilt:  api.newSecurityDetectionsGetStatusPrebuilt(),
		GetPrivilegesSpace: api.newSecurityDetectionsGetPrivilegesSpace(),
		ImportRule:         api.newSecurityDetectionsImportRules(),
		InstallPrebuilt:    api.newSecurityDetectionsInstallPrebuilt(),
		ListRules:          api.newSecurityDetectionsListRules(),
		ListTags:           api.newSecurityDetectionsListTags(),
		PatchRule:          api.newSecurityDetectionsPatchRule(),
		PreviewAlerts:      api.newSecurityDetectionsPreviewAlerts(),
		SearchAlerts:       api.newSecurityDetectionsSearchAlerts(),
		SetAlertStatus:     api.newSecurityDetectionsSetAlertStatus(),
		UpdateRule:         api.newSecurityDetectionsUpdateRule(),
		UpdateTags:         api.newSecurityDetectionsUpdateTags(),
	}

	api.SecurityEndpointManagement = SecurityEndpointManagement{
		GetActionStatus: api.newSecurityEndpointManagementGetActionStatus(),
		ListActions:     api.newSecurityEndpointManagementListActions(),
	}

	api.SecurityExceptions = SecurityExceptions{
		CreateItem:       api.newSecurityExceptionsCreateItem(),
		CreateItems:      api.newSecurityExceptionsCreateItems(),
		CreateList:       api.newSecurityExceptionsCreateList(),
		CreateSharedList: api.newSecurityExceptionsCreateSharedList(),
		DeleteItem:       api.newSecurityExceptionsDeleteItem(),
		DeleteList:       api.newSecurityExceptionsDeleteList(),
		DuplicateList:    api.newSecurityExceptionsDuplicateList(),
		ExportList:       api.newSecurityExceptionsExportList(),
		GetList:          api.newSecurityExceptionsGetList(),
		GetItem:          api.newSecurityExceptionsGetItem(),
		GetSummary:       api.newSecurityExceptionsGetSummary(),
		ImportList:       api.newSecurityExceptionsImportList(),
		ListItems:        api.newSecurityExceptionsListItems(),
		ListLists:        api.newSecurityExceptionsListLists(),
		UpdateItem:       api.newSecurityExceptionsUpdateItem(),
		UpdateList:       api.newSecurityExceptionsUpdateList(),
	}

	api.ShortURL = ShortURL{
		Create:  api.newShortURLCreate(),
		Delete:  api.newShortURLDelete(),
		Get:     api.newShortURLGet(),
		Resolve: api.newShortURLResolve(),
	}

	api.Spaces = Spaces{
		CopyObjects:            api.newSpacesCopyObjects(),
		Create:                 api.newSpacesCreate(),
		Delete:                 api.newSpacesDelete(),
		Get:                    api.newSpacesGet(),
		GetAll:                 api.newSpacesGetAll(),
		GetShareableReferences: api.newSpacesShareableReferences(),
		Update:                 api.newSpacesUpdate(),
		UpdateObjects:          api.newSpacesUpdateObjects(),
	}

	api.Status = Status{
		Get:         api.newStatusFunc(),
		GetRedacted: api.newStatusRedactedFunc(),
	}

	api.TaskManager = TaskManager{
		Health: api.newTaskManagerHealth(),
	}

	api.Uptime = Uptime{
		GetSettings:    api.newUptimeGetSettings(),
		UpdateSettings: api.newUptimeUpdateSettings(),
	}

	return api
}
