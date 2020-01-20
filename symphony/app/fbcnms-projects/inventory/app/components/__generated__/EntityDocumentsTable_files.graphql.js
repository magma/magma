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
type DocumentTable_files$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type EntityDocumentsTable_files$ref: FragmentReference;
declare export opaque type EntityDocumentsTable_files$fragmentType: EntityDocumentsTable_files$ref;
export type EntityDocumentsTable_files = $ReadOnlyArray<{|
  +$fragmentRefs: DocumentTable_files$ref,
  +$refType: EntityDocumentsTable_files$ref,
|}>;
export type EntityDocumentsTable_files$data = EntityDocumentsTable_files;
export type EntityDocumentsTable_files$key = $ReadOnlyArray<{
  +$data?: EntityDocumentsTable_files$data,
  +$fragmentRefs: EntityDocumentsTable_files$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "EntityDocumentsTable_files",
  "type": "File",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "FragmentSpread",
      "name": "DocumentTable_files",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'bf947c22e479b0150cc5e2861374aa18';
module.exports = node;
