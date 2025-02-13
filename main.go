package main

import (
	"encoding/json"
	"fmt"
	"github.com/robertoAraneda/go-fhir-validator/v1"
)

func main() {
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
