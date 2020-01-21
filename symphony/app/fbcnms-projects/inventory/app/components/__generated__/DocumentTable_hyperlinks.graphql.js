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
type HyperlinkTableRow_hyperlink$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type DocumentTable_hyperlinks$ref: FragmentReference;
declare export opaque type DocumentTable_hyperlinks$fragmentType: DocumentTable_hyperlinks$ref;
export type DocumentTable_hyperlinks = $ReadOnlyArray<{|
  +id: string,
  +category: ?string,
  +$fragmentRefs: HyperlinkTableRow_hyperlink$ref,
  +$refType: DocumentTable_hyperlinks$ref,
|}>;
export type DocumentTable_hyperlinks$data = DocumentTable_hyperlinks;
export type DocumentTable_hyperlinks$key = $ReadOnlyArray<{
  +$data?: DocumentTable_hyperlinks$data,
  +$fragmentRefs: DocumentTable_hyperlinks$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "DocumentTable_hyperlinks",
  "type": "Hyperlink",
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
      "name": "category",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "HyperlinkTableRow_hyperlink",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'f74e5e118c2770d772174b292fb4266e';
module.exports = node;
