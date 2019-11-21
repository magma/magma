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
type DocumentMenu_document$ref = any;
type ImageDialog_img$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ImageAttachment_img$ref: FragmentReference;
declare export opaque type ImageAttachment_img$fragmentType: ImageAttachment_img$ref;
export type ImageAttachment_img = {
  +id: string,
  +storeKey: ?string,
  +$fragmentRefs: DocumentMenu_document$ref & ImageDialog_img$ref,
};
export type ImageAttachment_img$data = ImageAttachment_img;
export type ImageAttachment_img$key = {
  +$data?: ImageAttachment_img$data,
  +$fragmentRefs: ImageAttachment_img$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ImageAttachment_img",
  "type": "File",
  "metadata": {
    "mask": false
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
      "name": "storeKey",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "DocumentMenu_document",
      "args": null
    },
    {
      "kind": "FragmentSpread",
      "name": "ImageDialog_img",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '67ad749bfc49740be43bf4cb43cfa561';
module.exports = node;
