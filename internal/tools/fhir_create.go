package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/langcare/langcare-mcp-fhir/internal/audit"
	"github.com/langcare/langcare-mcp-fhir/internal/fhir"
)

// FHIRCreateTool implements the fhir_create tool
type FHIRCreateTool struct {
	client      fhir.Client
	auditLogger *audit.AuditLogger
}

// NewFHIRCreateTool creates a new FHIR create tool
func NewFHIRCreateTool(client fhir.Client, auditLogger *audit.AuditLogger) *FHIRCreateTool {
	return &FHIRCreateTool{
		client:      client,
		auditLogger: auditLogger,
	}
}

// Name returns the tool name
func (t *FHIRCreateTool) Name() string {
	return "fhir_create"
}

// Description returns the tool description
func (t *FHIRCreateTool) Description() string {
	return "Create a new FHIR resource"
}

// InputSchema returns the JSON schema for tool inputs
func (t *FHIRCreateTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"resourceType": map[string]interface{}{
				"type":        "string",
				"description": "The FHIR resource type to create (e.g., Patient, Observation, Medication)",
			},
			"resource": map[string]interface{}{
				"type":        "object",
				"description": "The FHIR resource data to create",
			},
		},
		"required": []string{"resourceType", "resource"},
	}
}

// Execute executes the tool
func (t *FHIRCreateTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	resourceType, ok := args["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("resourceType is required and must be a string")
	}
	if err := validateResourceType(resourceType); err != nil {
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

	// Execute FHIR create
	created, err := t.client.Create(ctx, resourceType, resourceJSON)

	// Log PHI access
	if t.auditLogger != nil {
		t.auditLogger.LogPHIAccess(ctx, audit.PHIAccessEvent{
			Operation:    "create",
			ResourceType: resourceType,
			Status:       audit.GetStatusFromError(err),
			ErrorMessage: audit.GetErrorMessage(err),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %s: %w", resourceType, err)
	}

	// Parse to return as structured data
	var result map[string]interface{}
	if err := json.Unmarshal(created, &result); err != nil {
		return nil, fmt.Errorf("failed to parse created resource: %w", err)
	}

	return result, nil
}
