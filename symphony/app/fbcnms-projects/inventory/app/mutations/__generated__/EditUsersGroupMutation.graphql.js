/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash ddc2eb464d3a64b895d8df805799463e
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type EditUsersGroupInput = {|
  id: string,
  name?: ?string,
  description?: ?string,
  status?: ?UsersGroupStatus,
|};
export type EditUsersGroupMutationVariables = {|
  input: EditUsersGroupInput
|};
export type EditUsersGroupMutationResponse = {|
  +editUsersGroup: {|
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
export type EditUsersGroupMutation = {|
  variables: EditUsersGroupMutationVariables,
  response: EditUsersGroupMutationResponse,
|};
*/


/*
mutation EditUsersGroupMutation(
  $input: EditUsersGroupInput!
) {
  editUsersGroup(input: $input) {
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
    "type": "EditUsersGroupInput!",
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
    "name": "editUsersGroup",
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
    "name": "EditUsersGroupMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EditUsersGroupMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditUsersGroupMutation",
    "id": null,
    "text": "mutation EditUsersGroupMutation(\n  $input: EditUsersGroupInput!\n) {\n  editUsersGroup(input: $input) {\n    id\n    name\n    description\n    status\n    members {\n      id\n      authID\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '34078eb6c09ee4ce62f28cd37509f647';
module.exports = node;
