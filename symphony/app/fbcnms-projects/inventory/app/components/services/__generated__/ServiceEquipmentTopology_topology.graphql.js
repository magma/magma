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
type ForceNetworkTopology_topology$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceEquipmentTopology_topology$ref: FragmentReference;
declare export opaque type ServiceEquipmentTopology_topology$fragmentType: ServiceEquipmentTopology_topology$ref;
export type ServiceEquipmentTopology_topology = {|
  +nodes: $ReadOnlyArray<{|
    +id?: string,
    +name?: string,
  |}>,
  +$fragmentRefs: ForceNetworkTopology_topology$ref,
  +$refType: ServiceEquipmentTopology_topology$ref,
|};
export type ServiceEquipmentTopology_topology$data = ServiceEquipmentTopology_topology;
export type ServiceEquipmentTopology_topology$key = {
  +$data?: ServiceEquipmentTopology_topology$data,
  +$fragmentRefs: ServiceEquipmentTopology_topology$ref,
  ...
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
      "concreteType": null,
      "plural": true,
      "selections": [
        {
          "kind": "InlineFragment",
          "type": "Equipment",
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
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "ForceNetworkTopology_topology",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '3aacf7059c6ccea42895132c69cbb030';
module.exports = node;
