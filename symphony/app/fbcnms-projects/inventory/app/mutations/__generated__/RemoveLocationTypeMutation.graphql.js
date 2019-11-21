/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 5b8349f194975a27b593a68f28d9a386
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveLocationTypeMutationVariables = {|
  id: string
|};
export type RemoveLocationTypeMutationResponse = {|
  +removeLocationType: string
|};
export type RemoveLocationTypeMutation = {|
  variables: RemoveLocationTypeMutationVariables,
  response: RemoveLocationTypeMutationResponse,
|};
*/


/*
mutation RemoveLocationTypeMutation(
  $id: ID!
) {
  removeLocationType(id: $id)
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
    "name": "removeLocationType",
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
    "name": "RemoveLocationTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveLocationTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveLocationTypeMutation",
    "id": null,
    "text": "mutation RemoveLocationTypeMutation(\n  $id: ID!\n) {\n  removeLocationType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9df4715e9aebae6d3f5b3e7975c33de5';
module.exports = node;
