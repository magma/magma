/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4d236942a1465e14e097b8ad000c6036
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveLocationMutationVariables = {|
  id: string
|};
export type RemoveLocationMutationResponse = {|
  +removeLocation: string
|};
export type RemoveLocationMutation = {|
  variables: RemoveLocationMutationVariables,
  response: RemoveLocationMutationResponse,
|};
*/


/*
mutation RemoveLocationMutation(
  $id: ID!
) {
  removeLocation(id: $id)
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
    "name": "removeLocation",
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
    "name": "RemoveLocationMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveLocationMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveLocationMutation",
    "id": null,
    "text": "mutation RemoveLocationMutation(\n  $id: ID!\n) {\n  removeLocation(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '020ce82a52db80ee9ae4e4b370e756a7';
module.exports = node;
