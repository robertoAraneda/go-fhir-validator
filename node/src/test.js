const fhirpath = require(require.resolve("fhirpath", { paths: [__dirname] }));
const fhirpath_r4_model = require(require.resolve("fhirpath/fhir-context/r4", { paths: [__dirname] }));


try {
    const data =   {
        "div": "\u003cdiv xmlns=\"http://www.w3.org/1999/xhtml\"\u003e\n      \n      \u003cp\u003eHenry Levin the 7th\u003c/p\u003e\n    \n    \u003c/div\u003e",
            "status": "generated"
    }

    /*
      {
    "rootData": {
      "active": true,
      "birthDate": "1932-09-24",
      "gender": "male",
      "id": "xcda",
      "identifier": [
        {
          "period": {
            "end": "2014-02-01",
            "start": "2025-05-06"
          },
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
        "div": "\u003cdiv xmlns=\"http://www.w3.org/1999/xhtml\"\u003e\n      \n      \u003cp\u003eHenry Levin the 7th\u003c/p\u003e\n    \n    \u003c/div\u003e",
        "status": "generated"
      }
    },
    "data": {
      "div": "\u003cdiv xmlns=\"http://www.w3.org/1999/xhtml\"\u003e\n      \n      \u003cp\u003eHenry Levin the 7th\u003c/p\u003e\n    \n    \u003c/div\u003e",
      "status": "generated"
    },
    "path": "",
    "constraintExpression": "htmlChecks()",
    "constraintKey": "txt-1",
    "constraintHuman": "The narrative SHALL contain only the basic html formatting elements and attributes described in chapters 7-11 (except section 4 of chapter 9) and 15 of the HTML 4.0 standard, \u003ca\u003e elements (either name or href), images and internally contained style attributes",
    "constraintSeverity": "error",
    "constraintSource": "",
    "parentPath": "Patient.text"
  },
     */


    const constraint = "children().count()"

    const result = fhirpath.evaluate(
            data,
        {
            base: "Patient.identifier",
            expression: constraint
        },
            null,
            fhirpath_r4_model
        );
    console.log(JSON.stringify(result));
}catch (error) {
    console.log("catch")
    console.error("Error:", error.message);
    //process.exit(1);
}