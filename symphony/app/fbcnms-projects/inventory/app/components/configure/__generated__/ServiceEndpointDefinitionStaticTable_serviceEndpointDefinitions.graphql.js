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
declare export opaque type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref: FragmentReference;
declare export opaque type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$fragmentType: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref;
export type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +role: ?string,
  +index: number,
  +equipmentType: {|
    +id: string,
    +name: string,
  |},
  +$refType: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref,
|}>;
export type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$data = ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions;
export type ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$key = $ReadOnlyArray<{
  +$data?: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$data,
  +$fragmentRefs: ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions$ref,
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
  "name": "ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions",
  "type": "ServiceEndpointDefinition",
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
      "name": "role",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "index",
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
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '355047f79b4b79d16aa746b866d228ae';
module.exports = node;
