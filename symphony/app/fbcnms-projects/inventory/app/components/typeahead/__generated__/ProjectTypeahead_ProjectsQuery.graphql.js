/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 7191b6634e06abbdc90ae53a7963f340
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type ProjectFilterType = "PROJECT_NAME" | "%future added value";
export type ProjectFilterInput = {|
  filterType: ProjectFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
|};
export type ProjectTypeahead_ProjectsQueryVariables = {|
  limit?: ?number,
  filters: $ReadOnlyArray<ProjectFilterInput>,
|};
export type ProjectTypeahead_ProjectsQueryResponse = {|
  +projectSearch: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +type: {|
      +name: string
    |},
  |}>
|};
export type ProjectTypeahead_ProjectsQuery = {|
  variables: ProjectTypeahead_ProjectsQueryVariables,
  response: ProjectTypeahead_ProjectsQueryResponse,
|};
*/


/*
query ProjectTypeahead_ProjectsQuery(
  $limit: Int
  $filters: [ProjectFilterInput!]!
) {
  projectSearch(limit: $limit, filters: $filters) {
    id
    name
    type {
      name
      id
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "limit",
    "type": "Int",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "filters",
    "type": "[ProjectFilterInput!]!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "filters",
    "variableName": "filters"
  },
  {
    "kind": "Variable",
    "name": "limit",
    "variableName": "limit"
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
  "name": "name",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ProjectTypeahead_ProjectsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "projectSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Project",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "type",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectType",
            "plural": false,
            "selections": [
              (v3/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectTypeahead_ProjectsQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "projectSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Project",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "type",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectType",
            "plural": false,
            "selections": [
              (v3/*: any*/),
              (v2/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectTypeahead_ProjectsQuery",
    "id": null,
    "text": "query ProjectTypeahead_ProjectsQuery(\n  $limit: Int\n  $filters: [ProjectFilterInput!]!\n) {\n  projectSearch(limit: $limit, filters: $filters) {\n    id\n    name\n    type {\n      name\n      id\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '6041501fde832bf8b99f06d0c5d9103a';
module.exports = node;
