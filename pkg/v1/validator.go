package v1

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

var payload []*FhirPathPayload

// addOperationOutcome appends an issue to the OperationOutcome, with optional details about the incoming value type
func addOperationOutcome(outcome *OperationOutcome, code, diagnostic, location, details, severity string) {
	issue := IssueEntry{
		Severity:    severity,
		Code:        code,
		Diagnostics: diagnostic,
	}
	if location != "" {
		issue.Expression = []string{location}
	}

	if details != "" {
		issue.Details = &CodeableConcept{
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
			addOperationOutcome(outcome, "invalid", "Element has an empty path", parentPath, "Path is empty", "error")
			continue
		}

		validator(rootData, data, elements[i], rootSpec, spec, parentPath, outcome)
	}
}

func ValidateResource(data map[string]interface{}) (*OperationOutcome, error) {

	outcome := &OperationOutcome{ResourceType: "OperationOutcome"}
	payload = []*FhirPathPayload{}

	// extract the resource type
	resourceType, ok := data["resourceType"].(string)

	if !ok {
		return nil, fmt.Errorf("resource type not found")
	}

	// check if the resource type is valid
	if !contains(FhirR4ResourceTypes, resourceType) {
		// stringify the resource types
		resourceTypes := strings.Join(FhirR4ResourceTypes, ", ")
		addOperationOutcome(outcome, "error", fmt.Sprintf("Invalid resource type '%s'. Expected one of: %s", resourceType, resourceTypes), "", "Invalid resource type", "error")
		return outcome, nil
	}

	spec, ok := specLibraryData.Config[resourceType]
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

	//response, err := FhirPathValidator(rootData, specLibraryData, constraint.Expression)

	results, trace, err := FhirPathValidatorMultiple(payload)

	if err != nil {
		fmt.Printf("Error validating constraint %s\n", err)
		addOperationOutcome(outcome, "exception", fmt.Sprintf("Error validating constraint %s", err), "", "", "fatal")
		return outcome, nil
	}

	fmt.Printf("Results: %v\n", results)
	fmt.Printf("Trace: %v\n", trace)

	if len(outcome.Issue) == 0 && len(*results) == 0 {
		addOperationOutcome(outcome, "information", "Validation successful", "", "", "information")
	} else {
		for i := 0; i < len(*results); i++ {
			diagnostics := fmt.Sprintf("Failed constraint '%s'", (*results)[i].Key)
			code := "invariant"
			details := fmt.Sprintf("%s: %s", (*results)[i].Key, (*results)[i].Human)
			if (*results)[i].Source != "" {
				diagnostics = fmt.Sprintf("Failed constraint '%s' (source: %s)", (*results)[i].Key, (*results)[i].Source)
			}

			if (*results)[i].Key == "dom-6" {
				diagnostics = fmt.Sprintf("Failed constraint '%s' (source: %s)", (*results)[i].Key, (*results)[i].Source)
				code = "informational"

			}
			// generate better response "Failed constraint '%s'" if (*results)[i].Source is present, use it
			addOperationOutcome(outcome, code, diagnostics, (*results)[i].Path, details, (*results)[i].Severity)
		}
	}

	return outcome, nil
}

type FhirPathPayload struct {
	RootData             map[string]interface{} `json:"rootData"`
	Data                 map[string]interface{} `json:"data"`
	Path                 string                 `json:"path"`
	ConstraintExpression string                 `json:"constraintExpression"`
	ConstraintKey        string                 `json:"constraintKey"`
	ConstraintHuman      string                 `json:"constraintHuman"`
	ConstraintSeverity   string                 `json:"constraintSeverity"`
	ConstraintSource     string                 `json:"constraintSource"`
	ParentPath           string                 `json:"parentPath"`
}

