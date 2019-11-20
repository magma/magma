/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9235d6c3a917b84b519c4fc1fb43500b
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveServiceMutationVariables = {|
  id: string
|};
export type RemoveServiceMutationResponse = {|
  +removeService: string
|};
export type RemoveServiceMutation = {|
  variables: RemoveServiceMutationVariables,
  response: RemoveServiceMutationResponse,
|};
*/


/*
mutation RemoveServiceMutation(
  $id: ID!
) {
  removeService(id: $id)
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
    "name": "removeService",
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
    "name": "RemoveServiceMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveServiceMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveServiceMutation",
    "id": null,
    "text": "mutation RemoveServiceMutation(\n  $id: ID!\n) {\n  removeService(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '0e33b688caed9845743fa09141151ba7';
module.exports = node;
