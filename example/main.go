package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertoAraneda/go-fhir-validator/pkg/v1"
	"log"
)

func main() {
	// Initialize JSON data in memory
	_, err := v1.LoadData()
	if err != nil {
		log.Fatalf("Error loading JSON data: %v", err)
	}

	data, err := v1.ReadJSONFile("resource.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	outcome, err := v1.ValidateResource(data)

	if err != nil {
		fmt.Println(err)
		return
	}

	outcomeJSON, _ := json.MarshalIndent(outcome, "", "  ")
	fmt.Println(string(outcomeJSON))
}
