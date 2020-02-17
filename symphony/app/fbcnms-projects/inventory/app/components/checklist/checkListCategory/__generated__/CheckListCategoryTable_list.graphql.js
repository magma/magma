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
type CheckListCategoryItemsDialog_items$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListCategoryTable_list$ref: FragmentReference;
declare export opaque type CheckListCategoryTable_list$fragmentType: CheckListCategoryTable_list$ref;
export type CheckListCategoryTable_list = $ReadOnlyArray<{|
  +id: string,
  +title: string,
  +description: ?string,
  +checkList: $ReadOnlyArray<{|
    +checked: ?boolean,
    +$fragmentRefs: CheckListCategoryItemsDialog_items$ref,
  |}>,
  +$refType: CheckListCategoryTable_list$ref,
|}>;
export type CheckListCategoryTable_list$data = CheckListCategoryTable_list;
export type CheckListCategoryTable_list$key = $ReadOnlyArray<{
  +$data?: CheckListCategoryTable_list$data,
  +$fragmentRefs: CheckListCategoryTable_list$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListCategoryTable_list",
  "type": "CheckListCategory",
  "metadata": {
    "plural": true
  },
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
      "name": "title",
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
      "kind": "LinkedField",
      "alias": null,
      "name": "checkList",
      "storageKey": null,
      "args": null,
      "concreteType": "CheckListItem",
      "plural": true,
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
          "name": "CheckListCategoryItemsDialog_items",
          "args": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '5e0f7a136c99c45462629626f1c9088a';
module.exports = node;
