package kbapi

import "encoding/json"

// AgentPolicy defines model for agent_policy.
type AgentPolicy struct {
	AdvancedSettings *struct {
		AgentDownloadTargetDirectory      *interface{} `json:"agent_download_target_directory"`
		AgentDownloadTimeout              *interface{} `json:"agent_download_timeout"`
		AgentLimitsGoMaxProcs             *interface{} `json:"agent_limits_go_max_procs"`
		AgentLoggingFilesInterval         *interface{} `json:"agent_logging_files_interval"`
		AgentLoggingFilesKeepfiles        *interface{} `json:"agent_logging_files_keepfiles"`
		AgentLoggingFilesRotateeverybytes *interface{} `json:"agent_logging_files_rotateeverybytes"`
		AgentLoggingLevel                 *interface{} `json:"agent_logging_level"`
		AgentLoggingMetricsPeriod         *interface{} `json:"agent_logging_metrics_period"`
		AgentLoggingToFiles               *interface{} `json:"agent_logging_to_files"`
	} `json:"advanced_settings,omitempty"`
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	Agentless *struct {
		Resources *struct {
			Requests *struct {
				Cpu    *string `json:"cpu,omitempty"`
				Memory *string `json:"memory,omitempty"`
			} `json:"requests,omitempty"`
		} `json:"resources,omitempty"`
	} `json:"agentless,omitempty"`
	Agents            *float32 `json:"agents,omitempty"`
	DataOutputId      *string  `json:"data_output_id"`
	Description       *string  `json:"description,omitempty"`
	DownloadSourceId  *string  `json:"download_source_id"`
	FleetServerHostId *string  `json:"fleet_server_host_id"`

	// GlobalDataTags User defined data tags that are added to all of the inputs. The values can be strings or numbers.
	GlobalDataTags *[]struct {
		Name string `json:"name"`
		//TODO Fix json.RawMessage union
		Value json.RawMessage `json:"value"`
	} `json:"global_data_tags,omitempty"`
	HasFleetServer       *bool    `json:"has_fleet_server,omitempty"`
	Id                   string   `json:"id"`
	InactivityTimeout    *float32 `json:"inactivity_timeout,omitempty"`
	IsDefault            *bool    `json:"is_default,omitempty"`
	IsDefaultFleetServer *bool    `json:"is_default_fleet_server,omitempty"`
	IsManaged            bool     `json:"is_managed"`
	IsPreconfigured      *bool    `json:"is_preconfigured,omitempty"`

	// IsProtected Indicates whether the agent policy has tamper protection enabled. Default false.
	IsProtected bool `json:"is_protected"`

	// KeepMonitoringAlive When set to true, monitoring will be enabled but logs/metrics collection will be disabled
	KeepMonitoringAlive   *bool `json:"keep_monitoring_alive"`
	MonitoringDiagnostics *struct {
		Limit *struct {
			Burst    *float32 `json:"burst,omitempty"`
			Interval *string  `json:"interval,omitempty"`
		} `json:"limit,omitempty"`
		Uploader *struct {
			InitDur    *string  `json:"init_dur,omitempty"`
			MaxDur     *string  `json:"max_dur,omitempty"`
			MaxRetries *float32 `json:"max_retries,omitempty"`
		} `json:"uploader,omitempty"`
	} `json:"monitoring_diagnostics,omitempty"`
	MonitoringEnabled *[]string `json:"monitoring_enabled,omitempty"`
	MonitoringHttp    *struct {
		Buffer *struct {
			Enabled *bool `json:"enabled,omitempty"`
		} `json:"buffer,omitempty"`
		Enabled *bool    `json:"enabled,omitempty"`
		Host    *string  `json:"host,omitempty"`
		Port    *float32 `json:"port,omitempty"`
	} `json:"monitoring_http,omitempty"`
	MonitoringOutputId     *string `json:"monitoring_output_id"`
	MonitoringPprofEnabled *bool   `json:"monitoring_pprof_enabled,omitempty"`
	Name                   string  `json:"name"`
	Namespace              string  `json:"namespace"`

	// Overrides Override settings that are defined in the agent policy. Input settings cannot be overridden. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides        *map[string]interface{}     `json:"overrides"`
	PackagePolicies  *AgentPolicyPackagePolicies `json:"package_policies,omitempty"`
	RequiredVersions *[]struct {
		// Percentage Target percentage of agents to auto upgrade
		Percentage float32 `json:"percentage"`

		// Version Target version for automatic agent upgrade
		Version string `json:"version"`
	} `json:"required_versions"`
	Revision      float32   `json:"revision"`
	SchemaVersion *string   `json:"schema_version,omitempty"`
	SpaceIds      *[]string `json:"space_ids,omitempty"`
	Status        string    `json:"status"`

	// SupportsAgentless Indicates whether the agent policy supports agentless integrations.
	SupportsAgentless  *bool    `json:"supports_agentless"`
	UnenrollTimeout    *float32 `json:"unenroll_timeout,omitempty"`
	UnprivilegedAgents *float32 `json:"unprivileged_agents,omitempty"`
	UpdatedAt          string   `json:"updated_at"`
	UpdatedBy          string   `json:"updated_by"`
	Version            *string  `json:"version,omitempty"`
}

