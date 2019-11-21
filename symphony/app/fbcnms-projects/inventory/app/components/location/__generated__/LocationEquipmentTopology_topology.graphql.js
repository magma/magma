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
declare export opaque type LocationEquipmentTopology_topology$ref: FragmentReference;
declare export opaque type LocationEquipmentTopology_topology$fragmentType: LocationEquipmentTopology_topology$ref;
export type LocationEquipmentTopology_topology = {|
  +nodes: $ReadOnlyArray<{|
    +id: string,
    +name: string,
  |}>,
  +links: $ReadOnlyArray<{|
    +source: string,
    +target: string,
  |}>,
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
(node/*: any*/).hash = 'ac372f45e9fa24b0cbad523906bd6743';
module.exports = node;
