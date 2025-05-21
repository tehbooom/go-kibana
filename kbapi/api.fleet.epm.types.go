package kbapi

import "encoding/json"

type PackageListItem_Conditions struct {
	Elastic              *PackageListItem_Conditions_Elastic `json:"elastic,omitempty"`
	Kibana               *PackageListItem_Conditions_Kibana  `json:"kibana,omitempty"`
	AdditionalProperties map[string]interface{}              `json:"-"`
}

type PackageListItem_Conditions_Elastic struct {
	Capabilities         *[]string              `json:"capabilities,omitempty"`
	Subscription         *string                `json:"subscription,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_Conditions_Kibana struct {
	Version              *string                `json:"version,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_Discovery struct {
	Fields               *[]PackageListItem_Discovery_Fields_Item `json:"fields,omitempty"`
	AdditionalProperties map[string]interface{}                   `json:"-"`
}

type PackageListItem_Discovery_Fields_Item struct {
	Name                 string                 `json:"name"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_Icons_Item struct {
	DarkMode             *bool                  `json:"dark_mode,omitempty"`
	Path                 *string                `json:"path,omitempty"`
	Size                 *string                `json:"size,omitempty"`
	Src                  string                 `json:"src"`
	Title                *string                `json:"title,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo defines model for package_info.
type PackageInfoList struct {
	Categories           *[]string                         `json:"categories,omitempty"`
	Conditions           *PackageListItem_Conditions       `json:"conditions,omitempty"`
	DataStreams          *[]map[string]interface{}         `json:"data_streams,omitempty"`
	Description          *string                           `json:"description,omitempty"`
	Discovery            *PackageListItem_Discovery        `json:"discovery,omitempty"`
	Download             *string                           `json:"download,omitempty"`
	FormatVersion        *string                           `json:"format_version,omitempty"`
	Icons                *[]PackageListItem_Icons_Item     `json:"icons,omitempty"`
	Id                   string                            `json:"id"`
	InstallationInfo     *PackageListItem_InstallationInfo `json:"installationInfo,omitempty"`
	Integration          *string                           `json:"integration,omitempty"`
	Internal             *bool                             `json:"internal,omitempty"`
	LatestVersion        *string                           `json:"latestVersion,omitempty"`
	Name                 string                            `json:"name"`
	Owner                *PackageListItem_Owner            `json:"owner,omitempty"`
	Path                 *string                           `json:"path,omitempty"`
	PolicyTemplates      *[]map[string]interface{}         `json:"policy_templates,omitempty"`
	Readme               *string                           `json:"readme,omitempty"`
	Release              *string                           `json:"release,omitempty"`
	SignaturePath        *string                           `json:"signature_path,omitempty"`
	Source               *PackageListItem_Source           `json:"source,omitempty"`
	Status               *string                           `json:"status,omitempty"`
	Title                string                            `json:"title"`
	Type                 json.RawMessage                   `json:"type,omitempty"`
	Vars                 *[]map[string]interface{}         `json:"vars,omitempty"`
	Version              string                            `json:"version"`
	AdditionalProperties map[string]interface{}            `json:"-"`
}

type PackageListItem_InstallationInfo_InstalledEs_Item struct {
	Deferred             *bool                  `json:"deferred,omitempty"`
	Id                   string                 `json:"id"`
	Type                 string                 `json:"type"`
	Version              *string                `json:"version,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_InstallationInfo_ExperimentalDataStreamFeatures_Features struct {
	DocValueOnlyNumeric  *bool                  `json:"doc_value_only_numeric,omitempty"`
	DocValueOnlyOther    *bool                  `json:"doc_value_only_other,omitempty"`
	SyntheticSource      *bool                  `json:"synthetic_source,omitempty"`
	Tsdb                 *bool                  `json:"tsdb,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_InstallationInfo_ExperimentalDataStreamFeatures_Item struct {
	DataStream           string                                                                   `json:"data_stream"`
	Features             PackageListItem_InstallationInfo_ExperimentalDataStreamFeatures_Features `json:"features"`
	AdditionalProperties map[string]interface{}                                                   `json:"-"`
}

type PackageListItem_InstallationInfo_AdditionalSpacesInstalledKibana_Item struct {
	Id                   string                 `json:"id"`
	OriginId             *string                `json:"originId,omitempty"`
	Type                 json.RawMessage        `json:"type"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_InstallationInfo struct {
	AdditionalSpacesInstalledKibana *map[string][]PackageListItem_InstallationInfo_AdditionalSpacesInstalledKibana_Item `json:"additional_spaces_installed_kibana,omitempty"`
	CreatedAt                       *string                                                                             `json:"created_at,omitempty"`
	ExperimentalDataStreamFeatures  *[]PackageListItem_InstallationInfo_ExperimentalDataStreamFeatures_Item             `json:"experimental_data_stream_features,omitempty"`
	InstallFormatSchemaVersion      *string                                                                             `json:"install_format_schema_version,omitempty"`
	InstallSource                   string                                                                              `json:"install_source"`
	InstallStatus                   string                                                                              `json:"install_status"`
	InstalledEs                     []PackageListItem_InstallationInfo_InstalledEs_Item                                 `json:"installed_es"`
	InstalledKibana                 []PackageListItem_InstallationInfo_InstalledKibana_Item                             `json:"installed_kibana"`
	InstalledKibanaSpaceId          *string                                                                             `json:"installed_kibana_space_id,omitempty"`
	LatestExecutedState             *PackageListItem_InstallationInfo_LatestExecutedState                               `json:"latest_executed_state,omitempty"`
	LatestInstallFailedAttempts     *[]PackageListItem_InstallationInfo_LatestInstallFailedAttempts_Item                `json:"latest_install_failed_attempts,omitempty"`
	Name                            string                                                                              `json:"name"`
	Namespaces                      *[]string                                                                           `json:"namespaces,omitempty"`
	Type                            string                                                                              `json:"type"`
	UpdatedAt                       *string                                                                             `json:"updated_at,omitempty"`
	VerificationKeyId               *string                                                                             `json:"verification_key_id"`
	VerificationStatus              string                                                                              `json:"verification_status"`
	Version                         string                                                                              `json:"version"`
	AdditionalProperties            map[string]interface{}                                                              `json:"-"`
}

type PackageListItem_InstallationInfo_InstalledKibana_Item struct {
	Id                   string                 `json:"id"`
	OriginId             *string                `json:"originId,omitempty"`
	Type                 json.RawMessage        `json:"type"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_InstallationInfo_LatestInstallFailedAttempts_Item struct {
	CreatedAt            string                                                             `json:"created_at"`
	Error                PackageListItem_InstallationInfo_LatestInstallFailedAttempts_Error `json:"error"`
	TargetVersion        string                                                             `json:"target_version"`
	AdditionalProperties map[string]interface{}                                             `json:"-"`
}

type PackageListItem_InstallationInfo_LatestExecutedState struct {
	Error                *string                `json:"error,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	StartedAt            *string                `json:"started_at,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_Owner struct {
	Github               *string                `json:"github,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}
type PackageListItem_Source struct {
	License              string                 `json:"license"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type PackageListItem_InstallationInfo_LatestInstallFailedAttempts_Error struct {
	Message              string                 `json:"message"`
	Name                 string                 `json:"name"`
	Stack                *string                `json:"stack,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo defines model for package_info.
type PackageInfo struct {
	Agent *struct {
		Privileges *struct {
			Root *bool `json:"root,omitempty"`
		} `json:"privileges,omitempty"`
	} `json:"agent,omitempty"`
	AssetTags *[]struct {
		AssetIds   *[]string `json:"asset_ids,omitempty"`
		AssetTypes *[]string `json:"asset_types,omitempty"`
		Text       string    `json:"text"`
	} `json:"asset_tags,omitempty"`
	Assets               map[string]interface{}        `json:"assets"`
	Categories           *[]string                     `json:"categories,omitempty"`
	Conditions           *PackageInfo_Conditions       `json:"conditions,omitempty"`
	DataStreams          *[]map[string]interface{}     `json:"data_streams,omitempty"`
	Description          *string                       `json:"description,omitempty"`
	Discovery            *PackageInfo_Discovery        `json:"discovery,omitempty"`
	Download             *string                       `json:"download,omitempty"`
	Elasticsearch        *map[string]interface{}       `json:"elasticsearch,omitempty"`
	FormatVersion        *string                       `json:"format_version,omitempty"`
	Icons                *[]PackageInfo_Icons_Item     `json:"icons,omitempty"`
	InstallationInfo     *PackageInfo_InstallationInfo `json:"installationInfo,omitempty"`
	Internal             *bool                         `json:"internal,omitempty"`
	KeepPoliciesUpToDate *bool                         `json:"keepPoliciesUpToDate,omitempty"`
	LatestVersion        *string                       `json:"latestVersion,omitempty"`
	License              *string                       `json:"license,omitempty"`
	LicensePath          *string                       `json:"licensePath,omitempty"`
	Name                 string                        `json:"name"`
	Notice               *string                       `json:"notice,omitempty"`
	Owner                *PackageInfo_Owner            `json:"owner,omitempty"`
	Path                 *string                       `json:"path,omitempty"`
	PolicyTemplates      *[]map[string]interface{}     `json:"policy_templates,omitempty"`
	Readme               *string                       `json:"readme,omitempty"`
	Release              *string                       `json:"release,omitempty"`
	Screenshots          *[]struct {
		DarkMode *bool   `json:"dark_mode,omitempty"`
		Path     *string `json:"path,omitempty"`
		Size     *string `json:"size,omitempty"`
		Src      string  `json:"src"`
		Title    *string `json:"title,omitempty"`
		Type     *string `json:"type,omitempty"`
	} `json:"screenshots,omitempty"`
	SignaturePath        *string                   `json:"signature_path,omitempty"`
	Source               *PackageInfo_Source       `json:"source,omitempty"`
	Status               *string                   `json:"status,omitempty"`
	Title                string                    `json:"title"`
	Type                 *string                   `json:"type,omitempty"`
	Vars                 *[]map[string]interface{} `json:"vars,omitempty"`
	Version              string                    `json:"version"`
	AdditionalProperties map[string]interface{}    `json:"-"`
}

// PackageInfo_Conditions defines model for PackageInfo.Conditions.
type PackageInfo_Conditions struct {
	Elastic              *PackageInfo_Conditions_Elastic `json:"elastic,omitempty"`
	Kibana               *PackageInfo_Conditions_Kibana  `json:"kibana,omitempty"`
	AdditionalProperties map[string]interface{}          `json:"-"`
}

// PackageInfo_Conditions_Elastic defines model for PackageInfo.Conditions.Elastic.
type PackageInfo_Conditions_Elastic struct {
	Capabilities         *[]string              `json:"capabilities,omitempty"`
	Subscription         *string                `json:"subscription,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_Conditions_Kibana defines model for PackageInfo.Conditions.Kibana.
type PackageInfo_Conditions_Kibana struct {
	Version              *string                `json:"version,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_Discovery defines model for PackageInfo.Discovery.
type PackageInfo_Discovery struct {
	Fields               *[]PackageInfo_Discovery_Fields_Item `json:"fields,omitempty"`
	AdditionalProperties map[string]interface{}               `json:"-"`
}

// PackageInfo_Discovery_Fields_Item defines model for PackageInfo.Discovery.Fields.Item.
type PackageInfo_Discovery_Fields_Item struct {
	Name                 string                 `json:"name"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_Icons_Item defines model for package_info.icons.Item.
type PackageInfo_Icons_Item struct {
	DarkMode             *bool                  `json:"dark_mode,omitempty"`
	Path                 *string                `json:"path,omitempty"`
	Size                 *string                `json:"size,omitempty"`
	Src                  string                 `json:"src"`
	Title                *string                `json:"title,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo defines model for PackageInfo.InstallationInfo.
type PackageInfo_InstallationInfo struct {
	AdditionalSpacesInstalledKibana *map[string][]PackageInfo_InstallationInfo_AdditionalSpacesInstalledKibana_Item `json:"additional_spaces_installed_kibana,omitempty"`
	CreatedAt                       *string                                                                         `json:"created_at,omitempty"`
	ExperimentalDataStreamFeatures  *[]PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Item             `json:"experimental_data_stream_features,omitempty"`
	InstallFormatSchemaVersion      *string                                                                         `json:"install_format_schema_version,omitempty"`
	InstallSource                   string                                                                          `json:"install_source"`
	InstallStatus                   string                                                                          `json:"install_status"`
	InstalledES                     []PackageInfo_InstallationInfo_InstalledES_Item                                 `json:"installed_es"`
	InstalledKibana                 []PackageInfo_InstallationInfo_InstalledKibana_Item                             `json:"installed_kibana"`
	InstalledKibanaSpaceId          *string                                                                         `json:"installed_kibana_space_id,omitempty"`
	LatestExecutedState             *PackageInfo_InstallationInfo_LatestExecutedState                               `json:"latest_executed_state,omitempty"`
	LatestInstallFailedAttempts     *[]PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Item                `json:"latest_install_failed_attempts,omitempty"`
	Name                            string                                                                          `json:"name"`
	Namespaces                      *[]string                                                                       `json:"namespaces,omitempty"`
	Type                            string                                                                          `json:"type"`
	UpdatedAt                       *string                                                                         `json:"updated_at,omitempty"`
	VerificationKeyID               *string                                                                         `json:"verification_key_id"`
	VerificationStatus              string                                                                          `json:"verification_status"`
	Version                         string                                                                          `json:"version"`
	AdditionalProperties            map[string]interface{}                                                          `json:"-"`
}

// PackageInfo_InstallationInfo_AdditionalSpacesInstalledKibana_Item defines model for PackageInfo.InstallationInfo.AdditionalSpacesInstalledKibana.Item.
type PackageInfo_InstallationInfo_AdditionalSpacesInstalledKibana_Item struct {
	ID                   string                 `json:"id"`
	OriginID             *string                `json:"originId,omitempty"`
	Type                 *string                `json:"type"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Item defines model for PackageInfo.InstallationInfo.ExperimentalDataStreamFeatures.Item.
type PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Item struct {
	DataStream           string                                                               `json:"data_stream"`
	Features             PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Features `json:"features"`
	AdditionalProperties map[string]interface{}                                               `json:"-"`
}

// PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Features defines model for PackageInfo.InstallationInfo.ExperimentalDataStreamFeatures.Features.
type PackageInfo_InstallationInfo_ExperimentalDataStreamFeatures_Features struct {
	DocValueOnlyNumeric  *bool                  `json:"doc_value_only_numeric,omitempty"`
	DocValueOnlyOther    *bool                  `json:"doc_value_only_other,omitempty"`
	SyntheticSource      *bool                  `json:"synthetic_source,omitempty"`
	TSDB                 *bool                  `json:"tsdb,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo_InstalledES_Item defines model for PackageInfo.InstallationInfo.InstalledEs.Item.
type PackageInfo_InstallationInfo_InstalledES_Item struct {
	Deferred             *bool                  `json:"deferred,omitempty"`
	ID                   string                 `json:"id"`
	Type                 string                 `json:"type"`
	Version              *string                `json:"version,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo_InstalledKibana_Item defines model for PackageInfo.InstallationInfo.InstalledKibana.Item.
type PackageInfo_InstallationInfo_InstalledKibana_Item struct {
	ID                   string                 `json:"id"`
	OriginID             *string                `json:"originId,omitempty"`
	Type                 *string                `json:"type"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo_LatestExecutedState defines model for PackageInfo.InstallationInfo.LatestExecutedState.
type PackageInfo_InstallationInfo_LatestExecutedState struct {
	Error                *string                `json:"error,omitempty"`
	Name                 *string                `json:"name,omitempty"`
	StartedAt            *string                `json:"started_at,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Item defines model for PackageInfo.InstallationInfo.LatestInstallFailedAttempts.Item.
type PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Item struct {
	CreatedAt            string                                                         `json:"created_at"`
	Error                PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Error `json:"error"`
	TargetVersion        string                                                         `json:"target_version"`
	AdditionalProperties map[string]interface{}                                         `json:"-"`
}

// PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Error defines model for PackageInfo.InstallationInfo.LatestInstallFailedAttempts.Error.
type PackageInfo_InstallationInfo_LatestInstallFailedAttempts_Error struct {
	Message              string                 `json:"message"`
	Name                 string                 `json:"name"`
	Stack                *string                `json:"stack,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_Owner defines model for PackageInfo.Owner.
type PackageInfo_Owner struct {
	Github               *string                `json:"github,omitempty"`
	Type                 *string                `json:"type,omitempty"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

// PackageInfo_Source defines model for PackageInfo.Source.
type PackageInfo_Source struct {
	License              string                 `json:"license"`
	AdditionalProperties map[string]interface{} `json:"-"`
}

type Package struct {
	ID       string `json:"id"`
	OriginID string `json:"originId"`
	Type     string `json:"type"`
}
