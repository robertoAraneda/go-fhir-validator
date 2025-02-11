package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// FhirR4ResourceTypes contains all resource types in FHIR R4.
var FhirR4ResourceTypes = []string{
	"Account", "ActivityDefinition", "AdverseEvent", "AllergyIntolerance", "Appointment",
	"AppointmentResponse", "AuditEvent", "Basic", "Binary", "BiologicallyDerivedProduct",
	"BodyStructure", "Bundle", "CapabilityStatement", "CarePlan", "CareTeam",
	"CatalogEntry", "ChargeItem", "ChargeItemDefinition", "Claim", "ClaimResponse",
	"ClinicalImpression", "CodeSystem", "Communication", "CommunicationRequest", "CompartmentDefinition",
	"Composition", "ConceptMap", "Condition", "Consent", "Contract",
	"Coverage", "CoverageEligibilityRequest", "CoverageEligibilityResponse", "DetectedIssue", "Device",
	"DeviceDefinition", "DeviceMetric", "DeviceRequest", "DeviceUseStatement", "DiagnosticReport",
	"DocumentManifest", "DocumentReference", "EffectEvidenceSynthesis", "Encounter", "Endpoint",
	"EnrollmentRequest", "EnrollmentResponse", "EpisodeOfCare", "EventDefinition", "Evidence",
	"EvidenceVariable", "ExampleScenario", "ExplanationOfBenefit", "FamilyMemberHistory", "Flag",
	"Goal", "GraphDefinition", "Group", "GuidanceResponse", "HealthcareService",
	"ImagingStudy", "Immunization", "ImmunizationEvaluation", "ImmunizationRecommendation", "ImplementationGuide",
	"InsurancePlan", "Invoice", "Library", "Linkage", "List",
	"Location", "Measure", "MeasureReport", "Media", "Medication",
	"MedicationAdministration", "MedicationDispense", "MedicationKnowledge", "MedicationRequest", "MedicationStatement",
	"MedicinalProduct", "MedicinalProductAuthorization", "MedicinalProductContraindication", "MedicinalProductIndication", "MedicinalProductIngredient",
	"MedicinalProductInteraction", "MedicinalProductManufactured", "MedicinalProductPackaged", "MedicinalProductPharmaceutical", "MedicinalProductUndesirableEffect",
	"MessageDefinition", "MessageHeader", "MolecularSequence", "NamingSystem", "NutritionOrder",
	"Observation", "ObservationDefinition", "OperationDefinition", "OperationOutcome", "Organization",
	"OrganizationAffiliation", "Parameters", "Patient", "PaymentNotice", "PaymentReconciliation",
	"Person", "PlanDefinition", "Practitioner", "PractitionerRole", "Procedure",
	"Provenance", "Questionnaire", "QuestionnaireResponse", "RelatedPerson", "RequestGroup",
	"ResearchDefinition", "ResearchElementDefinition", "ResearchStudy", "ResearchSubject", "RiskAssessment",
	"RiskEvidenceSynthesis", "Schedule", "SearchParameter", "ServiceRequest", "Slot",
	"Specimen", "SpecimenDefinition", "StructureDefinition", "StructureMap", "Subscription",
	"Substance", "SubstanceNucleicAcid", "SubstancePolymer", "SubstanceProtein", "SubstanceReferenceInformation",
	"SubstanceSourceMaterial", "SubstanceSpecification", "SupplyDelivery", "SupplyRequest", "Task",
	"TerminologyCapabilities", "TestReport", "TestScript", "ValueSet", "VerificationResult",
	"VisionPrescription",
}

// addError appends an issue to the OperationOutcome, with optional details about the incoming value type
func addError(outcome *OperationOutcome, code, diagnostic, location, details string) {
	issue := IssueEntry{
		Severity:    "error",
		Code:        code,
		Diagnostics: diagnostic,
	}
	if location != "" {
		issue.Expression = []string{location}
	}

	if details != "" {
		issue.Details = CodeableConcept{
			Text: details,
		}
	}
	outcome.Issue = append(outcome.Issue, issue)
}

