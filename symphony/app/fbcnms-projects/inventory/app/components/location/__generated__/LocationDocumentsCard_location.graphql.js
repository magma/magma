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
type EntityDocumentsTable_files$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type LocationDocumentsCard_location$ref: FragmentReference;
declare export opaque type LocationDocumentsCard_location$fragmentType: LocationDocumentsCard_location$ref;
export type LocationDocumentsCard_location = {|
  +id: string,
  +images: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +files: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +$refType: LocationDocumentsCard_location$ref,
|};
export type LocationDocumentsCard_location$data = LocationDocumentsCard_location;
export type LocationDocumentsCard_location$key = {
  +$data?: LocationDocumentsCard_location$data,
  +$fragmentRefs: LocationDocumentsCard_location$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = [
  {
    "kind": "FragmentSpread",
    "name": "EntityDocumentsTable_files",
    "args": null
  }
];
return {
  "kind": "Fragment",
  "name": "LocationDocumentsCard_location",
  "type": "Location",
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
      "kind": "LinkedField",
      "alias": null,
      "name": "images",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v0/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "files",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v0/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '08cd49b5ce8161141e23fd84378c8fe1';
module.exports = node;
