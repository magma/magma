/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 6f056f173b0f862f6d7afb0d6f584c90
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UsersByAuthIDQueryVariables = {|
  authID: string
|};
export type UsersByAuthIDQueryResponse = {|
  +user: ?{|
    +id: string
  |}
|};
export type UsersByAuthIDQuery = {|
  variables: UsersByAuthIDQueryVariables,
  response: UsersByAuthIDQueryResponse,
|};
*/


/*
query UsersByAuthIDQuery(
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
    "name": "UsersByAuthIDQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "UsersByAuthIDQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "UsersByAuthIDQuery",
    "id": null,
    "text": "query UsersByAuthIDQuery(\n  $authID: String!\n) {\n  user(authID: $authID) {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a242fdd9ce8fb33496aab66755ad0186';
module.exports = node;