// validateElements applies a validation function to multiple elements.
func validateElements(
	rootData map[string]interface{},
	data map[string]interface{},
	elements []Element,
	rootSpec StructureDefinition,
	spec StructureDefinition,
	parentPath string,
	outcome *OperationOutcome,
	validator func(map[string]interface{}, map[string]interface{}, Element, StructureDefinition, StructureDefinition, string, *OperationOutcome),
) {
	if len(elements) == 0 {
		return // Exit early if there are no elements
	}

	// Use index-based iteration to prevent range aliasing issues
	for i := 0; i < len(elements); i++ {
		if elements[i].Path == "" {
			addError(outcome, "invalid", "Element has an empty path", parentPath, "")
			continue
		}

		validator(rootData, data, elements[i], rootSpec, spec, parentPath, outcome)
	}
}

func ValidateResource(data map[string]interface{}) (*OperationOutcome, error) {
	outcome := &OperationOutcome{ResourceType: "OperationOutcome"}

	// extract the resource type
	resourceType, ok := data["resourceType"].(string)

	if !ok {
		return nil, fmt.Errorf("resource type not found")
	}

	// check if the resource type is valid
	if !contains(FhirR4ResourceTypes, resourceType) {
		// stringify the resource types
		resourceTypes := strings.Join(FhirR4ResourceTypes, ", ")
		addError(outcome, "error", fmt.Sprintf("Invalid resource type '%s'. Expected one of: %s", resourceType, resourceTypes), "", "")
		return outcome, nil
	}

	spec, ok := resources.Load(resourceType)
	if !ok {
		return nil, fmt.Errorf("resource type '%s' not found in definitions", resourceType)
	}

	Validate(data, data, spec.(StructureDefinition), spec.(StructureDefinition), resourceType, outcome)

	// create a file with the payload
	payloadJSON, _ := json.MarshalIndent(payload, "", "  ")
	err := os.WriteFile("payload.json", payloadJSON, 0644)
	if err != nil {
		fmt.Printf("Error writing payload file %s\n", err)
	}

	//response, err := FhirPathValidator(rootData, data, constraint.Expression)

	_, err = FhirPathValidatorMultiple(payload)

	if err != nil {
		fmt.Printf("Error validating constraint %s\n", err)
	}

	if len(outcome.Issue) == 0 {
		addError(outcome, "information", "Validation successful", "", "")
	}

	return outcome, nil
}

type FhirPathPayload struct {
	RootData             map[string]interface{} `json:"rootData"`
	Data                 map[string]interface{} `json:"data"`
	ConstraintExpression string                 `json:"constraintExpression"`
	ConstraintKey        string                 `json:"constraintKey"`
	ConstraintHuman      string                 `json:"constraintHuman"`
	ParentPath           string                 `json:"parentPath"`
}

var payload []FhirPathPayload

