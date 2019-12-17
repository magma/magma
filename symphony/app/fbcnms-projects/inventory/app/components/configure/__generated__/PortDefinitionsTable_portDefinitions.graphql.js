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
declare export opaque type PortDefinitionsTable_portDefinitions$ref: FragmentReference;
declare export opaque type PortDefinitionsTable_portDefinitions$fragmentType: PortDefinitionsTable_portDefinitions$ref;
export type PortDefinitionsTable_portDefinitions = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
  +portType: ?{|
    +id: string,
    +name: string,
  |},
  +$refType: PortDefinitionsTable_portDefinitions$ref,
|}>;
export type PortDefinitionsTable_portDefinitions$data = PortDefinitionsTable_portDefinitions;
export type PortDefinitionsTable_portDefinitions$key = $ReadOnlyArray<{
  +$data?: PortDefinitionsTable_portDefinitions$data,
  +$fragmentRefs: PortDefinitionsTable_portDefinitions$ref,
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
  "name": "PortDefinitionsTable_portDefinitions",
  "type": "EquipmentPortDefinition",
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
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortType",
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
(node/*: any*/).hash = '0ee0d296004879c18c2ac5d4bd8f8181';
module.exports = node;
