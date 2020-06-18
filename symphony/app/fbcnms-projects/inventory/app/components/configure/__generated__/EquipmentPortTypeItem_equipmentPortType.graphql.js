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
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPortTypeItem_equipmentPortType$ref: FragmentReference;
declare export opaque type EquipmentPortTypeItem_equipmentPortType$fragmentType: EquipmentPortTypeItem_equipmentPortType$ref;
export type EquipmentPortTypeItem_equipmentPortType = {|
  +id: string,
  +name: string,
  +numberOfPortDefinitions: number,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: DynamicPropertyTypesGrid_propertyTypes$ref
  |}>,
  +linkPropertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: DynamicPropertyTypesGrid_propertyTypes$ref
  |}>,
  +$refType: EquipmentPortTypeItem_equipmentPortType$ref,
|};
export type EquipmentPortTypeItem_equipmentPortType$data = EquipmentPortTypeItem_equipmentPortType;
export type EquipmentPortTypeItem_equipmentPortType$key = {
  +$data?: EquipmentPortTypeItem_equipmentPortType$data,
  +$fragmentRefs: EquipmentPortTypeItem_equipmentPortType$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = [
  {
    "kind": "FragmentSpread",
    "name": "DynamicPropertyTypesGrid_propertyTypes",
    "args": null
  }
];
return {
  "kind": "Fragment",
  "name": "EquipmentPortTypeItem_equipmentPortType",
  "type": "EquipmentPortType",
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
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfPortDefinitions",
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
      "selections": (v0/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "linkPropertyTypes",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": true,
      "selections": (v0/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '64680fb091846a1ed120759b8199a89c';
module.exports = node;
