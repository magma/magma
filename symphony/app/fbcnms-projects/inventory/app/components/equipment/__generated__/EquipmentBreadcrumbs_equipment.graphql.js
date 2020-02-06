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
declare export opaque type EquipmentBreadcrumbs_equipment$ref: FragmentReference;
declare export opaque type EquipmentBreadcrumbs_equipment$fragmentType: EquipmentBreadcrumbs_equipment$ref;
export type EquipmentBreadcrumbs_equipment = {|
  +id: string,
  +name: string,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +locationHierarchy: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +locationType: {|
      +name: string
    |},
  |}>,
  +positionHierarchy: $ReadOnlyArray<{|
    +id: string,
    +definition: {|
      +id: string,
      +name: string,
      +visibleLabel: ?string,
    |},
    +parentEquipment: {|
      +id: string,
      +name: string,
      +equipmentType: {|
        +id: string,
        +name: string,
      |},
    |},
  |}>,
  +$refType: EquipmentBreadcrumbs_equipment$ref,
|};
export type EquipmentBreadcrumbs_equipment$data = EquipmentBreadcrumbs_equipment;
export type EquipmentBreadcrumbs_equipment$key = {
  +$data?: EquipmentBreadcrumbs_equipment$data,
  +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
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
},
v2 = {
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
};
return {
  "kind": "Fragment",
  "name": "EquipmentBreadcrumbs_equipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationHierarchy",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": true,
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
          "selections": [
            (v1/*: any*/)
          ]
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "positionHierarchy",
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
          "selections": [
            (v0/*: any*/),
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
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
            (v2/*: any*/)
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '16aecc58a50e73e65b4a6094e719a457';
module.exports = node;
