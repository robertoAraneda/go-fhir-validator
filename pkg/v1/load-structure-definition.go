package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// LibraryData holds the JSON content in memory
type LibraryData struct {
	Config map[string]interface{} `json:"config"`
}

var (
	once            sync.Once
	specLibraryData *LibraryData
)

func loadJSON(filePath string) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file %s: %v\n", filePath, err)
		}
	}(file)

	var content interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&content); err != nil {
		return nil, fmt.Errorf("failed to decode file %s: %w", filePath, err)
	}

	return content, nil
}

// LoadData loads JSON files into memory (singleton)
func LoadData() (*LibraryData, error) {
	var err error
	once.Do(func() {
		specLibraryData = &LibraryData{
			Config: make(map[string]interface{}),
		}

		err = filepath.WalkDir("spec", func(path string, d fs.DirEntry, walErr error) error {
			if walErr != nil {
				return walErr
			}

			if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
				fileName := filepath.Base(path)

				// Load JSON file
				jsonData, loadErr := loadJSON(path)
				if loadErr != nil {
					return loadErr
				}

				// access the JSON specLibraryData resourceType attribute
				switch jsonData.(type) {
				case map[string]interface{}:

					rawData := jsonData.(map[string]interface{})
					resourceType, ok := rawData["resourceType"].(string)
					if !ok {
						return fmt.Errorf("missing or invalid 'resourceType' in JSON")
					}

					fmt.Printf("ResourceType '%s\n'", resourceType)

					switch resourceType {
					case "StructureDefinition":
						var structureDef StructureDefinition
						jsonBytes, _ := json.Marshal(rawData) // Convert map to JSON
						if err := json.Unmarshal(jsonBytes, &structureDef); err != nil {
							return fmt.Errorf("failed to parse StructureDefinition: %v", err)
						}

						var key string
						if structureDef.Type == "Extension" && structureDef.ID != "Extension" {
							key = structureDef.URL
						} else {
							key = structureDef.ID
						}
						specLibraryData.Config[key] = structureDef

					case "ValueSet":

						var valueSet ValueSet
						jsonBytes, _ := json.Marshal(rawData) // Convert map to JSON
						if err := json.Unmarshal(jsonBytes, &valueSet); err != nil {
							return fmt.Errorf("failed to parse StructureDefinition: %v", err)
						}

						specLibraryData.Config[valueSet.URL] = valueSet

					case "CodeSystem":
						var codeSystem CodeSystem
						jsonBytes, _ := json.Marshal(rawData) // Convert map to JSON
						if err := json.Unmarshal(jsonBytes, &codeSystem); err != nil {
							return fmt.Errorf("failed to parse StructureDefinition: %v", err)
						}

						specLibraryData.Config[codeSystem.URL] = codeSystem
					default:
						return fmt.Errorf("unknown resourceType: %s", resourceType)
					}

				default:
					return fmt.Errorf("invalid JSON specLibraryData in file %s", fileName)
				}
			}

			return nil
		})

	})
	return specLibraryData, err
}

func GetSpec() (*LibraryData, error) {

	if specLibraryData == nil {
		return LoadData()
	}

	return specLibraryData, nil
}

func GetConfig() (map[string]interface{}, error) {
	d, err := GetSpec()
	if err != nil {
		return nil, err
	}

	return d.Config, nil
}

// ReadJSONFile reads and parses a JSON file into a map
func ReadJSONFile(filename string) (map[string]interface{}, error) {
	// print current path
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

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

	var specLibraryData map[string]interface{}
	err = json.Unmarshal(byteValue, &specLibraryData)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return specLibraryData, nil
}

/*
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


*/
