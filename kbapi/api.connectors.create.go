package kbapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TODO: Update the call
// ConnectorsCreateResponse wraps the response from a <todo> call
type ConnectorsCreateResponse struct {
	StatusCode int
	Body       *ConnectorsCreateResponseBody
	Error      interface{}
	RawBody    io.ReadCloser
}

type ConnectorsCreateResponseBody struct {
	// Id The identifier for the connector.
	ID     string                 `json:"id"`
	Config map[string]interface{} `json:"config"`
	// Name  The name of the rule.
	Name string `json:"name"`
	// ConnectorTypeId The connector type identifier.
	ConnectorTypeID string `json:"connector_type_id"`
	// IsDeprecated Indicates whether the connector is deprecated.
	IsDeprecated bool `json:"is_deprecated"`
	// IsMissingSecrets Indicates whether the connector is missing secrets.
	IsMissingSecrets bool `json:"is_missing_secrets"`
	// IsPreconfigured Indicates whether the connector is preconfigured. If true, the `config` and `is_missing_secrets` properties are omitted from the response.
	IsPreconfigured bool `json:"is_preconfigured"`
	// IsSystemActionType Indicates whether the connector is used for system actions.
	IsSystemActionType bool `json:"is_system_action_type"`
}

type ConnectorsCreateRequest struct {
	ID   string
	Body ConnectorsCreateRequestBody
}

// SetBedRock sets the AWS Bedrock configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetBedRock(config BedrockConfig, secrets BedrockSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON

	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON

	return nil
}

// SetCasesWebhook sets the Cases Webhook configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetCasesWebhook(config CasesWebhookConfig, secrets CasesWebhookSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetCrowdstrike sets the Crowdstrike configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetCrowdstrike(config CrowdstrikeConfig, secrets CrowdstrikeSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetD3security sets the D3 Security configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetD3security(config D3securityConfig, secrets D3securitySecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetEmail sets the Email configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetEmail(config EmailConfig, secrets EmailSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetGemini sets the Gemini configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetGemini(config GeminiConfig, secrets GeminiSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetGenaiAzure sets the Azure GenAI configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetGenaiAzure(config GenaiAzureConfig, secrets GenaiSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetGenaiOpenai sets the OpenAI GenAI configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetGenaiOpenai(config GenaiOpenaiConfig, secrets GenaiSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetIndex sets the Index configuration for a connector create request
func (body *ConnectorsCreateRequestBody) SetIndex(config IndexConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	return nil
}

// SetJira sets the Jira configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetJira(config JiraConfig, secrets JiraSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetOpsgenie sets the Opsgenie configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetOpsgenie(config OpsgenieConfig, secrets OpsgenieSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetPagerduty sets the PagerDuty configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetPagerduty(config PagerdutyConfig, secrets PagerdutySecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetResilient sets the Resilient configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetResilient(config ResilientConfig, secrets ResilientSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetSentinelone sets the SentinelOne configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetSentinelone(config SentineloneConfig, secrets SentineloneSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetServicenow sets the ServiceNow configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetServicenow(config ServicenowConfig, secrets ServicenowSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetServicenowItom sets the ServiceNow ITOM configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetServicenowItom(config ServicenowItomConfig, secrets ServicenowSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetSlackAPI sets the Slack API configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetSlackAPI(config SlackAPIConfig, secrets SlackAPISecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetSwimlane sets the Swimlane configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetSwimlane(config SwimlaneConfig, secrets SwimlaneSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetTeams sets the Microsoft Teams secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetTeams(secrets TeamsSecrets) error {
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetThehive sets the TheHive configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetThehive(config ThehiveConfig, secrets ThehiveSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetTines sets the Tines configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetTines(config TinesConfig, secrets TinesSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetTorq sets the Torq configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetTorq(config TorqConfig, secrets TorqSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetWebhook sets the Webhook configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetWebhook(config WebhookConfig, secrets WebhookSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

// SetXmatters sets the xMatters configuration and secrets for a connector create request
func (body *ConnectorsCreateRequestBody) SetXmatters(config XmattersConfig, secrets XmattersSecrets) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	body.Config = configJSON
	secretsJSON, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	body.Secrets = secretsJSON
	return nil
}

type ConnectorsCreateRequestBody struct {
	Name            string          `json:"name"`
	ConnectorTypeID string          `json:"connector_type_id"`
	Config          json.RawMessage `json:"config"`
	Secrets         json.RawMessage `json:"secrets,omitempty"`
}

// newConnectorsCreate returns a function that performs POST /api/actions/connector/{id} API requests
func (api *API) newConnectorsCreate() func(context.Context, *ConnectorsCreateRequest, ...RequestOption) (*ConnectorsCreateResponse, error) {
	return func(ctx context.Context, req *ConnectorsCreateRequest, opts ...RequestOption) (*ConnectorsCreateResponse, error) {
		if req == nil {
			return nil, fmt.Errorf("Request cannot be nil")
		}

		// Get instrumentation if available
		var instrument Instrumentation
		if i, ok := api.transport.(Instrumented); ok {
			instrument = i.InstrumentationEnabled()
		}

		// Start instrumentation span if available
		if instrument != nil {
			var newCtx context.Context
			newCtx = instrument.Start(ctx, "connectors.create")
			defer instrument.Close(newCtx)
			ctx = newCtx
		}

		path := fmt.Sprintf("/api/actions/connector/%s", req.ID)

		// Create HTTP request
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, nil)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		jsonBody, err := json.Marshal(req.Body)
		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		httpReq.Body = io.NopCloser(bytes.NewReader(jsonBody))
		httpReq.Header.Set("Content-Type", "application/json")

		// Apply all the functional options
		for _, opt := range opts {
			if err := opt(httpReq); err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
		}

		// Pre-request instrumentation
		if instrument != nil {
			instrument.BeforeRequest(httpReq, "connectors.create")
			if reader := instrument.RecordRequestBody(ctx, "connectors.create", httpReq.Body); reader != nil {
				httpReq.Body = reader
			}
		}

		// Execute request
		httpResp, err := api.transport.Perform(httpReq)

		if instrument != nil {
			instrument.AfterRequest(httpReq, "kibana", path)
		}

		if err != nil {
			if instrument != nil {
				instrument.RecordError(ctx, err)
			}
			return nil, err
		}

		// Prepare response
		resp := &ConnectorsCreateResponse{
			StatusCode: httpResp.StatusCode,
			RawBody:    httpResp.Body,
		}

		var result ConnectorsCreateResponseBody

		if httpResp.StatusCode == 200 {
			if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
				httpResp.Body.Close()
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, err
			}
			resp.Body = &result
			return resp, nil
		} else {
			// For all non-200 responses
			bodyBytes, err := io.ReadAll(httpResp.Body)
			httpResp.Body.Close()
			if err != nil {
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return nil, fmt.Errorf("failed to read response body: %v", err)
			}

			// Try to decode as JSON
			var errorObj interface{}
			if err := json.Unmarshal(bodyBytes, &errorObj); err == nil {
				resp.Error = errorObj

				errorMessage, _ := json.Marshal(errorObj)

				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, errorMessage)
			} else {
				// Not valid JSON
				resp.Error = string(bodyBytes)
				if instrument != nil {
					instrument.RecordError(ctx, err)
				}
				return resp, fmt.Errorf("HTTP Status Code %d: %s", httpResp.StatusCode, string(bodyBytes))
			}
		}
	}
}
