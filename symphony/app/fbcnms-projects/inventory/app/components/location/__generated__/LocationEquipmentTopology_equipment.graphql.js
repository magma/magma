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
declare export opaque type LocationEquipmentTopology_equipment$ref: FragmentReference;
declare export opaque type LocationEquipmentTopology_equipment$fragmentType: LocationEquipmentTopology_equipment$ref;
export type LocationEquipmentTopology_equipment = $ReadOnlyArray<{|
  +id: string,
  +$refType: LocationEquipmentTopology_equipment$ref,
|}>;
export type LocationEquipmentTopology_equipment$data = LocationEquipmentTopology_equipment;
export type LocationEquipmentTopology_equipment$key = $ReadOnlyArray<{
  +$data?: LocationEquipmentTopology_equipment$data,
  +$fragmentRefs: LocationEquipmentTopology_equipment$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "LocationEquipmentTopology_equipment",
  "type": "Equipment",
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
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '8d782dfc3b488a1005aaeef451b77d56';
module.exports = node;
