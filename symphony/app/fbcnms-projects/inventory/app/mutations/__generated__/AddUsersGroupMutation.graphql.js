/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash fa5e1f7acc53cb966fd855e0b6f844ce
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type AddUsersGroupInput = {|
  name: string,
  description?: ?string,
  members?: ?$ReadOnlyArray<string>,
|};
export type AddUsersGroupMutationVariables = {|
  input: AddUsersGroupInput
|};
export type AddUsersGroupMutationResponse = {|
  +addUsersGroup: {|
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
export type AddUsersGroupMutation = {|
  variables: AddUsersGroupMutationVariables,
  response: AddUsersGroupMutationResponse,
|};
*/


/*
mutation AddUsersGroupMutation(
  $input: AddUsersGroupInput!
) {
  addUsersGroup(input: $input) {
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
    "type": "AddUsersGroupInput!",
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
    "name": "addUsersGroup",
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
    "name": "AddUsersGroupMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddUsersGroupMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddUsersGroupMutation",
    "id": null,
    "text": "mutation AddUsersGroupMutation(\n  $input: AddUsersGroupInput!\n) {\n  addUsersGroup(input: $input) {\n    id\n    name\n    description\n    status\n    members {\n      id\n      authID\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a7b68e738f78bffc91dc722538502f62';
module.exports = node;
