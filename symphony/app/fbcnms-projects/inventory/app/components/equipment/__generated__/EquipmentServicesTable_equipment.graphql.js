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
declare export opaque type EquipmentServicesTable_equipment$ref: FragmentReference;
declare export opaque type EquipmentServicesTable_equipment$fragmentType: EquipmentServicesTable_equipment$ref;
export type EquipmentServicesTable_equipment = {|
  +id: string,
  +name: string,
  +services: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +externalId: ?string,
    +customer: ?{|
      +name: string
    |},
    +serviceType: {|
      +id: string,
      +name: string,
    |},
  |}>,
  +$refType: EquipmentServicesTable_equipment$ref,
|};
export type EquipmentServicesTable_equipment$data = EquipmentServicesTable_equipment;
export type EquipmentServicesTable_equipment$key = {
  +$data?: EquipmentServicesTable_equipment$data,
  +$fragmentRefs: EquipmentServicesTable_equipment$ref,
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
  "name": "EquipmentServicesTable_equipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "services",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
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
          "name": "customer",
          "storageKey": null,
          "args": null,
          "concreteType": "Customer",
          "plural": false,
          "selections": [
            (v1/*: any*/)
          ]
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "serviceType",
          "storageKey": null,
          "args": null,
          "concreteType": "ServiceType",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/)
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '26ffdbf9cc9e158c631da821bc4d0393';
module.exports = node;
