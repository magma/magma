/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash f54620de606fa6d71e919de1773e13db
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveServiceTypeMutationVariables = {|
  id: string
|};
export type RemoveServiceTypeMutationResponse = {|
  +removeServiceType: string
|};
export type RemoveServiceTypeMutation = {|
  variables: RemoveServiceTypeMutationVariables,
  response: RemoveServiceTypeMutationResponse,
|};
*/


/*
mutation RemoveServiceTypeMutation(
  $id: ID!
) {
  removeServiceType(id: $id)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "removeServiceType",
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      }
    ],
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "RemoveServiceTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveServiceTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveServiceTypeMutation",
    "id": null,
    "text": "mutation RemoveServiceTypeMutation(\n  $id: ID!\n) {\n  removeServiceType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '499330f58e6161c5f94f25863abdd098';
module.exports = node;
