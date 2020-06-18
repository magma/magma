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
type CheckListItem_item$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type BasicCheckListItemDefinition_item$ref: FragmentReference;
declare export opaque type BasicCheckListItemDefinition_item$fragmentType: BasicCheckListItemDefinition_item$ref;
export type BasicCheckListItemDefinition_item = {|
  +title: string,
  +checked: ?boolean,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: BasicCheckListItemDefinition_item$ref,
|};
export type BasicCheckListItemDefinition_item$data = BasicCheckListItemDefinition_item;
export type BasicCheckListItemDefinition_item$key = {
  +$data?: BasicCheckListItemDefinition_item$data,
  +$fragmentRefs: BasicCheckListItemDefinition_item$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "BasicCheckListItemDefinition_item",
  "type": "CheckListItem",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "title",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "checked",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "CheckListItem_item",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '6bb61023ad62f75200474c3512183886';
module.exports = node;