// AgentPolicyPackagePolicies This field is present only when retrieving a single agent policy, or when retrieving a list of agent policies with the ?full=true parameter
type AgentPolicyPackagePolicies = []struct {
	// AdditionalDatastreamsPermissions Additional datastream permissions, that will be added to the agent policy.
	AdditionalDatastreamsPermissions *[]string `json:"additional_datastreams_permissions"`
	Agents                           *float32  `json:"agents,omitempty"`
	CreatedAt                        string    `json:"created_at"`
	CreatedBy                        string    `json:"created_by"`

	// Description Package policy description
	Description   *string                                  `json:"description,omitempty"`
	Elasticsearch *AgentPolicyPackagePoliciesElasticsearch `json:"elasticsearch,omitempty"`
	Enabled       bool                                     `json:"enabled"`
	Id            string                                   `json:"id"`
	Inputs        AgentPolicyPackagePoliciesInputs         `json:"inputs"`
	IsManaged     *bool                                    `json:"is_managed,omitempty"`

	// Name Package policy name (should be unique)
	Name string `json:"name"`

	// Namespace The package policy namespace. Leave blank to inherit the agent policy's namespace.
	Namespace *string `json:"namespace,omitempty"`
	OutputId  *string `json:"output_id"`

	// Overrides Override settings that are defined in the package policy. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides *struct {
		Inputs *map[string]interface{} `json:"inputs,omitempty"`
	} `json:"overrides"`
	Package *struct {
		ExperimentalDataStreamFeatures *[]struct {
			DataStream string `json:"data_stream"`
			Features   struct {
				DocValueOnlyNumeric *bool `json:"doc_value_only_numeric,omitempty"`
				DocValueOnlyOther   *bool `json:"doc_value_only_other,omitempty"`
				SyntheticSource     *bool `json:"synthetic_source,omitempty"`
				Tsdb                *bool `json:"tsdb,omitempty"`
			} `json:"features"`
		} `json:"experimental_data_stream_features,omitempty"`

		// Name Package name
		Name         string  `json:"name"`
		RequiresRoot *bool   `json:"requires_root,omitempty"`
		Title        *string `json:"title,omitempty"`

		// Version Package version
		Version string `json:"version"`
	} `json:"package,omitempty"`

	// PolicyId Agent policy ID where that package policy will be added
	// Deprecated:
	PolicyId         *string   `json:"policy_id"`
	PolicyIds        *[]string `json:"policy_ids,omitempty"`
	Revision         float32   `json:"revision"`
	SecretReferences *[]struct {
		Id string `json:"id"`
	} `json:"secret_references,omitempty"`
	SpaceIds *[]string `json:"spaceIds,omitempty"`

	// SupportsAgentless Indicates whether the package policy belongs to an agentless agent policy.
	SupportsAgentless *bool  `json:"supports_agentless"`
	UpdatedAt         string `json:"updated_at"`
	UpdatedBy         string `json:"updated_by"`
	// TODO: Fix the union message
	Vars    *json.RawMessage `json:"vars,omitempty"`
	Version *string          `json:"version,omitempty"`
}

