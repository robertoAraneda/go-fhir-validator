package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

var datatypes sync.Map
var resources sync.Map
var codesystems sync.Map
var valuesets sync.Map
var extensions sync.Map
var definitions sync.Map

func init() {
	// List of directory-storage mappings
	mappings := []struct {
		dir     string
		storage *sync.Map
	}{
		{"datatypes", &definitions},
		{"datatypes", &datatypes},
		{"resources", &definitions},
		{"resources", &resources},
		{"extensions", &definitions},
		{"extensions", &extensions},
		{"codesystems", &codesystems},
		{"valuesets", &valuesets},
	}

	// Load definitions for each mapping
	for _, mapping := range mappings {
		LoadDefinitions(mapping.dir, mapping.storage)
	}
}

func main() {
	data, err := readJSONFile("resource.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	outcome, err := ValidateResource(data)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Validation failed:")
	outcomeJSON, _ := json.MarshalIndent(outcome, "", "  ")
	fmt.Println(string(outcomeJSON))
}
