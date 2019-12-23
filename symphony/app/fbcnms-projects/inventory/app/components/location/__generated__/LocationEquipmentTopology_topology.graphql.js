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
declare export opaque type LocationEquipmentTopology_topology$ref: FragmentReference;
declare export opaque type LocationEquipmentTopology_topology$fragmentType: LocationEquipmentTopology_topology$ref;
export type LocationEquipmentTopology_topology = {|
  +nodes: $ReadOnlyArray<{|
    +id?: string,
    +name?: string,
  |}>,
  +$fragmentRefs: ForceNetworkTopology_topology$ref,
  +$refType: LocationEquipmentTopology_topology$ref,
|};
export type LocationEquipmentTopology_topology$data = LocationEquipmentTopology_topology;
export type LocationEquipmentTopology_topology$key = {
  +$data?: LocationEquipmentTopology_topology$data,
  +$fragmentRefs: LocationEquipmentTopology_topology$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "LocationEquipmentTopology_topology",
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
(node/*: any*/).hash = '85e688a255ed451854dcd0a64bd6cded';
module.exports = node;
