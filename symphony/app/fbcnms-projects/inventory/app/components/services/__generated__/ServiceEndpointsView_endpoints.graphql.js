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
type EquipmentBreadcrumbs_equipment$ref = any;
export type ServiceEndpointRole = "CONSUMER" | "PROVIDER" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceEndpointsView_endpoints$ref: FragmentReference;
declare export opaque type ServiceEndpointsView_endpoints$fragmentType: ServiceEndpointsView_endpoints$ref;
export type ServiceEndpointsView_endpoints = $ReadOnlyArray<{|
  +id: string,
  +port: {|
    +parentEquipment: {|
      +name: string,
      +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
    |},
    +definition: {|
      +id: string,
      +name: string,
    |},
  |},
  +role: ServiceEndpointRole,
  +$refType: ServiceEndpointsView_endpoints$ref,
|}>;
export type ServiceEndpointsView_endpoints$data = ServiceEndpointsView_endpoints;
export type ServiceEndpointsView_endpoints$key = $ReadOnlyArray<{
  +$data?: ServiceEndpointsView_endpoints$data,
  +$fragmentRefs: ServiceEndpointsView_endpoints$ref,
}>;
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ServiceEndpointsView_endpoints",
  "type": "ServiceEndpoint",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "port",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPort",
      "plural": false,
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
            (v1/*: any*/),
            {
              "kind": "FragmentSpread",
              "name": "EquipmentBreadcrumbs_equipment",
              "args": null
            }
          ]
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "definition",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/)
          ]
        }
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "role",
      "args": null,
      "storageKey": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '6434214b1d518a116cccc1aa3dd9c192';
module.exports = node;
