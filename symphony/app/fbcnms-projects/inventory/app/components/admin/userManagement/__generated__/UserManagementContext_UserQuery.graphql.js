/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 5291309aee4f30261773a31fbcff12e7
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UserManagementContext_UserQueryVariables = {|
  authID: string
|};
export type UserManagementContext_UserQueryResponse = {|
  +user: ?{|
    +id: string
  |}
|};
export type UserManagementContext_UserQuery = {|
  variables: UserManagementContext_UserQueryVariables,
  response: UserManagementContext_UserQueryResponse,
|};
*/


/*
query UserManagementContext_UserQuery(
  $authID: String!
) {
  user(authID: $authID) {
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "authID",
    "type": "String!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "user",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "authID",
        "variableName": "authID"
      }
    ],
    "concreteType": "User",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "id",
        "args": null,
        "storageKey": null
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "UserManagementContext_UserQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "UserManagementContext_UserQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "UserManagementContext_UserQuery",
    "id": null,
    "text": "query UserManagementContext_UserQuery(\n  $authID: String!\n) {\n  user(authID: $authID) {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f98277cf4d6a1afbf4e1b674bb2e2012';
module.exports = node;
