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
declare export opaque type ServiceEndpointDefinitionTable_serviceEndpointDefinitions$ref: FragmentReference;
declare export opaque type ServiceEndpointDefinitionTable_serviceEndpointDefinitions$fragmentType: ServiceEndpointDefinitionTable_serviceEndpointDefinitions$ref;
export type ServiceEndpointDefinitionTable_serviceEndpointDefinitions = $ReadOnlyArray<{|
  +id: string,
  +index: number,
  +role: ?string,
  +name: string,
  +equipmentType: {|
    +name: string,
    +id: string,
  |},
  +$refType: ServiceEndpointDefinitionTable_serviceEndpointDefinitions$ref,
|}>;
export type ServiceEndpointDefinitionTable_serviceEndpointDefinitions$data = ServiceEndpointDefinitionTable_serviceEndpointDefinitions;
export type ServiceEndpointDefinitionTable_serviceEndpointDefinitions$key = $ReadOnlyArray<{
  +$data?: ServiceEndpointDefinitionTable_serviceEndpointDefinitions$data,
  +$fragmentRefs: ServiceEndpointDefinitionTable_serviceEndpointDefinitions$ref,
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
  "name": "ServiceEndpointDefinitionTable_serviceEndpointDefinitions",
  "type": "ServiceEndpointDefinition",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
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
      "name": "role",
      "args": null,
      "storageKey": null
    },
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": [
        (v1/*: any*/),
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '5e4e6ca6fca81e3fafa893676cb515fa';
module.exports = node;
