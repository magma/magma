/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 5c076afc2fb606058788fbac10d6e146
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type UsersGroupFilterType = "GROUP_NAME" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersGroupFilterInput = {|
  filterType: UsersGroupFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  maxDepth?: ?number,
|};
export type GroupSearchContextQueryVariables = {|
  filters: $ReadOnlyArray<UsersGroupFilterInput>
|};
export type GroupSearchContextQueryResponse = {|
  +usersGroupSearch: {|
    +usersGroups: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +description: ?string,
      +status: UsersGroupStatus,
      +members: $ReadOnlyArray<{|
        +id: string,
        +authID: string,
      |}>,
    |}>
  |}
|};
export type GroupSearchContextQuery = {|
  variables: GroupSearchContextQueryVariables,
  response: GroupSearchContextQueryResponse,
|};
*/


/*
query GroupSearchContextQuery(
  $filters: [UsersGroupFilterInput!]!
) {
  usersGroupSearch(filters: $filters) {
    usersGroups {
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
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "filters",
    "type": "[UsersGroupFilterInput!]!",
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
    "name": "usersGroupSearch",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filters",
        "variableName": "filters"
      }
    ],
    "concreteType": "UsersGroupSearchResult",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "usersGroups",
        "storageKey": null,
        "args": null,
        "concreteType": "UsersGroup",
        "plural": true,
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
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "GroupSearchContextQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "GroupSearchContextQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "GroupSearchContextQuery",
    "id": null,
    "text": "query GroupSearchContextQuery(\n  $filters: [UsersGroupFilterInput!]!\n) {\n  usersGroupSearch(filters: $filters) {\n    usersGroups {\n      id\n      name\n      description\n      status\n      members {\n        id\n        authID\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f5a72e6b8ac8dc0d0f4cd78ea3fc8d6d';
module.exports = node;