func Validate(rootData map[string]interface{}, data map[string]interface{}, rootSpec StructureDefinition, spec StructureDefinition, parentPath string, outcome *OperationOutcome) {

	fmt.Printf("Validating %s, %s, from %s\n", rootSpec.ID, spec.ID, parentPath)

	if len(data) == 0 {
		return // Exit early if no data to validate
	}

	elements := spec.Snapshot.Element
	backboneElementLookup := make(map[string]struct{}) // Efficient lookup for backbone elements
	topLevelElements := make([]Element, 0, len(elements))
	nestedBackboneElements := make([]Element, 0, len(elements))
	elementsWithVariableTypes := make([]Element, 0, len(elements))

	// Categorize elements into separate groups
	CategorizeElements(elements, backboneElementLookup, &topLevelElements, &nestedBackboneElements, &elementsWithVariableTypes)

	// Validate each category separately
	validateElements(rootData, data, topLevelElements, rootSpec, spec, parentPath, outcome, ValidateElement)
	validateElements(rootData, data, nestedBackboneElements, rootSpec, spec, parentPath, outcome, ValidateBackboneElement)
	validateElements(rootData, data, elementsWithVariableTypes, rootSpec, spec, parentPath, outcome, ValidateElementWithMultipleTypes)
	// ValidateUnderscoreFields(rootData, data, parentPath, rootSpec, spec, outcome)

	// find the constraints in the spec.Snapshot.Element when id is equal to spec.ID

	ele, matchingError := findMatchingElement(spec)
	if matchingError != nil {
		fmt.Printf("Error finding matching element %s\n", matchingError)
		return
	}

	if ele.Constraint == nil || len(ele.Constraint) == 0 {
		fmt.Printf("No constraints found for element: %s\n", ele.ID)
		return
	}

	for i := 0; i < len(ele.Constraint); i++ {
		constraint := ele.Constraint[i]
		var excludeFields = []string{"ele-1", "txt-1", "txt-2"}

		if contains(excludeFields, constraint.Key) {
			continue
		}

		payload = append(payload, FhirPathPayload{
			RootData:             rootData,
			Data:                 data,
			ConstraintExpression: constraint.Expression,
			ConstraintKey:        constraint.Key,
			ConstraintHuman:      constraint.Human,
			ParentPath:           parentPath,
		})

		//response, err := FhirPathValidator(rootData, data, constraint.Expression)

		/*
			response, err := FhirPathValidatorMultiple(payload)

			if err != nil {
				fmt.Printf("Error validating constraint %s\n", err)
				continue
			}

			if !response.IsValid {
				addError(outcome, "invalid", fmt.Sprintf("Field '%s' does not meet constraint '%s'", parentPath, constraint.Key), constraint.Expression, constraint.Human)

				if len(response.Urls) > 0 || len(response.Ids) > 0 {
					fmt.Printf("Constraint '%s' failed for field '%s'. URLs: %v, IDs: %v\n", constraint.Key, parentPath, response.Urls, response.Ids)
					addError(outcome, "invalid", fmt.Sprintf("Constraint '%s' failed for field '%s'. URLs: %v, IDs: %v", constraint.Key, parentPath, response.Urls, response.Ids), constraint.Expression, "")
				}
			}

		*/
	}
}

// CategorizeElements classifies elements into different types.
func CategorizeElements(
	elements []Element,
	backboneElementMap map[string]struct{},
	topLevelElements *[]Element,
	nestedBackboneElements *[]Element,
	multiTypeElements *[]Element,
) {
	for i := range elements {
		element := &elements[i]

		switch {
		case IsBackboneElement(*element):
			backboneElementMap[element.ID] = struct{}{}
		case IsMultipleType(*element):
			*multiTypeElements = append(*multiTypeElements, *element)
		case IsNestedInBackbone(element.Path, backboneElementMap):
			*nestedBackboneElements = append(*nestedBackboneElements, *element)
		default:
			if !strings.Contains(element.Path, "[x]") {
				*topLevelElements = append(*topLevelElements, *element)
			}
		}
	}
}

func ValidateElementWithMultipleTypes(rootData map[string]interface{}, data map[string]interface{}, element Element, rootSpec StructureDefinition, spec StructureDefinition, parentPath string, outcome *OperationOutcome) {
	fmt.Printf("Validating element with multiple types %s\n", element.Path)
}

// ValidateBackboneElement validates a single element against the specification
func ValidateBackboneElement(rootData map[string]interface{}, data map[string]interface{}, element Element, rootSpec StructureDefinition, spec StructureDefinition, parentPath string, outcome *OperationOutcome) {
	fmt.Printf("Validating backbone element %s\n", element.Path)
}

func findMatchingElement(spec StructureDefinition) (*Element, error) {
	for i := 0; i < len(spec.Snapshot.Element); i++ {
		element := spec.Snapshot.Element[i]
		if element.ID == spec.ID {
			return &element, nil
		}
	}

	return nil, fmt.Errorf("Element not found")
}

