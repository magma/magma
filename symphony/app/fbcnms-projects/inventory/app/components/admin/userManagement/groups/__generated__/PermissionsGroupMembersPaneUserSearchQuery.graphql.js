/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 553f067f5c58a1e59227f2e87ac9a452
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type UserFilterType = "USER_NAME" | "%future added value";
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UserFilterInput = {|
  filterType: UserFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  externalId?: ?string,
  name: string,
  type: PropertyKind,
  index?: ?number,
  category?: ?string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};
export type PermissionsGroupMembersPaneUserSearchQueryVariables = {|
  filters: $ReadOnlyArray<UserFilterInput>
|};
export type PermissionsGroupMembersPaneUserSearchQueryResponse = {|
  +userSearch: {|
    +users: $ReadOnlyArray<?{|
      +id: string,
      +authID: string,
      +firstName: string,
      +lastName: string,
      +email: string,
      +status: UserStatus,
      +role: UserRole,
      +groups: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
      |}>,
      +profilePhoto: ?{|
        +id: string,
        +fileName: string,
        +storeKey: ?string,
      |},
    |}>
  |}
|};
export type PermissionsGroupMembersPaneUserSearchQuery = {|
  variables: PermissionsGroupMembersPaneUserSearchQueryVariables,
  response: PermissionsGroupMembersPaneUserSearchQueryResponse,
|};
*/


/*
query PermissionsGroupMembersPaneUserSearchQuery(
  $filters: [UserFilterInput!]!
) {
  userSearch(filters: $filters) {
    users {
      id
      authID
      firstName
      lastName
      email
      status
      role
      groups {
        id
        name
      }
      profilePhoto {
        id
        fileName
        storeKey
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
    "type": "[UserFilterInput!]!",
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
    "name": "userSearch",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filters",
        "variableName": "filters"
      }
    ],
    "concreteType": "UserSearchResult",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "users",
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
            "name": "groups",
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
              }
            ]
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
              (v1/*: any*/),
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
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "PermissionsGroupMembersPaneUserSearchQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PermissionsGroupMembersPaneUserSearchQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PermissionsGroupMembersPaneUserSearchQuery",
    "id": null,
    "text": "query PermissionsGroupMembersPaneUserSearchQuery(\n  $filters: [UserFilterInput!]!\n) {\n  userSearch(filters: $filters) {\n    users {\n      id\n      authID\n      firstName\n      lastName\n      email\n      status\n      role\n      groups {\n        id\n        name\n      }\n      profilePhoto {\n        id\n        fileName\n        storeKey\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '01be0fd4e5f83a289935b03fdad95aa0';
module.exports = node;
