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
type DynamicPropertiesGrid_properties$ref = any;
type DynamicPropertiesGrid_propertyTypes$ref = any;
type PropertyFormField_property$ref = any;
type PropertyTypeFormField_propertyType$ref = any;
export type DiscoveryMethod = "INVENTORY" | "MANUAL" | "%future added value";
export type ServiceStatus = "DISCONNECTED" | "IN_SERVICE" | "MAINTENANCE" | "PENDING" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServicesView_service$ref: FragmentReference;
declare export opaque type ServicesView_service$fragmentType: ServicesView_service$ref;
export type ServicesView_service = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +externalId: ?string,
  +status: ServiceStatus,
  +customer: ?{|
    +id: string,
    +name: string,
  |},
  +serviceType: {|
    +id: string,
    +name: string,
    +discoveryMethod: DiscoveryMethod,
    +propertyTypes: $ReadOnlyArray<?{|
      +$fragmentRefs: PropertyTypeFormField_propertyType$ref & DynamicPropertiesGrid_propertyTypes$ref
    |}>,
  |},
  +properties: $ReadOnlyArray<?{|
    +$fragmentRefs: PropertyFormField_property$ref & DynamicPropertiesGrid_properties$ref
  |}>,
  +$refType: ServicesView_service$ref,
|}>;
export type ServicesView_service$data = ServicesView_service;
export type ServicesView_service$key = $ReadOnlyArray<{
  +$data?: ServicesView_service$data,
  +$fragmentRefs: ServicesView_service$ref,
  ...
}>;
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
};
return {
  "kind": "Fragment",
  "name": "ServicesView_service",
  "type": "Service",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "externalId",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "customer",
      "storageKey": null,
      "args": null,
      "concreteType": "Customer",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/)
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceType",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceType",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
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
            {
              "kind": "FragmentSpread",
              "name": "PropertyTypeFormField_propertyType",
              "args": null
            },
            {
              "kind": "FragmentSpread",
              "name": "DynamicPropertiesGrid_propertyTypes",
              "args": null
            }
          ]
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "properties",
      "storageKey": null,
      "args": null,
      "concreteType": "Property",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "PropertyFormField_property",
          "args": null
        },
        {
          "kind": "FragmentSpread",
          "name": "DynamicPropertiesGrid_properties",
          "args": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '082bbb9aa27247bc35148e31b2a04903';
module.exports = node;
