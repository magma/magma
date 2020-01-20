/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash a9fe8bb92b0cb474fbdcaa933318b03c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveWorkOrderMutationVariables = {|
  id: string
|};
export type RemoveWorkOrderMutationResponse = {|
  +removeWorkOrder: string
|};
export type RemoveWorkOrderMutation = {|
  variables: RemoveWorkOrderMutationVariables,
  response: RemoveWorkOrderMutationResponse,
|};
*/


/*
mutation RemoveWorkOrderMutation(
  $id: ID!
) {
  removeWorkOrder(id: $id)
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
    "name": "removeWorkOrder",
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
    "name": "RemoveWorkOrderMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveWorkOrderMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveWorkOrderMutation",
    "id": null,
    "text": "mutation RemoveWorkOrderMutation(\n  $id: ID!\n) {\n  removeWorkOrder(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '48a499f7eec8a47bc7b909821951fe45';
module.exports = node;
