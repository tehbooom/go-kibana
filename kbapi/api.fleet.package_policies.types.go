package kbapi

type PackagePolicy struct {
	// AdditionalDatastreamsPermissions Additional datastream permissions, that will be added to the agent policy.
	AdditionalDatastreamsPermissions *[]string `json:"additional_datastreams_permissions,omitempty"`
	Agents                           *float32  `json:"agents,omitempty"`
	CreatedAt                        string    `json:"created_at"`
	CreatedBy                        string    `json:"created_by"`

	// Description Package policy description
	Description   *string                      `json:"description,omitempty"`
	Elasticsearch *PackagePolicy_Elasticsearch `json:"elasticsearch,omitempty"`
	Enabled       bool                         `json:"enabled"`
	ID            string                       `json:"id"`

	// Inputs Package policy inputs (see integration documentation to know what inputs are available)
	Inputs    []PackagePolicyInput `json:"inputs"`
	IsManaged *bool                `json:"is_managed,omitempty"`

	// Name Package policy name (should be unique)
	Name string `json:"name"`

	// Namespace The package policy namespace. Leave blank to inherit the agent policy's namespace.
	Namespace *string `json:"namespace,omitempty"`
	OutputID  *string `json:"output_id"`

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
				TSDB                *bool `json:"tsdb,omitempty"`
			} `json:"features"`
		} `json:"experimental_data_stream_features,omitempty"`

		// Name Package name
		Name         string  `json:"name"`
		RequiresRoot *bool   `json:"requires_root,omitempty"`
		Title        *string `json:"title,omitempty"`

		// Version Package version
		Version string `json:"version"`
	} `json:"package,omitempty"`

	PolicyIDs        *[]string                 `json:"policy_ids,omitempty"`
	Revision         float32                   `json:"revision"`
	SecretReferences *[]PackagePolicySecretRef `json:"secret_references,omitempty"`
	SpaceIDs         *[]string                 `json:"spaceIds,omitempty"`

	// SupportsAgentless Indicates whether the package policy belongs to an agentless agent policy.
	SupportsAgentless *bool                   `json:"supports_agentless"`
	UpdatedAt         string                  `json:"updated_at"`
	UpdatedBy         string                  `json:"updated_by"`
	Vars              *map[string]interface{} `json:"vars,omitempty"`
	Version           *string                 `json:"version,omitempty"`
}

// PackagePolicy_Elasticsearch defines model for PackagePolicy.Elasticsearch.
type PackagePolicy_Elasticsearch struct {
	Privileges           *PackagePolicy_Elasticsearch_Privileges `json:"privileges,omitempty"`
	AdditionalProperties map[string]interface{}                  `json:"-"`
}

// PackagePolicy_Elasticsearch_Privileges defines model for PackagePolicy.Elasticsearch.Privileges.
type PackagePolicy_Elasticsearch_Privileges struct {
	Cluster              *[]string              `json:"cluster,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackagePolicyInput defines model for package_policy_input.
type PackagePolicyInput struct {
	Type           string `json:"type"`
	PolicyTemplate string `json:"policy_template"`
	// Enabled enable or disable that input, (default to true)
	Enabled *bool `json:"enabled,omitempty"`
	// Streams Input streams (see integration documentation to know what streams are available)
	Streams       []PackagePolicyInputStream `json:"streams,omitempty"`
	Vars          *map[string]interface{}    `json:"vars,omitempty"`
	KeepEnabled   bool                       `json:"keep_enabled,omitempty"`
	Config        map[string]interface{}     `json:"config,omitempty"`
	CompiledInput map[string]interface{}     `json:"compiled_input,omitempty"`
}

// PackagePolicyInputStream defines model for package_policy_input_stream.
type PackagePolicyInputStream struct {
	// Enabled enable or disable that stream, (default to true)
	Enabled    bool `json:"enabled"`
	DataStream *struct {
		Type    string `json:"type,omitempty"`
		Dataset string `json:"dataset,omitempty"`
	} `json:"data_stream,omitempty"`
	Vars           map[string]interface{} `json:"vars,omitempty"`
	ID             string                 `json:"id,omitempty"`
	CompiledStream map[string]interface{} `json:"compiled_stream,omitempty"`
}

// PackagePolicySecretRef defines model for package_policy_secret_ref.
type PackagePolicySecretRef struct {
	ID string `json:"id"`
}
