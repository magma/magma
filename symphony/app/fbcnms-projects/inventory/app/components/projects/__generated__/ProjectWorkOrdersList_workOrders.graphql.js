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
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ProjectWorkOrdersList_workOrders$ref: FragmentReference;
declare export opaque type ProjectWorkOrdersList_workOrders$fragmentType: ProjectWorkOrdersList_workOrders$ref;
export type ProjectWorkOrdersList_workOrders = $ReadOnlyArray<{|
  +id: string,
  +workOrderType: {|
    +name: string,
    +id: string,
  |},
  +name: string,
  +description: ?string,
  +ownerName: string,
  +creationDate: any,
  +installDate: ?any,
  +status: WorkOrderStatus,
  +priority: WorkOrderPriority,
  +$refType: ProjectWorkOrdersList_workOrders$ref,
|}>;
export type ProjectWorkOrdersList_workOrders$data = ProjectWorkOrdersList_workOrders;
export type ProjectWorkOrdersList_workOrders$key = $ReadOnlyArray<{
  +$data?: ProjectWorkOrdersList_workOrders$data,
  +$fragmentRefs: ProjectWorkOrdersList_workOrders$ref,
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
  "name": "ProjectWorkOrdersList_workOrders",
  "type": "WorkOrder",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "workOrderType",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkOrderType",
      "plural": false,
      "selections": [
        (v1/*: any*/),
        (v0/*: any*/)
      ]
    },
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "description",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "ownerName",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "creationDate",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "installDate",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "priority",
      "args": null,
      "storageKey": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b4bc6743797a4049c715993f77380b2a';
module.exports = node;
