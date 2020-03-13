/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 024a2038b0137c2caf3c55a776604fc4
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UsersSettings_UserQueryVariables = {|
  authID: string
|};
export type UsersSettings_UserQueryResponse = {|
  +user: ?{|
    +id: string
  |}
|};
export type UsersSettings_UserQuery = {|
  variables: UsersSettings_UserQueryVariables,
  response: UsersSettings_UserQueryResponse,
|};
*/


/*
query UsersSettings_UserQuery(
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
    "name": "UsersSettings_UserQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "UsersSettings_UserQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "UsersSettings_UserQuery",
    "id": null,
    "text": "query UsersSettings_UserQuery(\n  $authID: String!\n) {\n  user(authID: $authID) {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8bb8a3865034e4ca2fdbf08976b0f1ca';
module.exports = node;
