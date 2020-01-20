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
declare export opaque type BasicCheckListItemFilling_item$ref: FragmentReference;
declare export opaque type BasicCheckListItemFilling_item$fragmentType: BasicCheckListItemFilling_item$ref;
export type BasicCheckListItemFilling_item = {|
  +title: string,
  +checked: ?boolean,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: BasicCheckListItemFilling_item$ref,
|};
export type BasicCheckListItemFilling_item$data = BasicCheckListItemFilling_item;
export type BasicCheckListItemFilling_item$key = {
  +$data?: BasicCheckListItemFilling_item$data,
  +$fragmentRefs: BasicCheckListItemFilling_item$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "BasicCheckListItemFilling_item",
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
(node/*: any*/).hash = 'e911cf65a12669869d4ea70bfad2726c';
module.exports = node;
