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
type CheckListTable_list$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListCategoryItemsDialog_items$ref: FragmentReference;
declare export opaque type CheckListCategoryItemsDialog_items$fragmentType: CheckListCategoryItemsDialog_items$ref;
export type CheckListCategoryItemsDialog_items = $ReadOnlyArray<{|
  +checked: ?boolean,
  +$fragmentRefs: CheckListTable_list$ref,
  +$refType: CheckListCategoryItemsDialog_items$ref,
|}>;
export type CheckListCategoryItemsDialog_items$data = CheckListCategoryItemsDialog_items;
export type CheckListCategoryItemsDialog_items$key = $ReadOnlyArray<{
  +$data?: CheckListCategoryItemsDialog_items$data,
  +$fragmentRefs: CheckListCategoryItemsDialog_items$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListCategoryItemsDialog_items",
  "type": "CheckListItem",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "checked",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "CheckListTable_list",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'eaa431ed222d3b0cee90cf7a747db214';
module.exports = node;
