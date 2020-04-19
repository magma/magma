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
type PropertyTypeFormField_propertyType$ref = any;
type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceTypeItem_serviceType$ref: FragmentReference;
declare export opaque type ServiceTypeItem_serviceType$fragmentType: ServiceTypeItem_serviceType$ref;
export type ServiceTypeItem_serviceType = {|
  +id: string,
  +name: string,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: PropertyTypeFormField_propertyType$ref
  |}>,
  +endpointDefinitions: $ReadOnlyArray<?{|
    +$fragmentRefs: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref
  |}>,
  +numberOfServices: number,
  +$refType: ServiceTypeItem_serviceType$ref,
|};
export type ServiceTypeItem_serviceType$data = ServiceTypeItem_serviceType;
export type ServiceTypeItem_serviceType$key = {
  +$data?: ServiceTypeItem_serviceType$data,
  +$fragmentRefs: ServiceTypeItem_serviceType$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ServiceTypeItem_serviceType",
  "type": "ServiceType",
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
          "name": "PropertyTypeFormField_propertyType",
          "args": null
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
        {
          "kind": "FragmentSpread",
          "name": "ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions",
          "args": null
        }
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfServices",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '1e60a484ed66bed55f3de99645b11d03';
module.exports = node;
