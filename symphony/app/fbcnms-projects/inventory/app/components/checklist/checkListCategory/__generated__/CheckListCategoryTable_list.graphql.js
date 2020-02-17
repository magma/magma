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
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListCategoryTable_list$ref: FragmentReference;
declare export opaque type CheckListCategoryTable_list$fragmentType: CheckListCategoryTable_list$ref;
export type CheckListCategoryTable_list = $ReadOnlyArray<{|
  +id: string,
  +title: string,
  +description: ?string,
  +checkList: $ReadOnlyArray<{|
    +id: string,
    +title: string,
    +type: CheckListItemType,
    +index: ?number,
    +helpText: ?string,
    +enumValues: ?string,
    +stringValue: ?string,
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


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "CheckListCategoryTable_list",
  "type": "CheckListCategory",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
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
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "type",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "index",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "helpText",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "enumValues",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "stringValue",
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
          "name": "CheckListCategoryItemsDialog_items",
          "args": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f061bb5f566a2141726de2fc6642f89c';
module.exports = node;
