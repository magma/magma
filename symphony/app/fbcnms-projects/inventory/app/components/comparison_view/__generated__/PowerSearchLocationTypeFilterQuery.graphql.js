/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 5afd7bcc3a6a2b55ec9c8172c0936d11
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PowerSearchLocationTypeFilterQueryVariables = {||};
export type PowerSearchLocationTypeFilterQueryResponse = {|
  +locationTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type PowerSearchLocationTypeFilterQuery = {|
  variables: PowerSearchLocationTypeFilterQueryVariables,
  response: PowerSearchLocationTypeFilterQueryResponse,
|};
*/


/*
query PowerSearchLocationTypeFilterQuery {
  locationTypes(first: 50) {
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
    "name": "locationTypes",
    "storageKey": "locationTypes(first:50)",
    "args": [
      {
        "kind": "Literal",
        "name": "first",
        "value": 50
      }
    ],
    "concreteType": "LocationTypeConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "LocationTypeEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "LocationType",
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
    "name": "PowerSearchLocationTypeFilterQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchLocationTypeFilterQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchLocationTypeFilterQuery",
    "id": null,
    "text": "query PowerSearchLocationTypeFilterQuery {\n  locationTypes(first: 50) {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'cf67fcab81f4a29c18222aeebeb27e43';
module.exports = node;
