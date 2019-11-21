/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 31f155f153091284f030e18af90116db
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveEquipmentMutationVariables = {|
  id: string,
  work_order_id?: ?string,
|};
export type RemoveEquipmentMutationResponse = {|
  +removeEquipment: string
|};
export type RemoveEquipmentMutation = {|
  variables: RemoveEquipmentMutationVariables,
  response: RemoveEquipmentMutationResponse,
|};
*/


/*
mutation RemoveEquipmentMutation(
  $id: ID!
  $work_order_id: ID
) {
  removeEquipment(id: $id, workOrderId: $work_order_id)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "work_order_id",
    "type": "ID",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "removeEquipment",
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      },
      {
        "kind": "Variable",
        "name": "workOrderId",
        "variableName": "work_order_id"
      }
    ],
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "RemoveEquipmentMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveEquipmentMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveEquipmentMutation",
    "id": null,
    "text": "mutation RemoveEquipmentMutation(\n  $id: ID!\n  $work_order_id: ID\n) {\n  removeEquipment(id: $id, workOrderId: $work_order_id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '81fe5c2632e64bde78a3342c28d9258d';
module.exports = node;
