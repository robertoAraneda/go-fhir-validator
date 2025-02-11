const inputJSON = process.argv[2];


const fhirpath = require(require.resolve("fhirpath", { paths: [__dirname] }));
const fhirpath_r4_model = require(require.resolve("fhirpath/fhir-context/r4", { paths: [__dirname] }));

let tracefunction = function (x, label) {
    console.log(label + ":", JSON.stringify(x, null, 2));
};

try {

  // transform inputJSON to an array of objects
  const resources = JSON.parse(inputJSON);
  const accumulator = [];

  for (let i = 0; i < resources.length; i++) {
    const bundle = resources[i];

    let result = null;

    if (bundle.constraintExpression.startsWith("contained.where(")) {

        bundle.constraintExpression = `
contained.where((('#'+id in (%resource.descendants().reference | %resource.descendants().ofType(canonical) | %resource.descendants().ofType(uri) | %resource.descendants().ofType(url))) or descendants().where(reference = '#').exists() or descendants().where(ofType(canonical) = '#').exists() or descendants().where(ofType(canonical) = '#').exists()).not()).trace('unmatched', id).empty()
`;


        result = fhirpath.evaluate(
            bundle.data,
            bundle.constraintExpression,
            { resource: bundle.rootData },
            fhirpath_r4_model,
            {
              //  traceFn: tracefunction
            }
        );

        accumulator.push({
            result: result[0],
            key: bundle.constraintKey,
            path: bundle.parentPath,
            human: bundle.constraintHuman
        });

    }else{
         result = fhirpath.evaluate(
            bundle.data,
            bundle.constraintExpression,
            { rootResource: bundle.rootData },
            fhirpath_r4_model,
            {
               // traceFn: tracefunction
            }
        );

        accumulator.push({
            result: result[0],
            key: bundle.constraintKey,
            path: bundle.parentPath,
            human: bundle.constraintHuman
        });
    }
  }
    console.log("Result: " + JSON.stringify(accumulator, null, 2));
}catch (error) {
  console.error("Error:", error.message);
  process.exit(1);
}