/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 83336aad022f991548e08b68ab539e2b
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DownloadPythonPackageQueryVariables = {||};
export type DownloadPythonPackageQueryResponse = {|
  +pythonPackages: $ReadOnlyArray<{|
    +version: string,
    +whlFileKey: string,
    +uploadTime: any,
  |}>
|};
export type DownloadPythonPackageQuery = {|
  variables: DownloadPythonPackageQueryVariables,
  response: DownloadPythonPackageQueryResponse,
|};
*/


/*
query DownloadPythonPackageQuery {
  pythonPackages {
    version
    whlFileKey
    uploadTime
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "pythonPackages",
    "storageKey": null,
    "args": null,
    "concreteType": "PythonPackage",
    "plural": true,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "version",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "whlFileKey",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "uploadTime",
        "args": null,
        "storageKey": null
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "DownloadPythonPackageQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "DownloadPythonPackageQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "DownloadPythonPackageQuery",
    "id": null,
    "text": "query DownloadPythonPackageQuery {\n  pythonPackages {\n    version\n    whlFileKey\n    uploadTime\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c2ffc17589d8cfa0daad1826bd83f69c';
module.exports = node;
