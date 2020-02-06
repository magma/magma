/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 3bca0a33295020ff7587d39b593251ae
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PowerSearchLocationFilterQueryVariables = {|
  name: string,
  types?: ?$ReadOnlyArray<string>,
|};
export type PowerSearchLocationFilterQueryResponse = {|
  +locations: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type PowerSearchLocationFilterQuery = {|
  variables: PowerSearchLocationFilterQueryVariables,
  response: PowerSearchLocationFilterQueryResponse,
|};
*/


/*
query PowerSearchLocationFilterQuery(
  $name: String!
  $types: [ID!]
) {
  locations(name: $name, first: 10, types: $types) {
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
    "kind": "LocalArgument",
    "name": "name",
    "type": "String!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "types",
    "type": "[ID!]",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "locations",
    "storageKey": null,
    "args": [
      {
        "kind": "Literal",
        "name": "first",
        "value": 10
      },
      {
        "kind": "Variable",
        "name": "name",
        "variableName": "name"
      },
      {
        "kind": "Variable",
        "name": "types",
        "variableName": "types"
      }
    ],
    "concreteType": "LocationConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "LocationEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "Location",
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
    "name": "PowerSearchLocationFilterQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchLocationFilterQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchLocationFilterQuery",
    "id": null,
    "text": "query PowerSearchLocationFilterQuery(\n  $name: String!\n  $types: [ID!]\n) {\n  locations(name: $name, first: 10, types: $types) {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c2aaf778bfd1a263c4290d97b0253930';
module.exports = node;
