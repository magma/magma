/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 8b912a45f9d2d8295e39d82e6ffff8b4
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ActionID = "magma_reboot_node" | "%future added value";
export type TriggerID = "magma_alert" | "%future added value";
export type AddActionsRuleInput = {|
  name: string,
  triggerID: TriggerID,
  ruleActions: $ReadOnlyArray<?ActionsRuleActionInput>,
  ruleFilters: $ReadOnlyArray<?ActionsRuleFilterInput>,
|};
export type ActionsRuleActionInput = {|
  actionID: ActionID,
  data: string,
|};
export type ActionsRuleFilterInput = {|
  filterID: string,
  operatorID: string,
  data: string,
|};
export type AddActionsRuleMutationVariables = {|
  input: AddActionsRuleInput
|};
export type AddActionsRuleMutationResponse = {|
  +addActionsRule: {|
    +id: string,
    +name: string,
  |}
|};
export type AddActionsRuleMutation = {|
  variables: AddActionsRuleMutationVariables,
  response: AddActionsRuleMutationResponse,
|};
*/


/*
mutation AddActionsRuleMutation(
  $input: AddActionsRuleInput!
) {
  addActionsRule(input: $input) {
    id
    name
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddActionsRuleInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addActionsRule",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "ActionsRule",
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
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddActionsRuleMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddActionsRuleMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddActionsRuleMutation",
    "id": null,
    "text": "mutation AddActionsRuleMutation(\n  $input: AddActionsRuleInput!\n) {\n  addActionsRule(input: $input) {\n    id\n    name\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ba56c51c8719dde6ebb3357a2f349757';
module.exports = node;
