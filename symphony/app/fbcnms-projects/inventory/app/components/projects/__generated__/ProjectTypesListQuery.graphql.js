/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 13e4f432dcb7618918fb5f938f6438b3
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ProjectTypesListQueryVariables = {||};
export type ProjectTypesListQueryResponse = {|
  +projectTypes: ?{|
    +edges: ?$ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type ProjectTypesListQuery = {|
  variables: ProjectTypesListQueryVariables,
  response: ProjectTypesListQueryResponse,
|};
*/


/*
query ProjectTypesListQuery {
  projectTypes(first: 50) {
    edges {
      node {
        id
        name
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "projectTypes",
    "storageKey": "projectTypes(first:50)",
    "args": [
      {
        "kind": "Literal",
        "name": "first",
        "value": 50
      }
    ],
    "concreteType": "ProjectTypeConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "ProjectTypeEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectType",
            "plural": false,
            "selections": [
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "id",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "name",
                "args": null,
                "storageKey": null
              }
            ]
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
    "name": "ProjectTypesListQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectTypesListQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectTypesListQuery",
    "id": null,
    "text": "query ProjectTypesListQuery {\n  projectTypes(first: 50) {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7ed39a2ca143bf8f9cd64150d796c71b';
module.exports = node;
