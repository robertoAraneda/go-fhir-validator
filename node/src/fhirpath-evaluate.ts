import * as fhirpath from "fhirpath";
import * as fhirpath_r4_model from "fhirpath/fhir-context/r4";

const inputJSON: string = process.argv[2];

const traceFunction = (x: unknown, label: string): void => {
    console.log(`${label}:`, JSON.stringify(x, null, 2));
};

interface Bundle {
    constraintExpression: string;
    data: any;
    rootData: any;
    constraintKey: string;
    parentPath: string;
    constraintHuman: string;
    constraintSource: string;
    constraintSeverity: string;
}

interface ResponseBundle {
    result: boolean;
    key: string;
    path: string;
    human: string;
    source: string;
    severity: string;
}

try {
    // Transform inputJSON to an array of objects
    const resources: Bundle[] = JSON.parse(inputJSON);
    const accumulator: ResponseBundle[] = [];


    for (const bundle of resources) {
        let result: boolean[] = [];
        const resourceType = bundle.data.resourceType as string;
        console.log("parentPath:", bundle.parentPath);

        if (bundle.constraintExpression.startsWith("contained.where(")) {
            bundle.constraintExpression = `
contained.where((('#'+id in (%resource.descendants().reference | %resource.descendants().ofType(canonical) | %resource.descendants().ofType(uri) | %resource.descendants().ofType(url))) or descendants().where(reference = '#').exists() or descendants().where(ofType(canonical) = '#').exists() or descendants().where(ofType(canonical) = '#').exists()).not()).trace('unmatched', id).empty()
`;

            result = fhirpath.evaluate(
                bundle.data,
                resourceType ? bundle.constraintExpression : { base: bundle.parentPath, expression: bundle.constraintExpression },
                { resource: bundle.rootData },
                fhirpath_r4_model
            ) as boolean[];
        } else {
            result = fhirpath.evaluate(
                bundle.data,
                resourceType ? bundle.constraintExpression : { base: bundle.parentPath, expression: bundle.constraintExpression },
                { rootResource: bundle.rootData },
                fhirpath_r4_model
            ) as boolean[];
        }

        accumulator.push({
            result: result[0],
            key: bundle.constraintKey,
            path: bundle.parentPath,
            human: bundle.constraintHuman,
            source: bundle.constraintSource,
            severity: bundle.constraintSeverity,
        });
    }

    console.log("Result:", JSON.stringify(accumulator, null, 2));
} catch (error: any) {
    console.error("Error:", error.message);
    //process.exit(1);
}
