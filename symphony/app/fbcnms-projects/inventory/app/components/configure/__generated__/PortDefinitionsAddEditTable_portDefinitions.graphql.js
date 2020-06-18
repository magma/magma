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
declare export opaque type PortDefinitionsAddEditTable_portDefinitions$ref: FragmentReference;
declare export opaque type PortDefinitionsAddEditTable_portDefinitions$fragmentType: PortDefinitionsAddEditTable_portDefinitions$ref;
export type PortDefinitionsAddEditTable_portDefinitions = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
  +portType: ?{|
    +id: string,
    +name: string,
  |},
  +$refType: PortDefinitionsAddEditTable_portDefinitions$ref,
|}>;
export type PortDefinitionsAddEditTable_portDefinitions$data = PortDefinitionsAddEditTable_portDefinitions;
export type PortDefinitionsAddEditTable_portDefinitions$key = $ReadOnlyArray<{
  +$data?: PortDefinitionsAddEditTable_portDefinitions$data,
  +$fragmentRefs: PortDefinitionsAddEditTable_portDefinitions$ref,
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
  "name": "PortDefinitionsAddEditTable_portDefinitions",
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
(node/*: any*/).hash = '02bde9d17bd7bb5430914f693cdd659b';
module.exports = node;
