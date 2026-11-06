package config_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/langcare/langcare-mcp-fhir/internal/config"
)

func TestLoad_GenericAuth(t *testing.T) {
	t.Parallel()
	const fixtureValue = "synthetic-config-value"
	tests := []struct {
		name      string
		authType  string
		omit      string
		style     string
		wantError string
	}{
		{name: "valid OAuth2", authType: "oauth2_client_credentials"},
		{name: "missing token_url", authType: "oauth2_client_credentials", omit: "token_url", wantError: "token_url is required"},
		{name: "missing client_id", authType: "oauth2_client_credentials", omit: "client_id", wantError: "client_id is required"},
		{name: "missing client_secret", authType: "oauth2_client_credentials", omit: "client_secret", wantError: "client_secret is required"},
		{name: "invalid auth_style", authType: "oauth2_client_credentials", style: "invalid", wantError: "auth_style must be one of"},
		{name: "auth_style params", authType: "oauth2_client_credentials", style: "params"},
		{name: "auth_style header", authType: "oauth2_client_credentials", style: "header"},
		{name: "auth_style auto", authType: "oauth2_client_credentials", style: "auto"},
		{name: "auth_style empty", authType: "oauth2_client_credentials"},
		{name: "bearer", authType: "bearer"},
		{name: "basic", authType: "basic"},
		{name: "none", authType: "none"},
		{name: "empty type defaults to none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var auth strings.Builder
			fmt.Fprintf(&auth, "    type: %q\n", tt.authType)
			switch tt.authType {
			case "oauth2_client_credentials":
				for _, field := range []struct{ name, value string }{
					{"token_url", "https://auth.example.com/token"},
					{"client_id", fixtureValue},
					{"client_secret", fixtureValue},
				} {
					if field.name != tt.omit {
						fmt.Fprintf(&auth, "    %s: %q\n", field.name, field.value)
					}
				}
				fmt.Fprintf(&auth, "    auth_style: %q\n", tt.style)
			case "bearer":
				fmt.Fprintf(&auth, "    token: %q\n", fixtureValue)
			case "basic":
				fmt.Fprintf(&auth, "    username: %q\n    password: %q\n", fixtureValue, fixtureValue)
			}
			yaml := "fhir_server:\n  provider: generic\n  base_url: https://fhir.example.com\n  auth:\n" + auth.String() + "transport:\n  stdio: true\n"
			file, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
			if err != nil {
				t.Fatalf("create config: %v", err)
			}
			defer file.Close()
			if _, err := file.WriteString(yaml); err != nil {
				t.Fatalf("write config: %v", err)
			}
			if err := file.Close(); err != nil {
				t.Fatalf("close config: %v", err)
			}
			cfg, err := config.Load(file.Name())
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("Load error = %v, want containing %q", err, tt.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			wantType := tt.authType
			if wantType == "" {
				wantType = "none"
			}
			if cfg.FHIRServer.Auth.Type != wantType {
				t.Errorf("auth type = %q, want %q", cfg.FHIRServer.Auth.Type, wantType)
			}
			if cfg.FHIRServer.Auth.AuthStyle != tt.style {
				t.Errorf("auth style = %q, want %q", cfg.FHIRServer.Auth.AuthStyle, tt.style)
			}
		})
	}
}
