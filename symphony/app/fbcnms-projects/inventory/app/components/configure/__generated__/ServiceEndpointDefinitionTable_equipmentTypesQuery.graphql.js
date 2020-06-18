/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4de2c231d4ba0e5deba98a1ada4e0935
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ServiceEndpointDefinitionTable_equipmentTypesQueryVariables = {||};
export type ServiceEndpointDefinitionTable_equipmentTypesQueryResponse = {|
  +equipmentTypes: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type ServiceEndpointDefinitionTable_equipmentTypesQuery = {|
  variables: ServiceEndpointDefinitionTable_equipmentTypesQueryVariables,
  response: ServiceEndpointDefinitionTable_equipmentTypesQueryResponse,
|};
*/


/*
query ServiceEndpointDefinitionTable_equipmentTypesQuery {
  equipmentTypes(first: 500) {
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
    "concreteType": "EquipmentTypeEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentType",
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
    "name": "ServiceEndpointDefinitionTable_equipmentTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentTypes",
        "name": "__ServiceEndpointDefinitionTable_equipmentTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceEndpointDefinitionTable_equipmentTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentTypes",
        "storageKey": "equipmentTypes(first:500)",
        "args": (v1/*: any*/),
        "concreteType": "EquipmentTypeConnection",
        "plural": false,
        "selections": (v0/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "equipmentTypes",
        "args": (v1/*: any*/),
        "handle": "connection",
        "key": "ServiceEndpointDefinitionTable_equipmentTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceEndpointDefinitionTable_equipmentTypesQuery",
    "id": null,
    "text": "query ServiceEndpointDefinitionTable_equipmentTypesQuery {\n  equipmentTypes(first: 500) {\n    edges {\n      node {\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "equipmentTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '88a844a623c430a24a1944839e78684d';
module.exports = node;
