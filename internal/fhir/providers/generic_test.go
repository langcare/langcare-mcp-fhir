package providers

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestGenericProvider_OAuth2ClientCredentials(t *testing.T) {
	t.Parallel()
	const fixtureValue = "synthetic-client-value"
	tests := []struct {
		name         string
		style        string
		endpoint     map[string]string
		refresh      bool
		stringExpiry bool
		unauthorized bool
	}{
		{name: "token acquisition and FHIR Read and Search"},
		{name: "auth style params", style: "params"},
		{name: "auth style header", style: "header"},
		{name: "endpoint params", endpoint: map[string]string{"audience": "https://fhir.example.com"}},
		{name: "token refresh after expiry", refresh: true},
		{name: "string expires_in", stringExpiry: true},
		{name: "token endpoint 401", style: "params", unauthorized: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var tokenCalls, fhirCalls atomic.Int32
			tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				call := tokenCalls.Add(1)
				if r.Method != http.MethodPost {
					t.Errorf("token method = %s, want POST", r.Method)
				}
				if err := r.ParseForm(); err != nil {
					t.Errorf("parse token form: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if got := r.PostForm.Get("grant_type"); got != "client_credentials" {
					t.Errorf("grant_type = %q, want client_credentials", got)
				}
				if tt.style == "params" {
					if r.PostForm.Get("client_id") != fixtureValue || r.PostForm.Get("client_secret") != fixtureValue {
						t.Error("client credentials missing from form body")
					}
					if r.Header.Get("Authorization") != "" {
						t.Error("params style unexpectedly sent Authorization header")
					}
				} else {
					username, password, ok := r.BasicAuth()
					if !ok || username != fixtureValue || password != fixtureValue {
						t.Error("client credentials missing from Basic auth header")
					}
					if r.PostForm.Get("client_id") != "" || r.PostForm.Get("client_secret") != "" {
						t.Error("header style unexpectedly sent credentials in form body")
					}
				}
				for key, want := range tt.endpoint {
					if got := r.PostForm.Get(key); got != want {
						t.Errorf("endpoint param %s = %q, want %q", key, got, want)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				if tt.unauthorized {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(w, `{"error":"invalid_client"}`)
					return
				}
				expiry := "3600"
				if tt.refresh && call == 1 {
					expiry = "1"
				} else if tt.stringExpiry {
					expiry = `"3600"`
				}
				fmt.Fprintf(w, `{"access_token":"t%d","token_type":"Bearer","expires_in":%s}`, call, expiry)
			}))
			defer tokenServer.Close()

			fhirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				call := fhirCalls.Add(1)
				want := "Bearer t1"
				if tt.refresh && call == 2 {
					want = "Bearer t2"
				}
				if got := r.Header.Get("Authorization"); got != want {
					t.Errorf("FHIR Authorization = %q, want %q", got, want)
				}
				w.Header().Set("Content-Type", "application/fhir+json")
				switch r.URL.RequestURI() {
				case "/Patient/example":
					fmt.Fprint(w, `{"resourceType":"Patient","id":"example"}`)
				case "/Patient?name=example":
					fmt.Fprint(w, `{"resourceType":"Bundle","entry":[]}`)
				default:
					t.Errorf("unexpected FHIR request: %s", r.URL.RequestURI())
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer fhirServer.Close()

			provider, err := NewGenericProvider(&GenericConfig{
				BaseURL: fhirServer.URL, AuthType: "oauth2_client_credentials",
				TokenURL: tokenServer.URL, ClientID: fixtureValue, ClientSecret: fixtureValue,
				AuthStyle: tt.style, EndpointParams: tt.endpoint,
			}, slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatalf("NewGenericProvider: %v", err)
			}
			for i, operation := range []string{"Read", "Search"} {
				if tt.refresh && i == 1 {
					time.Sleep(1100 * time.Millisecond)
				}
				var body []byte
				if operation == "Read" {
					body, err = provider.Read(context.Background(), "Patient", "example")
				} else {
					body, err = provider.Search(context.Background(), "Patient", "name=example")
				}
				if tt.unauthorized {
					if err == nil {
						t.Fatalf("%s succeeded, want token endpoint error", operation)
					}
					var retrieveErr *oauth2.RetrieveError
					if !errors.As(err, &retrieveErr) || retrieveErr.Response.StatusCode != http.StatusUnauthorized {
						t.Fatalf("%s error does not wrap OAuth2 HTTP 401: %v", operation, err)
					}
					for _, fragment := range []string{"request failed:", "invalid_client"} {
						if !strings.Contains(err.Error(), fragment) {
							t.Errorf("%s error %q missing %q", operation, err, fragment)
						}
					}
					continue
				}
				if err != nil {
					t.Fatalf("%s: %v", operation, err)
				}
				want := `{"resourceType":"Patient","id":"example"}`
				if operation == "Search" {
					want = `{"resourceType":"Bundle","entry":[]}`
				}
				if string(body) != want {
					t.Errorf("%s response = %s, want %s", operation, body, want)
				}
			}
			wantTokenCalls, wantFHIRCalls := int32(1), int32(2)
			if tt.refresh || tt.unauthorized {
				wantTokenCalls = 2
			}
			if tt.unauthorized {
				wantFHIRCalls = 0
			}
			if got := tokenCalls.Load(); got != wantTokenCalls {
				t.Errorf("token calls = %d, want %d", got, wantTokenCalls)
			}
			if got := fhirCalls.Load(); got != wantFHIRCalls {
				t.Errorf("FHIR calls = %d, want %d", got, wantFHIRCalls)
			}
		})
	}
}

func TestGenericProvider_LegacyAuth(t *testing.T) {
	t.Parallel()
	const fixtureValue = "synthetic-legacy-value"
	tests := []struct {
		name string
		want string
	}{
		{name: "bearer", want: "Bearer " + fixtureValue},
		{name: "basic", want: "Basic " + base64.StdEncoding.EncodeToString([]byte(fixtureValue+":"+fixtureValue))},
		{name: "none"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var calls atomic.Int32
			fhirServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if got := r.Header.Get("Authorization"); got != tt.want {
					t.Errorf("Authorization = %q, want %q", got, tt.want)
				}
				w.Header().Set("Content-Type", "application/fhir+json")
				fmt.Fprint(w, `{"resourceType":"Patient","id":"example"}`)
			}))
			defer fhirServer.Close()
			cfg := &GenericConfig{BaseURL: fhirServer.URL, AuthType: tt.name}
			switch tt.name {
			case "bearer":
				cfg.Token = fixtureValue
			case "basic":
				cfg.Username, cfg.Password = fixtureValue, fixtureValue
			}
			provider, err := NewGenericProvider(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatalf("NewGenericProvider: %v", err)
			}
			if _, err := provider.Read(context.Background(), "Patient", "example"); err != nil {
				t.Fatalf("Read: %v", err)
			}
			if _, err := provider.Search(context.Background(), "Patient", "name=example"); err != nil {
				t.Fatalf("Search: %v", err)
			}
			if got := calls.Load(); got != 2 {
				t.Errorf("FHIR calls = %d, want 2", got)
			}
		})
	}
}
