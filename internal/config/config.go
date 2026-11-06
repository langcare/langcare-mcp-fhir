package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// GetMCPAuthTokens returns the list of valid MCP auth tokens
func (c *Config) GetMCPAuthTokens() []string {
	// Try config file first
	tokens := c.Security.MCPAuthTokens

	// Fallback to environment variable
	if tokens == "" {
		tokens = os.Getenv("MCP_AUTH_TOKENS")
	}

	// Split by comma and trim whitespace
	if tokens == "" {
		return []string{}
	}

	var result []string
	for _, token := range strings.Split(tokens, ",") {
		trimmed := strings.TrimSpace(token)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// Config represents the application configuration
type Config struct {
	FHIRServer FHIRServerConfig `yaml:"fhir_server"`
	MCP        MCPConfig        `yaml:"mcp"`
	Transport  TransportConfig  `yaml:"transport"`
	Logging    LoggingConfig    `yaml:"logging"`
	Security   SecurityConfig   `yaml:"security"`
}

// FHIRServerConfig contains FHIR server connection settings
type FHIRServerConfig struct {
	Provider string             `yaml:"provider"` // gcp, epic, cerner, openemr, generic
	BaseURL  string             `yaml:"base_url"`
	Auth     AuthConfig         `yaml:"auth"`
	GCP      GCPConfig          `yaml:"gcp"`
	EPIC     EPICOAuthConfig    `yaml:"epic"`
	Cerner   CernerOAuthConfig  `yaml:"cerner"`
	OpenEMR  OpenEMROAuthConfig `yaml:"openemr"`
	Timeout  time.Duration      `yaml:"timeout"`
	Headers  map[string]string  `yaml:"headers"`
}

// AuthConfig contains authentication settings for generic provider
type AuthConfig struct {
	Type           string            `yaml:"type"` // bearer, basic, none, oauth2_client_credentials
	Token          string            `yaml:"token"`
	Username       string            `yaml:"username"`
	Password       string            `yaml:"password"`
	TokenURL       string            `yaml:"token_url"`
	ClientID       string            `yaml:"client_id"`
	ClientSecret   string            `yaml:"client_secret"`
	Scopes         []string          `yaml:"scopes"`
	EndpointParams map[string]string `yaml:"endpoint_params"`
	AuthStyle      string            `yaml:"auth_style"` // params, header, auto (default: auto)
}

// GCPConfig contains GCP Healthcare API settings
type GCPConfig struct {
	ProjectID      string `yaml:"project_id"`
	Location       string `yaml:"location"`
	DatasetID      string `yaml:"dataset_id"`
	FHIRStoreID    string `yaml:"fhir_store_id"`
	CredentialPath string `yaml:"credential_path"` // Path to service account JSON
}

// EPICOAuthConfig contains EPIC OAuth2 settings
type EPICOAuthConfig struct {
	ClientID       string   `yaml:"client_id"`
	PrivateKeyPath string   `yaml:"private_key_path"` // Path to RSA private key PEM
	TokenURL       string   `yaml:"token_url"`
	Scopes         []string `yaml:"scopes"`
}

// OpenEMROAuthConfig contains OpenEMR SMART on FHIR Backend Services settings
// (private_key_jwt with RS384, per the SMART on FHIR Backend Services spec).
type OpenEMROAuthConfig struct {
	ClientID       string   `yaml:"client_id"`
	PrivateKeyPath string   `yaml:"private_key_path"` // Path to RSA private key PEM
	TokenURL       string   `yaml:"token_url"`        // e.g. https://host/oauth2/default/token
	KeyID          string   `yaml:"key_id"`           // JWK kid registered with OpenEMR
	Scopes         []string `yaml:"scopes"`           // system/* scopes
}

// CernerOAuthConfig contains Cerner OAuth2 settings
type CernerOAuthConfig struct {
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	TokenURL     string   `yaml:"token_url"`
	Scopes       []string `yaml:"scopes"`
}

// MCPConfig contains MCP server settings
type MCPConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// TransportConfig contains transport layer settings
type TransportConfig struct {
	Stdio bool       `yaml:"stdio"`
	HTTP  HTTPConfig `yaml:"http"`
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	MCPAuthTokens string          `yaml:"mcp_auth_tokens"` // Comma-separated list of valid MCP tokens
	RateLimit     RateLimitConfig `yaml:"rate_limit"`
}

// RateLimitConfig contains rate limiting settings
type RateLimitConfig struct {
	Enabled bool `yaml:"enabled"`
	Rate    int  `yaml:"rate"`  // requests per second
	Burst   int  `yaml:"burst"` // burst size
}

// HTTPConfig contains HTTP transport settings
type HTTPConfig struct {
	Enabled bool      `yaml:"enabled"`
	Port    int       `yaml:"port"`
	TLS     TLSConfig `yaml:"tls"`
}

// TLSConfig contains TLS settings
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level    string `yaml:"level"`  // debug, info, warn, error
	Format   string `yaml:"format"` // json, text
	ScrubPHI bool   `yaml:"scrub_phi"`
}

// Load reads configuration from a YAML file and environment variables
func Load(path string) (*Config, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the YAML content
	expandedData := os.ExpandEnv(string(data))

	// Parse YAML
	var cfg Config
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	cfg.setDefaults()

	// Validate
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default values for optional fields
func (c *Config) setDefaults() {
	if c.FHIRServer.Timeout == 0 {
		c.FHIRServer.Timeout = 30 * time.Second
	}
	if c.MCP.Name == "" {
		c.MCP.Name = "LangCare MCP FHIR Server"
	}
	if c.MCP.Version == "" {
		c.MCP.Version = "2.5.0"
	}
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}
	if c.Transport.HTTP.Port == 0 {
		c.Transport.HTTP.Port = 8080
	}
	// Default to PHI scrubbing enabled for safety
	if !c.Logging.ScrubPHI {
		c.Logging.ScrubPHI = true
	}
	// Rate limiting defaults
	if c.Security.RateLimit.Rate == 0 {
		c.Security.RateLimit.Rate = 100
	}
	if c.Security.RateLimit.Burst == 0 {
		c.Security.RateLimit.Burst = 200
	}
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	// Validate provider type
	if c.FHIRServer.Provider == "" {
		c.FHIRServer.Provider = "generic" // Default to generic for backwards compatibility
	}

	validProviders := map[string]bool{
		"gcp":     true,
		"epic":    true,
		"cerner":  true,
		"openemr": true,
		"generic": true,
	}
	if !validProviders[c.FHIRServer.Provider] {
		return fmt.Errorf("fhir_server.provider must be one of: gcp, epic, cerner, openemr, generic")
	}

	// Validate provider-specific configuration
	switch c.FHIRServer.Provider {
	case "gcp":
		// BaseURL is optional for GCP (will be constructed)
		if c.FHIRServer.GCP.ProjectID == "" {
			return fmt.Errorf("fhir_server.gcp.project_id is required for GCP provider")
		}
		if c.FHIRServer.GCP.Location == "" {
			return fmt.Errorf("fhir_server.gcp.location is required for GCP provider")
		}
		if c.FHIRServer.GCP.DatasetID == "" {
			return fmt.Errorf("fhir_server.gcp.dataset_id is required for GCP provider")
		}
		if c.FHIRServer.GCP.FHIRStoreID == "" {
			return fmt.Errorf("fhir_server.gcp.fhir_store_id is required for GCP provider")
		}

	case "epic":
		if c.FHIRServer.BaseURL == "" {
			return fmt.Errorf("fhir_server.base_url is required for EPIC provider")
		}
		if c.FHIRServer.EPIC.ClientID == "" {
			return fmt.Errorf("fhir_server.epic.client_id is required for EPIC provider")
		}
		if c.FHIRServer.EPIC.PrivateKeyPath == "" {
			return fmt.Errorf("fhir_server.epic.private_key_path is required for EPIC provider")
		}
		if c.FHIRServer.EPIC.TokenURL == "" {
			return fmt.Errorf("fhir_server.epic.token_url is required for EPIC provider")
		}

	case "openemr":
		if c.FHIRServer.BaseURL == "" {
			return fmt.Errorf("fhir_server.base_url is required for OpenEMR provider")
		}
		if c.FHIRServer.OpenEMR.ClientID == "" {
			return fmt.Errorf("fhir_server.openemr.client_id is required for OpenEMR provider")
		}
		if c.FHIRServer.OpenEMR.PrivateKeyPath == "" {
			return fmt.Errorf("fhir_server.openemr.private_key_path is required for OpenEMR provider")
		}
		if c.FHIRServer.OpenEMR.TokenURL == "" {
			return fmt.Errorf("fhir_server.openemr.token_url is required for OpenEMR provider")
		}

	case "cerner":
		if c.FHIRServer.BaseURL == "" {
			return fmt.Errorf("fhir_server.base_url is required for Cerner provider")
		}
		if c.FHIRServer.Cerner.ClientID == "" {
			return fmt.Errorf("fhir_server.cerner.client_id is required for Cerner provider")
		}
		if c.FHIRServer.Cerner.ClientSecret == "" {
			return fmt.Errorf("fhir_server.cerner.client_secret is required for Cerner provider")
		}
		if c.FHIRServer.Cerner.TokenURL == "" {
			return fmt.Errorf("fhir_server.cerner.token_url is required for Cerner provider")
		}

	case "generic":
		if c.FHIRServer.BaseURL == "" {
			return fmt.Errorf("fhir_server.base_url is required")
		}
		if !strings.HasPrefix(c.FHIRServer.BaseURL, "http://") && !strings.HasPrefix(c.FHIRServer.BaseURL, "https://") {
			return fmt.Errorf("fhir_server.base_url must start with http:// or https://")
		}

		// Validate auth type for generic provider
		validAuthTypes := map[string]bool{
			"bearer":                    true,
			"basic":                     true,
			"none":                      true,
			"oauth2_client_credentials": true,
		}
		if c.FHIRServer.Auth.Type == "" {
			c.FHIRServer.Auth.Type = "none" // Default
		}
		if !validAuthTypes[c.FHIRServer.Auth.Type] {
			return fmt.Errorf("fhir_server.auth.type must be one of: bearer, basic, none, oauth2_client_credentials")
		}

		// Validate auth credentials based on type
		switch c.FHIRServer.Auth.Type {
		case "bearer":
			if c.FHIRServer.Auth.Token == "" {
				return fmt.Errorf("fhir_server.auth.token is required when auth type is bearer")
			}
		case "basic":
			if c.FHIRServer.Auth.Username == "" || c.FHIRServer.Auth.Password == "" {
				return fmt.Errorf("fhir_server.auth.username and password are required when auth type is basic")
			}
		case "oauth2_client_credentials":
			if c.FHIRServer.Auth.TokenURL == "" {
				return fmt.Errorf("fhir_server.auth.token_url is required when auth type is oauth2_client_credentials")
			}
			if c.FHIRServer.Auth.ClientID == "" {
				return fmt.Errorf("fhir_server.auth.client_id is required when auth type is oauth2_client_credentials")
			}
			if c.FHIRServer.Auth.ClientSecret == "" {
				return fmt.Errorf("fhir_server.auth.client_secret is required when auth type is oauth2_client_credentials")
			}
			validAuthStyles := map[string]bool{"params": true, "header": true, "auto": true, "": true}
			if !validAuthStyles[c.FHIRServer.Auth.AuthStyle] {
				return fmt.Errorf("fhir_server.auth.auth_style must be one of: params, header, auto")
			}
		}
	}

	// Validate transport - at least one must be enabled
	if !c.Transport.Stdio && !c.Transport.HTTP.Enabled {
		return fmt.Errorf("at least one transport (stdio or http) must be enabled")
	}

	// Validate MCP auth tokens for HTTP transport
	if c.Transport.HTTP.Enabled && c.Security.MCPAuthTokens == "" {
		// Check environment variable as fallback
		if envTokens := os.Getenv("MCP_AUTH_TOKENS"); envTokens == "" {
			return fmt.Errorf("security.mcp_auth_tokens is required for HTTP transport (or set MCP_AUTH_TOKENS env var)")
		}
	}

	// Validate TLS config if enabled
	if c.Transport.HTTP.TLS.Enabled {
		if c.Transport.HTTP.TLS.CertFile == "" || c.Transport.HTTP.TLS.KeyFile == "" {
			return fmt.Errorf("transport.http.tls.cert_file and key_file are required when TLS is enabled")
		}
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error")
	}

	// Validate log format
	if c.Logging.Format != "json" && c.Logging.Format != "text" {
		return fmt.Errorf("logging.format must be either json or text")
	}

	return nil
}
