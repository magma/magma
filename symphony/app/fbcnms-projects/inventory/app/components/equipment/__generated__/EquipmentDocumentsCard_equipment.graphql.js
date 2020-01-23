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
type EntityDocumentsTable_hyperlinks$ref = any;
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
  +hyperlinks: $ReadOnlyArray<{|
    +$fragmentRefs: EntityDocumentsTable_hyperlinks$ref
  |}>,
  +$refType: EquipmentDocumentsCard_equipment$ref,
|};
export type EquipmentDocumentsCard_equipment$data = EquipmentDocumentsCard_equipment;
export type EquipmentDocumentsCard_equipment$key = {
  +$data?: EquipmentDocumentsCard_equipment$data,
  +$fragmentRefs: EquipmentDocumentsCard_equipment$ref,
  ...
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
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "hyperlinks",
      "storageKey": null,
      "args": null,
      "concreteType": "Hyperlink",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "EntityDocumentsTable_hyperlinks",
          "args": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '47fa4ebbcddc93a1562ca029f897dc4e';
module.exports = node;
