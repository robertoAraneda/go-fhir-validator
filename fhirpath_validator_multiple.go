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

type ResponseFhirPathValidatorMultiple struct {
	IsValid bool     `json:"isValid"`
	Urls    []string `json:"urls,omitempty"`
	Ids     []string `json:"ids,omitempty"`
}

type ValidationResult struct {
	Result bool   `json:"result"`
	Key    string `json:"key"`
	Path   string `json:"path"`
	Human  string `json:"human"`
}

// TraceData estructura para almacenar los valores de TRACE
type TraceData struct {
	URL       []string `json:"url"`
	IDs       []string `json:"ids"`
	Unmatched []string `json:"unmatched"`
}

// cleanValues limpia comillas dobles y escapadas en los valores extraídos
func cleanValues(raw []string) []string {
	if len(raw) < 2 {
		return nil
	}

	// Remove unnecessary newlines and split by commas
	rawValues := strings.Split(strings.ReplaceAll(raw[1], "\n", ""), ",")
	var result []string
	for _, v := range rawValues {
		cleaned := strings.TrimSpace(v)
		cleaned = strings.Trim(cleaned, `"\`)
		result = append(result, cleaned)
	}
	return result
}

func FhirPathValidatorMultiple(array []FhirPathPayload) (*interface{}, error) {

	// Convert the FHIR resource to JSON
	resourceJSON, err := json.Marshal(array)
	if err != nil {
		log.Fatalf("Error converting resource to JSON: %v", err)
		return nil, err
	}

	path := filepath.Join("node/dist", "fhirpath-evaluate.js")

	// print with scaped characters
	fmt.Printf("%q\n", string(resourceJSON))

	// Step 3: Execute the Node.js script
	cmd := exec.Command("node", path, string(resourceJSON))

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

	// Regex to extract the JSON array inside "Result:"
	re := regexp.MustCompile(`Result:\s*(\[[\s\S]*\])`)
	matches := re.FindStringSubmatch(rawOutput)

	if len(matches) < 2 {
		fmt.Println("No 'Result' array found in input.")
	}

	resultJSON := matches[1]

	// Parse JSON into a slice of ValidationResult
	var results []ValidationResult
	err = json.Unmarshal([]byte(resultJSON), &results)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Filter only results where "result" is false
	var failedResults []ValidationResult
	for _, res := range results {
		if !res.Result {
			failedResults = append(failedResults, res)
		}
	}

	// Print filtered results
	fmt.Println("Extracted Failed Results:")
	for _, res := range failedResults {
		fmt.Printf("- Key: %s, Path: %s, Human: %s\n", res.Key, res.Path, res.Human)
	}

	trace := TraceData{}
	// Regex mejorada para capturar listas de valores correctamente
	tracePatterns := map[string]*regexp.Regexp{
		"url":       regexp.MustCompile(`(?m)TRACE:\[url\]\s*\[\s*([\s\S]*?)\s*\]`),
		"ids":       regexp.MustCompile(`(?m)TRACE:\[ids\]\s*\[\s*([\s\S]*?)\s*\]`),
		"unmatched": regexp.MustCompile(`(?m)TRACE:\[unmatched\]\s*\[\s*([\s\S]*?)\s*\]`),
	}

	// Extraer datos usando las expresiones regulares
	for key, pattern := range tracePatterns {
		traceMatches := pattern.FindStringSubmatch(rawOutput)
		if len(traceMatches) > 1 {
			// Intentar parsear como JSON array (corrige comillas escapadas)
			var values []string
			values = cleanValues(traceMatches)
			if err == nil {
				// Asignar valores al struct solo si se parsea bien
				switch key {
				case "url":
					trace.URL = values
				case "ids":
					trace.IDs = values
				case "unmatched":
					trace.Unmatched = values
				}
			}
		}
	}

	// Convertir a JSON para mostrarlo de forma estructurada
	traceJSON, _ := json.MarshalIndent(trace, "", "  ")
	fmt.Println("Extracted TRACE Data:")
	fmt.Println(string(traceJSON))

	/*
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

	*/

	return nil, nil
}
