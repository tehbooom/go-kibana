package kbapi

import "encoding/json"

// BedrockConfig Defines properties for connectors when type is `.bedrock`.
type BedrockConfig struct {
	// APIURL The Amazon Bedrock request URL.
	APIURL string `json:"apiURL"`

	// DefaultModel The generative artificial intelligence model for Amazon Bedrock to use. Current support is for the Anthropic Claude models.
	DefaultModel *string `json:"defaultModel,omitempty"`
}

// BedrockSecrets Defines secrets for connectors when type is `.bedrock`.
type BedrockSecrets struct {
	// AccessKey The AWS access key for authentication.
	AccessKey string `json:"accessKey"`

	// Secret The AWS secret for authentication.
	Secret string `json:"secret"`
}

// CasesWebhookConfig Defines properties for connectors when type is `.cases-webhook`.
type CasesWebhookConfig struct {
	// AuthType The type of authentication to use: basic, SSL, or none.
	AuthType *string `json:"authType"`

	// CA A base64 encoded version of the certificate authority file that the connector can trust to sign and validate certificates. This option is available for all authentication types.
	CA *string `json:"ca,omitempty"`

	// CertType If the `authType` is `webhook-authentication-ssl`, specifies whether the certificate authentication data is in a CRT and key file format or a PFX file format.
	CertType *string `json:"certType,omitempty"`

	// CreateCommentJSON A JSON payload sent to the create comment URL to create a case comment. You can use variables to add Kibana Cases data to the payload. The required variable is `case.comment`. Due to Mustache template variables (the text enclosed in triple braces, for example, `{{{case.title}}}`), the JSON is not validated when you create the connector. The JSON is validated once the Mustache variables have been placed when the REST method runs. Manually ensure that the JSON is valid, disregarding the Mustache variables, so the later validation will pass.
	CreateCommentJSON *string `json:"createCommentJSON,omitempty"`

	// CreateCommentMethod The REST API HTTP request method to create a case comment in the third-party system. Valid values are `patch`, `post`, and `put`.
	CreateCommentMethod *string `json:"createCommentMethod,omitempty"`

	// CreateCommentURL The REST API URL to create a case comment by ID in the third-party system. You can use a variable to add the external system ID to the URL. If you are using the `xpack.actions.allowedHosts setting`, add the hostname to the allowed hosts.
	CreateCommentURL *string `json:"createCommentURL,omitempty"`

	// CreateIncidentJSON A JSON payload sent to the create case URL to create a case. You can use variables to add case data to the payload. Required variables are `case.title` and `case.description`. Due to Mustache template variables (which is the text enclosed in triple braces, for example, `{{{case.title}}}`), the JSON is not validated when you create the connector. The JSON is validated after the Mustache variables have been placed when REST method runs. Manually ensure that the JSON is valid to avoid future validation errors; disregard Mustache variables during your review.
	CreateIncidentJSON string `json:"createIncidentJSON"`

	// CreateIncidentMethod The REST API HTTP request method to create a case in the third-party system. Valid values are `patch`, `post`, and `put`.
	CreateIncidentMethod *string `json:"createIncidentMethod,omitempty"`

	// CreateIncidentResponseKey The JSON key in the create external case response that contains the case ID.
	CreateIncidentResponseKey string `json:"createIncidentResponseKey"`

	// CreateIncidentURL The REST API URL to create a case in the third-party system. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	CreateIncidentURL string `json:"createIncidentURL"`

	// GetIncidentResponseExternalTitleKey The JSON key in get external case response that contains the case title.
	GetIncidentResponseExternalTitleKey string `json:"getIncidentResponseExternalTitleKey"`

	// GetIncidentURL The REST API URL to get the case by ID from the third-party system. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts. You can use a variable to add the external system ID to the URL. Due to Mustache template variables (the text enclosed in triple braces, for example, `{{{case.title}}}`), the JSON is not validated when you create the connector. The JSON is validated after the Mustache variables have been placed when REST method runs. Manually ensure that the JSON is valid, disregarding the Mustache variables, so the later validation will pass.
	GetIncidentURL string `json:"getIncidentURL"`

	// HasAuth If true, a username and password for login type authentication must be provided.
	HasAuth *bool `json:"hasAuth,omitempty"`

	// Headers A set of key-value pairs sent as headers with the request URLs for the create case, update case, get case, and create comment methods.
	Headers *string `json:"headers,omitempty"`

	// UpdateIncidentJSON The JSON payload sent to the update case URL to update the case. You can use variables to add Kibana Cases data to the payload. Required variables are `case.title` and `case.description`. Due to Mustache template variables (which is the text enclosed in triple braces, for example, `{{{case.title}}}`), the JSON is not validated when you create the connector. The JSON is validated after the Mustache variables have been placed when REST method runs. Manually ensure that the JSON is valid to avoid future validation errors; disregard Mustache variables during your review.
	UpdateIncidentJSON string `json:"updateIncidentJSON"`

	// UpdateIncidentMethod The REST API HTTP request method to update the case in the third-party system. Valid values are `patch`, `post`, and `put`.
	UpdateIncidentMethod *string `json:"updateIncidentMethod,omitempty"`

	// UpdateIncidentURL The REST API URL to update the case by ID in the third-party system. You can use a variable to add the external system ID to the URL. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	UpdateIncidentURL string `json:"updateIncidentURL"`

	// VerificationMode Controls the verification of certificates. Use `full` to validate that the certificate has an issue date within the `not_before` and `not_after` dates, chains to a trusted certificate authority (CA), and has a hostname or IP address that matches the names within the certificate. Use `certificate` to validate the certificate and verify that it is signed by a trusted authority; this option does not check the certificate hostname. Use `none` to skip certificate validation.
	VerificationMode *string `json:"verificationMode,omitempty"`

	// ViewIncidentURL The URL to view the case in the external system. You can use variables to add the external system ID or external system title to the URL.
	ViewIncidentURL string `json:"viewIncidentURL"`
}

