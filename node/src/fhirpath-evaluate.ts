import * as fhirpath from "fhirpath";
import * as fhirpath_r4_model from "fhirpath/fhir-context/r4";

const inputJSON: string = process.argv[2];

const traceFunction = (x: unknown, label: string): void => {
    console.log(`${label}:`, JSON.stringify(x, null, 2));
};

try {
    // Transform inputJSON to an array of objects
    const resources: any[] = JSON.parse(inputJSON);
    const accumulator: Array<{ result: boolean; key: string; path: string; human: string }> = [];

    for (const bundle of resources) {
        let result: boolean[] = [];

        if (bundle.constraintExpression.startsWith("contained.where(")) {
            bundle.constraintExpression = `
contained.where((('#'+id in (%resource.descendants().reference | %resource.descendants().ofType(canonical) | %resource.descendants().ofType(uri) | %resource.descendants().ofType(url))) or descendants().where(reference = '#').exists() or descendants().where(ofType(canonical) = '#').exists() or descendants().where(ofType(canonical) = '#').exists()).not()).trace('unmatched', id).empty()
`;

            result = fhirpath.evaluate(
                bundle.data,
                bundle.constraintExpression,
                { resource: bundle.rootData },
                fhirpath_r4_model
            ) as boolean[];
        } else {
            result = fhirpath.evaluate(
                bundle.data,
                bundle.constraintExpression,
                { rootResource: bundle.rootData },
                fhirpath_r4_model
            ) as boolean[];
        }

        accumulator.push({
            result: result[0],
            key: bundle.constraintKey,
            path: bundle.parentPath,
            human: bundle.constraintHuman
        });
    }

    console.log("Result:", JSON.stringify(accumulator, null, 2));
} catch (error: any) {
    console.error("Error:", error.message);
    process.exit(1);
}
