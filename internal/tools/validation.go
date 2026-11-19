package tools

import (
	"fmt"
	"regexp"
)

// These values come from the agent (and so, indirectly, from the LLM) and are
// placed directly into the FHIR request path, e.g. {base}/{resourceType}/{id}.
// Without validation a value such as "1/../../Observation/2" or
// "Patient?_count=1000" changes which resource is requested, while the audit
// log still records the original resourceType and id.

// resourceTypePattern matches FHIR resource type names, which are PascalCase
// letters only (e.g. Patient, MedicationRequest).
var resourceTypePattern = regexp.MustCompile(`^[A-Z][A-Za-z]{0,63}$`)

// resourceIDPattern matches the FHIR "id" datatype: 1-64 characters from
// [A-Za-z0-9-.]. See https://hl7.org/fhir/R4/datatypes.html#id
var resourceIDPattern = regexp.MustCompile(`^[A-Za-z0-9\-.]{1,64}$`)

// validateResourceType checks that resourceType is a single, well-formed FHIR
// resource type name and cannot alter the request path.
func validateResourceType(resourceType string) error {
	if !resourceTypePattern.MatchString(resourceType) {
		return fmt.Errorf("invalid resourceType %q: must be a FHIR resource type name such as Patient or Observation", resourceType)
	}
	return nil
}

// validateResourceID checks that id is a valid FHIR logical id and cannot
// alter the request path.
func validateResourceID(id string) error {
	// "." and ".." match the character set but are path dot-segments.
	if !resourceIDPattern.MatchString(id) || id == "." || id == ".." {
		return fmt.Errorf("invalid id %q: must be 1-64 characters of letters, digits, '-' or '.'", id)
	}
	return nil
}
