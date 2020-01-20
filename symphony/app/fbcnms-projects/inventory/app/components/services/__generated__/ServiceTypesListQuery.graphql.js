/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 563c9d66f57e677e769c2096a5e3f300
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ServiceTypesListQueryVariables = {||};
export type ServiceTypesListQueryResponse = {|
  +serviceTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type ServiceTypesListQuery = {|
  variables: ServiceTypesListQueryVariables,
  response: ServiceTypesListQueryResponse,
|};
*/


/*
query ServiceTypesListQuery {
  serviceTypes(first: 50) {
    edges {
      node {
        id
        name
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "edges",
    "storageKey": null,
    "args": null,
    "concreteType": "ServiceTypeEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceType",
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
            "name": "__typename",
            "args": null,
            "storageKey": null
          }
        ]
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "cursor",
        "args": null,
        "storageKey": null
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "pageInfo",
    "storageKey": null,
    "args": null,
    "concreteType": "PageInfo",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "endCursor",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "hasNextPage",
        "args": null,
        "storageKey": null
      }
    ]
  }
],
v1 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ServiceTypesListQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "serviceTypes",
        "name": "__ServiceTypesListQuery_serviceTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceTypesListQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceTypes",
        "storageKey": "serviceTypes(first:50)",
        "args": (v1/*: any*/),
        "concreteType": "ServiceTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "serviceTypes",
        "args": (v1/*: any*/),
        "handle": "connection",
        "key": "ServiceTypesListQuery_serviceTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceTypesListQuery",
    "id": null,
    "text": "query ServiceTypesListQuery {\n  serviceTypes(first: 50) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "serviceTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '88e4d8dffdb018ebff8f12f4a5c45e17';
module.exports = node;
