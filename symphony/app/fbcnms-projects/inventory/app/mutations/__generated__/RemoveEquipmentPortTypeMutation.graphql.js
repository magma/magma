/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 24e4af8962784fc2704217c26d2b39b9
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveEquipmentPortTypeMutationVariables = {|
  id: string
|};
export type RemoveEquipmentPortTypeMutationResponse = {|
  +removeEquipmentPortType: string
|};
export type RemoveEquipmentPortTypeMutation = {|
  variables: RemoveEquipmentPortTypeMutationVariables,
  response: RemoveEquipmentPortTypeMutationResponse,
|};
*/


/*
mutation RemoveEquipmentPortTypeMutation(
  $id: ID!
) {
  removeEquipmentPortType(id: $id)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "removeEquipmentPortType",
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      }
    ],
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "RemoveEquipmentPortTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveEquipmentPortTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveEquipmentPortTypeMutation",
    "id": null,
    "text": "mutation RemoveEquipmentPortTypeMutation(\n  $id: ID!\n) {\n  removeEquipmentPortType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f14522f936273971e9cd0691f19f241e';
module.exports = node;
