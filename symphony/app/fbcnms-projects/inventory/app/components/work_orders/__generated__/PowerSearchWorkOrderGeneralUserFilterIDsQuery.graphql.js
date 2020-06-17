/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1bf9d6968c0d42910914b6c000562591
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PowerSearchWorkOrderGeneralUserFilterIDsQueryVariables = {|
  id: string
|};
export type PowerSearchWorkOrderGeneralUserFilterIDsQueryResponse = {|
  +node: ?{|
    +id?: string,
    +email?: string,
  |}
|};
export type PowerSearchWorkOrderGeneralUserFilterIDsQuery = {|
  variables: PowerSearchWorkOrderGeneralUserFilterIDsQueryVariables,
  response: PowerSearchWorkOrderGeneralUserFilterIDsQueryResponse,
|};
*/


/*
query PowerSearchWorkOrderGeneralUserFilterIDsQuery(
  $id: ID!
) {
  node(id: $id) {
    __typename
    ... on User {
      id
      email
    }
    id
  }
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
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "email",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "PowerSearchWorkOrderGeneralUserFilterIDsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "User",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchWorkOrderGeneralUserFilterIDsQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "User",
            "selections": [
              (v3/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchWorkOrderGeneralUserFilterIDsQuery",
    "id": null,
    "text": "query PowerSearchWorkOrderGeneralUserFilterIDsQuery(\n  $id: ID!\n) {\n  node(id: $id) {\n    __typename\n    ... on User {\n      id\n      email\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '46680cb3247e4632e08a63153e09df71';
module.exports = node;
