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
declare export opaque type EquipmentDocumentsCard_equipment$ref: FragmentReference;
declare export opaque type EquipmentDocumentsCard_equipment$fragmentType: EquipmentDocumentsCard_equipment$ref;
export type EquipmentDocumentsCard_equipment = {|
  +id: string,
  +images: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +files: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +$refType: EquipmentDocumentsCard_equipment$ref,
|};
export type EquipmentDocumentsCard_equipment$data = EquipmentDocumentsCard_equipment;
export type EquipmentDocumentsCard_equipment$key = {
  +$data?: EquipmentDocumentsCard_equipment$data,
  +$fragmentRefs: EquipmentDocumentsCard_equipment$ref,
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
  "name": "EquipmentDocumentsCard_equipment",
  "type": "Equipment",
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
(node/*: any*/).hash = '48d00bc9a7faafee190884542a0489b3';
module.exports = node;
