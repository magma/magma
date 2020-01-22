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
type WorkOrderDetailsPaneEquipmentItem_equipment$ref = any;
type WorkOrderDetailsPaneLinkItem_link$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type WorkOrderDetailsPane_workOrder$ref: FragmentReference;
declare export opaque type WorkOrderDetailsPane_workOrder$fragmentType: WorkOrderDetailsPane_workOrder$ref;
export type WorkOrderDetailsPane_workOrder = {|
  +id: string,
  +name: string,
  +equipmentToAdd: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: WorkOrderDetailsPaneEquipmentItem_equipment$ref,
  |}>,
  +equipmentToRemove: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: WorkOrderDetailsPaneEquipmentItem_equipment$ref,
  |}>,
  +linksToAdd: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: WorkOrderDetailsPaneLinkItem_link$ref,
  |}>,
  +linksToRemove: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: WorkOrderDetailsPaneLinkItem_link$ref,
  |}>,
  +$refType: WorkOrderDetailsPane_workOrder$ref,
|};
export type WorkOrderDetailsPane_workOrder$data = WorkOrderDetailsPane_workOrder;
export type WorkOrderDetailsPane_workOrder$key = {
  +$data?: WorkOrderDetailsPane_workOrder$data,
  +$fragmentRefs: WorkOrderDetailsPane_workOrder$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = [
  (v0/*: any*/),
  {
    "kind": "FragmentSpread",
    "name": "WorkOrderDetailsPaneEquipmentItem_equipment",
    "args": null
  }
],
v2 = [
  (v0/*: any*/),
  {
    "kind": "FragmentSpread",
    "name": "WorkOrderDetailsPaneLinkItem_link",
    "args": null
  }
];
return {
  "kind": "Fragment",
  "name": "WorkOrderDetailsPane_workOrder",
  "type": "WorkOrder",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "name",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentToAdd",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentToRemove",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "linksToAdd",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": true,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "linksToRemove",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": true,
      "selections": (v2/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '4c85915884fc2d9d9f8c4fbee32cc2a4';
module.exports = node;
