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
declare export opaque type ServiceEquipmentTopology_endpoints$ref: FragmentReference;
declare export opaque type ServiceEquipmentTopology_endpoints$fragmentType: ServiceEquipmentTopology_endpoints$ref;
export type ServiceEquipmentTopology_endpoints = $ReadOnlyArray<{|
  +definition: {|
    +role: ?string
  |},
  +equipment: {|
    +id: string,
    +positionHierarchy: $ReadOnlyArray<{|
      +parentEquipment: {|
        +id: string
      |}
    |}>,
  |},
  +$refType: ServiceEquipmentTopology_endpoints$ref,
|}>;
export type ServiceEquipmentTopology_endpoints$data = ServiceEquipmentTopology_endpoints;
export type ServiceEquipmentTopology_endpoints$key = $ReadOnlyArray<{
  +$data?: ServiceEquipmentTopology_endpoints$data,
  +$fragmentRefs: ServiceEquipmentTopology_endpoints$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ServiceEquipmentTopology_endpoints",
  "type": "ServiceEndpoint",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "definition",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceEndpointDefinition",
      "plural": false,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "role",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipment",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "positionHierarchy",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPosition",
          "plural": true,
          "selections": [
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "parentEquipment",
              "storageKey": null,
              "args": null,
              "concreteType": "Equipment",
              "plural": false,
              "selections": [
                (v0/*: any*/)
              ]
            }
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '085a9519bd88793b015ab955f716fb5f';
module.exports = node;