// CasesWebhookSecrets defines model for cases_webhook_secrets.
type CasesWebhookSecrets struct {
	// Crt If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-crt-key`, it is a base64 encoded version of the CRT or CERT file.
	CRT *string `json:"crt,omitempty"`

	// Key If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-crt-key`, it is a base64 encoded version of the KEY file.
	Key *string `json:"key,omitempty"`

	// Password The password for HTTP basic authentication. If `hasAuth` is set to `true` and and `authType` is `webhook-authentication-basic`, this property is required.
	Password *string `json:"password,omitempty"`

	// Pfx If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-pfx`, it is a base64 encoded version of the PFX or P12 file.
	PFX *string `json:"pfx,omitempty"`

	// User The username for HTTP basic authentication. If `hasAuth` is set to `true` and `authType` is `webhook-authentication-basic`, this property is required.
	User *string `json:"user,omitempty"`
}

// CrowdstrikeConfig Defines config properties for connectors when type is `.crowdstrike`.
type CrowdstrikeConfig struct {
	// URL The CrowdStrike tenant URL. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	URL string `json:"url"`
}

// CrowdstrikeSecrets Defines secrets for connectors when type is `.crowdstrike`.
type CrowdstrikeSecrets struct {
	// ClientID The CrowdStrike API client identifier.
	ClientID string `json:"clientID"`

	// ClientSecret The CrowdStrike API client secret to authenticate the `clientID`.
	ClientSecret string `json:"clientSecret"`
}

// D3securityConfig Defines properties for connectors when type is `.d3security`.
type D3securityConfig struct {
	// URL The D3 Security API request URL. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	URL string `json:"url"`
}

// D3securitySecrets Defines secrets for connectors when type is `.d3security`.
type D3securitySecrets struct {
	// Token The D3 Security token.
	Token string `json:"token"`
}

// EmailConfig Defines properties for connectors when type is `.email`.
type EmailConfig struct {
	// ClientID The client identifier, which is a part of OAuth 2.0 client credentials authentication, in GUID format. If `service` is `exchange_server`, this property is required.
	ClientID *string `json:"clientID"`

	// From The from address for all emails sent by the connector. It must be specified in `user@host-name` format.
	From string `json:"from"`

	// HasAuth Specifies whether a user and password are required inside the secrets configuration.
	HasAuth *bool `json:"hasAuth,omitempty"`

	// Host The host name of the service provider. If the `service` is `elastic_cloud` (for Elastic Cloud notifications) or one of Nodemailer's well-known email service providers, this property is ignored. If `service` is `other`, this property must be defined.
	Host          *string `json:"host,omitempty"`
	OauthTokenURL *string `json:"oauthTokenUrl"`

	// Port The port to connect to on the service provider. If the `service` is `elastic_cloud` (for Elastic Cloud notifications) or one of Nodemailer's well-known email service providers, this property is ignored. If `service` is `other`, this property must be defined.
	Port *int `json:"port,omitempty"`

	// Secure Specifies whether the connection to the service provider will use TLS. If the `service` is `elastic_cloud` (for Elastic Cloud notifications) or one of Nodemailer's well-known email service providers, this property is ignored.
	Secure *bool `json:"secure,omitempty"`

	// Service The name of the email service.
	Service *string `json:"service,omitempty"`

	// TenantID The tenant identifier, which is part of OAuth 2.0 client credentials authentication, in GUID format. If `service` is `exchange_server`, this property is required.
	TenantID *string `json:"tenantID"`
}

// EmailSecrets Defines secrets for connectors when type is `.email`.
type EmailSecrets struct {
	// ClientSecret The Microsoft Exchange Client secret for OAuth 2.0 client credentials authentication. It must be URL-encoded. If `service` is `exchange_server`, this property is required.
	ClientSecret *string `json:"clientSecret,omitempty"`

	// Password The password for HTTP basic authentication. If `hasAuth` is set to `true`, this property is required.
	Password *string `json:"password,omitempty"`

	// User The username for HTTP basic authentication. If `hasAuth` is set to `true`, this property is required.
	User *string `json:"user,omitempty"`
}

// GeminiConfig Defines properties for connectors when type is `.gemini`.
type GeminiConfig struct {
	// APIURL The Google Gemini request URL.
	APIURL string `json:"apiURL"`

	// DefaultModel The generative artificial intelligence model for Google Gemini to use.
	DefaultModel *string `json:"defaultModel,omitempty"`

	// GcpProjectID The Google ProjectID that has Vertex AI endpoint enabled.
	GcpProjectID string `json:"gcpProjectID"`

	// GcpRegion The GCP region where the Vertex AI endpoint enabled.
	GcpRegion string `json:"gcpRegion"`
}

// GeminiSecrets Defines secrets for connectors when type is `.gemini`.
type GeminiSecrets struct {
	// CredentialsJSON The service account credentials JSON file. The service account should have Vertex AI user IAM role assigned to it.
	CredentialsJSON string `json:"credentialsJSON"`
}

// GenaiAzureConfig Defines properties for connectors when type is `.gen-ai` and the API provider is `Azure OpenAI`.
type GenaiAzureConfig struct {
	// APIProvider The OpenAI API provider.
	APIProvider string `json:"apiProvider"`

	// APIURL The OpenAI API endpoint.
	APIURL string `json:"apiURL"`
}

