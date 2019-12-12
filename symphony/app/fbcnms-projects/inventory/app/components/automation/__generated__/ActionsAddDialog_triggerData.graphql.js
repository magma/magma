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
type ActionRow_data$ref = any;
type TriggerFilterRow_data$ref = any;
export type TriggerID = "magma_alert" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ActionsAddDialog_triggerData$ref: FragmentReference;
declare export opaque type ActionsAddDialog_triggerData$fragmentType: ActionsAddDialog_triggerData$ref;
export type ActionsAddDialog_triggerData = {|
  +triggerID: TriggerID,
  +description: string,
  +$fragmentRefs: ActionRow_data$ref & TriggerFilterRow_data$ref,
  +$refType: ActionsAddDialog_triggerData$ref,
|};
export type ActionsAddDialog_triggerData$data = ActionsAddDialog_triggerData;
export type ActionsAddDialog_triggerData$key = {
  +$data?: ActionsAddDialog_triggerData$data,
  +$fragmentRefs: ActionsAddDialog_triggerData$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ActionsAddDialog_triggerData",
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
      "kind": "ScalarField",
      "alias": null,
      "name": "description",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "ActionRow_data",
      "args": null
    },
    {
      "kind": "FragmentSpread",
      "name": "TriggerFilterRow_data",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '0af12b6eb810778e0e590238e2c67a25';
module.exports = node;
