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
declare export opaque type FreeTextCheckListItemDefinition_item$ref: FragmentReference;
declare export opaque type FreeTextCheckListItemDefinition_item$fragmentType: FreeTextCheckListItemDefinition_item$ref;
export type FreeTextCheckListItemDefinition_item = {|
  +title: string,
  +helpText: ?string,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: FreeTextCheckListItemDefinition_item$ref,
|};
export type FreeTextCheckListItemDefinition_item$data = FreeTextCheckListItemDefinition_item;
export type FreeTextCheckListItemDefinition_item$key = {
  +$data?: FreeTextCheckListItemDefinition_item$data,
  +$fragmentRefs: FreeTextCheckListItemDefinition_item$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "FreeTextCheckListItemDefinition_item",
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
      "name": "helpText",
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
(node/*: any*/).hash = '1f50cfff233bf20b8ff2d472fdce0014';
module.exports = node;
