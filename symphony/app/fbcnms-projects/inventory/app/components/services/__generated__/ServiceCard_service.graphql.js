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
type ServiceEquipmentTopology_terminationPoints$ref = any;
type ServiceEquipmentTopology_topology$ref = any;
type ServicePanel_service$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceCard_service$ref: FragmentReference;
declare export opaque type ServiceCard_service$fragmentType: ServiceCard_service$ref;
export type ServiceCard_service = {|
  +id: string,
  +name: string,
  +terminationPoints: $ReadOnlyArray<?{|
    +$fragmentRefs: ServiceEquipmentTopology_terminationPoints$ref
  |}>,
  +topology: {|
    +$fragmentRefs: ServiceEquipmentTopology_topology$ref
  |},
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
(node/*: any*/).hash = '8282fa4a93fe9c919d3f42b9dc90f1b2';
module.exports = node;
