/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0c51d2c9df4aa6d55e8a6b4969a73c92
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PowerSearchPortDefinitionFilterQueryVariables = {||};
export type PowerSearchPortDefinitionFilterQueryResponse = {|
  +equipmentPortDefinitions: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type PowerSearchPortDefinitionFilterQuery = {|
  variables: PowerSearchPortDefinitionFilterQueryVariables,
  response: PowerSearchPortDefinitionFilterQueryResponse,
|};
*/


/*
query PowerSearchPortDefinitionFilterQuery {
  equipmentPortDefinitions {
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
    "name": "equipmentPortDefinitions",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPortDefinitionConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortDefinitionEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortDefinition",
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
    "name": "PowerSearchPortDefinitionFilterQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchPortDefinitionFilterQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchPortDefinitionFilterQuery",
    "id": null,
    "text": "query PowerSearchPortDefinitionFilterQuery {\n  equipmentPortDefinitions {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7e8942b0ece87929d960954bb040fc3a';
module.exports = node;
