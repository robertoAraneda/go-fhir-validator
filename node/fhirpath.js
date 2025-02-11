const inputJSON = process.argv[2];

const fhirExpression = process.argv[3];

const rootResource = process.argv[4];

const fhirpath = require(require.resolve("fhirpath", { paths: [__dirname] }));

let tracefunction = function (x, label) {
  console.log("Trace output [" + label + "]: ", x);
};

try {
  const resource = JSON.parse(inputJSON);
  const root = JSON.parse(rootResource);

  console.log("Expression: ", fhirExpression);

/*
  fhirpath.evaluate({ "answer": { "valueQuantity": ""}},
      { "base": "QuestionnaireResponse.item",
        "expression": "answer.value = 2 year"},
      null, null);

*/

  const result = fhirpath.evaluate(resource, fhirExpression, {
    rootResource: root,
    // traceFn: tracefunction,
    //resolveInternalTypes: false
  });
  console.log(JSON.stringify(result));
}catch (error) {
  console.log("catch")
  console.error("Error:", error.message);
  process.exit(1);
}