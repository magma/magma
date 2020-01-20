/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1c7dad697c8ec8505a4fcc0be7fa5e3b
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ActionsListCard_actionsRule$ref = any;
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
export type EditActionsRuleMutationVariables = {|
  id: string,
  input: AddActionsRuleInput,
|};
export type EditActionsRuleMutationResponse = {|
  +editActionsRule: ?{|
    +$fragmentRefs: ActionsListCard_actionsRule$ref
  |}
|};
export type EditActionsRuleMutation = {|
  variables: EditActionsRuleMutationVariables,
  response: EditActionsRuleMutationResponse,
|};
*/


/*
mutation EditActionsRuleMutation(
  $id: ID!
  $input: AddActionsRuleInput!
) {
  editActionsRule(id: $id, input: $input) {
    ...ActionsListCard_actionsRule
    id
  }
}

fragment ActionRow_data on ActionsTrigger {
  triggerID
  supportedActions {
    actionID
    dataType
    description
  }
}

fragment ActionsAddDialog_triggerData on ActionsTrigger {
  triggerID
  description
  ...ActionRow_data
  ...TriggerFilterRow_data
}

fragment ActionsListCard_actionsRule on ActionsRule {
  id
  name
  trigger {
    description
    ...ActionsAddDialog_triggerData
    id
  }
  ruleActions {
    actionID
    data
  }
  ruleFilters {
    filterID
    operatorID
    data
  }
}

fragment TriggerFilterOperator_data on ActionsFilter {
  supportedOperators {
    operatorID
    description
    dataType
  }
}

fragment TriggerFilterRow_data on ActionsTrigger {
  triggerID
  supportedFilters {
    filterID
    description
    supportedOperators {
      operatorID
    }
    ...TriggerFilterOperator_data
  }
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
    "name": "input",
    "type": "AddActionsRuleInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  },
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "actionID",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "dataType",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "filterID",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "operatorID",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "data",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditActionsRuleMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editActionsRule",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ActionsRule",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "ActionsListCard_actionsRule",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EditActionsRuleMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editActionsRule",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ActionsRule",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "name",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "trigger",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsTrigger",
            "plural": false,
            "selections": [
              (v3/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "triggerID",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "supportedActions",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsAction",
                "plural": true,
                "selections": [
                  (v4/*: any*/),
                  (v5/*: any*/),
                  (v3/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "supportedFilters",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsFilter",
                "plural": true,
                "selections": [
                  (v6/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "supportedOperators",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ActionsOperator",
                    "plural": true,
                    "selections": [
                      (v7/*: any*/),
                      (v3/*: any*/),
                      (v5/*: any*/)
                    ]
                  }
                ]
              },
              (v2/*: any*/)
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "ruleActions",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsRuleAction",
            "plural": true,
            "selections": [
              (v4/*: any*/),
              (v8/*: any*/)
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "ruleFilters",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsRuleFilter",
            "plural": true,
            "selections": [
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditActionsRuleMutation",
    "id": null,
    "text": "mutation EditActionsRuleMutation(\n  $id: ID!\n  $input: AddActionsRuleInput!\n) {\n  editActionsRule(id: $id, input: $input) {\n    ...ActionsListCard_actionsRule\n    id\n  }\n}\n\nfragment ActionRow_data on ActionsTrigger {\n  triggerID\n  supportedActions {\n    actionID\n    dataType\n    description\n  }\n}\n\nfragment ActionsAddDialog_triggerData on ActionsTrigger {\n  triggerID\n  description\n  ...ActionRow_data\n  ...TriggerFilterRow_data\n}\n\nfragment ActionsListCard_actionsRule on ActionsRule {\n  id\n  name\n  trigger {\n    description\n    ...ActionsAddDialog_triggerData\n    id\n  }\n  ruleActions {\n    actionID\n    data\n  }\n  ruleFilters {\n    filterID\n    operatorID\n    data\n  }\n}\n\nfragment TriggerFilterOperator_data on ActionsFilter {\n  supportedOperators {\n    operatorID\n    description\n    dataType\n  }\n}\n\nfragment TriggerFilterRow_data on ActionsTrigger {\n  triggerID\n  supportedFilters {\n    filterID\n    description\n    supportedOperators {\n      operatorID\n    }\n    ...TriggerFilterOperator_data\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '521657dc4718a97d75660361d8d42456';
module.exports = node;
