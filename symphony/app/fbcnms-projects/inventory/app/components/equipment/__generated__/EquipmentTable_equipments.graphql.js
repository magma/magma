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
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentTable_equipments$ref: FragmentReference;
declare export opaque type EquipmentTable_equipments$fragmentType: EquipmentTable_equipments$ref;
export type EquipmentTable_equipments = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +futureState: ?FutureState,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +workOrder: ?{|
    +id: string,
    +status: WorkOrderStatus,
  |},
  +device: ?{|
    +up: ?boolean
  |},
  +services: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +$refType: EquipmentTable_equipments$ref,
|}>;
export type EquipmentTable_equipments$data = EquipmentTable_equipments;
export type EquipmentTable_equipments$key = $ReadOnlyArray<{
  +$data?: EquipmentTable_equipments$data,
  +$fragmentRefs: EquipmentTable_equipments$ref,
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
  "name": "EquipmentTable_equipments",
  "type": "Equipment",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "futureState",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/)
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "workOrder",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkOrder",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "status",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "device",
      "storageKey": null,
      "args": null,
      "concreteType": "Device",
      "plural": false,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "up",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "services",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": true,
      "selections": [
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ad1708fe40398b5cab77735cbd8a6417';
module.exports = node;
