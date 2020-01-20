/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash c5b27cf50db273c5629fc4bb1456bb17
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type WorkOrderTypesListQueryVariables = {||};
export type WorkOrderTypesListQueryResponse = {|
  +workOrderTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type WorkOrderTypesListQuery = {|
  variables: WorkOrderTypesListQueryVariables,
  response: WorkOrderTypesListQueryResponse,
|};
*/


/*
query WorkOrderTypesListQuery {
  workOrderTypes(first: 50) {
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
    "concreteType": "WorkOrderTypeEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderType",
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
    "name": "WorkOrderTypesListQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrderTypes",
        "name": "__WorkOrderTypesListQuery_workOrderTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "WorkOrderTypesListQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderTypes",
        "storageKey": "workOrderTypes(first:50)",
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "workOrderTypes",
        "args": (v1/*: any*/),
        "handle": "connection",
        "key": "WorkOrderTypesListQuery_workOrderTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrderTypesListQuery",
    "id": null,
    "text": "query WorkOrderTypesListQuery {\n  workOrderTypes(first: 50) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "workOrderTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '6996261d6c40d4af29a061aef08666f1';
module.exports = node;
