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
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListItem_item$ref: FragmentReference;
declare export opaque type CheckListItem_item$fragmentType: CheckListItem_item$ref;
export type CheckListItem_item = {|
  +id: string,
  +title: string,
  +type: CheckListItemType,
  +index: ?number,
  +helpText: ?string,
  +enumValues: ?string,
  +stringValue: ?string,
  +checked: ?boolean,
  +$refType: CheckListItem_item$ref,
|};
export type CheckListItem_item$data = CheckListItem_item;
export type CheckListItem_item$key = {
  +$data?: CheckListItem_item$data,
  +$fragmentRefs: CheckListItem_item$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CheckListItem_item",
  "type": "CheckListItem",
  "metadata": null,
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
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '4627620778b0360e1335dc3c4811e918';
module.exports = node;