// GenaiOpenaiConfig Defines properties for connectors when type is `.gen-ai` and the API provider is `OpenAI`.
type GenaiOpenaiConfig struct {
	// APIProvider The OpenAI API provider.
	APIProvider string `json:"apiProvider"`

	// APIURL The OpenAI API endpoint.
	APIURL string `json:"apiURL"`

	// DefaultModel The default model to use for requests.
	DefaultModel *string `json:"defaultModel,omitempty"`
}

// GenaiSecrets Defines secrets for connectors when type is `.gen-ai`.
type GenaiSecrets struct {
	// APIKey The OpenAI API key.
	APIKey *string `json:"apiKey,omitempty"`
}

// IndexConfig Defines properties for connectors when type is `.index`.
type IndexConfig struct {
	// ExecutionTimeField A field that indicates when the document was indexed.
	ExecutionTimeField *string `json:"executionTimeField"`

	// Index The Elasticsearch index to be written to.
	Index string `json:"index"`

	// Refresh The refresh policy for the write request, which affects when changes are made visible to search. Refer to the refresh setting for Elasticsearch document APIs.
	Refresh *bool `json:"refresh,omitempty"`
}

// JiraConfig Defines properties for connectors when type is `.jira`.
type JiraConfig struct {
	// APIURL The Jira instance URL.
	APIURL string `json:"apiURL"`

	// ProjectKey The Jira project key.
	ProjectKey string `json:"projectKey"`
}

// JiraSecrets Defines secrets for connectors when type is `.jira`.
type JiraSecrets struct {
	// APIToken The Jira API authentication token for HTTP basic authentication.
	APIToken string `json:"apiToken"`

	// Email The account email for HTTP Basic authentication.
	Email string `json:"email"`
}

// Key If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-crt-key`, it is a base64 encoded version of the KEY file.
type Key = string

// OpsgenieConfig Defines properties for connectors when type is `.opsgenie`.
type OpsgenieConfig struct {
	// APIURL The Opsgenie URL. For example, `https://api.opsgenie.com` or `https://api.eu.opsgenie.com`. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	APIURL string `json:"apiURL"`
}

// OpsgenieSecrets Defines secrets for connectors when type is `.opsgenie`.
type OpsgenieSecrets struct {
	// APIKey The Opsgenie API authentication key for HTTP Basic authentication.
	APIKey string `json:"apiKey"`
}

// PagerdutyConfig Defines properties for connectors when type is `.pagerduty`.
type PagerdutyConfig struct {
	// APIURL The PagerDuty event URL.
	APIURL *string `json:"apiURL"`
}

// PagerdutySecrets Defines secrets for connectors when type is `.pagerduty`.
type PagerdutySecrets struct {
	// RoutingKey A 32 character PagerDuty Integration Key for an integration on a service.
	RoutingKey string `json:"routingKey"`
}

// ResilientConfig Defines properties for connectors when type is `.resilient`.
type ResilientConfig struct {
	// APIURL The IBM Resilient instance URL.
	APIURL string `json:"apiURL"`

	// OrgID The IBM Resilient organization ID.
	OrgID string `json:"orgID"`
}

// ResilientSecrets Defines secrets for connectors when type is `.resilient`.
type ResilientSecrets struct {
	// APIKeyID The authentication key ID for HTTP Basic authentication.
	APIKeyID string `json:"apiKeyID"`

	// APIKeySecret The authentication key secret for HTTP Basic authentication.
	APIKeySecret string `json:"apiKeySecret"`
}

// SentineloneConfig Defines properties for connectors when type is `.sentinelone`.
type SentineloneConfig struct {
	// URL The SentinelOne tenant URL. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	URL string `json:"url"`
}

// SentineloneSecrets Defines secrets for connectors when type is `.sentinelone`.
type SentineloneSecrets struct {
	// Token The A SentinelOne API token.
	Token string `json:"token"`
}

// ServicenowConfig Defines properties for connectors when type is `.servicenow`.
type ServicenowConfig struct {
	// APIURL The ServiceNow instance URL.
	APIURL string `json:"apiURL"`

	// ClientID The client ID assigned to your OAuth application. This property is required when `isOAuth` is `true`.
	ClientID *string `json:"clientID,omitempty"`

	// IsOAuth The type of authentication to use. The default value is false, which means basic authentication is used instead of open authorization (OAuth).
	IsOAuth *bool `json:"isOAuth,omitempty"`

	// JWTKeyID The key identifier assigned to the JWT verifier map of your OAuth application. This property is required when `isOAuth` is `true`.
	JWTKeyID *string `json:"jwtKeyID,omitempty"`

	// UserIDentifierValue The identifier to use for OAuth authentication. This identifier should be the user field you selected when you created an OAuth JWT API endpoint for external clients in your ServiceNow instance. For example, if the selected user field is `Email`, the user identifier should be the user's email address. This property is required when `isOAuth` is `true`.
	UserIDentifierValue *string `json:"userIDentifierValue,omitempty"`

	// UsesTableAPI Determines whether the connector uses the Table API or the Import Set API. This property is supported only for ServiceNow ITSM and ServiceNow SecOps connectors.  NOTE: If this property is set to `false`, the Elastic application should be installed in ServiceNow.
	UsesTableAPI *bool `json:"usesTableAPI,omitempty"`
}

