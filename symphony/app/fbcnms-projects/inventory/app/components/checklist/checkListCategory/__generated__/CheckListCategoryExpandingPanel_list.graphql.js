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
type CheckListCategoryTable_list$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListCategoryExpandingPanel_list$ref: FragmentReference;
declare export opaque type CheckListCategoryExpandingPanel_list$fragmentType: CheckListCategoryExpandingPanel_list$ref;
export type CheckListCategoryExpandingPanel_list = $ReadOnlyArray<{|
  +$fragmentRefs: CheckListCategoryTable_list$ref,
  +$refType: CheckListCategoryExpandingPanel_list$ref,
|}>;
export type CheckListCategoryExpandingPanel_list$data = CheckListCategoryExpandingPanel_list;
export type CheckListCategoryExpandingPanel_list$key = $ReadOnlyArray<{
  +$data?: CheckListCategoryExpandingPanel_list$data,
  +$fragmentRefs: CheckListCategoryExpandingPanel_list$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListCategoryExpandingPanel_list",
  "type": "CheckListCategory",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "FragmentSpread",
      "name": "CheckListCategoryTable_list",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '3cd9ca1fd31f6f622c38b3209c1ad1e2';
module.exports = node;
