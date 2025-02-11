# FHIRPath Multiple Validator

This project implements a FHIRPath validator using a Node.js script to validate multiple FHIR resources. The Go code handles the execution of the Node.js script and processes the validation results.

## Prerequisites
- [Node.js](https://nodejs.org/)
- [npm](https://www.npmjs.com/)
- [Go](https://go.dev/)

## Installation

1. Clone the repository:
   ```sh
   git clone <repository-url>
   cd <repository-directory>
   ```
2. Install Node.js dependencies:
   ```sh
   cd node
   npm install
   ```
3. Ensure you have Go installed and set up.

## Usage

1. Place your FHIR resources in a JSON file, e.g., `resource.json`.
2. Run the Go application:
   ```sh
   go run main.go
   ```

## Project Structure

- `main.go`: Entry point of the Go application. It reads the FHIR resources from a JSON file and calls the validation function.
- `fhirpath_validator_multiple.go`: Contains the `FhirPathValidatorMultiple` function that executes the Node.js script and processes the validation results.
- `node/fhirpath_multiple.js`: Node.js script that performs the FHIRPath validation.
- `node/README.md`: This README file.

## Example

To validate a FHIR resource:
1. Ensure your `resource.json` file is correctly formatted.
2. Run the Go application:
   ```sh
    go run main.go validator.go utils.go load-structure-definition.go struct.go fhirpath_validator_multiple.go
   ```
3. The output will display the validation results, including any failed validations and extracted trace data.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.