// ValidateElement validates a single element against the specification
func ValidateElement(rootData map[string]interface{}, data map[string]interface{}, element Element, rootSpec StructureDefinition, spec StructureDefinition, parentPath string, outcome *OperationOutcome) {

	fieldName := strings.TrimPrefix(element.Path, spec.ID+".")

	fullPath := joinPath(parentPath, fieldName)
	if underscorePart := extractUnderscorePart(element.ID); underscorePart != "" {
		fullPath = joinPath(parentPath, underscorePart)
	}

	childData := data[fieldName]

	if childData == nil {
		if element.Min > 0 {
			addError(outcome, "required", fmt.Sprintf("Field '%s' is required", fullPath), fullPath, "")
		}
		return
	}

	ValidateField(rootData, childData, element, fullPath, rootSpec, spec, outcome, false)

}

// ValidateField validates a single field against the specification
func ValidateField(rootData map[string]interface{}, value interface{}, element Element, fullPath string, rootSpec StructureDefinition, spec StructureDefinition, outcome *OperationOutcome, comingFromArray bool) {

	if IsArrayElement(element) && !comingFromArray {
		// Validate each element in the array
		ValidateArray(rootData, value, element, fullPath, rootSpec, spec, outcome)
	} else {
		// Validate the single value
		ValidateValue(rootData, value, element, fullPath, rootSpec, spec, outcome)
	}
}

// ValidateArray validates an array field against the specification
func ValidateArray(
	rootData map[string]interface{},
	value interface{},
	element Element,
	fullPath string,
	rootSpec StructureDefinition,
	spec StructureDefinition,
	outcome *OperationOutcome,
) {

	// Ensure the value is an array
	array, ok := value.([]interface{})
	if !ok {
		addError(outcome, "invalid", fmt.Sprintf("Field '%s' must be an array", fullPath), fullPath, "")
		return
	}

	length := len(array)

	// Validate minItems
	if length < element.Min {
		addError(outcome, "required", fmt.Sprintf("Field '%s' has too few items: minimum is %d. Found %d elements", fullPath, element.Min, length), fullPath, "")
	}

	// Validate maxItems
	maxItems, isUnlimited := ParseMaxItems(element.Max)
	if !isUnlimited && length > maxItems {
		addError(outcome, "invalid", fmt.Sprintf("Field '%s' has too many items: maximum is %d. Found %d elements", fullPath, maxItems, length), fullPath, "")
	}

	// Validate each element in the array
	for i, item := range array {
		itemPath := fmt.Sprintf("%s[%d]", fullPath, i)

		// check if item is a string

		switch v := item.(type) {
		case string, map[string]interface{}:
			ValidateField(rootData, v, element, itemPath, rootSpec, spec, outcome, true)
		default:
			addError(outcome, "invalid", fmt.Sprintf("Field '%s' must be a single value", fullPath), fullPath, "")
		}
	}
}

// ValidateValue validates a single value against the specification
func ValidateValue(
	rootData map[string]interface{},
	value interface{},
	element Element,
	fullPath string,
	rootSpec StructureDefinition,
	spec StructureDefinition,
	outcome *OperationOutcome,
) {

	if value == nil {
		addError(outcome, "invalid", fmt.Sprintf("All children of '%s' must be present", fullPath), fullPath, "")
		return
	}

	switch v := value.(type) {
	case []interface{}:
		addError(outcome, "invalid", fmt.Sprintf("Field '%s' must be a single value", fullPath), fullPath, "")
	case map[string]interface{}:
		// Validate nested object
		ValidateComplexType(rootData, v, element.Type[0].Code, fullPath, rootSpec, spec, outcome)
	case string:
		// Validate primitive type
		ValidatePrimitiveType(v, element.Type[0].Code, fullPath, rootSpec, spec, outcome)
	default:
		// Additional type checks can be added here if needed
	}
}

// ValidateComplexType validates nested complex types like Address, Organization, etc.
// It handles both single objects and slices of complex types.
func ValidateComplexType(rootData map[string]interface{}, value interface{}, typeCode, path string, rootSpec, spec StructureDefinition, outcome *OperationOutcome) {

	// Load the structure definition for the type
	nestedSpec, found := definitions.Load(typeCode)
	if !found {
		addError(outcome, "invalid", fmt.Sprintf("No structure definition found for type '%s'", typeCode), path, "")
		return
	}

	// Ensure correct type assertion
	specDefinition, valid := nestedSpec.(StructureDefinition)
	if !valid {
		addError(outcome, "invalid", fmt.Sprintf("Invalid structure definition for type '%s'", typeCode), path, "")
		return
	}

	Validate(rootData, value.(map[string]interface{}), rootSpec, specDefinition, path, outcome)
}

