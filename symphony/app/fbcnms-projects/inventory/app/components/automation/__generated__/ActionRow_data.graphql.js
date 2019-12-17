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
export type ActionID = "magma_reboot_node" | "%future added value";
export type ActionsDataType = "string" | "stringArray" | "%future added value";
export type TriggerID = "magma_alert" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ActionRow_data$ref: FragmentReference;
declare export opaque type ActionRow_data$fragmentType: ActionRow_data$ref;
export type ActionRow_data = {|
  +triggerID: TriggerID,
  +supportedActions: $ReadOnlyArray<?{|
    +actionID: ActionID,
    +dataType: ActionsDataType,
    +description: string,
  |}>,
  +$refType: ActionRow_data$ref,
|};
export type ActionRow_data$data = ActionRow_data;
export type ActionRow_data$key = {
  +$data?: ActionRow_data$data,
  +$fragmentRefs: ActionRow_data$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ActionRow_data",
  "type": "ActionsTrigger",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
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
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "actionID",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "dataType",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "description",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'ee5d5f84b4912c3c7b63fda08f43b60a';
module.exports = node;
