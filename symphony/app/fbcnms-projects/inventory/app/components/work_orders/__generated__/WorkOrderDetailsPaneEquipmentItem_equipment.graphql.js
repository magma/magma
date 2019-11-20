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
declare export opaque type WorkOrderDetailsPaneEquipmentItem_equipment$ref: FragmentReference;
declare export opaque type WorkOrderDetailsPaneEquipmentItem_equipment$fragmentType: WorkOrderDetailsPaneEquipmentItem_equipment$ref;
export type WorkOrderDetailsPaneEquipmentItem_equipment = {|
  +id: string,
  +name: string,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +parentLocation: ?{|
    +id: string,
    +name: string,
    +locationType: {|
      +id: string,
      +name: string,
    |},
  |},
  +parentPosition: ?{|
    +id: string,
    +definition: {|
      +name: string,
      +visibleLabel: ?string,
    |},
    +parentEquipment: {|
      +id: string,
      +name: string,
    |},
  |},
  +$refType: WorkOrderDetailsPaneEquipmentItem_equipment$ref,
|};
export type WorkOrderDetailsPaneEquipmentItem_equipment$data = WorkOrderDetailsPaneEquipmentItem_equipment;
export type WorkOrderDetailsPaneEquipmentItem_equipment$key = {
  +$data?: WorkOrderDetailsPaneEquipmentItem_equipment$data,
  +$fragmentRefs: WorkOrderDetailsPaneEquipmentItem_equipment$ref,
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
},
v2 = [
  (v0/*: any*/),
  (v1/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "WorkOrderDetailsPaneEquipmentItem_equipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "parentLocation",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "locationType",
          "storageKey": null,
          "args": null,
          "concreteType": "LocationType",
          "plural": false,
          "selections": (v2/*: any*/)
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "parentPosition",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPosition",
      "plural": false,
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
            (v1/*: any*/),
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
          "name": "parentEquipment",
          "storageKey": null,
          "args": null,
          "concreteType": "Equipment",
          "plural": false,
          "selections": (v2/*: any*/)
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '50466f4314e67740e00fa80dc786d8e8';
module.exports = node;