// ServicenowItomConfig Defines properties for connectors when type is `.servicenow-itom`.
type ServicenowItomConfig struct {
	// APIURL The ServiceNow instance URL.
	APIURL string `json:"apiURL"`

	// ClientID The client ID assigned to your OAuth application. This property is required when `isOAuth` is `true`.
	ClientID *string `json:"clientID,omitempty"`

	// IsOAuth The type of authentication to use. The default value is false, which means basic authentication is used instead of open authorization (OAuth).
	IsOAuth *bool `json:"isOAuth,omitempty"`

	// JWTKeyID The key identifier assigned to the JWT verifier map of your OAuth application. This property is required when `isOAuth` is `true`.
	JWTKeyID *string `json:"jwtKeyID,omitempty"`

	// UserIDentifierValue The identifier to use for OAuth authentication. This identifier should be the user field you selected when you created an OAuth JWT API endpoint for external clients in your ServiceNow instance. For example, if the selected user field is `Email`, the user identifier should be the user's email address. This property is required when `isOAuth` is `true`.
	UserIDentifierValue *string `json:"userIDentifierValue,omitempty"`
}

// ServicenowSecrets Defines secrets for connectors when type is `.servicenow`, `.servicenow-sir`, or `.servicenow-itom`.
type ServicenowSecrets struct {
	// ClientSecret The client secret assigned to your OAuth application. This property is required when `isOAuth` is `true`.
	ClientSecret *string `json:"clientSecret,omitempty"`

	// Password The password for HTTP basic authentication. This property is required when `isOAuth` is `false`.
	Password *string `json:"password,omitempty"`

	// PrivateKey The RSA private key that you created for use in ServiceNow. This property is required when `isOAuth` is `true`.
	PrivateKey *string `json:"privateKey,omitempty"`

	// PrivateKeyPassword The password for the RSA private key. This property is required when `isOAuth` is `true` and you set a password on your private key.
	PrivateKeyPassword *string `json:"privateKeyPassword,omitempty"`

	// Username The username for HTTP basic authentication. This property is required when `isOAuth` is `false`.
	Username *string `json:"username,omitempty"`
}

// SlackAPIConfig Defines properties for connectors when type is `.slack_api`.
type SlackAPIConfig struct {
	// AllowedChannels A list of valid Slack channels.
	AllowedChannels *[]struct {
		// ID The Slack channel ID.
		ID string `json:"id"`

		// Name The Slack channel name.
		Name string `json:"name"`
	} `json:"allowedChannels,omitempty"`
}

// SlackAPISecrets Defines secrets for connectors when type is `.slack`.
type SlackAPISecrets struct {
	// Token Slack bot user OAuth token.
	Token string `json:"token"`
}

// SwimlaneConfig Defines properties for connectors when type is `.swimlane`.
type SwimlaneConfig struct {
	// APIURL The Swimlane instance URL.
	APIURL string `json:"apiURL"`

	// AppID The Swimlane application ID.
	AppID string `json:"appID"`

	// ConnectorType The type of connector. Valid values are `all`, `alerts`, and `cases`.
	ConnectorType string `json:"connectorType"`

	// Mappings The field mapping.
	Mappings *struct {
		// AlertIDConfig Mapping for the alert ID.
		AlertIDConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"alertIDConfig,omitempty"`

		// CaseIDConfig Mapping for the case ID.
		CaseIDConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"caseIDConfig,omitempty"`

		// CaseNameConfig Mapping for the case name.
		CaseNameConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"caseNameConfig,omitempty"`

		// CommentsConfig Mapping for the case comments.
		CommentsConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"commentsConfig,omitempty"`

		// DescriptionConfig Mapping for the case description.
		DescriptionConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"descriptionConfig,omitempty"`

		// RuleNameConfig Mapping for the name of the alert's rule.
		RuleNameConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"ruleNameConfig,omitempty"`

		// SeverityConfig Mapping for the severity.
		SeverityConfig *struct {
			// FieldType The type of field in Swimlane.
			FieldType string `json:"fieldType"`

			// ID The identifier for the field in Swimlane.
			ID string `json:"id"`

			// Key The key for the field in Swimlane.
			Key string `json:"key"`

			// Name The name of the field in Swimlane.
			Name string `json:"name"`
		} `json:"severityConfig,omitempty"`
	} `json:"mappings,omitempty"`
}

// SwimlaneSecrets Defines secrets for connectors when type is `.swimlane`.
type SwimlaneSecrets struct {
	// APIToken Swimlane API authentication token.
	APIToken *string `json:"apiToken,omitempty"`
}

// TeamsSecrets Defines secrets for connectors when type is `.teams`.
type TeamsSecrets struct {
	// WebhookURL The URL of the incoming webhook. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	WebhookURL string `json:"webhookUrl"`
}

// ThehiveConfig Defines configuration properties for connectors when type is `.thehive`.
type ThehiveConfig struct {
	// Organisation The organisation in TheHive that will contain the alerts or cases. By default, the connector uses the default organisation of the user account that created the API key.
	Organisation *string `json:"organisation,omitempty"`

	// URL The instance URL in TheHive. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	URL string `json:"url"`
}

// ThehiveSecrets Defines secrets for connectors when type is `.thehive`.
type ThehiveSecrets struct {
	// APIKey The API key for authentication in TheHive.
	APIKey string `json:"apiKey"`
}

// TinesConfig Defines properties for connectors when type is `.tines`.
type TinesConfig struct {
	// URL The Tines tenant URL. If you are using the `xpack.actions.allowedHosts` setting, make sure this hostname is added to the allowed hosts.
	URL string `json:"url"`
}

// TinesSecrets Defines secrets for connectors when type is `.tines`.
type TinesSecrets struct {
	// Email The email used to sign in to Tines.
	Email string `json:"email"`

	// Token The Tines API token.
	Token string `json:"token"`
}

// TorqConfig Defines properties for connectors when type is `.torq`.
type TorqConfig struct {
	// WebhookIntegrationURL The endpoint URL of the Elastic Security integration in Torq.
	WebhookIntegrationURL string `json:"webhookIntegrationUrl"`
}

