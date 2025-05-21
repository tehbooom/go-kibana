package kbapi

type Role struct {
	Name              string                 `json:"name"`
	Kibana            []KibanaPermission     `json:"kibana"`
	Metadata          Metadata               `json:"metadata"`
	Description       string                 `json:"description"`
	Elasticsearch     ElasticsearchPrivilege `json:"elasticsearch"`
	TransientMetadata TransientMetadata      `json:"transient_metadata"`
}

type KibanaPermission struct {
	Base    []string               `json:"base"`
	Spaces  []string               `json:"spaces"`
	Feature map[string]interface{} `json:"feature"`
}

type Metadata struct {
	Version int `json:"version"`
}

type ElasticsearchPrivilege struct {
	RunAs         []string         `json:"run_as,omitempty"`
	Cluster       []string         `json:"cluster"`
	Indices       []IndexPrivilege `json:"indices"`
	RemoteCluster []RemoteCluster  `json:"remote_cluster,omitempty"`
	RemoteIndices []RemoteIndices  `json:"remote_indices,omitempty"`
}

type RemoteCluster struct {
	Clusters   []string `json:"clusters"`
	Privileges []string `json:"privileges"`
}

type RemoteIndices struct {
	Names                  []string      `json:"names"`
	Clusters               []string      `json:"clusters"`
	Privileges             []string      `json:"privileges"`
	FieldSecurity          FieldSecurity `json:"field_security,omitempty"`
	Query                  string        `json:"query,omitempty"`
	AllowRestrictedIndices bool          `json:"allow_restricted_indices,omitempty"`
}

type IndexPrivilege struct {
	Names                  []string      `json:"names,omitempty"`
	Query                  string        `json:"query,omitempty"`
	Privileges             []string      `json:"privileges"`
	FieldSecurity          FieldSecurity `json:"field_security,omitempty"`
	AllowRestrictedIndices bool          `json:"allow_restricted_indices,omitempty"`
}

type FieldSecurity struct {
	Grant []string `json:"grant"`
}

type TransientMetadata struct {
	Enabled bool `json:"enabled"`
}
