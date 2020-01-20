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
declare export opaque type ForceNetworkTopology_topology$ref: FragmentReference;
declare export opaque type ForceNetworkTopology_topology$fragmentType: ForceNetworkTopology_topology$ref;
export type ForceNetworkTopology_topology = {|
  +nodes: $ReadOnlyArray<{|
    +id: string
  |}>,
  +links: $ReadOnlyArray<{|
    +source: {|
      +id: string
    |},
    +target: {|
      +id: string
    |},
  |}>,
  +$refType: ForceNetworkTopology_topology$ref,
|};
export type ForceNetworkTopology_topology$data = ForceNetworkTopology_topology;
export type ForceNetworkTopology_topology$key = {
  +$data?: ForceNetworkTopology_topology$data,
  +$fragmentRefs: ForceNetworkTopology_topology$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "id",
    "args": null,
    "storageKey": null
  }
];
return {
  "kind": "Fragment",
  "name": "ForceNetworkTopology_topology",
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
      "selections": (v0/*: any*/)
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
          "kind": "LinkedField",
          "alias": null,
          "name": "source",
          "storageKey": null,
          "args": null,
          "concreteType": null,
          "plural": false,
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "target",
          "storageKey": null,
          "args": null,
          "concreteType": null,
          "plural": false,
          "selections": (v0/*: any*/)
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '63b721aff87697366221cbc6250df26e';
module.exports = node;