// TorqSecrets Defines secrets for connectors when type is `.torq`.
type TorqSecrets struct {
	// Token The secret of the webhook authentication header.
	Token string `json:"token"`
}

// WebhookConfig Defines properties for connectors when type is `.webhook`.
type WebhookConfig struct {
	// AuthType The type of authentication to use: basic, SSL, or none.
	AuthType *string `json:"authType"`

	// CA A base64 encoded version of the certificate authority file that the connector can trust to sign and validate certificates. This option is available for all authentication types.
	CA *string `json:"ca,omitempty"`

	// CertType If the `authType` is `webhook-authentication-ssl`, specifies whether the certificate authentication data is in a CRT and key file format or a PFX file format.
	CertType *string `json:"certType,omitempty"`

	// HasAuth If true, a username and password for login type authentication must be provided.
	HasAuth *bool `json:"hasAuth,omitempty"`

	// Headers A set of key-value pairs sent as headers with the request.
	Headers *map[string]interface{} `json:"headers"`

	// Method The HTTP request method, either `post` or `put`.
	Method *string `json:"method,omitempty"`

	// URL The request URL. If you are using the `xpack.actions.allowedHosts` setting, add the hostname to the allowed hosts.
	URL *string `json:"url,omitempty"`

	// VerificationMode Controls the verification of certificates. Use `full` to validate that the certificate has an issue date within the `not_before` and `not_after` dates, chains to a trusted certificate authority (CA), and has a hostname or IP address that matches the names within the certificate. Use `certificate` to validate the certificate and verify that it is signed by a trusted authority; this option does not check the certificate hostname. Use `none` to skip certificate validation.
	VerificationMode *string `json:"verificationMode,omitempty"`
}

// WebhookSecrets Defines secrets for connectors when type is `.webhook`.
type WebhookSecrets struct {
	// Crt If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-crt-key`, it is a base64 encoded version of the CRT or CERT file.
	Crt *string `json:"crt,omitempty"`

	// Key If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-crt-key`, it is a base64 encoded version of the KEY file.
	Key *string `json:"key,omitempty"`

	// Password The password for HTTP basic authentication or the passphrase for the SSL certificate files. If `hasAuth` is set to `true` and `authType` is `webhook-authentication-basic`, this property is required.
	Password *string `json:"password,omitempty"`

	// Pfx If `authType` is `webhook-authentication-ssl` and `certType` is `ssl-pfx`, it is a base64 encoded version of the PFX or P12 file.
	Pfx *string `json:"pfx,omitempty"`

	// User The username for HTTP basic authentication. If `hasAuth` is set to `true`  and `authType` is `webhook-authentication-basic`, this property is required.
	User *string `json:"user,omitempty"`
}

// XmattersConfig Defines properties for connectors when type is `.xmatters`.
type XmattersConfig struct {
	// ConfigURL The request URL for the Elastic Alerts trigger in xMatters. It is applicable only when `usesBasic` is `true`.
	ConfigURL *string `json:"configUrl"`

	// UsesBasic Specifies whether the connector uses HTTP basic authentication (`true`) or URL authentication (`false`).
	UsesBasic *bool `json:"usesBasic,omitempty"`
}

// XmattersSecrets Defines secrets for connectors when type is `.xmatters`.
type XmattersSecrets struct {
	// Password A user name for HTTP basic authentication. It is applicable only when `usesBasic` is `true`.
	Password *string `json:"password,omitempty"`

	// SecretsURL The request URL for the Elastic Alerts trigger in xMatters with the API key included in the URL. It is applicable only when `usesBasic` is `false`.
	SecretsURL *string `json:"secretsUrl,omitempty"`

	// User A password for HTTP basic authentication. It is applicable only when `usesBasic` is `true`.
	User *string `json:"user,omitempty"`
}

// RunAcknowledgeResolvePagerduty Test an action that acknowledges or resolves a PagerDuty alert.
type RunAcknowledgeResolvePagerduty struct {
	// DedupKey The deduplication key for the PagerDuty alert.
	DedupKey string `json:"dedupKey"`

	// EventAction The type of event.
	EventAction string `json:"eventAction"`
}

