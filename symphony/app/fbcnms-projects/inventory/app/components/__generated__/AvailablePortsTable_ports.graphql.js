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
import type { FragmentReference } from "relay-runtime";
declare export opaque type AvailablePortsTable_ports$ref: FragmentReference;
declare export opaque type AvailablePortsTable_ports$fragmentType: AvailablePortsTable_ports$ref;
export type AvailablePortsTable_ports = $ReadOnlyArray<{|
  +id: string,
  +parentEquipment: {|
    +id: string,
    +name: string,
    +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
  |},
  +definition: {|
    +id: string,
    +name: string,
    +portType: ?{|
      +name: string
    |},
    +visibleLabel: ?string,
  |},
  +$refType: AvailablePortsTable_ports$ref,
|}>;
export type AvailablePortsTable_ports$data = AvailablePortsTable_ports;
export type AvailablePortsTable_ports$key = $ReadOnlyArray<{
  +$data?: AvailablePortsTable_ports$data,
  +$fragmentRefs: AvailablePortsTable_ports$ref,
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
  "name": "AvailablePortsTable_ports",
  "type": "EquipmentPort",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
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
        {
          "kind": "FragmentSpread",
          "name": "EquipmentBreadcrumbs_equipment",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "definition",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortDefinition",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "portType",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortType",
          "plural": false,
          "selections": [
            (v1/*: any*/)
          ]
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "visibleLabel",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '28fb7ac76ca11ecf2ff60dae1869a25b';
module.exports = node;