// AgentPolicy_PackagePolicies_1_Elasticsearch_Privileges defines model for AgentPolicy.PackagePolicies.1.Elasticsearch.Privileges.
type AgentPolicyPackagePoliciesElasticsearchPrivileges struct {
	Cluster              *[]string              `json:"cluster,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// AgentPolicy_PackagePolicies_1_Elasticsearch defines model for AgentPolicy.PackagePolicies.1.Elasticsearch.
type AgentPolicyPackagePoliciesElasticsearch struct {
	Privileges           *AgentPolicyPackagePoliciesElasticsearchPrivileges `json:"privileges,omitempty"`
	AdditionalProperties map[string]interface{}                             `json:"-"`
}

// AgentPolicyPackagePoliciesInputs defines model for .
type AgentPolicyPackagePoliciesInputs = []struct {
	CompiledInput interface{} `json:"compiled_input"`

	// Config Package variable (see integration documentation for more information)
	Config *map[string]struct {
		Frozen *bool       `json:"frozen,omitempty"`
		Type   *string     `json:"type,omitempty"`
		Value  interface{} `json:"value"`
	} `json:"config,omitempty"`
	Enabled        bool    `json:"enabled"`
	Id             *string `json:"id,omitempty"`
	KeepEnabled    *bool   `json:"keep_enabled,omitempty"`
	PolicyTemplate *string `json:"policy_template,omitempty"`
	Streams        []struct {
		CompiledStream interface{} `json:"compiled_stream"`

		// Config Package variable (see integration documentation for more information)
		Config *map[string]struct {
			Frozen *bool       `json:"frozen,omitempty"`
			Type   *string     `json:"type,omitempty"`
			Value  interface{} `json:"value"`
		} `json:"config,omitempty"`
		DataStream struct {
			Dataset       string `json:"dataset"`
			Elasticsearch *struct {
				DynamicDataset   *bool `json:"dynamic_dataset,omitempty"`
				DynamicNamespace *bool `json:"dynamic_namespace,omitempty"`
				Privileges       *struct {
					Indices *[]string `json:"indices,omitempty"`
				} `json:"privileges,omitempty"`
			} `json:"elasticsearch,omitempty"`
			Type string `json:"type"`
		} `json:"data_stream"`
		Enabled     bool    `json:"enabled"`
		Id          *string `json:"id,omitempty"`
		KeepEnabled *bool   `json:"keep_enabled,omitempty"`
		Release     *string `json:"release,omitempty"`

		// Vars Package variable (see integration documentation for more information)
		Vars *map[string]struct {
			Frozen *bool       `json:"frozen,omitempty"`
			Type   *string     `json:"type,omitempty"`
			Value  interface{} `json:"value"`
		} `json:"vars,omitempty"`
	} `json:"streams"`
	Type string `json:"type"`

	// Vars Package variable (see integration documentation for more information)
	Vars *map[string]struct {
		Frozen *bool       `json:"frozen,omitempty"`
		Type   *string     `json:"type,omitempty"`
		Value  interface{} `json:"value"`
	} `json:"vars,omitempty"`
}

type FleetBulkGetAgentPoliciesResponseBody struct {
	Items []struct {
		AdvancedSettings *struct {
			AgentDownloadTargetDirectory      *interface{} `json:"agent_download_target_directory"`
			AgentDownloadTimeout              *interface{} `json:"agent_download_timeout"`
			AgentLimitsGoMaxProcs             *interface{} `json:"agent_limits_go_max_procs"`
			AgentLoggingFilesInterval         *interface{} `json:"agent_logging_files_interval"`
			AgentLoggingFilesKeepfiles        *interface{} `json:"agent_logging_files_keepfiles"`
			AgentLoggingFilesRotateeverybytes *interface{} `json:"agent_logging_files_rotateeverybytes"`
			AgentLoggingLevel                 *interface{} `json:"agent_logging_level"`
			AgentLoggingMetricsPeriod         *interface{} `json:"agent_logging_metrics_period"`
			AgentLoggingToFiles               *interface{} `json:"agent_logging_to_files"`
		} `json:"advanced_settings,omitempty"`
		AgentFeatures *[]struct {
			Enabled bool   `json:"enabled"`
			Name    string `json:"name"`
		} `json:"agent_features,omitempty"`
		Agentless *struct {
			Resources *struct {
				Requests *struct {
					Cpu    *string `json:"cpu,omitempty"`
					Memory *string `json:"memory,omitempty"`
				} `json:"requests,omitempty"`
			} `json:"resources,omitempty"`
		} `json:"agentless,omitempty"`
		Agents            *float32 `json:"agents,omitempty"`
		DataOutputId      *string  `json:"data_output_id"`
		Description       *string  `json:"description,omitempty"`
		DownloadSourceId  *string  `json:"download_source_id"`
		FleetServerHostId *string  `json:"fleet_server_host_id"`

		// GlobalDataTags User defined data tags that are added to all of the inputs. The values can be strings or numbers.
		GlobalDataTags *[]struct {
			Name  string          `json:"name"`
			Value json.RawMessage `json:"value"`
		} `json:"global_data_tags,omitempty"`
		HasFleetServer       *bool    `json:"has_fleet_server,omitempty"`
		Id                   string   `json:"id"`
		InactivityTimeout    *float32 `json:"inactivity_timeout,omitempty"`
		IsDefault            *bool    `json:"is_default,omitempty"`
		IsDefaultFleetServer *bool    `json:"is_default_fleet_server,omitempty"`
		IsManaged            bool     `json:"is_managed"`
		IsPreconfigured      *bool    `json:"is_preconfigured,omitempty"`

		// IsProtected Indicates whether the agent policy has tamper protection enabled. Default false.
		IsProtected bool `json:"is_protected"`

		// KeepMonitoringAlive When set to true, monitoring will be enabled but logs/metrics collection will be disabled
		KeepMonitoringAlive   *bool `json:"keep_monitoring_alive"`
		MonitoringDiagnostics *struct {
			Limit *struct {
				Burst    *float32 `json:"burst,omitempty"`
				Interval *string  `json:"interval,omitempty"`
			} `json:"limit,omitempty"`
			Uploader *struct {
				InitDur    *string  `json:"init_dur,omitempty"`
				MaxDur     *string  `json:"max_dur,omitempty"`
				MaxRetries *float32 `json:"max_retries,omitempty"`
			} `json:"uploader,omitempty"`
		} `json:"monitoring_diagnostics,omitempty"`
		MonitoringEnabled *[]string `json:"monitoring_enabled,omitempty"`
		MonitoringHttp    *struct {
			Buffer *struct {
				Enabled *bool `json:"enabled,omitempty"`
			} `json:"buffer,omitempty"`
			Enabled *bool    `json:"enabled,omitempty"`
			Host    *string  `json:"host,omitempty"`
			Port    *float32 `json:"port,omitempty"`
		} `json:"monitoring_http,omitempty"`
		MonitoringOutputId     *string `json:"monitoring_output_id"`
		MonitoringPprofEnabled *bool   `json:"monitoring_pprof_enabled,omitempty"`
		Name                   string  `json:"name"`
		Namespace              string  `json:"namespace"`

		// Overrides Override settings that are defined in the agent policy. Input settings cannot be overridden. The override option should be used only in unusual circumstances and not as a routine procedure.
		Overrides        *map[string]interface{}     `json:"overrides"`
		PackagePolicies  *AgentPolicyPackagePolicies `json:"package_policies,omitempty"`
		RequiredVersions *[]struct {
			// Percentage Target percentage of agents to auto upgrade
			Percentage float32 `json:"percentage"`

			// Version Target version for automatic agent upgrade
			Version string `json:"version"`
		} `json:"required_versions"`
		Revision      float32   `json:"revision"`
		SchemaVersion *string   `json:"schema_version,omitempty"`
		SpaceIds      *[]string `json:"space_ids,omitempty"`
		Status        string    `json:"status"`

		// SupportsAgentless Indicates whether the agent policy supports agentless integrations.
		SupportsAgentless  *bool    `json:"supports_agentless"`
		UnenrollTimeout    *float32 `json:"unenroll_timeout,omitempty"`
		UnprivilegedAgents *float32 `json:"unprivileged_agents,omitempty"`
		UpdatedAt          string   `json:"updated_at"`
		UpdatedBy          string   `json:"updated_by"`
		Version            *string  `json:"version,omitempty"`
	} `json:"items"`
}

type FleetCreateAgentPolicyRequestBody struct {
	AdvancedSettings *struct {
		AgentDownloadTargetDirectory      *interface{} `json:"agent_download_target_directory"`
		AgentDownloadTimeout              *interface{} `json:"agent_download_timeout"`
		AgentLimitsGoMaxProcs             *interface{} `json:"agent_limits_go_max_procs"`
		AgentLoggingFilesInterval         *interface{} `json:"agent_logging_files_interval"`
		AgentLoggingFilesKeepfiles        *interface{} `json:"agent_logging_files_keepfiles"`
		AgentLoggingFilesRotateeverybytes *interface{} `json:"agent_logging_files_rotateeverybytes"`
		AgentLoggingLevel                 *interface{} `json:"agent_logging_level"`
		AgentLoggingMetricsPeriod         *interface{} `json:"agent_logging_metrics_period"`
		AgentLoggingToFiles               *interface{} `json:"agent_logging_to_files"`
	} `json:"advanced_settings,omitempty"`
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	Agentless *struct {
		Resources *struct {
			Requests *struct {
				Cpu    *string `json:"cpu,omitempty"`
				Memory *string `json:"memory,omitempty"`
			} `json:"requests,omitempty"`
		} `json:"resources,omitempty"`
	} `json:"agentless,omitempty"`
	DataOutputId      *string `json:"data_output_id"`
	Description       *string `json:"description,omitempty"`
	DownloadSourceId  *string `json:"download_source_id"`
	FleetServerHostId *string `json:"fleet_server_host_id"`
	Force             *bool   `json:"force,omitempty"`

	// GlobalDataTags User defined data tags that are added to all of the inputs. The values can be strings or numbers.
	GlobalDataTags *[]struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
	} `json:"global_data_tags,omitempty"`
	HasFleetServer       *bool    `json:"has_fleet_server,omitempty"`
	Id                   *string  `json:"id,omitempty"`
	InactivityTimeout    *float32 `json:"inactivity_timeout,omitempty"`
	IsDefault            *bool    `json:"is_default,omitempty"`
	IsDefaultFleetServer *bool    `json:"is_default_fleet_server,omitempty"`
	IsManaged            *bool    `json:"is_managed,omitempty"`
	IsProtected          *bool    `json:"is_protected,omitempty"`

	// KeepMonitoringAlive When set to true, monitoring will be enabled but logs/metrics collection will be disabled
	KeepMonitoringAlive   *bool `json:"keep_monitoring_alive,omitempty"`
	MonitoringDiagnostics *struct {
		Limit *struct {
			Burst    *float32 `json:"burst,omitempty"`
			Interval *string  `json:"interval,omitempty"`
		} `json:"limit,omitempty"`
		Uploader *struct {
			InitDur    *string  `json:"init_dur,omitempty"`
			MaxDur     *string  `json:"max_dur,omitempty"`
			MaxRetries *float32 `json:"max_retries,omitempty"`
		} `json:"uploader,omitempty"`
	} `json:"monitoring_diagnostics,omitempty"`
	MonitoringEnabled *[]string `json:"monitoring_enabled,omitempty"`
	MonitoringHttp    *struct {
		Buffer *struct {
			Enabled *bool `json:"enabled,omitempty"`
		} `json:"buffer,omitempty"`
		Enabled *bool    `json:"enabled,omitempty"`
		Host    *string  `json:"host,omitempty"`
		Port    *float32 `json:"port,omitempty"`
	} `json:"monitoring_http,omitempty"`
	MonitoringOutputId     *string `json:"monitoring_output_id"`
	MonitoringPprofEnabled *bool   `json:"monitoring_pprof_enabled,omitempty"`
	Name                   string  `json:"name"`
	Namespace              string  `json:"namespace"`

	// Overrides Override settings that are defined in the agent policy. Input settings cannot be overridden. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides        *map[string]interface{} `json:"overrides,omitempty"`
	RequiredVersions *[]struct {
		// Percentage Target percentage of agents to auto upgrade
		Percentage float32 `json:"percentage"`

		// Version Target version for automatic agent upgrade
		Version string `json:"version"`
	} `json:"required_versions,omitempty"`
	SpaceIds *[]string `json:"space_ids,omitempty"`

	// SupportsAgentless Indicates whether the agent policy supports agentless integrations.
	SupportsAgentless *bool    `json:"supports_agentless,omitempty"`
	UnenrollTimeout   *float32 `json:"unenroll_timeout,omitempty"`
}

type PutFleetAgentPolicyRequestBody struct {
	AdvancedSettings *struct {
		AgentDownloadTargetDirectory      *interface{} `json:"agent_download_target_directory"`
		AgentDownloadTimeout              *interface{} `json:"agent_download_timeout"`
		AgentLimitsGoMaxProcs             *interface{} `json:"agent_limits_go_max_procs"`
		AgentLoggingFilesInterval         *interface{} `json:"agent_logging_files_interval"`
		AgentLoggingFilesKeepfiles        *interface{} `json:"agent_logging_files_keepfiles"`
		AgentLoggingFilesRotateeverybytes *interface{} `json:"agent_logging_files_rotateeverybytes"`
		AgentLoggingLevel                 *interface{} `json:"agent_logging_level"`
		AgentLoggingMetricsPeriod         *interface{} `json:"agent_logging_metrics_period"`
		AgentLoggingToFiles               *interface{} `json:"agent_logging_to_files"`
	} `json:"advanced_settings,omitempty"`
	AgentFeatures *[]struct {
		Enabled bool   `json:"enabled"`
		Name    string `json:"name"`
	} `json:"agent_features,omitempty"`
	Agentless *struct {
		Resources *struct {
			Requests *struct {
				Cpu    *string `json:"cpu,omitempty"`
				Memory *string `json:"memory,omitempty"`
			} `json:"requests,omitempty"`
		} `json:"resources,omitempty"`
	} `json:"agentless,omitempty"`
	BumpRevision      *bool   `json:"bumpRevision,omitempty"`
	DataOutputId      *string `json:"data_output_id"`
	Description       *string `json:"description,omitempty"`
	DownloadSourceId  *string `json:"download_source_id"`
	FleetServerHostId *string `json:"fleet_server_host_id"`
	Force             *bool   `json:"force,omitempty"`

	// GlobalDataTags User defined data tags that are added to all of the inputs. The values can be strings or numbers.
	GlobalDataTags *[]struct {
		Name string `json:"name"`
		// TODO fix json.RawMessage union
		Value json.RawMessage `json:"value"`
	} `json:"global_data_tags,omitempty"`
	HasFleetServer       *bool    `json:"has_fleet_server,omitempty"`
	Id                   *string  `json:"id,omitempty"`
	InactivityTimeout    *float32 `json:"inactivity_timeout,omitempty"`
	IsDefault            *bool    `json:"is_default,omitempty"`
	IsDefaultFleetServer *bool    `json:"is_default_fleet_server,omitempty"`
	IsManaged            *bool    `json:"is_managed,omitempty"`
	IsProtected          *bool    `json:"is_protected,omitempty"`

	// KeepMonitoringAlive When set to true, monitoring will be enabled but logs/metrics collection will be disabled
	KeepMonitoringAlive   *bool `json:"keep_monitoring_alive,omitempty"`
	MonitoringDiagnostics *struct {
		Limit *struct {
			Burst    *float32 `json:"burst,omitempty"`
			Interval *string  `json:"interval,omitempty"`
		} `json:"limit,omitempty"`
		Uploader *struct {
			InitDur    *string  `json:"init_dur,omitempty"`
			MaxDur     *string  `json:"max_dur,omitempty"`
			MaxRetries *float32 `json:"max_retries,omitempty"`
		} `json:"uploader,omitempty"`
	} `json:"monitoring_diagnostics,omitempty"`
	MonitoringEnabled *[]string `json:"monitoring_enabled,omitempty"`
	MonitoringHttp    *struct {
		Buffer *struct {
			Enabled *bool `json:"enabled,omitempty"`
		} `json:"buffer,omitempty"`
		Enabled *bool    `json:"enabled,omitempty"`
		Host    *string  `json:"host,omitempty"`
		Port    *float32 `json:"port,omitempty"`
	} `json:"monitoring_http,omitempty"`
	MonitoringOutputId     *string `json:"monitoring_output_id"`
	MonitoringPprofEnabled *bool   `json:"monitoring_pprof_enabled,omitempty"`
	Name                   string  `json:"name"`
	Namespace              string  `json:"namespace"`

	// Overrides Override settings that are defined in the agent policy. Input settings cannot be overridden. The override option should be used only in unusual circumstances and not as a routine procedure.
	Overrides        *map[string]interface{} `json:"overrides,omitempty"`
	RequiredVersions *[]struct {
		// Percentage Target percentage of agents to auto upgrade
		Percentage float32 `json:"percentage"`

		// Version Target version for automatic agent upgrade
		Version string `json:"version"`
	} `json:"required_versions,omitempty"`
	SpaceIds *[]string `json:"space_ids,omitempty"`

	// SupportsAgentless Indicates whether the agent policy supports agentless integrations.
	SupportsAgentless *bool    `json:"supports_agentless,omitempty"`
	UnenrollTimeout   *float32 `json:"unenroll_timeout,omitempty"`
}