// RunAddevent The `addEvent` subaction for ServiceNow ITOM connectors.
type RunAddevent struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`

	// SubActionParams The set of configuration properties for the action.
	SubActionParams *struct {
		// AdditionalInfo Additional information about the event.
		AdditionalInfo *string `json:"additional_info,omitempty"`

		// Description The details about the event.
		Description *string `json:"description,omitempty"`

		// EventClass A specific instance of the source.
		EventClass *string `json:"event_class,omitempty"`

		// MessageKey All actions sharing this key are associated with the same ServiceNow alert. The default value is `<rule ID>:<alert instance ID>`.
		MessageKey *string `json:"message_key,omitempty"`

		// MetricName The name of the metric.
		MetricName *string `json:"metric_name,omitempty"`

		// Node The host that the event was triggered for.
		Node *string `json:"node,omitempty"`

		// Resource The name of the resource.
		Resource *string `json:"resource,omitempty"`

		// Severity The severity of the event.
		Severity *string `json:"severity,omitempty"`

		// Source The name of the event source type.
		Source *string `json:"source,omitempty"`

		// TimeOfEvent The time of the event.
		TimeOfEvent *string `json:"time_of_event,omitempty"`

		// Type The type of event.
		Type *string `json:"type,omitempty"`
	} `json:"subActionParams,omitempty"`
}

// RunClosealert The `closeAlert` subaction for Opsgenie connectors.
type RunClosealert struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// Alias The unique identifier used for alert deduplication in Opsgenie. The alias must match the value used when creating the alert.
		Alias string `json:"alias"`

		// Note Additional information for the alert.
		Note *string `json:"note,omitempty"`

		// Source The display name for the source of the alert.
		Source *string `json:"source,omitempty"`

		// User The display name for the owner.
		User *string `json:"user,omitempty"`
	} `json:"subActionParams"`
}

// RunCloseincident The `closeIncident` subaction for ServiceNow ITSM connectors.
type RunCloseincident struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		Incident RunCloseincident_SubActionParams_Incident `json:"incident"`
	} `json:"subActionParams"`
}

// RunCloseincident_SubActionParams_Incident defines model for RunCloseincident.SubActionParams.Incident.
type RunCloseincident_SubActionParams_Incident struct {
	// CorrelationId An identifier that is assigned to the incident when it is created by the connector.
	// NOTE: If you use the default value and the rule generates multiple alerts that use the same alert IDs,
	// the latest open incident for this correlation ID is closed unless you specify the external ID.
	CorrelationId *string `json:"correlation_id"`

	// ExternalId The unique identifier (`incidentId`) for the incident in ServiceNow.
	ExternalId *string `json:"externalId"`
}

// RunCreatealert The `createAlert` subaction for Opsgenie and TheHive connectors.
type RunCreatealert struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// Actions The custom actions available to the alert in Opsgenie connectors.
		Actions *[]string `json:"actions,omitempty"`

		// Alias The unique identifier used for alert deduplication in Opsgenie.
		Alias *string `json:"alias,omitempty"`

		// Description A description that provides detailed information about the alert.
		Description *string `json:"description,omitempty"`

		// Details The custom properties of the alert in Opsgenie connectors.
		Details *map[string]interface{} `json:"details,omitempty"`

		// Entity The domain of the alert in Opsgenie connectors. For example, the application or server name.
		Entity *string `json:"entity,omitempty"`

		// Message The alert message in Opsgenie connectors.
		Message *string `json:"message,omitempty"`

		// Note Additional information for the alert in Opsgenie connectors.
		Note *string `json:"note,omitempty"`

		// Priority The priority level for the alert in Opsgenie connectors.
		Priority *string `json:"priority,omitempty"`

		// Responders The entities to receive notifications about the alert in Opsgenie connectors. If `type` is `user`, either `id` or `username` is required. If `type` is `team`, either `id` or `name` is required.
		Responders *[]struct {
			// Id The identifier for the entity.
			Id *string `json:"id,omitempty"`

			// Name The name of the entity.
			Name *string `json:"name,omitempty"`

			// Type The type of responders, in this case `escalation`.
			Type *string `json:"type,omitempty"`

			// Username A valid email address for the user.
			Username *string `json:"username,omitempty"`
		} `json:"responders,omitempty"`

		// Severity The severity of the incident for TheHive connectors. The value ranges from 1 (low) to 4 (critical) with a default value of 2 (medium).
		Severity *int `json:"severity,omitempty"`

		// Source The display name for the source of the alert in Opsgenie and TheHive connectors.
		Source *string `json:"source,omitempty"`

		// SourceRef A source reference for the alert in TheHive connectors.
		SourceRef *string `json:"sourceRef,omitempty"`

		// Tags The tags for the alert in Opsgenie and TheHive connectors.
		Tags *[]string `json:"tags,omitempty"`

		// Title A title for the incident for TheHive connectors. It is used for searching the contents of the knowledge base.
		Title *string `json:"title,omitempty"`

		// Tlp The traffic light protocol designation for the incident in TheHive connectors. Valid values include: 0 (clear), 1 (green), 2 (amber), 3 (amber and strict), and 4 (red).
		Tlp *int `json:"tlp,omitempty"`

		// Type The type of alert in TheHive connectors.
		Type *string `json:"type,omitempty"`

		// User The display name for the owner.
		User *string `json:"user,omitempty"`

		// VisibleTo The teams and users that the alert will be visible to without sending a notification. Only one of `id`, `name`, or `username` is required.
		VisibleTo *[]struct {
			// Id The identifier for the entity.
			Id *string `json:"id,omitempty"`

			// Name The name of the entity.
			Name *string `json:"name,omitempty"`

			// Type Valid values are `team` and `user`.
			Type string `json:"type"`

			// Username The user name. This property is required only when the `type` is `user`.
			Username *string `json:"username,omitempty"`
		} `json:"visibleTo,omitempty"`
	} `json:"subActionParams"`
}

// RunDocuments Test an action that indexes a document into Elasticsearch.
type RunDocuments struct {
	// Documents The documents in JSON format for index connectors.
	Documents []map[string]interface{} `json:"documents"`
}

// RunFieldsbyissuetype The `fieldsByIssueType` subaction for Jira connectors.
type RunFieldsbyissuetype struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// Id The Jira issue type identifier.
		Id string `json:"id"`
	} `json:"subActionParams"`
}

// RunGetchoices The `getChoices` subaction for ServiceNow ITOM, ServiceNow ITSM, and ServiceNow SecOps connectors.
type RunGetchoices struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`

	// SubActionParams The set of configuration properties for the action.
	SubActionParams struct {
		// Fields An array of fields.
		Fields []string `json:"fields"`
	} `json:"subActionParams"`
}

// RunGetfields The `getFields` subaction for Jira, ServiceNow ITSM, and ServiceNow SecOps connectors.
type RunGetfields struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`
}

// RunGetincident The `getIncident` subaction for Jira, ServiceNow ITSM, and ServiceNow SecOps connectors.
type RunGetincident struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// ExternalId The Jira, ServiceNow ITSM, or ServiceNow SecOps issue identifier.
		ExternalId string `json:"externalId"`
	} `json:"subActionParams"`
}

// RunIssue The `issue` subaction for Jira connectors.
type RunIssue struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams *struct {
		// Id The Jira issue identifier.
		Id string `json:"id"`
	} `json:"subActionParams,omitempty"`
}

// RunIssues The `issues` subaction for Jira connectors.
type RunIssues struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// Title The title of the Jira issue.
		Title string `json:"title"`
	} `json:"subActionParams"`
}

// RunIssuetypes The `issueTypes` subaction for Jira connectors.
type RunIssuetypes struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`
}

