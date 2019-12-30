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
declare export opaque type PropertyFormField_property$ref: FragmentReference;
declare export opaque type PropertyFormField_property$fragmentType: PropertyFormField_property$ref;
export type PropertyFormField_property = {|
  +id: string,
  +propertyType: {|
    +id: string,
    +name: string,
    +type: PropertyKind,
    +isEditable: ?boolean,
    +isMandatory: ?boolean,
    +isInstanceProperty: ?boolean,
    +stringValue: ?string,
  |},
  +stringValue: ?string,
  +intValue: ?number,
  +floatValue: ?number,
  +booleanValue: ?boolean,
  +latitudeValue: ?number,
  +longitudeValue: ?number,
  +rangeFromValue: ?number,
  +rangeToValue: ?number,
  +equipmentValue: ?{|
    +id: string,
    +name: string,
  |},
  +locationValue: ?{|
    +id: string,
    +name: string,
  |},
  +serviceValue: ?{|
    +id: string,
    +name: string,
  |},
  +$refType: PropertyFormField_property$ref,
|};
export type PropertyFormField_property$data = PropertyFormField_property;
export type PropertyFormField_property$key = {
  +$data?: PropertyFormField_property$data,
  +$fragmentRefs: PropertyFormField_property$ref,
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
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v3 = [
  (v0/*: any*/),
  (v1/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "PropertyFormField_property",
  "type": "Property",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
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
        },
        (v2/*: any*/)
      ]
    },
    (v2/*: any*/),
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
      "name": "floatValue",
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
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v3/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v3/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": false,
      "selections": (v3/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '67524c0dbba3cf1a4b79dffefa3b014f';
module.exports = node;
