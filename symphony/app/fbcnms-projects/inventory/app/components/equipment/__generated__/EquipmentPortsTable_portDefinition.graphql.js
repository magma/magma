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
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPortsTable_portDefinition$ref: FragmentReference;
declare export opaque type EquipmentPortsTable_portDefinition$fragmentType: EquipmentPortsTable_portDefinition$ref;
export type EquipmentPortsTable_portDefinition = {|
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
  +portType: ?{|
    +id: string,
    +name: string,
    +propertyTypes: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +type: PropertyKind,
      +nodeType: ?string,
      +index: ?number,
      +stringValue: ?string,
      +intValue: ?number,
      +booleanValue: ?boolean,
      +floatValue: ?number,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +isEditable: ?boolean,
      +isInstanceProperty: ?boolean,
      +isMandatory: ?boolean,
      +category: ?string,
      +isDeleted: ?boolean,
    |}>,
    +linkPropertyTypes: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +type: PropertyKind,
      +nodeType: ?string,
      +index: ?number,
      +stringValue: ?string,
      +intValue: ?number,
      +booleanValue: ?boolean,
      +floatValue: ?number,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +isEditable: ?boolean,
      +isInstanceProperty: ?boolean,
      +isMandatory: ?boolean,
      +category: ?string,
      +isDeleted: ?boolean,
    |}>,
  |},
  +$refType: EquipmentPortsTable_portDefinition$ref,
|};
export type EquipmentPortsTable_portDefinition$data = EquipmentPortsTable_portDefinition;
export type EquipmentPortsTable_portDefinition$key = {
  +$data?: EquipmentPortsTable_portDefinition$data,
  +$fragmentRefs: EquipmentPortsTable_portDefinition$ref,
  ...
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
v3 = [
  (v0/*: any*/),
  (v1/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "type",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "nodeType",
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
    "name": "rangeFromValue",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "rangeToValue",
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
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "category",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isDeleted",
    "args": null,
    "storageKey": null
  }
];
return {
  "kind": "Fragment",
  "name": "EquipmentPortsTable_portDefinition",
  "type": "EquipmentPortDefinition",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    (v2/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "visibleLabel",
      "args": null,
      "storageKey": null
    },
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
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "propertyTypes",
          "storageKey": null,
          "args": null,
          "concreteType": "PropertyType",
          "plural": true,
          "selections": (v3/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "linkPropertyTypes",
          "storageKey": null,
          "args": null,
          "concreteType": "PropertyType",
          "plural": true,
          "selections": (v3/*: any*/)
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c3402192a79fc61b03fc615e09549a3e';
module.exports = node;
