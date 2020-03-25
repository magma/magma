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
declare export opaque type FreeTextCheckListItemFilling_item$ref: FragmentReference;
declare export opaque type FreeTextCheckListItemFilling_item$fragmentType: FreeTextCheckListItemFilling_item$ref;
export type FreeTextCheckListItemFilling_item = {|
  +title: string,
  +helpText: ?string,
  +stringValue: ?string,
  +checked: ?boolean,
  +$fragmentRefs: CheckListItem_item$ref,
  +$refType: FreeTextCheckListItemFilling_item$ref,
|};
export type FreeTextCheckListItemFilling_item$data = FreeTextCheckListItemFilling_item;
export type FreeTextCheckListItemFilling_item$key = {
  +$data?: FreeTextCheckListItemFilling_item$data,
  +$fragmentRefs: FreeTextCheckListItemFilling_item$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "FreeTextCheckListItemFilling_item",
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
      "name": "CheckListItem_item",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '42299731165e04e04bcba59bed45764e';
module.exports = node;
