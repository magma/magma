/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash dad11eac08f2caf8099186c65b2afc91
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
  workOrders(first: 50, showCompleted: false) {
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
var v0 = {
  "kind": "Literal",
  "name": "showCompleted",
  "value": false
},
v1 = [
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
v2 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50
  },
  (v0/*: any*/)
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
        "storageKey": "__WorkOrdersPane_workOrders_connection(showCompleted:false)",
        "args": [
          (v0/*: any*/)
        ],
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": (v1/*: any*/)
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
        "storageKey": "workOrders(first:50,showCompleted:false)",
        "args": (v2/*: any*/),
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": (v1/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "workOrders",
        "args": (v2/*: any*/),
        "handle": "connection",
        "key": "WorkOrdersPane_workOrders",
        "filters": [
          "showCompleted"
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrdersPaneQuery",
    "id": null,
    "text": "query WorkOrdersPaneQuery {\n  workOrders(first: 50, showCompleted: false) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
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
(node/*: any*/).hash = 'c78f283ec8ed6da0a8611d4519f220cf';
module.exports = node;
