package main

import "time"

type StructureDefinition struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id"`
	Text         Text         `json:"text"`
	Extension    []Extension  `json:"extension"`
	URL          string       `json:"url"`
	Version      string       `json:"version"`
	Name         string       `json:"name"`
	Status       string       `json:"status"`
	Date         time.Time    `json:"date"`
	Publisher    string       `json:"publisher"`
	Contact      []Contact    `json:"contact"`
	Description  string       `json:"description"`
	FHIRVersion  string       `json:"fhirVersion"`
	Mapping      []Mapping    `json:"mapping"`
	Kind         string       `json:"kind"`
	Abstract     bool         `json:"abstract"`
	Type         string       `json:"type"`
	Snapshot     *Snapshot    `json:"snapshot"`
	Differential Differential `json:"differential"`
}

// ValueSet representa un ValueSet de FHIR
type ValueSet struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id"`
	Meta         Meta         `json:"meta"`
	Text         Narrative    `json:"text"`
	Extension    []Extension  `json:"extension"`
	URL          string       `json:"url"`
	Identifier   []Identifier `json:"identifier"`
	Version      string       `json:"version"`
	Name         string       `json:"name"`
	Title        string       `json:"title"`
	Status       string       `json:"status"`
	Experimental bool         `json:"experimental"`
	Date         time.Time    `json:"date"`
	Publisher    string       `json:"publisher"`
	Contact      []Contact    `json:"contact"`
	Description  string       `json:"description"`
	Immutable    bool         `json:"immutable"`
	Compose      Compose      `json:"compose"`
}

// CodeSystem representa un sistema de códigos en FHIR
type CodeSystem struct {
	ResourceType  string       `json:"resourceType"`
	ID            string       `json:"id"`
	Meta          Meta         `json:"meta"`
	Text          Narrative    `json:"text"`
	Extension     []Extension  `json:"extension"`
	URL           string       `json:"url"`
	Identifier    []Identifier `json:"identifier"`
	Version       string       `json:"version"`
	Name          string       `json:"name"`
	Title         string       `json:"title"`
	Status        string       `json:"status"`
	Experimental  bool         `json:"experimental"`
	Date          time.Time    `json:"date"`
	Publisher     string       `json:"publisher"`
	Contact       []Contact    `json:"contact"`
	Description   string       `json:"description"`
	CaseSensitive bool         `json:"caseSensitive"`
	ValueSet      string       `json:"valueSet"`
	Content       string       `json:"content"`
	Concept       []Concept    `json:"concept"`
}

// OperationOutcome represents a FHIR OperationOutcome resource
type OperationOutcome struct {
	ResourceType string       `json:"resourceType"`
	Issue        []IssueEntry `json:"issue"`
}

// IssueEntry represents an individual issue in OperationOutcome
type IssueEntry struct {
	Severity    string           `json:"severity"` // "error", "warning", etc.
	Code        string           `json:"code"`     // "invalid", "required", etc.
	Details     *CodeableConcept `json:"details,omitempty"`
	Diagnostics string           `json:"diagnostics,omitempty"` // Additional diagnostic information
	Expression  []string         `json:"expression,omitempty"`  // FHIRPath expression
	Location    string           `json:"location,omitempty"`    // Path to the field causing the issue
}

// Concept representa un concepto en el CodeSystem
type Concept struct {
	Code       string    `json:"code"`
	Display    string    `json:"display"`
	Definition string    `json:"definition"`
	Concept    []Concept `json:"concept,omitempty"`
}

// Meta representa la metadata del ValueSet
type Meta struct {
	LastUpdated time.Time `json:"lastUpdated"`
	Profile     []string  `json:"profile"`
}

// Narrative representa el texto narrativo generado para el recurso
type Narrative struct {
	Status string `json:"status"`
	Div    string `json:"div"`
}

// Extension representa las extensiones en FHIR
type Extension struct {
	URL         string `json:"url"`
	ValueCode   string `json:"valueCode,omitempty"`
	ValueInt    int    `json:"valueInteger,omitempty"`
	ValueString string `json:"valueString,omitempty"`
	ValueUrl    string `json:"valueUrl,omitempty"`
}

// Identifier representa un identificador del recurso
type Identifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

// Contact representa la información de contacto del publisher
type Contact struct {
	Telecom []Telecom `json:"telecom"`
}

// Telecom representa información de telecomunicación (URL o email)
type Telecom struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

// Compose representa la composición del ValueSet
type Compose struct {
	Include []Include `json:"include"`
}

// Include representa los sistemas de código incluidos en el ValueSet
type Include struct {
	System  string    `json:"system"`
	Concept []Concept `json:"concept"`
	Filter  []struct {
		Property string `json:"property"`
		Op       string `json:"op"`
		Value    string `json:"value"`
	} `json:"filter"`
}

type Text struct {
	Status string `json:"status"`
	Div    string `json:"div"`
}

type Mapping struct {
	Identity string `json:"identity"`
	URI      string `json:"uri"`
	Name     string `json:"name"`
}

type Snapshot struct {
	Element []Element `json:"element"`
}

type Differential struct {
	Element []Element `json:"element"`
}

type Element struct {
	ID         string           `json:"id"`
	Extension  []Extension      `json:"extension,omitempty"`
	Path       string           `json:"path"`
	Short      string           `json:"short"`
	Definition string           `json:"definition"`
	Min        int              `json:"min"`
	Max        string           `json:"max"`
	Base       Base             `json:"base"`
	Type       []Type           `json:"type"`
	Condition  []string         `json:"condition,omitempty"`
	Constraint []Constraint     `json:"constraint,omitempty"`
	IsModifier bool             `json:"isModifier"`
	IsSummary  bool             `json:"isSummary"`
	Binding    *Binding         `json:"binding,omitempty"`
	Mapping    []ElementMapping `json:"mapping,omitempty"`
}

type Binding struct {
	Description string      `json:"description"`
	Strength    string      `json:"strength"`
	ValueSet    string      `json:"valueSet"`
	Extension   []Extension `json:"extension,omitempty"`
}

type Type struct {
	Code          string      `json:"code"`
	Extension     []Extension `json:"extension,omitempty"`
	TargetProfile []string    `json:"targetProfile,omitempty"`
}

type Base struct {
	Path string `json:"path"`
	Min  int    `json:"min"`
	Max  string `json:"max"`
}

type Constraint struct {
	Key        string `json:"key"`
	Severity   string `json:"severity"`
	Human      string `json:"human"`
	Expression string `json:"expression"`
	XPath      string `json:"xpath"`
	Source     string `json:"source"`
}

type ElementMapping struct {
	Identity string `json:"identity"`
	Map      string `json:"map"`
}

type Coding struct {
	System  string `json:"system,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}
