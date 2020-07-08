/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash a01160079bc62ddcb9aa66db83068ea0
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type WorkOrdersPaneQueryVariables = {||};
export type WorkOrdersPaneQueryResponse = {|
  +workOrders: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type WorkOrdersPaneQuery = {|
  variables: WorkOrdersPaneQueryVariables,
  response: WorkOrdersPaneQueryResponse,
|};
*/


/*
query WorkOrdersPaneQuery {
  workOrders(first: 500) {
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
    "concreteType": "WorkOrderEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrder",
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
    "value": 500
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrdersPaneQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrders",
        "name": "__WorkOrdersPane_workOrders_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "WorkOrdersPaneQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrders",
        "storageKey": "workOrders(first:500)",
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "workOrders",
        "args": (v1/*: any*/),
        "handle": "connection",
        "key": "WorkOrdersPane_workOrders",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrdersPaneQuery",
    "id": null,
    "text": "query WorkOrdersPaneQuery {\n  workOrders(first: 500) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "workOrders"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '0e030cf6ac3d385f5de8a99532eb512a';
module.exports = node;
