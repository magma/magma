/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 89bb11ae6070831b7dcdcd7f9a605eaf
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_OR_EQUAL_THAN" | "DATE_GREATER_THAN" | "DATE_LESS_OR_EQUAL_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type UserFilterType = "USER_NAME" | "USER_STATUS" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UserFilterInput = {|
  filterType: UserFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  statusValue?: ?UserStatus,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  externalId?: ?string,
  name: string,
  type: PropertyKind,
  nodeType?: ?string,
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
export type PowerSearchWorkOrderGeneralUserFilter_userQueryVariables = {|
  filters: $ReadOnlyArray<UserFilterInput>
|};
export type PowerSearchWorkOrderGeneralUserFilter_userQueryResponse = {|
  +userSearch: {|
    +users: $ReadOnlyArray<?{|
      +id: string,
      +email: string,
    |}>
  |}
|};
export type PowerSearchWorkOrderGeneralUserFilter_userQuery = {|
  variables: PowerSearchWorkOrderGeneralUserFilter_userQueryVariables,
  response: PowerSearchWorkOrderGeneralUserFilter_userQueryResponse,
|};
*/


/*
query PowerSearchWorkOrderGeneralUserFilter_userQuery(
  $filters: [UserFilterInput!]!
) {
  userSearch(limit: 10, filters: $filters) {
    users {
      id
      email
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
v1 = [
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
      },
      {
        "kind": "Literal",
        "name": "limit",
        "value": 10
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
            "name": "email",
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
    "name": "PowerSearchWorkOrderGeneralUserFilter_userQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchWorkOrderGeneralUserFilter_userQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchWorkOrderGeneralUserFilter_userQuery",
    "id": null,
    "text": "query PowerSearchWorkOrderGeneralUserFilter_userQuery(\n  $filters: [UserFilterInput!]!\n) {\n  userSearch(limit: 10, filters: $filters) {\n    users {\n      id\n      email\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b13e2c314465152f9bd6e17750487648';
module.exports = node;
