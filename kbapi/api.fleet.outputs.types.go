package kbapi

type FleetOutputsResponseBodyItem struct {
	AllowEdit            *[]string             `json:"allow_edit,omitempty"`
	ID                   string                `json:"id"`
	Name                 string                `json:"name"`
	Type                 string                `json:"type"`
	Hosts                []string              `json:"hosts"`
	IsDefault            bool                  `json:"is_default"`
	IsDefaultMonitoring  bool                  `json:"is_default_monitoring"`
	IsInternal           bool                  `json:"is_internal,omitempty"`
	IsPreconfigured      bool                  `json:"is_preconfigured,omitempty"`
	ConfigYAML           string                `json:"config_yaml,omitempty"`
	CaSHA256             string                `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint string                `json:"ca_trusted_fingerprint,omitempty"`
	SSL                  *SSL                  `json:"ssl,omitempty"`
	ProxyID              *string               `json:"proxy_id,omitempty"`
	Shipper              *NewOutputShipper     `json:"shipper,omitempty"`
	Preset               *string               `json:"preset,omitempty"`
	Secrets              *map[string]SecretRef `json:"secrets,omitempty"`
	Version              *string               `json:"version,omitempty"`
	ClientID             *string               `json:"client_id,omitempty"`
	Compression          *string               `json:"compression,omitempty"`
	CompressionLevel     *interface{}          `json:"compression_level"`
	AuthType             *string               `json:"auth_type,omitempty"`
	ConnectionType       *interface{}          `json:"connection_type"`
	BrokerTimeout        *int                  `json:"broker_timeout,omitempty"`
	Hash                 *struct {
		Hash   *string `json:"hash,omitempty"`
		Random *bool   `json:"random,omitempty"`
	} `json:"hash,omitempty"`
	Headers      []Header      `json:"headers,omitempty"`
	Username     string        `json:"username,omitempty"`
	SASL         *SASL         `json:"sasl,omitempty"`
	Partition    string        `json:"partition,omitempty"`
	ServiceToken *string       `json:"service_token"`
	Random       *RandomConfig `json:"random,omitempty"`
	Topic        string        `json:"topic,omitempty"`
	Timeout      int           `json:"timeout,omitempty"`
	RequiredAcks int           `json:"required_acks,omitempty"`
	Password     interface{}   `json:"password,omitempty"`
}

type SSL struct {
	Certificate            string   `json:"certificate,omitempty"`
	CertificateAuthorities []string `json:"certificate_authorities,omitempty"`
	VerificationMode       string   `json:"verification_mode,omitempty"`
	Key                    string   `json:"key,omitempty"`
}

type SecretRef struct {
	ID string `json:"id"`
}

type SASL struct {
	Mechanism string `json:"mechanism,omitempty"`
}

type RandomConfig struct {
	GroupEvents int `json:"group_events,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KafkaOutput struct {
	AllowEdit *[]string `json:"allow_edit,omitempty"`
	// Values are none, user_pass, ssl, or kerberos.
	AuthType             string   `json:"auth_type"`
	BrokerTimeout        *float32 `json:"broker_timeout,omitempty"`
	CaSHA256             *string  `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string  `json:"ca_trusted_fingerprint,omitempty"`
	ClientID             *string  `json:"client_id,omitempty"`
	Compression          *string  `json:"compression,omitempty"`
	CompressionLevel     *float32 `json:"compression_level,omitempty"`
	ConfigYaml           *string  `json:"config_yaml"`
	// Values are plaintext or encryption
	ConnectionType string `json:"connection_type"`
	Hash           *struct {
		Hash   *string `json:"hash,omitempty"`
		Random *bool   `json:"random,omitempty"`
	} `json:"hash,omitempty"`
	Headers *[]struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"headers,omitempty"`
	// Format "host:port" without protocol.
	Hosts               []string `json:"hosts"`
	ID                  *string  `json:"id,omitempty"`
	IsDefault           *bool    `json:"is_default,omitempty"`
	IsDefaultMonitoring *bool    `json:"is_default_monitoring,omitempty"`
	IsInternal          *bool    `json:"is_internal,omitempty"`
	IsPreconfigured     *bool    `json:"is_preconfigured,omitempty"`
	Key                 *string  `json:"key,omitempty"`
	Name                string   `json:"name"`
	Partition           *string  `json:"partition,omitempty"`
	Password            *string  `json:"password"`
	ProxyID             *string  `json:"proxy_id,omitempty"`
	Random              *struct {
		GroupEvents *float32 `json:"group_events,omitempty"`
	} `json:"random,omitempty"`
	RequiredAcks *int `json:"required_acks,omitempty"`
	RoundRobin   *struct {
		GroupEvents *float32 `json:"group_events,omitempty"`
	} `json:"round_robin,omitempty"`
	SASL *struct {
		Mechanism *string `json:"mechanism,omitempty"`
	} `json:"sasl"`
	Secrets *struct {
		Password *struct {
			ID string `json:"id,omitempty"`
		} `json:"password,omitempty"`
		SSL *struct {
			Key *struct {
				ID string `json:"id,omitempty"`
			} `json:"key"`
		} `json:"ssl,omitempty"`
	} `json:"secrets,omitempty"`
	Shipper  *NewOutputShipper `json:"shipper,omitempty"`
	SSL      *NewOutputSSL     `json:"ssl,omitempty"`
	Timeout  *float32          `json:"timeout,omitempty"`
	Topic    *string           `json:"topic,omitempty"`
	Type     string            `json:"type"`
	Username *string           `json:"username"`
	Version  *string           `json:"version,omitempty"`
}

