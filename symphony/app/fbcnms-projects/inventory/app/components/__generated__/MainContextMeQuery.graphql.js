/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash b8dc6dfaa9688e2c238f79a5db0e086a
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type MainContextMeQueryVariables = {||};
export type MainContextMeQueryResponse = {|
  +me: ?{|
    +user: ?{|
      +id: string,
      +authID: string,
      +email: string,
      +firstName: string,
      +lastName: string,
    |},
    +permissions: {|
      +canWrite: boolean,
      +adminPolicy: {|
        +access: {|
          +isAllowed: PermissionValue
        |}
      |},
    |},
  |}
|};
export type MainContextMeQuery = {|
  variables: MainContextMeQueryVariables,
  response: MainContextMeQueryResponse,
|};
*/


/*
query MainContextMeQuery {
  me {
    user {
      id
      authID
      email
      firstName
      lastName
    }
    permissions {
      canWrite
      adminPolicy {
        access {
          isAllowed
        }
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "me",
    "storageKey": null,
    "args": null,
    "concreteType": "Viewer",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "user",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "id",
            "args": null,
            "storageKey": null
          },
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
            "name": "email",
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
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "permissions",
        "storageKey": null,
        "args": null,
        "concreteType": "PermissionSettings",
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "canWrite",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "adminPolicy",
            "storageKey": null,
            "args": null,
            "concreteType": "AdministrativePolicy",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "access",
                "storageKey": null,
                "args": null,
                "concreteType": "BasicPermissionRule",
                "plural": false,
                "selections": [
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "isAllowed",
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
    "name": "MainContextMeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "MainContextMeQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "MainContextMeQuery",
    "id": null,
    "text": "query MainContextMeQuery {\n  me {\n    user {\n      id\n      authID\n      email\n      firstName\n      lastName\n    }\n    permissions {\n      canWrite\n      adminPolicy {\n        access {\n          isAllowed\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f1a977ba665a002410f12c3df9e89af4';
module.exports = node;
