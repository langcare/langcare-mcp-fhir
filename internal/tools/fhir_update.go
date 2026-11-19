package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/langcare/langcare-mcp-fhir/internal/audit"
	"github.com/langcare/langcare-mcp-fhir/internal/fhir"
)

// FHIRUpdateTool implements the fhir_update tool
type FHIRUpdateTool struct {
	client      fhir.Client
	auditLogger *audit.AuditLogger
}

// NewFHIRUpdateTool creates a new FHIR update tool
func NewFHIRUpdateTool(client fhir.Client, auditLogger *audit.AuditLogger) *FHIRUpdateTool {
	return &FHIRUpdateTool{
		client:      client,
		auditLogger: auditLogger,
	}
}

// Name returns the tool name
func (t *FHIRUpdateTool) Name() string {
	return "fhir_update"
}

// Description returns the tool description
func (t *FHIRUpdateTool) Description() string {
	return "Update an existing FHIR resource"
}

// InputSchema returns the JSON schema for tool inputs
func (t *FHIRUpdateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"resourceType": map[string]interface{}{
				"type":        "string",
				"description": "The FHIR resource type to update (e.g., Patient, Observation, Medication)",
			},
			"id": map[string]interface{}{
				"type":        "string",
				"description": "The resource ID",
			},
			"resource": map[string]interface{}{
				"type":        "object",
				"description": "The updated FHIR resource data",
			},
		},
		"required": []string{"resourceType", "id", "resource"},
	}
}

// Execute executes the tool
func (t *FHIRUpdateTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	resourceData, ok := args["resource"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("resource is required and must be an object")
	}

	// Convert resource to JSON
	resourceJSON, err := json.Marshal(resourceData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	// Execute FHIR update
	updated, err := t.client.Update(ctx, resourceType, id, resourceJSON)

	// Log PHI access
	if t.auditLogger != nil {
		t.auditLogger.LogPHIAccess(ctx, audit.PHIAccessEvent{
			Operation:    "update",
			ResourceType: resourceType,
			ResourceID:   id,
			Status:       audit.GetStatusFromError(err),
			ErrorMessage: audit.GetErrorMessage(err),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update %s/%s: %w", resourceType, id, err)
	}

	// Parse to return as structured data
	var result map[string]interface{}
	if err := json.Unmarshal(updated, &result); err != nil {
		return nil, fmt.Errorf("failed to parse updated resource: %w", err)
	}

	return result, nil
}
