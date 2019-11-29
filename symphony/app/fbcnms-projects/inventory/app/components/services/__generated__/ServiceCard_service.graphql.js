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
type ServiceEquipmentTopology_terminationPoints$ref = any;
type ServiceEquipmentTopology_topology$ref = any;
type ServiceLinksView_links$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceCard_service$ref: FragmentReference;
declare export opaque type ServiceCard_service$fragmentType: ServiceCard_service$ref;
export type ServiceCard_service = {|
  +id: string,
  +name: string,
  +externalId: ?string,
  +customer: ?{|
    +name: string
  |},
  +serviceType: {|
    +id: string,
    +name: string,
    +propertyTypes: $ReadOnlyArray<?{|
      +$fragmentRefs: PropertyTypeFormField_propertyType$ref & DynamicPropertiesGrid_propertyTypes$ref
    |}>,
  |},
  +properties: $ReadOnlyArray<?{|
    +$fragmentRefs: PropertyFormField_property$ref & DynamicPropertiesGrid_properties$ref
  |}>,
  +links: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: ServiceLinksView_links$ref,
  |}>,
  +terminationPoints: $ReadOnlyArray<?{|
    +$fragmentRefs: ServiceEquipmentTopology_terminationPoints$ref
  |}>,
  +topology: {|
    +$fragmentRefs: ServiceEquipmentTopology_topology$ref
  |},
  +$refType: ServiceCard_service$ref,
|};
export type ServiceCard_service$data = ServiceCard_service;
export type ServiceCard_service$key = {
  +$data?: ServiceCard_service$data,
  +$fragmentRefs: ServiceCard_service$ref,
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
};
return {
  "kind": "Fragment",
  "name": "ServiceCard_service",
  "type": "Service",
  "metadata": null,
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
      "kind": "LinkedField",
      "alias": null,
      "name": "customer",
      "storageKey": null,
      "args": null,
      "concreteType": "Customer",
      "plural": false,
      "selections": [
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
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "links",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "FragmentSpread",
          "name": "ServiceLinksView_links",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "terminationPoints",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "ServiceEquipmentTopology_terminationPoints",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "topology",
      "storageKey": null,
      "args": null,
      "concreteType": "NetworkTopology",
      "plural": false,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "ServiceEquipmentTopology_topology",
          "args": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd7a8b953b12099588b0981685472c65c';
module.exports = node;
