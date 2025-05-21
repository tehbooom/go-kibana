package kbapi

type FleetBinaryDownloadResponseItem struct {
	Host      string `json:"host"`
	ID        string `json:"id"`
	IsDefault *bool  `json:"is_default,omitempty"`
	Name      string `json:"name"`
	// ProxyId The ID of the proxy to use for this download source. See the proxies API for more information.
	ProxyID *string `json:"proxy_id"`
	// Secrets *struct {
	// 	SSL *struct {
	// 		Key *string `json:"key,omitempty"`
	// 	} `json:"ssl,omitempty"`
	// } `json:"secrets,omitempty"`
	// SSL *struct {
	// 	Certificate            *string   `json:"certificate,omitempty"`
	// 	CertificateAuthorities *[]string `json:"certificate_authorities,omitempty"`
	// 	Key                    *string   `json:"key,omitempty"`
	// } `json:"ssl,omitempty"`
}
