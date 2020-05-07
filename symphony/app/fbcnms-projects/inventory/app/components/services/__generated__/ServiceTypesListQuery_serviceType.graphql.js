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
export type DiscoveryMethod = "INVENTORY" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceTypesListQuery_serviceType$ref: FragmentReference;
declare export opaque type ServiceTypesListQuery_serviceType$fragmentType: ServiceTypesListQuery_serviceType$ref;
export type ServiceTypesListQuery_serviceType = {|
  +id: string,
  +name: string,
  +discoveryMethod: ?DiscoveryMethod,
  +$refType: ServiceTypesListQuery_serviceType$ref,
|};
export type ServiceTypesListQuery_serviceType$data = ServiceTypesListQuery_serviceType;
export type ServiceTypesListQuery_serviceType$key = {
  +$data?: ServiceTypesListQuery_serviceType$data,
  +$fragmentRefs: ServiceTypesListQuery_serviceType$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ServiceTypesListQuery_serviceType",
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
      "kind": "ScalarField",
      "alias": null,
      "name": "discoveryMethod",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '746b4b7d14876a363b38950f2f7ab94d';
module.exports = node;
