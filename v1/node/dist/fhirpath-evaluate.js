"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
const fhirpath = __importStar(require("fhirpath"));
const fhirpath_r4_model = __importStar(require("fhirpath/fhir-context/r4"));
const inputJSON = process.argv[2];
const traceFunction = (x, label) => {
    console.log(`${label}:`, JSON.stringify(x, null, 2));
};
try {
    // Transform inputJSON to an array of objects
    const resources = JSON.parse(inputJSON);
    const accumulator = [];
    for (const bundle of resources) {
        let result = [];
        const resourceType = bundle.data.resourceType;
        console.log("parentPath:", bundle.parentPath);
        if (bundle.constraintExpression.startsWith("contained.where(")) {
            bundle.constraintExpression = `
contained.where((('#'+id in (%resource.descendants().reference | %resource.descendants().ofType(canonical) | %resource.descendants().ofType(uri) | %resource.descendants().ofType(url))) or descendants().where(reference = '#').exists() or descendants().where(ofType(canonical) = '#').exists() or descendants().where(ofType(canonical) = '#').exists()).not()).trace('unmatched', id).empty()
`;
            result = fhirpath.evaluate(bundle.data, resourceType ? bundle.constraintExpression : { base: bundle.parentPath, expression: bundle.constraintExpression }, { resource: bundle.rootData }, fhirpath_r4_model);
        }
        else {
            result = fhirpath.evaluate(bundle.data, resourceType ? bundle.constraintExpression : { base: bundle.parentPath, expression: bundle.constraintExpression }, { rootResource: bundle.rootData }, fhirpath_r4_model);
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
}
catch (error) {
    console.error("Error:", error.message);
    //process.exit(1);
}
