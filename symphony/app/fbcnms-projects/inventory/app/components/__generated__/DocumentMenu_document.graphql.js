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
export type FileType = "FILE" | "IMAGE" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type DocumentMenu_document$ref: FragmentReference;
declare export opaque type DocumentMenu_document$fragmentType: DocumentMenu_document$ref;
export type DocumentMenu_document = {|
  +id: string,
  +fileName: string,
  +storeKey: ?string,
  +fileType: ?FileType,
  +$refType: DocumentMenu_document$ref,
|};
export type DocumentMenu_document$data = DocumentMenu_document;
export type DocumentMenu_document$key = {
  +$data?: DocumentMenu_document$data,
  +$fragmentRefs: DocumentMenu_document$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "DocumentMenu_document",
  "type": "File",
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
      "name": "fileName",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "storeKey",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "fileType",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'e3bd50bdcbb80190e36c2e009e390b52';
module.exports = node;
