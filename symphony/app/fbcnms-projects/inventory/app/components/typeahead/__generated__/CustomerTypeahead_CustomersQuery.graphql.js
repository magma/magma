/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash f3cdf937dd6146cc44ffb9ebb07a76ba
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type CustomerTypeahead_CustomersQueryVariables = {|
  limit?: ?number
|};
export type CustomerTypeahead_CustomersQueryResponse = {|
  +customers: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +externalId: ?string,
      |}
    |}>
  |}
|};
export type CustomerTypeahead_CustomersQuery = {|
  variables: CustomerTypeahead_CustomersQueryVariables,
  response: CustomerTypeahead_CustomersQueryResponse,
|};
*/


/*
query CustomerTypeahead_CustomersQuery(
  $limit: Int
) {
  customers(first: $limit) {
    edges {
      node {
        id
        name
        externalId
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "limit",
    "type": "Int",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "customers",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "first",
        "variableName": "limit"
      }
    ],
    "concreteType": "CustomerConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "CustomerEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "Customer",
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
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "externalId",
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
    "name": "CustomerTypeahead_CustomersQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "CustomerTypeahead_CustomersQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "CustomerTypeahead_CustomersQuery",
    "id": null,
    "text": "query CustomerTypeahead_CustomersQuery(\n  $limit: Int\n) {\n  customers(first: $limit) {\n    edges {\n      node {\n        id\n        name\n        externalId\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '18df0202c887393e4158102ab7e9ba4c';
module.exports = node;
