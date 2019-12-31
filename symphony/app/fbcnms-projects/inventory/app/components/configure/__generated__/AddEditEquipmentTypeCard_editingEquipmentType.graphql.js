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
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type AddEditEquipmentTypeCard_editingEquipmentType$ref: FragmentReference;
declare export opaque type AddEditEquipmentTypeCard_editingEquipmentType$fragmentType: AddEditEquipmentTypeCard_editingEquipmentType$ref;
export type AddEditEquipmentTypeCard_editingEquipmentType = {|
  +id: string,
  +name: string,
  +propertyTypes: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +type: PropertyKind,
    +index: ?number,
    +stringValue: ?string,
    +intValue: ?number,
    +booleanValue: ?boolean,
    +floatValue: ?number,
    +latitudeValue: ?number,
    +longitudeValue: ?number,
    +isEditable: ?boolean,
    +isInstanceProperty: ?boolean,
    +isMandatory: ?boolean,
  |}>,
  +positionDefinitions: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +index: ?number,
    +visibleLabel: ?string,
  |}>,
  +portDefinitions: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +index: ?number,
    +visibleLabel: ?string,
    +portType: ?{|
      +id: string,
      +name: string,
    |},
  |}>,
  +numberOfEquipment: number,
  +$refType: AddEditEquipmentTypeCard_editingEquipmentType$ref,
|};
export type AddEditEquipmentTypeCard_editingEquipmentType$data = AddEditEquipmentTypeCard_editingEquipmentType;
export type AddEditEquipmentTypeCard_editingEquipmentType$key = {
  +$data?: AddEditEquipmentTypeCard_editingEquipmentType$data,
  +$fragmentRefs: AddEditEquipmentTypeCard_editingEquipmentType$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "AddEditEquipmentTypeCard_editingEquipmentType",
  "type": "EquipmentType",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyTypes",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "type",
          "args": null,
          "storageKey": null
        },
        (v2/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "stringValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "intValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "booleanValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "floatValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "latitudeValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "longitudeValue",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "isEditable",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "isInstanceProperty",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "isMandatory",
          "args": null,
          "storageKey": null
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
        (v0/*: any*/),
        (v1/*: any*/),
        (v2/*: any*/),
        (v3/*: any*/)
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
        (v0/*: any*/),
        (v1/*: any*/),
        (v2/*: any*/),
        (v3/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "portType",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortType",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/)
          ]
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
})();
// prettier-ignore
(node/*: any*/).hash = '6aa03df6f5bf29d5f15ae954fe87d0b6';
module.exports = node;
