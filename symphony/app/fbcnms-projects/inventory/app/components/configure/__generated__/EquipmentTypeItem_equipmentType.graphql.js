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
type DynamicPropertyTypesGrid_propertyTypes$ref = any;
type PortDefinitionsTable_portDefinitions$ref = any;
type PositionDefinitionsTable_positionDefinitions$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentTypeItem_equipmentType$ref: FragmentReference;
declare export opaque type EquipmentTypeItem_equipmentType$fragmentType: EquipmentTypeItem_equipmentType$ref;
export type EquipmentTypeItem_equipmentType = {|
  +id: string,
  +name: string,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: DynamicPropertyTypesGrid_propertyTypes$ref
  |}>,
  +positionDefinitions: $ReadOnlyArray<?{|
    +$fragmentRefs: PositionDefinitionsTable_positionDefinitions$ref
  |}>,
  +portDefinitions: $ReadOnlyArray<?{|
    +$fragmentRefs: PortDefinitionsTable_portDefinitions$ref
  |}>,
  +numberOfEquipment: number,
  +$refType: EquipmentTypeItem_equipmentType$ref,
|};
export type EquipmentTypeItem_equipmentType$data = EquipmentTypeItem_equipmentType;
export type EquipmentTypeItem_equipmentType$key = {
  +$data?: EquipmentTypeItem_equipmentType$data,
  +$fragmentRefs: EquipmentTypeItem_equipmentType$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "EquipmentTypeItem_equipmentType",
  "type": "EquipmentType",
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
      "kind": "ScalarField",
      "alias": null,
      "name": "name",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyTypes",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "DynamicPropertyTypesGrid_propertyTypes",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "positionDefinitions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPositionDefinition",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "PositionDefinitionsTable_positionDefinitions",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portDefinitions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortDefinition",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "PortDefinitionsTable_portDefinitions",
          "args": null
        }
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfEquipment",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'c31a356dc2022425956e57facd0ae246';
module.exports = node;
