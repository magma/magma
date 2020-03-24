/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 6c31bac93a99ba3a5ae6ac9194f40266
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersView_UsersQueryVariables = {||};
export type UsersView_UsersQueryResponse = {|
  +users: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +authID: string,
        +firstName: string,
        +lastName: string,
        +email: string,
        +status: UserStatus,
        +role: UserRole,
        +profilePhoto: ?{|
          +id: string,
          +fileName: string,
          +storeKey: ?string,
        |},
      |}
    |}>
  |}
|};
export type UsersView_UsersQuery = {|
  variables: UsersView_UsersQueryVariables,
  response: UsersView_UsersQueryResponse,
|};
*/


/*
query UsersView_UsersQuery {
  users {
    edges {
      node {
        id
        authID
        firstName
        lastName
        email
        status
        role
        profilePhoto {
          id
          fileName
          storeKey
        }
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "users",
    "storageKey": null,
    "args": null,
    "concreteType": "UserConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "UserEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "User",
            "plural": false,
            "selections": [
              (v0/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "authID",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "firstName",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "lastName",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "email",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "status",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "role",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "profilePhoto",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "fileName",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "storeKey",
                    "args": null,
                    "storageKey": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "UsersView_UsersQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "UsersView_UsersQuery",
    "argumentDefinitions": [],
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "UsersView_UsersQuery",
    "id": null,
    "text": "query UsersView_UsersQuery {\n  users {\n    edges {\n      node {\n        id\n        authID\n        firstName\n        lastName\n        email\n        status\n        role\n        profilePhoto {\n          id\n          fileName\n          storeKey\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'eac99844a47b8bd8a4904545c8de9856';
module.exports = node;