func Validate(rootData map[string]interface{}, data map[string]interface{}, rootSpec StructureDefinition, spec StructureDefinition, parentPath string, outcome *OperationOutcome) {

	fmt.Printf("Validating %s, %s, from %s\n", rootSpec.ID, spec.ID, parentPath)

	if len(data) == 0 {
		return // Exit early if no specLibraryData to validate
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
	// ValidateUnderscoreFields(rootData, specLibraryData, parentPath, rootSpec, specLibraryData, outcome)

	// find the constraints in the specLibraryData.Snapshot.Element when id is equal to specLibraryData.ID

	constraints, _ := findMatchingElementDos(data, spec)

	// Map to track unique payload entries (key: constraintKey + parentPath)
	payloadSet := make(map[string]bool)

	// TODO: fix this not getting all constrains for all elements.
	for i := 0; i < len(*constraints); i++ {
		constraint := (*constraints)[i]

		payloadKey := fmt.Sprintf("%s|%s", constraint.Key, parentPath)
		if _, exists := payloadSet[payloadKey]; exists {
			continue // Skip duplicates
		}
		payloadSet[payloadKey] = true

		var item = FhirPathPayload{
			RootData:             rootData,
			Data:                 data,
			ConstraintExpression: constraint.Expression,
			ConstraintKey:        constraint.Key,
			ConstraintHuman:      constraint.Human,
			ConstraintSeverity:   constraint.Severity,
			ConstraintSource:     constraint.Source,
			ParentPath:           parentPath,
		}

		payload = append(payload, &item)
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

func findMatchingElementDos(data map[string]interface{}, spec StructureDefinition) (*[]Constraint, error) {
	fmt.Printf("Finding constraints for %s\n", spec.ID)
	fmt.Printf("Data: %v\n", data)
	var constraints []Constraint
	// Iterate over the keys in the specLibraryData map
	for key := range data {
		// Construct the expected ID to match with specLibraryData elements
		expectedID := spec.ID + "." + key

		// Search for a match in the specLibraryData's Snapshot.Element array
		for _, element := range spec.Snapshot.Element {

			if spec.ID == element.ID {
				fmt.Printf("Element %s\n", element.ID)
				for _, constraint := range element.Constraint {
					var skippedKeys = []string{"txt-1", "txt-2", "ele-1"}
					if contains(skippedKeys, constraint.Key) {
						continue
					}

					fmt.Printf("Constraint %s\n", constraint.Key)

					constraints = append(constraints, constraint)
				}
			}

			if element.ID == expectedID {
				fmt.Println("Element found")
				fmt.Printf("Element %s\n", element.ID)
				for _, constraint := range element.Constraint {
					var skippedKeys = []string{"txt-1", "txt-2", "ele-1"}
					if contains(skippedKeys, constraint.Key) {
						continue
					}

					fmt.Printf("Constraint %s\n", constraint.Key)

					constraints = append(constraints, constraint)
				}

				break
			}
		}
	}

	return &constraints, fmt.Errorf("no constraints found")
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
			addOperationOutcome(outcome, "required", fmt.Sprintf("Field '%s' is required", fullPath), fullPath, "Field is required", "error")
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
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("Field '%s' must be an array", fullPath), fullPath, "Field must be an array", "error")
		return
	}

	length := len(array)

	// Validate minItems
	if length < element.Min {
		addOperationOutcome(outcome, "required", fmt.Sprintf("Field '%s' has too few items: minimum is %d. Found %d elements", fullPath, element.Min, length), fullPath, "Field has too few items", "error")
	}

	// Validate maxItems
	maxItems, isUnlimited := ParseMaxItems(element.Max)
	if !isUnlimited && length > maxItems {
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("Field '%s' has too many items: maximum is %d. Found %d elements", fullPath, maxItems, length), fullPath, "Field has too many items", "error")
	}

	// Validate each element in the array
	for i, item := range array {
		itemPath := fmt.Sprintf("%s[%d]", fullPath, i)

		// check if item is a string

		switch v := item.(type) {
		case string, map[string]interface{}:
			ValidateField(rootData, v, element, itemPath, rootSpec, spec, outcome, true)
		default:
			addOperationOutcome(outcome, "invalid", fmt.Sprintf("Field '%s' must be a single value", fullPath), fullPath, "Field must be a single value", "error")
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
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("All children of '%s' must be present", fullPath), fullPath, "Field must be present", "error")
		return
	}

	switch v := value.(type) {
	case []interface{}:
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("Field '%s' must be a single value", fullPath), fullPath, "Field must be a single value", "error")
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
	nestedSpec, found := specLibraryData.Config[typeCode]
	if !found {
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("No structure definition found for type '%s'", typeCode), path, "No structure definition found", "error")
		return
	}

	// Ensure correct type assertion
	specDefinition, valid := nestedSpec.(StructureDefinition)
	if !valid {
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("Invalid structure definition for type '%s'", typeCode), path, "Invalid structure definition", "error")
		return
	}

	Validate(rootData, value.(map[string]interface{}), rootSpec, specDefinition, path, outcome)
}

// ValidatePrimitiveType validates a FHIR primitive type against its expected regex pattern.
func ValidatePrimitiveType(value string, typeCode, path string, rootSpec StructureDefinition, spec StructureDefinition, outcome *OperationOutcome) {

	if typeCode == "http://hl7.org/fhirpath/System.String" {
		typeCode = "string" // Normalize FHIRPath string type
	}

	definition, found := specLibraryData.Config[typeCode]
	if !found {
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("No definition found for type '%s'", typeCode), path, "No definition found", "error")
		return
	}

	// Extract the value element definition from the snapshot
	var valueElement *Element
	if valueElement = ExtractValueElementID(definition.(StructureDefinition).ID, definition.(StructureDefinition).Snapshot); valueElement == nil {
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("No value element found for '%s'", path), path, "No value element found", "error")
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
			addOperationOutcome(outcome, "invalid", fmt.Sprintf("No regex pattern found for '%s'", path), path, "No regex pattern found", "error")
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
		addOperationOutcome(outcome, "invalid", fmt.Sprintf("Field '%s' does not match the expected pattern: %s", path, regex), path, "Field does not match the expected pattern", "error")
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
