/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash ecf91b236fa569770b84820fc827982e
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DeleteUsersGroupMutationVariables = {|
  id: string
|};
export type DeleteUsersGroupMutationResponse = {|
  +deleteUsersGroup: boolean
|};
export type DeleteUsersGroupMutation = {|
  variables: DeleteUsersGroupMutationVariables,
  response: DeleteUsersGroupMutationResponse,
|};
*/


/*
mutation DeleteUsersGroupMutation(
  $id: ID!
) {
  deleteUsersGroup(id: $id)
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
    "name": "deleteUsersGroup",
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
    "name": "DeleteUsersGroupMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "DeleteUsersGroupMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeleteUsersGroupMutation",
    "id": null,
    "text": "mutation DeleteUsersGroupMutation(\n  $id: ID!\n) {\n  deleteUsersGroup(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c5029d4e245f8583c35fafee471a4157';
module.exports = node;
