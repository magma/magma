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
declare export opaque type EquipmentPropertiesCard_position$ref: FragmentReference;
declare export opaque type EquipmentPropertiesCard_position$fragmentType: EquipmentPropertiesCard_position$ref;
export type EquipmentPropertiesCard_position = {
  +id: string,
  +definition: {
    +id: string,
    +name: string,
    +index: ?number,
    +visibleLabel: ?string,
    ...
  },
  +attachedEquipment: ?{
    +id: string,
    +name: string,
    +futureState: ?FutureState,
    +workOrder: ?{
      +id: string,
      +status: WorkOrderStatus,
      ...
    },
    ...
  },
  ...
};
export type EquipmentPropertiesCard_position$data = EquipmentPropertiesCard_position;
export type EquipmentPropertiesCard_position$key = {
  +$data?: EquipmentPropertiesCard_position$data,
  +$fragmentRefs: EquipmentPropertiesCard_position$ref,
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
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "EquipmentPropertiesCard_position",
  "type": "EquipmentPosition",
  "metadata": {
    "mask": false
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "definition",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPositionDefinition",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "index",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "visibleLabel",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "attachedEquipment",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
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
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd1e6f2ef9ef3182e98188066491a9856';
module.exports = node;
