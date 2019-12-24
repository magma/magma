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
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type PowerSearchEquipmentResultsTable_equipment$ref: FragmentReference;
declare export opaque type PowerSearchEquipmentResultsTable_equipment$fragmentType: PowerSearchEquipmentResultsTable_equipment$ref;
export type PowerSearchEquipmentResultsTable_equipment = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +futureState: ?FutureState,
  +externalId: ?string,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +workOrder: ?{|
    +id: string,
    +status: WorkOrderStatus,
  |},
  +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
  +$refType: PowerSearchEquipmentResultsTable_equipment$ref,
|}>;
export type PowerSearchEquipmentResultsTable_equipment$data = PowerSearchEquipmentResultsTable_equipment;
export type PowerSearchEquipmentResultsTable_equipment$key = $ReadOnlyArray<{
  +$data?: PowerSearchEquipmentResultsTable_equipment$data,
  +$fragmentRefs: PowerSearchEquipmentResultsTable_equipment$ref,
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
  "name": "PowerSearchEquipmentResultsTable_equipment",
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
      "kind": "ScalarField",
      "alias": null,
      "name": "externalId",
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
      "kind": "FragmentSpread",
      "name": "EquipmentBreadcrumbs_equipment",
      "args": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'dbd45296395196af507b9270ca136663';
module.exports = node;
