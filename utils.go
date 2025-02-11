package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// IsArrayElement determines if the given element should be treated as an array based on its max cardinality.
func IsArrayElement(element Element) bool {
	_, baseMaxIsWildcard := ParseMaxItems(element.Base.Max)
	maxItems, maxIsWildcard := ParseMaxItems(element.Max)

	return baseMaxIsWildcard || maxIsWildcard || maxItems > 1
}

// ParseMaxItems parses the 'max' value from the definition and checks if it is a wildcard.
func ParseMaxItems(max string) (int, bool) {
	if max == "*" {
		return 0, true // Wildcard means no upper limit
	}

	// Convert string to integer safely
	maxInt, err := strconv.Atoi(max)
	if err != nil {
		return 1, false // Default to max=1 if not explicitly defined
	}

	return maxInt, false
}

func ToStruct[T any](value interface{}) (*T, error) {
	var result T

	// If value is already of type T, return it directly
	if v, ok := value.(T); ok {
		return &v, nil
	}

	// Marshal into JSON and unmarshal into target type
	data, err := json.Marshal(value)
	if err != nil {
		return &result, fmt.Errorf("error marshaling value: %v", err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return &result, fmt.Errorf("error unmarshaling value: %v", err)
	}

	return &result, nil
}

// extractUnderscorePart removes everything before the first "_" and keeps the underscore-prefixed part
func extractUnderscorePart(id string) string {
	parts := strings.Split(id, ".")
	for _, part := range parts {
		if strings.HasPrefix(part, "_") {
			return part
		}
	}
	return "" // Return empty string if no underscore-prefixed part is found
}

// contains checks if a given resourceType exists in the list of valid FHIR R4 resource types.
func contains(resourceTypes []string, resourceType string) bool {
	for i := 0; i < len(resourceTypes); i++ {
		if resourceTypes[i] == resourceType {
			return true
		}
	}
	return false
}

func joinPath(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}

func IsMultipleType(element Element) bool {
	if strings.Contains(element.Path, "[x]") {
		return true
	}
	return false
}

// IsBackboneElement Determines if an element is of type BackboneElement
func IsBackboneElement(element Element) bool {

	if len(element.Type) == 0 {
		return false
	}

	for _, t := range element.Type {
		if t.Code == "BackboneElement" {
			return true
		}
	}
	return false
}

// IsNestedInBackbone checks if an element belongs to a backbone element.
func IsNestedInBackbone(elementPath string, backboneElementMap map[string]struct{}) bool {
	for id := range backboneElementMap {
		if strings.HasPrefix(elementPath, id+".") {
			return true
		}
	}
	return false
}
