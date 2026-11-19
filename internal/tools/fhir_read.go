package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/langcare/langcare-mcp-fhir/internal/audit"
	"github.com/langcare/langcare-mcp-fhir/internal/fhir"
)

// FHIRReadTool implements the fhir_read tool
type FHIRReadTool struct {
	client      fhir.Client
	auditLogger *audit.AuditLogger
}

// NewFHIRReadTool creates a new FHIR read tool
func NewFHIRReadTool(client fhir.Client, auditLogger *audit.AuditLogger) *FHIRReadTool {
	return &FHIRReadTool{
		client:      client,
		auditLogger: auditLogger,
	}
}

// Name returns the tool name
func (t *FHIRReadTool) Name() string {
	return "fhir_read"
}

// Description returns the tool description
func (t *FHIRReadTool) Description() string {
	return "Read a FHIR resource by type and ID"
}

// InputSchema returns the JSON schema for tool inputs
func (t *FHIRReadTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"resourceType": map[string]interface{}{
				"type":        "string",
				"description": "The FHIR resource type (e.g., Patient, Observation, Medication)",
			},
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The resource ID",
			},
		},
		"required": []string{"resourceType", "id"},
	}
}

// Execute executes the tool
func (t *FHIRReadTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	resourceType, ok := args["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("resourceType is required and must be a string")
	}
	if err := validateResourceType(resourceType); err != nil {
		return nil, err
	}

	id, ok := args["id"].(string)
	if !ok {
		return nil, fmt.Errorf("id is required and must be a string")
	}
	if err := validateResourceID(id); err != nil {
		return nil, err
	}

	// Execute FHIR read
	resource, err := t.client.Read(ctx, resourceType, id)

	// Log PHI access
	if t.auditLogger != nil {
		t.auditLogger.LogPHIAccess(ctx, audit.PHIAccessEvent{
			Operation:    "read",
			ResourceType: resourceType,
			ResourceID:   id,
			Status:       audit.GetStatusFromError(err),
			ErrorMessage: audit.GetErrorMessage(err),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read %s/%s: %w", resourceType, id, err)
	}

	// Parse to return as structured data
	var result map[string]interface{}
	if err := json.Unmarshal(resource, &result); err != nil {
		return nil, fmt.Errorf("failed to parse resource: %w", err)
	}

	return result, nil
}
