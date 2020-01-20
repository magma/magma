/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0c8f467d878da619303e63b30f7cf458
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DeleteFloorPlanMutationVariables = {|
  id: string
|};
export type DeleteFloorPlanMutationResponse = {|
  +deleteFloorPlan: boolean
|};
export type DeleteFloorPlanMutation = {|
  variables: DeleteFloorPlanMutationVariables,
  response: DeleteFloorPlanMutationResponse,
|};
*/


/*
mutation DeleteFloorPlanMutation(
  $id: ID!
) {
  deleteFloorPlan(id: $id)
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
    "name": "deleteFloorPlan",
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
    "name": "DeleteFloorPlanMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "DeleteFloorPlanMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeleteFloorPlanMutation",
    "id": null,
    "text": "mutation DeleteFloorPlanMutation(\n  $id: ID!\n) {\n  deleteFloorPlan(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '398cf21649438ad9454a4cdcb7c81c89';
module.exports = node;
