/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0798b383924849432e2ffa211550dacd
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveEquipmentTypeMutationVariables = {|
  id: string
|};
export type RemoveEquipmentTypeMutationResponse = {|
  +removeEquipmentType: string
|};
export type RemoveEquipmentTypeMutation = {|
  variables: RemoveEquipmentTypeMutationVariables,
  response: RemoveEquipmentTypeMutationResponse,
|};
*/


/*
mutation RemoveEquipmentTypeMutation(
  $id: ID!
) {
  removeEquipmentType(id: $id)
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
    "name": "removeEquipmentType",
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
    "name": "RemoveEquipmentTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveEquipmentTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveEquipmentTypeMutation",
    "id": null,
    "text": "mutation RemoveEquipmentTypeMutation(\n  $id: ID!\n) {\n  removeEquipmentType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '11b09437ff00a0fb11b78dcc441c853f';
module.exports = node;
