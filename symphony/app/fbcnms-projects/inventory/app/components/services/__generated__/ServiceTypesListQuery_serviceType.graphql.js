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
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceTypesListQuery_serviceType$ref: FragmentReference;
declare export opaque type ServiceTypesListQuery_serviceType$fragmentType: ServiceTypesListQuery_serviceType$ref;
export type ServiceTypesListQuery_serviceType = {|
  +id: string,
  +name: string,
  +$refType: ServiceTypesListQuery_serviceType$ref,
|};
export type ServiceTypesListQuery_serviceType$data = ServiceTypesListQuery_serviceType;
export type ServiceTypesListQuery_serviceType$key = {
  +$data?: ServiceTypesListQuery_serviceType$data,
  +$fragmentRefs: ServiceTypesListQuery_serviceType$ref,
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
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '2228caffa61db69903ecdcc99299f54e';
module.exports = node;
