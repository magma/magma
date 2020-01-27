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
import type { FragmentReference } from "relay-runtime";
declare export opaque type ImageDialog_img$ref: FragmentReference;
declare export opaque type ImageDialog_img$fragmentType: ImageDialog_img$ref;
export type ImageDialog_img = {|
  +storeKey: ?string,
  +fileName: string,
  +$refType: ImageDialog_img$ref,
|};
export type ImageDialog_img$data = ImageDialog_img;
export type ImageDialog_img$key = {
  +$data?: ImageDialog_img$data,
  +$fragmentRefs: ImageDialog_img$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ImageDialog_img",
  "type": "File",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
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
      "name": "fileName",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '9df3ae53271a85ffc0bd704420104cc5';
module.exports = node;
