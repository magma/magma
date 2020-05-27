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
export type DiscoveryMethod = "INVENTORY" | "MANUAL" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type AddEditServiceTypeCard_editingServiceType$ref: FragmentReference;
declare export opaque type AddEditServiceTypeCard_editingServiceType$fragmentType: AddEditServiceTypeCard_editingServiceType$ref;
export type AddEditServiceTypeCard_editingServiceType = {|
  +id: string,
  +name: string,
  +numberOfServices: number,
  +discoveryMethod: DiscoveryMethod,
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
    +isMandatory: ?boolean,
    +isInstanceProperty: ?boolean,
  |}>,
  +endpointDefinitions: $ReadOnlyArray<{|
    +id: string,
    +index: number,
    +role: ?string,
    +name: string,
    +equipmentType: {|
      +name: string,
      +id: string,
    |},
  |}>,
  +$refType: AddEditServiceTypeCard_editingServiceType$ref,
|};
export type AddEditServiceTypeCard_editingServiceType$data = AddEditServiceTypeCard_editingServiceType;
export type AddEditServiceTypeCard_editingServiceType$key = {
  +$data?: AddEditServiceTypeCard_editingServiceType$data,
  +$fragmentRefs: AddEditServiceTypeCard_editingServiceType$ref,
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
};
return {
  "kind": "Fragment",
  "name": "AddEditServiceTypeCard_editingServiceType",
  "type": "ServiceType",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfServices",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "discoveryMethod",
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
          "name": "isMandatory",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "isInstanceProperty",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "endpointDefinitions",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceEndpointDefinition",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        (v2/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "role",
          "args": null,
          "storageKey": null
        },
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "equipmentType",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentType",
          "plural": false,
          "selections": [
            (v1/*: any*/),
            (v0/*: any*/)
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a30bbf291d8bcee70f5da0491ec834a7';
module.exports = node;
