package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ResponseFhirPathValidator struct {
	IsValid bool     `json:"isValid"`
	Urls    []string `json:"urls,omitempty"`
	Ids     []string `json:"ids,omitempty"`
}

func FhirPathValidator(rootResource map[string]interface{}, fhirResource interface{}, expression string) (*ResponseFhirPathValidator, error) {

	// Convert the FHIR resource to JSON
	resourceJSON, err := json.Marshal(fhirResource)
	if err != nil {
		log.Fatalf("Error converting resource to JSON: %v", err)
		return nil, err
	}

	// Convert the root resource to JSON
	rootResourceJSON, err := json.Marshal(rootResource)
	if err != nil {
		log.Fatalf("Error converting root resource to JSON: %v", err)
		return nil, err
	}

	path := filepath.Join("", "fhirpath.js")

	// Step 3: Execute the Node.js script
	cmd := exec.Command("node", path, string(resourceJSON), expression, string(rootResourceJSON))

	// Capture output
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Node.js error: %s\n", stderr.String())
		return nil, fmt.Errorf("execution error: %w", err)
	}

	// Debug: Print raw output
	rawOutput := out.String()

	fmt.Printf("Raw Output: %s\n", rawOutput)

	re := regexp.MustCompile(`\[(true|false)\]\s*$`)
	matches := re.FindString(rawOutput)

	// Parse JSON response
	var result []bool
	err = json.Unmarshal([]byte(matches), &result)
	if err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w\nRaw Output: %s", err, rawOutput)
	}

	// Regular expressions for TRACE:[url] and TRACE:[ids]
	urlRegex := regexp.MustCompile(`(?m)TRACE:\[url\]\s*\[\s*([\s\S]*?)\s*\]`)
	idsRegex := regexp.MustCompile(`(?m)TRACE:\[ids\]\s*\[\s*([\s\S]*?)\s*\]`)

	// Extract URL and IDs from the output
	urlMatches := urlRegex.FindStringSubmatch(rawOutput)
	idsMatches := idsRegex.FindStringSubmatch(rawOutput)

	// Helper function to clean up extracted values
	cleanValues := func(match []string) []string {
		if len(match) < 2 {
			return nil
		}

		// Remove unnecessary newlines and split by commas
		rawValues := strings.Split(strings.ReplaceAll(match[1], "\n", ""), ",")
		var result []string
		for _, v := range rawValues {
			cleaned := strings.TrimSpace(v)
			cleaned = strings.Trim(cleaned, `"\`)
			result = append(result, cleaned)
		}
		return result
	}

	urlValues := cleanValues(urlMatches)

	fmt.Printf("URL Values: %v\n", urlValues)
	idsValues := cleanValues(idsMatches)

	fmt.Printf("IDs Values: %v\n", idsValues)

	return &ResponseFhirPathValidator{
		IsValid: result[0],
		Urls:    urlValues,
		Ids:     idsValues,
	}, nil
}
