/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 78123c64d00d8ae4c635b25c77f55ed8
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ProjectsMap_projects$ref = any;
type ProjectsTableView_projects$ref = any;
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type ProjectFilterType = "PROJECT_NAME" | "%future added value";
export type ProjectFilterInput = {|
  filterType: ProjectFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
|};
export type ProjectComparisonViewQueryRendererSearchQueryVariables = {|
  limit?: ?number,
  filters: $ReadOnlyArray<ProjectFilterInput>,
|};
export type ProjectComparisonViewQueryRendererSearchQueryResponse = {|
  +projectSearch: $ReadOnlyArray<?{|
    +$fragmentRefs: ProjectsTableView_projects$ref & ProjectsMap_projects$ref
  |}>
|};
export type ProjectComparisonViewQueryRendererSearchQuery = {|
  variables: ProjectComparisonViewQueryRendererSearchQueryVariables,
  response: ProjectComparisonViewQueryRendererSearchQueryResponse,
|};
*/


/*
query ProjectComparisonViewQueryRendererSearchQuery(
  $limit: Int
  $filters: [ProjectFilterInput!]!
) {
  projectSearch(limit: $limit, filters: $filters) {
    ...ProjectsTableView_projects
    ...ProjectsMap_projects
    id
  }
}

fragment ProjectsMap_projects on Project {
  id
  name
  location {
    id
    name
    latitude
    longitude
  }
  numberOfWorkOrders
}

fragment ProjectsTableView_projects on Project {
  id
  name
  creator
  location {
    id
    name
  }
  type {
    id
    name
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
    "name": "ProjectComparisonViewQueryRendererSearchQuery",
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
          {
            "kind": "FragmentSpread",
            "name": "ProjectsTableView_projects",
            "args": null
          },
          {
            "kind": "FragmentSpread",
            "name": "ProjectsMap_projects",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectComparisonViewQueryRendererSearchQuery",
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
            "kind": "ScalarField",
            "alias": null,
            "name": "creator",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "location",
            "storageKey": null,
            "args": null,
            "concreteType": "Location",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "latitude",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "longitude",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "type",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectType",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/)
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numberOfWorkOrders",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectComparisonViewQueryRendererSearchQuery",
    "id": null,
    "text": "query ProjectComparisonViewQueryRendererSearchQuery(\n  $limit: Int\n  $filters: [ProjectFilterInput!]!\n) {\n  projectSearch(limit: $limit, filters: $filters) {\n    ...ProjectsTableView_projects\n    ...ProjectsMap_projects\n    id\n  }\n}\n\nfragment ProjectsMap_projects on Project {\n  id\n  name\n  location {\n    id\n    name\n    latitude\n    longitude\n  }\n  numberOfWorkOrders\n}\n\nfragment ProjectsTableView_projects on Project {\n  id\n  name\n  creator\n  location {\n    id\n    name\n  }\n  type {\n    id\n    name\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '649095d4b9b56d7a5c136ad0a1feb5d3';
module.exports = node;
