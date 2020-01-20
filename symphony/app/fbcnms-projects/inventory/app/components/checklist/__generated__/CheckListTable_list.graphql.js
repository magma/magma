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
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListTable_list$ref: FragmentReference;
declare export opaque type CheckListTable_list$fragmentType: CheckListTable_list$ref;
export type CheckListTable_list = $ReadOnlyArray<{|
  +id: string,
  +index: ?number,
  +type: CheckListItemType,
  +title: string,
  +checked: ?boolean,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: CheckListTable_list$ref,
|}>;
export type CheckListTable_list$data = CheckListTable_list;
export type CheckListTable_list$key = $ReadOnlyArray<{
  +$data?: CheckListTable_list$data,
  +$fragmentRefs: CheckListTable_list$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListTable_list",
  "type": "CheckListItem",
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
      "name": "index",
      "args": null,
      "storageKey": null
    },
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
(node/*: any*/).hash = '2b2e334a3aed7484b01b617b45b34e8b';
module.exports = node;
