/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash afc2ff6d4b0aa213298bf9c5d0bf47cb
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PortDefinitionsAddEditTable__equipmentPortTypesQueryVariables = {||};
export type PortDefinitionsAddEditTable__equipmentPortTypesQueryResponse = {|
  +equipmentPortTypes: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type PortDefinitionsAddEditTable__equipmentPortTypesQuery = {|
  variables: PortDefinitionsAddEditTable__equipmentPortTypesQueryVariables,
  response: PortDefinitionsAddEditTable__equipmentPortTypesQueryResponse,
|};
*/


/*
query PortDefinitionsAddEditTable__equipmentPortTypesQuery {
  equipmentPortTypes(first: 500) {
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
    "concreteType": "EquipmentPortTypeEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortType",
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
    "name": "PortDefinitionsAddEditTable__equipmentPortTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentPortTypes",
        "name": "__PortDefinitionsTable_equipmentPortTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "PortDefinitionsAddEditTable__equipmentPortTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentPortTypes",
        "storageKey": "equipmentPortTypes(first:500)",
        "args": (v1/*: any*/),
        "concreteType": "EquipmentPortTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "equipmentPortTypes",
        "args": (v1/*: any*/),
        "handle": "connection",
        "key": "PortDefinitionsTable_equipmentPortTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "PortDefinitionsAddEditTable__equipmentPortTypesQuery",
    "id": null,
    "text": "query PortDefinitionsAddEditTable__equipmentPortTypesQuery {\n  equipmentPortTypes(first: 500) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "equipmentPortTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8387908919eda242a725beb795eb547a';
module.exports = node;
