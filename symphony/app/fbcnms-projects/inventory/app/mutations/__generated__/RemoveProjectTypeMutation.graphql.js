/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 42ce89779e4aa26b4f8371bbee88bdc0
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveProjectTypeMutationVariables = {|
  id: string
|};
export type RemoveProjectTypeMutationResponse = {|
  +deleteProjectType: boolean
|};
export type RemoveProjectTypeMutation = {|
  variables: RemoveProjectTypeMutationVariables,
  response: RemoveProjectTypeMutationResponse,
|};
*/


/*
mutation RemoveProjectTypeMutation(
  $id: ID!
) {
  deleteProjectType(id: $id)
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
    "name": "deleteProjectType",
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
    "name": "RemoveProjectTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveProjectTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveProjectTypeMutation",
    "id": null,
    "text": "mutation RemoveProjectTypeMutation(\n  $id: ID!\n) {\n  deleteProjectType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9ef2120c80fef6893e62a26dce1feb23';
module.exports = node;
