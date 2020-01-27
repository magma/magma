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
type DocumentTable_hyperlinks$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type EntityDocumentsTable_hyperlinks$ref: FragmentReference;
declare export opaque type EntityDocumentsTable_hyperlinks$fragmentType: EntityDocumentsTable_hyperlinks$ref;
export type EntityDocumentsTable_hyperlinks = $ReadOnlyArray<{|
  +$fragmentRefs: DocumentTable_hyperlinks$ref,
  +$refType: EntityDocumentsTable_hyperlinks$ref,
|}>;
export type EntityDocumentsTable_hyperlinks$data = EntityDocumentsTable_hyperlinks;
export type EntityDocumentsTable_hyperlinks$key = $ReadOnlyArray<{
  +$data?: EntityDocumentsTable_hyperlinks$data,
  +$fragmentRefs: EntityDocumentsTable_hyperlinks$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "EntityDocumentsTable_hyperlinks",
  "type": "Hyperlink",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "FragmentSpread",
      "name": "DocumentTable_hyperlinks",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'dff4323de5cb8d1ba32bfd2518647a88';
module.exports = node;
