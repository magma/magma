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
declare export opaque type ServiceEquipmentTopology_topology$ref: FragmentReference;
declare export opaque type ServiceEquipmentTopology_topology$fragmentType: ServiceEquipmentTopology_topology$ref;
export type ServiceEquipmentTopology_topology = {|
  +nodes: $ReadOnlyArray<{|
    +id: string,
    +name: string,
  |}>,
  +links: $ReadOnlyArray<{|
    +source: string,
    +target: string,
  |}>,
  +$refType: ServiceEquipmentTopology_topology$ref,
|};
export type ServiceEquipmentTopology_topology$data = ServiceEquipmentTopology_topology;
export type ServiceEquipmentTopology_topology$key = {
  +$data?: ServiceEquipmentTopology_topology$data,
  +$fragmentRefs: ServiceEquipmentTopology_topology$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ServiceEquipmentTopology_topology",
  "type": "NetworkTopology",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "nodes",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
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
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "links",
      "storageKey": null,
      "args": null,
      "concreteType": "TopologyLink",
      "plural": true,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "source",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "target",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '63c26d79e91218be27e4d2ced941d0ff';
module.exports = node;