// RunMessageEmail Test an action that sends an email message. There must be at least one recipient in `to`, `cc`, or `bcc`.
type RunMessageEmail struct {
	// Bcc A list of "blind carbon copy" email addresses. Addresses can be specified in `user@host-name` format or in name `<user@host-name>` format
	Bcc *[]string `json:"bcc,omitempty"`

	// Cc A list of "carbon copy" email addresses. Addresses can be specified in `user@host-name` format or in name `<user@host-name>` format
	Cc *[]string `json:"cc,omitempty"`

	// Message The email message text. Markdown format is supported.
	Message string `json:"message"`

	// Subject The subject line of the email.
	Subject string `json:"subject"`

	// To A list of email addresses. Addresses can be specified in `user@host-name` format or in name `<user@host-name>` format.
	To *[]string `json:"to,omitempty"`
}

// RunMessageServerlog Test an action that writes an entry to the Kibana server log.
type RunMessageServerlog struct {
	// Level The log level of the message for server log connectors.
	Level *string `json:"level,omitempty"`

	// Message The message for server log connectors.
	Message string `json:"message"`
}

// RunMessageSlack Test an action that sends a message to Slack. It is applicable only when the connector type is `.slack`.
type RunMessageSlack struct {
	// Message The Slack message text, which cannot contain Markdown, images, or other advanced formatting.
	Message string `json:"message"`
}

