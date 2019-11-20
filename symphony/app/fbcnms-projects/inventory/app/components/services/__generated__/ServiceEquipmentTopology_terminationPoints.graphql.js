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
declare export opaque type ServiceEquipmentTopology_terminationPoints$ref: FragmentReference;
declare export opaque type ServiceEquipmentTopology_terminationPoints$fragmentType: ServiceEquipmentTopology_terminationPoints$ref;
export type ServiceEquipmentTopology_terminationPoints = $ReadOnlyArray<{|
  +id: string,
  +$refType: ServiceEquipmentTopology_terminationPoints$ref,
|}>;
export type ServiceEquipmentTopology_terminationPoints$data = ServiceEquipmentTopology_terminationPoints;
export type ServiceEquipmentTopology_terminationPoints$key = $ReadOnlyArray<{
  +$data?: ServiceEquipmentTopology_terminationPoints$data,
  +$fragmentRefs: ServiceEquipmentTopology_terminationPoints$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ServiceEquipmentTopology_terminationPoints",
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
(node/*: any*/).hash = '0dbb364d6cbc900000ed6a6d6140b7ce';
module.exports = node;
