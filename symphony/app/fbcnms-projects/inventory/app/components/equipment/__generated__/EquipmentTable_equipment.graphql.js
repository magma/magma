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
declare export opaque type EquipmentTable_equipment$ref: FragmentReference;
declare export opaque type EquipmentTable_equipment$fragmentType: EquipmentTable_equipment$ref;
export type EquipmentTable_equipment = $ReadOnlyArray<{|
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
  +$refType: EquipmentTable_equipment$ref,
|}>;
export type EquipmentTable_equipment$data = EquipmentTable_equipment;
export type EquipmentTable_equipment$key = $ReadOnlyArray<{
  +$data?: EquipmentTable_equipment$data,
  +$fragmentRefs: EquipmentTable_equipment$ref,
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
  "name": "EquipmentTable_equipment",
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
(node/*: any*/).hash = 'cd50b22725719d42497777ab6583e2e5';
module.exports = node;
