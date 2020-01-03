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
type AddToEquipmentDialog_parentEquipment$ref = any;
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPositionsGrid_equipment$ref: FragmentReference;
declare export opaque type EquipmentPositionsGrid_equipment$fragmentType: EquipmentPositionsGrid_equipment$ref;
export type EquipmentPositionsGrid_equipment = {|
  +id: string,
  +positions: $ReadOnlyArray<?{|
    +id: string,
    +definition: {|
      +id: string,
      +name: string,
      +index: ?number,
      +visibleLabel: ?string,
    |},
    +attachedEquipment: ?{|
      +id: string,
      +name: string,
      +futureState: ?FutureState,
      +services: $ReadOnlyArray<?{|
        +id: string
      |}>,
    |},
    +parentEquipment: {|
      +id: string
    |},
  |}>,
  +equipmentType: {|
    +positionDefinitions: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +index: ?number,
      +visibleLabel: ?string,
    |}>
  |},
  +$fragmentRefs: AddToEquipmentDialog_parentEquipment$ref,
  +$refType: EquipmentPositionsGrid_equipment$ref,
|};
export type EquipmentPositionsGrid_equipment$data = EquipmentPositionsGrid_equipment;
export type EquipmentPositionsGrid_equipment$key = {
  +$data?: EquipmentPositionsGrid_equipment$data,
  +$fragmentRefs: EquipmentPositionsGrid_equipment$ref,
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
],
v3 = [
  (v0/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "EquipmentPositionsGrid_equipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "positions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPosition",
      "plural": true,
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
          "selections": (v2/*: any*/)
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
              "name": "services",
              "storageKey": null,
              "args": null,
              "concreteType": "Service",
              "plural": true,
              "selections": (v3/*: any*/)
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
          "selections": (v3/*: any*/)
        }
      ]
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
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "positionDefinitions",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPositionDefinition",
          "plural": true,
          "selections": (v2/*: any*/)
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "AddToEquipmentDialog_parentEquipment",
      "args": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c4ed56ad25227e14dccf883ee79b3e2d';
module.exports = node;
