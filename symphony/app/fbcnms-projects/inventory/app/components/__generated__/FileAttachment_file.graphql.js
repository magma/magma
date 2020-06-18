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
type ImageDialog_img$ref = any;
export type FileType = "FILE" | "IMAGE" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type FileAttachment_file$ref: FragmentReference;
declare export opaque type FileAttachment_file$fragmentType: FileAttachment_file$ref;
export type FileAttachment_file = {|
  +id: string,
  +fileName: string,
  +sizeInBytes: ?number,
  +uploaded: ?any,
  +fileType: ?FileType,
  +storeKey: ?string,
  +category: ?string,
  +$fragmentRefs: ImageDialog_img$ref,
  +$refType: FileAttachment_file$ref,
|};
export type FileAttachment_file$data = FileAttachment_file;
export type FileAttachment_file$key = {
  +$data?: FileAttachment_file$data,
  +$fragmentRefs: FileAttachment_file$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "FileAttachment_file",
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
      "name": "sizeInBytes",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "uploaded",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "fileType",
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
      "name": "category",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "ImageDialog_img",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '281b6befd7d441674861d38032feb4e5';
module.exports = node;
