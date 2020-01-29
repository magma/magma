/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 568578d67ed861129a870a9054b6a2ed
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveActionsRuleMutationVariables = {|
  id: string
|};
export type RemoveActionsRuleMutationResponse = {|
  +removeActionsRule: boolean
|};
export type RemoveActionsRuleMutation = {|
  variables: RemoveActionsRuleMutationVariables,
  response: RemoveActionsRuleMutationResponse,
|};
*/


/*
mutation RemoveActionsRuleMutation(
  $id: ID!
) {
  removeActionsRule(id: $id)
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
    "name": "removeActionsRule",
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
    "name": "RemoveActionsRuleMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveActionsRuleMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveActionsRuleMutation",
    "id": null,
    "text": "mutation RemoveActionsRuleMutation(\n  $id: ID!\n) {\n  removeActionsRule(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'bf9c9df600267a517d02f476465e2a8f';
module.exports = node;
