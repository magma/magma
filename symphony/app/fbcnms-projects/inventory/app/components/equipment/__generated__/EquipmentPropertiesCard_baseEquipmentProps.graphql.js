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
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPropertiesCard_baseEquipmentProps$ref: FragmentReference;
declare export opaque type EquipmentPropertiesCard_baseEquipmentProps$fragmentType: EquipmentPropertiesCard_baseEquipmentProps$ref;
export type EquipmentPropertiesCard_baseEquipmentProps = {
  +id: string,
  +name: string,
  +futureState: ?FutureState,
  +parentLocation: ?{
    +id: string,
    +name: string,
    ...
  },
  ...
};
export type EquipmentPropertiesCard_baseEquipmentProps$data = EquipmentPropertiesCard_baseEquipmentProps;
export type EquipmentPropertiesCard_baseEquipmentProps$key = {
  +$data?: EquipmentPropertiesCard_baseEquipmentProps$data,
  +$fragmentRefs: EquipmentPropertiesCard_baseEquipmentProps$ref,
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
  "name": "EquipmentPropertiesCard_baseEquipmentProps",
  "type": "Equipment",
  "metadata": {
    "mask": false
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
      "name": "parentLocation",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '38c108dd662662f36827a889ffc2706f';
module.exports = node;
