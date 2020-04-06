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
declare export opaque type CheckListTableDefinition_list$ref: FragmentReference;
declare export opaque type CheckListTableDefinition_list$fragmentType: CheckListTableDefinition_list$ref;
export type CheckListTableDefinition_list = $ReadOnlyArray<{|
  +id: string,
  +type: CheckListItemType,
  +index: ?number,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: CheckListTableDefinition_list$ref,
|}>;
export type CheckListTableDefinition_list$data = CheckListTableDefinition_list;
export type CheckListTableDefinition_list$key = $ReadOnlyArray<{
  +$data?: CheckListTableDefinition_list$data,
  +$fragmentRefs: CheckListTableDefinition_list$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListTableDefinition_list",
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
      "kind": "FragmentSpread",
      "name": "CheckListItem_item",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'e152c66dc9bf157465de5bdef97712e9';
module.exports = node;
