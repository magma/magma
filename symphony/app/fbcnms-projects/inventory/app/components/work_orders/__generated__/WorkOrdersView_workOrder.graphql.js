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
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type WorkOrdersView_workOrder$ref: FragmentReference;
declare export opaque type WorkOrdersView_workOrder$fragmentType: WorkOrdersView_workOrder$ref;
export type WorkOrdersView_workOrder = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +description: ?string,
  +ownerName: string,
  +creationDate: any,
  +installDate: ?any,
  +status: WorkOrderStatus,
  +assignee: ?string,
  +location: ?{|
    +id: string,
    +name: string,
  |},
  +workOrderType: {|
    +id: string,
    +name: string,
  |},
  +project: ?{|
    +id: string,
    +name: string,
  |},
  +$refType: WorkOrdersView_workOrder$ref,
|}>;
export type WorkOrdersView_workOrder$data = WorkOrdersView_workOrder;
export type WorkOrdersView_workOrder$key = $ReadOnlyArray<{
  +$data?: WorkOrdersView_workOrder$data,
  +$fragmentRefs: WorkOrdersView_workOrder$ref,
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
},
v2 = [
  (v0/*: any*/),
  (v1/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "WorkOrdersView_workOrder",
  "type": "WorkOrder",
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
      "name": "assignee",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "location",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "workOrderType",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkOrderType",
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "project",
      "storageKey": null,
      "args": null,
      "concreteType": "Project",
      "plural": false,
      "selections": (v2/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b420cc68c865cc34e5985dbd20c78f64';
module.exports = node;
