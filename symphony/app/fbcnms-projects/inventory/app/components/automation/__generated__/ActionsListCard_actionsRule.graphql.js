/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
type ActionsAddDialog_triggerData$ref = any;
export type ActionID = "magma_reboot_node" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ActionsListCard_actionsRule$ref: FragmentReference;
declare export opaque type ActionsListCard_actionsRule$fragmentType: ActionsListCard_actionsRule$ref;
export type ActionsListCard_actionsRule = {|
  +id: string,
  +name: string,
  +trigger: {|
    +description: string,
    +$fragmentRefs: ActionsAddDialog_triggerData$ref,
  |},
  +ruleActions: $ReadOnlyArray<?{|
    +actionID: ActionID,
    +data: string,
  |}>,
  +ruleFilters: $ReadOnlyArray<?{|
    +filterID: ?string,
    +operatorID: ?string,
    +data: string,
  |}>,
  +$refType: ActionsListCard_actionsRule$ref,
|};
export type ActionsListCard_actionsRule$data = ActionsListCard_actionsRule;
export type ActionsListCard_actionsRule$key = {
  +$data?: ActionsListCard_actionsRule$data,
  +$fragmentRefs: ActionsListCard_actionsRule$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "data",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ActionsListCard_actionsRule",
  "type": "ActionsRule",
  "metadata": null,
  "argumentDefinitions": [],
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
      "kind": "LinkedField",
      "alias": null,
      "name": "trigger",
      "storageKey": null,
      "args": null,
      "concreteType": "ActionsTrigger",
      "plural": false,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "description",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "FragmentSpread",
          "name": "ActionsAddDialog_triggerData",
          "args": null
        }
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
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "actionID",
          "args": null,
          "storageKey": null
        },
        (v0/*: any*/)
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
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "filterID",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "operatorID",
          "args": null,
          "storageKey": null
        },
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '5c259632c00bc3ceaa2d85fde01b1664';
module.exports = node;