type NewOutputShipper struct {
	CompressionLevel            *float32 `json:"compression_level"`
	DiskQueueCompressionEnabled *bool    `json:"disk_queue_compression_enabled"`
	DiskQueueEnabled            *bool    `json:"disk_queue_enabled"`
	DiskQueueEncryptionEnabled  *bool    `json:"disk_queue_encryption_enabled"`
	DiskQueueMaxSize            *float32 `json:"disk_queue_max_size"`
	DiskQueuePath               *string  `json:"disk_queue_path"`
	Loadbalance                 *bool    `json:"loadbalance"`
	MaxBatchBytes               *float32 `json:"max_batch_bytes"`
	MemQueueEvents              *float32 `json:"mem_queue_events"`
	QueueFlushTimeout           *float32 `json:"queue_flush_timeout"`
}

// NewOutputSsl defines model for new_output_ssl.
type NewOutputSSL struct {
	Certificate            *string   `json:"certificate,omitempty"`
	CertificateAuthorities *[]string `json:"certificate_authorities,omitempty"`
	Key                    *string   `json:"key,omitempty"`
	VerificationMode       *string   `json:"verification_mode,omitempty"`
}

type LogstashOutput struct {
	AllowEdit            *[]string `json:"allow_edit,omitempty"`
	CaSHA256             *string   `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string   `json:"ca_trusted_fingerprint,omitempty"`
	ConfigYaml           string    `json:"config_yaml"`
	Hosts                []string  `json:"hosts"`
	ID                   *string   `json:"id,omitempty"`
	IsDefault            *bool     `json:"is_default,omitempty"`
	IsDefaultMonitoring  *bool     `json:"is_default_monitoring,omitempty"`
	IsInternal           *bool     `json:"is_internal,omitempty"`
	IsPreconfigured      *bool     `json:"is_preconfigured,omitempty"`
	Name                 string    `json:"name"`
	ProxyID              *string   `json:"proxy_id,omitempty"`
	Secrets              *struct {
		SSL *struct {
			Key *struct {
				ID string `json:"id,omitempty"`
			} `json:"key,omitempty"`
		} `json:"ssl,omitempty"`
	} `json:"secrets,omitempty"`
	Shipper *NewOutputShipper `json:"shipper,omitempty"`
	SSL     *NewOutputSSL     `json:"ssl,omitempty"`
	Type    string            `json:"type"`
}

type RemoteElasticsearchOutput struct {
	AllowEdit            *[]string `json:"allow_edit,omitempty"`
	CaSHA256             *string   `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string   `json:"ca_trusted_fingerprint,omitempty"`
	ConfigYaml           *string   `json:"config_yaml"`
	Hosts                []string  `json:"hosts"`
	ID                   *string   `json:"id,omitempty"`
	IsDefault            *bool     `json:"is_default,omitempty"`
	IsDefaultMonitoring  *bool     `json:"is_default_monitoring,omitempty"`
	IsInternal           *bool     `json:"is_internal,omitempty"`
	IsPreconfigured      *bool     `json:"is_preconfigured,omitempty"`
	KibanaAPIKey         *string   `json:"kibana_api_key,omitempty"`
	KibanaURL            *string   `json:"kibana_url,omitempty"`
	Name                 string    `json:"name"`
	Preset               *string   `json:"preset,omitempty"`
	ProxyID              *string   `json:"proxy_id,omitempty"`
	Secrets              *struct {
		KibanaApiKey *struct {
			ID string `json:"id,omitempty"`
		} `json:"kibana_api_key,omitempty"`
		ServiceToken *struct {
			ID string `json:"id,omitempty"`
		} `json:"service_token,omitempty"`
		Ssl *struct {
			Key *struct {
				ID string `json:"id,omitempty"`
			} `json:"key,omitempty"`
		} `json:"ssl,omitempty"`
	} `json:"secrets,omitempty"`
	ServiceToken     *string           `json:"service_token"`
	Shipper          *NewOutputShipper `json:"shipper,omitempty"`
	Ssl              *NewOutputSSL     `json:"ssl,omitempty"`
	SyncIntegrations *bool             `json:"sync_integrations,omitempty"`
	Type             string            `json:"type"`
}

type ElasticsearchOutput struct {
	AllowEdit            *[]string `json:"allow_edit,omitempty"`
	CaSHA256             *string   `json:"ca_sha256,omitempty"`
	CaTrustedFingerprint *string   `json:"ca_trusted_fingerprint,omitempty"`
	ConfigYaml           *string   `json:"config_yaml"`
	Hosts                []string  `json:"hosts"`
	ID                   *string   `json:"id,omitempty"`
	IsDefault            *bool     `json:"is_default,omitempty"`
	IsDefaultMonitoring  *bool     `json:"is_default_monitoring,omitempty"`
	IsInternal           *bool     `json:"is_internal,omitempty"`
	IsPreconfigured      *bool     `json:"is_preconfigured,omitempty"`
	Name                 string    `json:"name"`
	Preset               *string   `json:"preset,omitempty"`
	ProxyID              *string   `json:"proxy_id,omitempty"`
	Secrets              *struct {
		SSL *struct {
			Key *struct {
				ID string `json:"id,omitempty"`
			} `json:"key,omitempty"`
		} `json:"ssl,omitempty"`
	} `json:"secrets,omitempty"`
	Shipper *NewOutputShipper `json:"shipper,omitempty"`
	SSL     *NewOutputSSL     `json:"ssl,omitempty"`
	Type    string            `json:"type"`
}
