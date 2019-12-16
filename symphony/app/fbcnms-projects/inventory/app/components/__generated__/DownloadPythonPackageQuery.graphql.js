/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 75d9df408cd382faa9469a98dbea6740
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DownloadPythonPackageQueryVariables = {||};
export type DownloadPythonPackageQueryResponse = {|
  +latestPythonPackage: ?{|
    +lastPythonPackage: ?{|
      +version: string,
      +whlFileKey: string,
    |}
  |}
|};
export type DownloadPythonPackageQuery = {|
  variables: DownloadPythonPackageQueryVariables,
  response: DownloadPythonPackageQueryResponse,
|};
*/


/*
query DownloadPythonPackageQuery {
  latestPythonPackage {
    lastPythonPackage {
      version
      whlFileKey
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "latestPythonPackage",
    "storageKey": null,
    "args": null,
    "concreteType": "LatestPythonPackageResult",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "lastPythonPackage",
        "storageKey": null,
        "args": null,
        "concreteType": "PythonPackage",
        "plural": false,
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
          }
        ]
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
    "text": "query DownloadPythonPackageQuery {\n  latestPythonPackage {\n    lastPythonPackage {\n      version\n      whlFileKey\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '684edf2b73ab7adf9b106bc0d215835b';
module.exports = node;
