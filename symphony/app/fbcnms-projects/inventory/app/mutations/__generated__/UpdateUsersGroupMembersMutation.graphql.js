/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 99dc18f7131a6223330b4bdbab750fbd
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UpdateUsersGroupMembersInput = {|
  id: string,
  addUserIds: $ReadOnlyArray<string>,
  removeUserIds: $ReadOnlyArray<string>,
|};
export type UpdateUsersGroupMembersMutationVariables = {|
  input: UpdateUsersGroupMembersInput
|};
export type UpdateUsersGroupMembersMutationResponse = {|
  +updateUsersGroupMembers: {|
    +id: string,
    +name: string,
    +description: ?string,
    +status: UsersGroupStatus,
    +members: $ReadOnlyArray<{|
      +id: string,
      +authID: string,
    |}>,
  |}
|};
export type UpdateUsersGroupMembersMutation = {|
  variables: UpdateUsersGroupMembersMutationVariables,
  response: UpdateUsersGroupMembersMutationResponse,
|};
*/


/*
mutation UpdateUsersGroupMembersMutation(
  $input: UpdateUsersGroupMembersInput!
) {
  updateUsersGroupMembers(input: $input) {
    id
    name
    description
    status
    members {
      id
      authID
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "UpdateUsersGroupMembersInput!",
    "defaultValue": null
  }
],
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v2 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "updateUsersGroupMembers",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "UsersGroup",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "name",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "description",
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
        "kind": "LinkedField",
        "alias": null,
        "name": "members",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "authID",
            "args": null,
            "storageKey": null
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
    "name": "UpdateUsersGroupMembersMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "UpdateUsersGroupMembersMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "UpdateUsersGroupMembersMutation",
    "id": null,
    "text": "mutation UpdateUsersGroupMembersMutation(\n  $input: UpdateUsersGroupMembersInput!\n) {\n  updateUsersGroupMembers(input: $input) {\n    id\n    name\n    description\n    status\n    members {\n      id\n      authID\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '32e7db4a603c20ec9d3c60bb7a9f51c0';
module.exports = node;
