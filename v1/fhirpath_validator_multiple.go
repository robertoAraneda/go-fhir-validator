package v1

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
	Result   bool   `json:"result"`
	Key      string `json:"key"`
	Path     string `json:"path"`
	Human    string `json:"human"`
	Source   string `json:"source"`
	Severity string `json:"severity"`
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

func FhirPathValidatorMultiple(array []FhirPathPayload) (*[]ValidationResult, *TraceData, error) {

	// Convert the FHIR resource to JSON
	resourceJSON, err := json.Marshal(array)
	if err != nil {
		log.Fatalf("Error converting resource to JSON: %v", err)
		return nil, nil, err
	}

	path := filepath.Join("v1/node/dist", "fhirpath-evaluate.js")

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
		return nil, nil, fmt.Errorf("execution error: %w", err)
	}

	// Debug: Print raw output
	rawOutput := out.String()

	fmt.Printf("Raw Output: %s\n", rawOutput)

	// Regex to extract the JSON array inside "Result:"
	re := regexp.MustCompile(`Result:\s*(\[[\s\S]*\])`)
	matches := re.FindStringSubmatch(rawOutput)

	if len(matches) < 2 {
		fmt.Println("No 'Result' array found in input.")
		return nil, nil, fmt.Errorf("no 'Result' array found in input")
	}

	resultJSON := matches[1]

	// Parse JSON into a slice of ValidationResult
	var results []ValidationResult
	err = json.Unmarshal([]byte(resultJSON), &results)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, nil, fmt.Errorf("error parsing JSON: %w", err)
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

	return &failedResults, &trace, nil
}
