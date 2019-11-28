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
import type { FragmentReference } from "relay-runtime";
declare export opaque type PowerSearchLinkFirstEquipmentResultsTable_equipment$ref: FragmentReference;
declare export opaque type PowerSearchLinkFirstEquipmentResultsTable_equipment$fragmentType: PowerSearchLinkFirstEquipmentResultsTable_equipment$ref;
export type PowerSearchLinkFirstEquipmentResultsTable_equipment = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +futureState: ?FutureState,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
  +$refType: PowerSearchLinkFirstEquipmentResultsTable_equipment$ref,
|}>;
export type PowerSearchLinkFirstEquipmentResultsTable_equipment$data = PowerSearchLinkFirstEquipmentResultsTable_equipment;
export type PowerSearchLinkFirstEquipmentResultsTable_equipment$key = $ReadOnlyArray<{
  +$data?: PowerSearchLinkFirstEquipmentResultsTable_equipment$data,
  +$fragmentRefs: PowerSearchLinkFirstEquipmentResultsTable_equipment$ref,
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
  "name": "PowerSearchLinkFirstEquipmentResultsTable_equipment",
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
      "kind": "FragmentSpread",
      "name": "EquipmentBreadcrumbs_equipment",
      "args": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '1716ce00a6510aae7310a382182ce067';
module.exports = node;
