package tools

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/langcare/langcare-mcp-fhir/internal/fhir/providers"
)

func TestValidateResourceType(t *testing.T) {
	valid := []string{"Patient", "Observation", "MedicationRequest", "AllergyIntolerance", "DocumentReference"}
	for _, rt := range valid {
		if err := validateResourceType(rt); err != nil {
			t.Errorf("validateResourceType(%q) = %v, want nil", rt, err)
		}
	}

	invalid := []string{
		"",
		"patient",                    // FHIR resource types start with an uppercase letter
		"Patient/123",                // embeds an id
		"Patient/../Observation",     // path traversal
		"../../datasets/other/fhir",  // escapes the configured base URL
		"Patient?_count=1000",        // injects a query string
		"Patient#frag",               // fragment truncates the rest of the URL
		"Patient Observation",        // whitespace
		"$everything",                // operations are not resource types
		"Patient%2F..%2FObservation", // encoded traversal
		strings.Repeat("A", 65),
	}
	for _, rt := range invalid {
		if err := validateResourceType(rt); err == nil {
			t.Errorf("validateResourceType(%q) = nil, want error", rt)
		}
	}
}

func TestValidateResourceID(t *testing.T) {
	valid := []string{"123", "example", "a1b2-c3d4", "pat.1", "eXu8Gq9ZsH4tq7fYbVpd1aA3", strings.Repeat("a", 64)}
	for _, id := range valid {
		if err := validateResourceID(id); err != nil {
			t.Errorf("validateResourceID(%q) = %v, want nil", id, err)
		}
	}

	invalid := []string{
		"",
		"123/_history/1",        // extra path segments
		"1/../../Observation/2", // path traversal
		"..",                    // dot segment on its own
		".",                     // dot segment on its own
		"1?_format=xml",         // injects a query string
		"1#x",                   // fragment
		"1 2",                   // whitespace
		"1%2F..%2F2",            // encoded traversal
		strings.Repeat("a", 65),
	}
	for _, id := range invalid {
		if err := validateResourceID(id); err == nil {
			t.Errorf("validateResourceID(%q) = nil, want error", id)
		}
	}
}

// recordingServer is a fake FHIR server that records every path it is asked for.
type recordingServer struct {
	mu    sync.Mutex
	paths []string
}

func (s *recordingServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		s.mu.Lock()
		s.paths = append(s.paths, r.URL.RawPath+"|"+r.URL.Path+"|"+r.URL.RawQuery)
		s.mu.Unlock()
		w.Header().Set("Content-Type", "application/fhir+json")
		_, _ = w.Write([]byte(`{"resourceType":"Bundle","type":"searchset"}`))
	})
}

func (s *recordingServer) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.paths)
}

func newTestProvider(t *testing.T, baseURL string) providers.Provider {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	p, err := providers.NewGenericProvider(&providers.GenericConfig{BaseURL: baseURL, AuthType: "none"}, logger)
	if err != nil {
		t.Fatalf("NewGenericProvider: %v", err)
	}
	return p
}

// TestToolsRejectCraftedPathParams checks that malformed resourceType/id values
// are rejected before any request reaches the FHIR server.
func TestToolsRejectCraftedPathParams(t *testing.T) {
	srv := &recordingServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	client := newTestProvider(t, ts.URL+"/fhir")
	ctx := context.Background()
	resource := map[string]interface{}{"resourceType": "Patient"}

	cases := []struct {
		name string
		tool Tool
		args map[string]interface{}
	}{
		{"read: traversal in id", NewFHIRReadTool(client, nil), map[string]interface{}{"resourceType": "Patient", "id": "1/../../Observation/2"}},
		{"read: query in id", NewFHIRReadTool(client, nil), map[string]interface{}{"resourceType": "Patient", "id": "1?_format=xml"}},
		{"read: traversal in type", NewFHIRReadTool(client, nil), map[string]interface{}{"resourceType": "../../other/fhir/Patient", "id": "1"}},
		{"search: traversal in type", NewFHIRSearchTool(client, nil), map[string]interface{}{"resourceType": "Patient/../Observation", "queryParams": "code=1234"}},
		{"create: id in type", NewFHIRCreateTool(client, nil), map[string]interface{}{"resourceType": "Patient/123", "resource": resource}},
		{"update: traversal in id", NewFHIRUpdateTool(client, nil), map[string]interface{}{"resourceType": "Patient", "id": "1/../2", "resource": resource}},
		{"update: bad type", NewFHIRUpdateTool(client, nil), map[string]interface{}{"resourceType": "patient", "id": "1", "resource": resource}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := srv.count()
			if _, err := tc.tool.Execute(ctx, tc.args); err == nil {
				t.Fatalf("Execute(%v) = nil error, want validation error", tc.args)
			}
			if after := srv.count(); after != before {
				t.Fatalf("FHIR server received %d request(s) for invalid input %v; want 0", after-before, tc.args)
			}
		})
	}
}

// TestToolsAcceptValidParams makes sure ordinary calls still go through unchanged.
func TestToolsAcceptValidParams(t *testing.T) {
	srv := &recordingServer{}
	ts := httptest.NewServer(srv.handler())
	defer ts.Close()

	client := newTestProvider(t, ts.URL+"/fhir")
	ctx := context.Background()
	resource := map[string]interface{}{"resourceType": "Patient"}

	cases := []struct {
		name     string
		tool     Tool
		args     map[string]interface{}
		wantPath string
	}{
		{"read", NewFHIRReadTool(client, nil), map[string]interface{}{"resourceType": "Patient", "id": "example-1.2"}, "|/fhir/Patient/example-1.2|"},
		{"search", NewFHIRSearchTool(client, nil), map[string]interface{}{"resourceType": "Observation", "queryParams": "patient=123&code=8867-4"}, "|/fhir/Observation|patient=123&code=8867-4"},
		{"create", NewFHIRCreateTool(client, nil), map[string]interface{}{"resourceType": "Patient", "resource": resource}, "|/fhir/Patient|"},
		{"update", NewFHIRUpdateTool(client, nil), map[string]interface{}{"resourceType": "Patient", "id": "123", "resource": resource}, "|/fhir/Patient/123|"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.tool.Execute(ctx, tc.args); err != nil {
				t.Fatalf("Execute(%v) = %v, want nil", tc.args, err)
			}
			srv.mu.Lock()
			got := srv.paths[len(srv.paths)-1]
			srv.mu.Unlock()
			if got != tc.wantPath {
				t.Fatalf("server saw %q, want %q", got, tc.wantPath)
			}
		})
	}
}
