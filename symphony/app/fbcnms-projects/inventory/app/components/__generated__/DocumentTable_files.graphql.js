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
type FileAttachment_file$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type DocumentTable_files$ref: FragmentReference;
declare export opaque type DocumentTable_files$fragmentType: DocumentTable_files$ref;
export type DocumentTable_files = $ReadOnlyArray<{|
  +id: string,
  +fileName: string,
  +category: ?string,
  +$fragmentRefs: FileAttachment_file$ref,
  +$refType: DocumentTable_files$ref,
|}>;
export type DocumentTable_files$data = DocumentTable_files;
export type DocumentTable_files$key = $ReadOnlyArray<{
  +$data?: DocumentTable_files$data,
  +$fragmentRefs: DocumentTable_files$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "DocumentTable_files",
  "type": "File",
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
      "name": "fileName",
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
      "name": "FileAttachment_file",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '171a55fcd66fd996e20dd78e3d3db780';
module.exports = node;
