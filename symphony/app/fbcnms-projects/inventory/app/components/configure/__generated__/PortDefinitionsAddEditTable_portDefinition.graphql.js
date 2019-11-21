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
declare export opaque type PortDefinitionsAddEditTable_portDefinition$ref: FragmentReference;
declare export opaque type PortDefinitionsAddEditTable_portDefinition$fragmentType: PortDefinitionsAddEditTable_portDefinition$ref;
export type PortDefinitionsAddEditTable_portDefinition = {|
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
  +type: string,
  +portType: ?{|
    +id: string,
    +name: string,
  |},
  +$refType: PortDefinitionsAddEditTable_portDefinition$ref,
|};
export type PortDefinitionsAddEditTable_portDefinition$data = PortDefinitionsAddEditTable_portDefinition;
export type PortDefinitionsAddEditTable_portDefinition$key = {
  +$data?: PortDefinitionsAddEditTable_portDefinition$data,
  +$fragmentRefs: PortDefinitionsAddEditTable_portDefinition$ref,
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
  "name": "PortDefinitionsAddEditTable_portDefinition",
  "type": "EquipmentPortDefinition",
  "metadata": null,
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
      "kind": "ScalarField",
      "alias": null,
      "name": "type",
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
(node/*: any*/).hash = 'cb451e164a86b3be37ecb901b872cea8';
module.exports = node;
