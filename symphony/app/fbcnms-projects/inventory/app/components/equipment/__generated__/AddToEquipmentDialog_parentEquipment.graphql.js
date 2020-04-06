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
declare export opaque type AddToEquipmentDialog_parentEquipment$ref: FragmentReference;
declare export opaque type AddToEquipmentDialog_parentEquipment$fragmentType: AddToEquipmentDialog_parentEquipment$ref;
export type AddToEquipmentDialog_parentEquipment = {|
  +id: string,
  +locationHierarchy: $ReadOnlyArray<{|
    +id: string
  |}>,
  +$refType: AddToEquipmentDialog_parentEquipment$ref,
|};
export type AddToEquipmentDialog_parentEquipment$data = AddToEquipmentDialog_parentEquipment;
export type AddToEquipmentDialog_parentEquipment$key = {
  +$data?: AddToEquipmentDialog_parentEquipment$data,
  +$fragmentRefs: AddToEquipmentDialog_parentEquipment$ref,
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
};
return {
  "kind": "Fragment",
  "name": "AddToEquipmentDialog_parentEquipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationHierarchy",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": true,
      "selections": [
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '5f1bb222d3448e6445d42e84470b6d5c';
module.exports = node;
