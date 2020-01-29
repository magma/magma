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
declare export opaque type CheckListTableFilling_list$ref: FragmentReference;
declare export opaque type CheckListTableFilling_list$fragmentType: CheckListTableFilling_list$ref;
export type CheckListTableFilling_list = $ReadOnlyArray<{|
  +id: string,
  +checked: ?boolean,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: CheckListTableFilling_list$ref,
|}>;
export type CheckListTableFilling_list$data = CheckListTableFilling_list;
export type CheckListTableFilling_list$key = $ReadOnlyArray<{
  +$data?: CheckListTableFilling_list$data,
  +$fragmentRefs: CheckListTableFilling_list$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListTableFilling_list",
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
(node/*: any*/).hash = '25c5ac03c595ad1f243a9fc17f713890';
module.exports = node;
