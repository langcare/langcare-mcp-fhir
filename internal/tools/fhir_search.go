package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/langcare/langcare-mcp-fhir/internal/audit"
	"github.com/langcare/langcare-mcp-fhir/internal/fhir"
)

// FHIRSearchTool implements the fhir_search tool
type FHIRSearchTool struct {
	client      fhir.Client
	auditLogger *audit.AuditLogger
}

// NewFHIRSearchTool creates a new FHIR search tool
func NewFHIRSearchTool(client fhir.Client, auditLogger *audit.AuditLogger) *FHIRSearchTool {
	return &FHIRSearchTool{
		client:      client,
		auditLogger: auditLogger,
	}
}

// Name returns the tool name
func (t *FHIRSearchTool) Name() string {
	return "fhir_search"
}

// Description returns the tool description
func (t *FHIRSearchTool) Description() string {
	return "Search FHIR resources with query parameters"
}

// InputSchema returns the JSON schema for tool inputs
func (t *FHIRSearchTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"resourceType": map[string]interface{}{
				"type":        "string",
				"description": "The FHIR resource type to search (e.g., Patient, Observation, Medication)",
			},
			"queryParams": map[string]interface{}{
				"type":        "string",
				"description": "FHIR search query parameters (e.g., 'name=John&birthdate=gt1990-01-01')",
			},
		},
		"required": []string{"resourceType"},
	}
}

// Execute executes the tool
func (t *FHIRSearchTool) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Parse arguments
	resourceType, ok := args["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("resourceType is required and must be a string")
	}
	if err := validateResourceType(resourceType); err != nil {
		return nil, err
	}

	queryParams := ""
	if qp, ok := args["queryParams"].(string); ok {
		queryParams = qp
	}

	// Execute FHIR search
	bundle, err := t.client.Search(ctx, resourceType, queryParams)

	// Log PHI access
	if t.auditLogger != nil {
		t.auditLogger.LogPHIAccess(ctx, audit.PHIAccessEvent{
			Operation:    "search",
			ResourceType: resourceType,
			QueryParams:  queryParams,
			Status:       audit.GetStatusFromError(err),
			ErrorMessage: audit.GetErrorMessage(err),
		})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search %s: %w", resourceType, err)
	}

	// Parse to return as structured data
	var result map[string]interface{}
	if err := json.Unmarshal(bundle, &result); err != nil {
		return nil, fmt.Errorf("failed to parse bundle: %w", err)
	}

	return result, nil
}