// RunPostmessage Test an action that sends a message to Slack. It is applicable only when the connector type is `.slack_api`.
type RunPostmessage struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`

	// SubActionParams The set of configuration properties for the action.
	SubActionParams struct {
		// ChannelIds The Slack channel identifier, which must be one of the `allowedChannels` in the connector configuration.
		ChannelIds *[]string `json:"channelIds,omitempty"`

		// Channels The name of a channel that your Slack app has access to.
		// Deprecated:
		Channels *[]string `json:"channels,omitempty"`

		// Text The Slack message text. If it is a Slack webhook connector, the text cannot contain Markdown, images, or other advanced formatting. If it is a Slack web API connector, it can contain either plain text or block kit messages.
		Text *string `json:"text,omitempty"`
	} `json:"subActionParams"`
}

// RunPushtoservice The `pushToService` subaction for Jira, ServiceNow ITSM, ServiceNow SecOps, Swimlane, TheHive, and Webhook - Case Management connectors.
type RunPushtoservice struct {
	// SubAction The action to test.
	SubAction string `json:"subAction"`

	// SubActionParams The set of configuration properties for the action.
	SubActionParams struct {
		// Comments Additional information that is sent to Jira, ServiceNow ITSM, ServiceNow SecOps, Swimlane, or TheHive.
		Comments *[]struct {
			// Comment A comment related to the incident. For example, describe how to troubleshoot the issue.
			Comment *string `json:"comment,omitempty"`

			// CommentId A unique identifier for the comment.
			CommentId *int `json:"commentId,omitempty"`
		} `json:"comments,omitempty"`

		// Incident Information necessary to create or update a Jira, ServiceNow ITSM, ServiveNow SecOps, Swimlane, or TheHive incident.
		Incident *struct {
			// AdditionalFields Additional fields for ServiceNow ITSM and ServiveNow SecOps connectors. The fields must exist in the Elastic ServiceNow application and must be specified in JSON format.
			AdditionalFields *string `json:"additional_fields"`

			// AlertId The alert identifier for Swimlane connectors.
			AlertId *string `json:"alertId,omitempty"`

			// CaseId The case identifier for the incident for Swimlane connectors.
			CaseId *string `json:"caseId,omitempty"`

			// CaseName The case name for the incident for Swimlane connectors.
			CaseName *string `json:"caseName,omitempty"`

			// Category The category of the incident for ServiceNow ITSM and ServiceNow SecOps connectors.
			Category *string `json:"category,omitempty"`

			// CorrelationDisplay A descriptive label of the alert for correlation purposes for ServiceNow ITSM and ServiceNow SecOps connectors.
			CorrelationDisplay *string `json:"correlation_display,omitempty"`

			// CorrelationId The correlation identifier for the security incident for ServiceNow ITSM and ServiveNow SecOps connectors. Connectors using the same correlation ID are associated with the same ServiceNow incident. This value determines whether a new ServiceNow incident is created or an existing one is updated. Modifying this value is optional; if not modified, the rule ID and alert ID are combined as `{{ruleID}}:{{alert ID}}` to form the correlation ID value in ServiceNow. The maximum character length for this value is 100 characters. NOTE: Using the default configuration of `{{ruleID}}:{{alert ID}}` ensures that ServiceNow creates a separate incident record for every generated alert that uses a unique alert ID. If the rule generates multiple alerts that use the same alert IDs, ServiceNow creates and continually updates a single incident record for the alert.
			CorrelationId *string `json:"correlation_id,omitempty"`

			// Description The description of the incident for Jira, ServiceNow ITSM, ServiceNow SecOps, Swimlane, TheHive, and Webhook - Case Management connectors.
			Description *string `json:"description,omitempty"`

			// DestIp A list of destination IP addresses related to the security incident for ServiceNow SecOps connectors. The IPs are added as observables to the security incident.
			// TODO: FIX RAW MESSAGE as either string or []string
			DestIp *json.RawMessage `json:"dest_ip,omitempty"`

			// ExternalId The Jira, ServiceNow ITSM, or ServiceNow SecOps issue identifier. If present, the incident is updated. Otherwise, a new incident is created.
			ExternalId *string `json:"externalId,omitempty"`

			// Id The external case identifier for Webhook - Case Management connectors.
			Id *string `json:"id,omitempty"`

			// Impact The impact of the incident for ServiceNow ITSM connectors.
			Impact *string `json:"impact,omitempty"`

			// IssueType The type of incident for Jira connectors. For example, 10006. To obtain the list of valid values, set `subAction` to `issueTypes`.
			IssueType *int `json:"issueType,omitempty"`

			// Labels The labels for the incident for Jira connectors. NOTE: Labels cannot contain spaces.
			Labels *[]string `json:"labels,omitempty"`

			// MalwareHash A list of malware hashes related to the security incident for ServiceNow SecOps connectors. The hashes are added as observables to the security incident.
			// TODO: FIX RAW MESSAGE as either string or []string
			MalwareHash *json.RawMessage `json:"malware_hash,omitempty"`

			// MalwareUrl A list of malware URLs related to the security incident for ServiceNow SecOps connectors. The URLs are added as observables to the security incident.
			MalwareUrl *string `json:"malware_url,omitempty"`

			// OtherFields Custom field identifiers and their values for Jira connectors.
			OtherFields *map[string]interface{} `json:"otherFields,omitempty"`

			// Parent The ID or key of the parent issue for Jira connectors. Applies only to `Sub-task` types of issues.
			Parent *string `json:"parent,omitempty"`

			// Priority The priority of the incident in Jira and ServiceNow SecOps connectors.
			Priority *string `json:"priority,omitempty"`

			// RuleName The rule name for Swimlane connectors.
			RuleName *string `json:"ruleName,omitempty"`

			// Severity The severity of the incident for ServiceNow ITSM, Swimlane, and TheHive connectors. In TheHive connectors, the severity value ranges from 1 (low) to 4 (critical) with a default value of 2 (medium).
			Severity *int `json:"severity,omitempty"`

			// ShortDescription A short description of the incident for ServiceNow ITSM and ServiceNow SecOps connectors. It is used for searching the contents of the knowledge base.
			ShortDescription *string `json:"short_description,omitempty"`

			// SourceIp A list of source IP addresses related to the security incident for ServiceNow SecOps connectors. The IPs are added as observables to the security incident.
			// TODO: FIX RAW MESSAGE as either string or []string
			SourceIp *json.RawMessage `json:"source_ip,omitempty"`

			// Status The status of the incident for Webhook - Case Management connectors.
			Status *string `json:"status,omitempty"`

			// Subcategory The subcategory of the incident for ServiceNow ITSM and ServiceNow SecOps connectors.
			Subcategory *string `json:"subcategory,omitempty"`

			// Summary A summary of the incident for Jira connectors.
			Summary *string `json:"summary,omitempty"`

			// Tags A list of tags for TheHive and Webhook - Case Management connectors.
			Tags *[]string `json:"tags,omitempty"`

			// Title A title for the incident for Jira, TheHive, and Webhook - Case Management connectors. It is used for searching the contents of the knowledge base.
			Title *string `json:"title,omitempty"`

			// Tlp The traffic light protocol designation for the incident in TheHive connectors. Valid values include: 0 (clear), 1 (green), 2 (amber), 3 (amber and strict), and 4 (red).
			Tlp *int `json:"tlp,omitempty"`

			// Urgency The urgency of the incident for ServiceNow ITSM connectors.
			Urgency *string `json:"urgency,omitempty"`
		} `json:"incident,omitempty"`
	} `json:"subActionParams"`
}

// RunTriggerPagerduty Test an action that triggers a PagerDuty alert.
type RunTriggerPagerduty struct {
	// Class The class or type of the event.
	Class *string `json:"class,omitempty"`

	// Component The component of the source machine that is responsible for the event.
	Component *string `json:"component,omitempty"`

	// CustomDetails Additional details to add to the event.
	CustomDetails *map[string]interface{} `json:"customDetails,omitempty"`

	// DedupKey All actions sharing this key will be associated with the same PagerDuty alert. This value is used to correlate trigger and resolution.
	DedupKey *string `json:"dedupKey,omitempty"`

	// EventAction The type of event.
	EventAction string `json:"eventAction"`

	// Group The logical grouping of components of a service.
	Group *string `json:"group,omitempty"`

	// Links A list of links to add to the event.
	Links *[]struct {
		// Href The URL for the link.
		Href *string `json:"href,omitempty"`

		// Text A plain text description of the purpose of the link.
		Text *string `json:"text,omitempty"`
	} `json:"links,omitempty"`

	// Severity The severity of the event on the affected system.
	Severity *string `json:"severity,omitempty"`

	// Source The affected system, such as a hostname or fully qualified domain name. Defaults to the Kibana saved object id of the action.
	Source *string `json:"source,omitempty"`

	// Summary A summery of the event.
	Summary *string `json:"summary,omitempty"`

	// Timestamp An ISO-8601 timestamp that indicates when the event was detected or generated.
	Timestamp *string `json:"timestamp,omitempty"`
}

// RunValidchannelid Retrieves information about a valid Slack channel identifier. It is applicable only when the connector type is `.slack_api`.
type RunValidchannelid struct {
	// SubAction The action to test.
	SubAction       string `json:"subAction"`
	SubActionParams struct {
		// ChannelId The Slack channel identifier.
		ChannelId string `json:"channelId"`
	} `json:"subActionParams"`
}
