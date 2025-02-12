# FHIRPath Multiple Validator (TypeScript)

This project implements a FHIRPath validator using a TypeScript script to validate multiple FHIR resources. The TypeScript code executes the validation using the `fhirpath` library and processes the validation results.

## Prerequisites
- [Node.js](https://nodejs.org/)
- [npm](https://www.npmjs.com/)
- [TypeScript](https://www.typescriptlang.org/)

## Installation

1. Install dependencies:
   ```sh
   npm install
   ```
2. Compile TypeScript to JavaScript:
    ```sh
    npm run build
    ```

## Usage

1. Create a JSON file named `data-to-validate.json` and populate it with the following content:
```json
   [{
    "rootData": {
      "active": true,
      "birthDate": "1932-09-24",
      "gender": "male",
      "id": "xcda",
      "identifier": [
        {
          "system": "urn:oid:2.16.840.1.113883.19.5",
          "type": {
            "coding": [
              {
                "code": "MR",
                "system": "http://terminology.hl7.org/CodeSystem/v2-0203"
              }
            ]
          },
          "use": "usual",
          "value": "12345"
        }
      ],
      "managingOrganization": {
        "display": "Good Health Clinic",
        "reference": "Organization/2.16.840.1.113883.19.5"
      },
      "name": [
        {
          "family": "Levin",
          "given": [
            "Henry"
          ]
        }
      ],
      "resourceType": "Patient",
      "text": {
        "div": "<div xmlns=\"http://www.w3.org/1999/xhtml\"><p>Henry Levin the 7th</p></div>",
        "status": "generated"
      }
    },
    "data": {
      "display": "Good Health Clinic",
      "reference": "Organization/2.16.840.1.113883.19.5"
    },
    "constraintExpression": "reference.startsWith('#').not() or (reference.substring(1).trace('url') in %rootResource.contained.id.trace('ids'))",
    "constraintKey": "ref-1",
    "constraintHuman": "SHALL have a contained resource if a local reference is provided",
    "parentPath": "Patient.managingOrganization"
}]
```

2. Run the TypeScript application:
   ```sh
   node dist/fhirpath-evaluate.js <json-string>
   ```

## Project Structure

- `src/fhirpath-evaluate.ts`: Main TypeScript script for validation.
- `dist/fhirpath-evaluate.js`: Compiled JavaScript file after TypeScript compilation.
- `package.json`: Project dependencies and scripts.
- `tsconfig.json`: TypeScript configuration file.

## Example

To validate a FHIR resource:
1. Run the application:
   ```sh
   node dist/fhirpath-evaluate.js "[{\"rootData\": {\"active\": true, \"birthDate\": \"1932-09-24\", \"gender\": \"male\", \"id\": \"xcda\", \"identifier\": [{\"system\": \"urn:oid:2.16.840.1.113883.19.5\", \"type\": {\"coding\": [{\"code\": \"MR\", \"system\": \"http://terminology.hl7.org/CodeSystem/v2-0203\"}]}, \"use\": \"usual\", \"value\": \"12345\"}], \"managingOrganization\": {\"display\": \"Good Health Clinic\", \"reference\": \"Organization/2.16.840.1.113883.19.5\"}, \"name\": [{\"family\": \"Levin\", \"given\": [\"Henry\"]}], \"resourceType\": \"Patient\", \"text\": {\"div\": \"<div xmlns=\\\\\\\"http://www.w3.org/1999/xhtml\\\\\\\">\\n      \\n      <p>Henry Levin the 7th</p>\\n    \\n    </div>\", \"status\": \"generated\"}}, \"data\": {\"display\": \"Good Health Clinic\", \"reference\": \"Organization/2.16.840.1.113883.19.5\"}, \"constraintExpression\": \"reference.startsWith('#').not() or (reference.substring(1).trace('url') in %rootResource.contained.id.trace('ids'))\", \"constraintKey\": \"ref-1\", \"constraintHuman\": \"SHALL have a contained resource if a local reference is provided\", \"parentPath\": \"Patient.managingOrganization\"}]"
   ```
2. The output will display the validation results, including any failed validations and extracted trace data.
3. Example output:
```text
TRACE:[url] [
 "rganization/2.16.840.1.113883.19.5"
]
TRACE:[ids] []
Result: [
  {
    "result": true,
    "key": "ref-1",
    "path": "Patient.managingOrganization",
    "human": "SHALL have a contained resource if a local reference is provided"
  }
]
```
   