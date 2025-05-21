package kbapi

// ServerHost defines model for server_host.
type FleetServerHostItem struct {
	HostURLs        []string `json:"host_urls"`
	ID              string   `json:"id"`
	IsDefault       *bool    `json:"is_default,omitempty"`
	IsInternal      *bool    `json:"is_internal,omitempty"`
	IsPreconfigured *bool    `json:"is_preconfigured,omitempty"`
	Name            string   `json:"name"`
	ProxyID         *string  `json:"proxy_id"`
}
