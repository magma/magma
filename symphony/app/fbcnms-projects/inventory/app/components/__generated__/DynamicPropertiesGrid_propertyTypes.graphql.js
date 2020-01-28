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
declare export opaque type DynamicPropertiesGrid_propertyTypes$ref: FragmentReference;
declare export opaque type DynamicPropertiesGrid_propertyTypes$fragmentType: DynamicPropertiesGrid_propertyTypes$ref;
export type DynamicPropertiesGrid_propertyTypes = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +index: ?number,
  +isInstanceProperty: ?boolean,
  +type: PropertyKind,
  +stringValue: ?string,
  +intValue: ?number,
  +booleanValue: ?boolean,
  +latitudeValue: ?number,
  +longitudeValue: ?number,
  +rangeFromValue: ?number,
  +rangeToValue: ?number,
  +floatValue: ?number,
  +$refType: DynamicPropertiesGrid_propertyTypes$ref,
|}>;
export type DynamicPropertiesGrid_propertyTypes$data = DynamicPropertiesGrid_propertyTypes;
export type DynamicPropertiesGrid_propertyTypes$key = $ReadOnlyArray<{
  +$data?: DynamicPropertiesGrid_propertyTypes$data,
  +$fragmentRefs: DynamicPropertiesGrid_propertyTypes$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "DynamicPropertiesGrid_propertyTypes",
  "type": "PropertyType",
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
      "name": "index",
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
      "name": "type",
      "args": null,
      "storageKey": null
    },
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
      "name": "floatValue",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '352ce07fece51150af9945c3ea537d3f';
module.exports = node;
