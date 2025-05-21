package kbapi

import "encoding/json"

type FleetProxiesResponseItem struct {
	Certificate            *string                 `json:"certificate,omitempty"`
	CertificateAuthorities *string                 `json:"certificate_authorities,omitempty"`
	CertificateKey         *string                 `json:"certificate_key,omitempty"`
	ID                     string                  `json:"id"`
	IsPreconfigured        *bool                   `json:"is_preconfigured,omitempty"`
	Name                   string                  `json:"name"`
	ProxyHeaders           *map[string]interface{} `json:"proxy_headers"`
	URL                    string                  `json:"url"`
}

type FleetProxiesRequestBody struct {
	Certificate            *string         `json:"certificate,omitempty"`
	CertificateAuthorities *string         `json:"certificate_authorities,omitempty"`
	CertificateKey         *string         `json:"certificate_key,omitempty"`
	ID                     *string         `json:"id,omitempty"`
	IsPreconfigured        *bool           `json:"is_preconfigured,omitempty"`
	Name                   string          `json:"name"`
	ProxyHeaders           json.RawMessage `json:"proxy_headers,omitempty"`
	URL                    string          `json:"url"`
}

// SetProxyHeaders sets all headers at once, replacing any existing headers
func (body *FleetProxiesRequestBody) SetProxyHeaders(headers map[string]interface{}) error {
	data, err := json.Marshal(headers)
	if err != nil {
		return err
	}
	body.ProxyHeaders = data
	return nil
}

// GetProxyHeaders returns the current proxy headers as a map
func (body *FleetProxiesRequestBody) GetProxyHeaders() (map[string]interface{}, error) {
	if len(body.ProxyHeaders) == 0 {
		return make(map[string]interface{}), nil
	}

	var headersMap map[string]interface{}
	if err := json.Unmarshal(body.ProxyHeaders, &headersMap); err != nil {
		return nil, err
	}
	return headersMap, nil
}
