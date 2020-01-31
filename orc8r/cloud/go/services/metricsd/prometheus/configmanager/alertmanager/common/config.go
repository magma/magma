package common

import (
	amconfig "github.com/prometheus/common/config"
)

// HTTPConfig is a copy of prometheus/common/config.HTTPClientConfig with
// `Secret` fields replaced with strings to enable marshaling without obfuscation
type HTTPConfig struct {
	// The HTTP basic authentication credentials for the targets.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty" json:"basic_auth,omitempty"`
	// The bearer token for the targets.
	BearerToken string `yaml:"bearer_token,omitempty" json:"bearer_token,omitempty"`
	// The bearer token file for the targets.

	// TODO: Support file storage
	//BearerTokenFile string `yaml:"bearer_token_file,omitempty"`
	// HTTP proxy server to use to connect to the targets.
	ProxyURL *amconfig.URL `yaml:"proxy_url,omitempty" json:"proxy_url,omitempty"`
	// TLSConfig to use to connect to the targets.
	TLSConfig TLSConfig `yaml:"tls_config,omitempty" json:"tls_config,omitempty"`
}

// BasicAuth is a copy of prometheus/common/config.BasicAuth with `Secret`
// fields replaced with strings to enable marshaling without obfuscation
type BasicAuth struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`

	// TODO: Support file storage
	//PasswordFile string `yaml:"password_file,omitempty"`
}

// TLSConfig is a copy of prometheus/common/config.TLSConfig without file fields
// since storing files is not supported by alertmanager-configurer yet
type TLSConfig struct {
	// TODO: Support file storage
	//// The CA cert to use for the targets.
	//CAFile string `yaml:"ca_file,omitempty"`
	//// The client cert file for the targets.
	//CertFile string `yaml:"cert_file,omitempty"`
	////The client key file for the targets.
	//KeyFile string `yaml:"key_file,omitempty"`

	// Used to verify the hostname for the targets.
	ServerName string `yaml:"server_name,omitempty" json:"server_name,omitempty"`
	// Disable target certificate validation.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify" json:"insecure_skip_verify,omitempty"`
}
