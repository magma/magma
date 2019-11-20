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
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceTypeItem_serviceType$ref: FragmentReference;
declare export opaque type ServiceTypeItem_serviceType$fragmentType: ServiceTypeItem_serviceType$ref;
export type ServiceTypeItem_serviceType = {|
  +id: string,
  +name: string,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: PropertyTypeFormField_propertyType$ref
  |}>,
  +numberOfServices: number,
  +$refType: ServiceTypeItem_serviceType$ref,
|};
export type ServiceTypeItem_serviceType$data = ServiceTypeItem_serviceType;
export type ServiceTypeItem_serviceType$key = {
  +$data?: ServiceTypeItem_serviceType$data,
  +$fragmentRefs: ServiceTypeItem_serviceType$ref,
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
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfServices",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '64aa4c97eaf3183bfa12c3cd93b03c4b';
module.exports = node;
