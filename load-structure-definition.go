package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// LoadDefinition loads a StructureDefinition, ValueSet, or CodeSystem from a file
func LoadDefinition(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var rawData map[string]interface{}
	if err := json.Unmarshal(byteValue, &rawData); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	resourceType, ok := rawData["resourceType"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'resourceType' in JSON")
	}
	switch resourceType {
	case "StructureDefinition":
		var definition StructureDefinition
		if err := json.Unmarshal(byteValue, &definition); err != nil {
			return nil, fmt.Errorf("error parsing StructureDefinition: %v", err)
		}
		return definition, nil
	case "ValueSet":
		var valueSet ValueSet
		if err := json.Unmarshal(byteValue, &valueSet); err != nil {
			return nil, fmt.Errorf("error parsing ValueSet: %v", err)
		}
		return valueSet, nil
	case "CodeSystem":
		var codeSystem CodeSystem
		if err := json.Unmarshal(byteValue, &codeSystem); err != nil {
			return nil, fmt.Errorf("error parsing CodeSystem: %v", err)
		}
		return codeSystem, nil
	default:
		return nil, fmt.Errorf("unknown resourceType: %s", resourceType)
	}
}

// readJSONFile reads and parses a JSON file into a map
func readJSONFile(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(file)

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return data, nil
}

// LoadDefinitions loads all FHIR definitions from a directory and stores them in `storage`
func LoadDefinitions(dir string, storage *sync.Map) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Error reading directory '%s': %v\n", dir, err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		definition, err := LoadDefinition(filePath)
		if err != nil {
			fmt.Printf("Error loading structure definition from '%s': %v\n", file.Name(), err)
			continue
		}

		// Store the definition in `storage`
		StoreDefinition(storage, definition)
	}
}

// StoreDefinition stores the FHIR structure definition in sync.Map based on its type
func StoreDefinition(storage *sync.Map, definition interface{}) {
	var key string

	switch v := definition.(type) {
	case CodeSystem:
		key = v.URL

	case ValueSet:
		// key = fmt.Sprintf("%s|%s", v.URL, v.Version)
		key = v.URL

	case StructureDefinition:
		if v.Type == "Extension" && v.ID != "Extension" {
			key = v.URL
		} else {
			key = v.ID
		}

	default:
		fmt.Println("Error: Unsupported FHIR resource type")
		return
	}

	storage.Store(key, definition)
}