// ValidatePrimitiveType validates a FHIR primitive type against its expected regex pattern.
func ValidatePrimitiveType(value string, typeCode, path string, rootSpec StructureDefinition, spec StructureDefinition, outcome *OperationOutcome) {

	if typeCode == "http://hl7.org/fhirpath/System.String" {
		typeCode = "string" // Normalize FHIRPath string type
	}

	definition, found := definitions.Load(typeCode)
	if !found {
		addError(outcome, "invalid", fmt.Sprintf("No definition found for type '%s'", typeCode), path, "")
		return
	}

	// Extract the value element definition from the snapshot
	var valueElement *Element
	if valueElement = ExtractValueElementID(definition.(StructureDefinition).ID, definition.(StructureDefinition).Snapshot); valueElement == nil {
		addError(outcome, "invalid", fmt.Sprintf("No value element found for '%s'", path), path, "")
		return
	}

	// Extract regex pattern
	regex := ExtractRegexFromElement(*valueElement)

	// If no regex is found, try to infer from FHIR extensions
	if regex == "" {
		fhirType := getFhirTypeFromExtensions(valueElement.Type)
		if fhirType == "string" {
			regex = "[ \\r\\n\\t\\S]+" // Default regex for string
		} else {
			addError(outcome, "invalid", fmt.Sprintf("No regex pattern found for '%s'", path), path, "")
			return
		}
	}

	// Validate value against regex pattern
	ValidateRegex(value, regex, path, outcome)
}

func ValidateRegex(value string, regex string, path string, outcome *OperationOutcome) {
	// Perform regex validation
	re := regexp.MustCompile("^" + regex + "$") // Add start and end anchors
	strValue := fmt.Sprintf("%v", value)

	if !re.MatchString(strValue) {
		addError(outcome, "invalid", fmt.Sprintf("Field '%s' does not match the expected pattern: %s", path, regex), path, "")
	}

}

// ExtractRegexFromElement extracts the regex pattern from an Element object.
func ExtractRegexFromElement(element Element) string {
	if len(element.Type) == 0 {
		return "" // No types available, return empty regex
	}

	for i := 0; i < len(element.Type); i++ {
		for j := 0; j < len(element.Type[i].Extension); j++ {
			if element.Type[i].Extension[j].URL == "http://hl7.org/fhir/StructureDefinition/regex" {
				return element.Type[i].Extension[j].ValueString // Return first found regex
			}
		}
	}

	return "" // No regex found
}

// ExtractValueElementID searches for an element with ID `rootId + ".value"` in the snapshot.
func ExtractValueElementID(rootId string, snapshot *Snapshot) *Element {

	if snapshot == nil || len(snapshot.Element) == 0 {
		fmt.Printf("⚠️ Warning: Snapshot is empty or nil for rootId '%s'\n", rootId)
		return nil // Return empty element if snapshot is missing
	}

	targetID := rootId + ".value"

	for i := 0; i < len(snapshot.Element); i++ {
		if snapshot.Element[i].ID == targetID {
			return &snapshot.Element[i] // Return first matching element
		}
	}

	fmt.Printf("⚠️ Warning: No value element found for rootId '%s'\n", rootId)
	return nil // Return empty struct if no match is found
}

// getFHIRTypeFromExtensions extracts the FHIR type from element extensions
func getFhirTypeFromExtensions(types []Type) string {
	for i := 0; i < len(types); i++ {
		for j := 0; j < len(types[i].Extension); j++ {
			if types[i].Extension[j].URL == "http://hl7.org/fhir/StructureDefinition/structuredefinition-fhir-type" {
				return types[i].Extension[j].ValueUrl
			}
		}
	}
	return ""
}
