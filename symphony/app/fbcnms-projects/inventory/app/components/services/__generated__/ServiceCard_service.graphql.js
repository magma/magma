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
type ServiceDetailsPanel_service$ref = any;
type ServiceEquipmentTopology_endpoints$ref = any;
type ServiceEquipmentTopology_topology$ref = any;
type ServicePanel_service$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceCard_service$ref: FragmentReference;
declare export opaque type ServiceCard_service$fragmentType: ServiceCard_service$ref;
export type ServiceCard_service = {|
  +id: string,
  +name: string,
  +topology: {|
    +$fragmentRefs: ServiceEquipmentTopology_topology$ref
  |},
  +endpoints: $ReadOnlyArray<?{|
    +$fragmentRefs: ServiceEquipmentTopology_endpoints$ref
  |}>,
  +$fragmentRefs: ServiceDetailsPanel_service$ref & ServicePanel_service$ref,
  +$refType: ServiceCard_service$ref,
|};
export type ServiceCard_service$data = ServiceCard_service;
export type ServiceCard_service$key = {
  +$data?: ServiceCard_service$data,
  +$fragmentRefs: ServiceCard_service$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ServiceCard_service",
  "type": "Service",
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
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "endpoints",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceEndpoint",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "ServiceEquipmentTopology_endpoints",
          "args": null
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "ServiceDetailsPanel_service",
      "args": null
    },
    {
      "kind": "FragmentSpread",
      "name": "ServicePanel_service",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'b365b307bdc31d3d737c3f7f1b6d33fe';
module.exports = node;
