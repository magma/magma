/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 02651cef4b097e16fc9d3b37b3071342
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_OR_EQUAL_THAN" | "DATE_GREATER_THAN" | "DATE_LESS_OR_EQUAL_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
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
  +projects: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +type: {|
          +name: string
        |},
      |}
    |}>
  |}
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
  projects(first: $limit, filterBy: $filters) {
    edges {
      node {
        id
        name
        type {
          name
          id
        }
      }
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
    "name": "filterBy",
    "variableName": "filters"
  },
  {
    "kind": "Variable",
    "name": "first",
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
        "name": "projects",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ProjectConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "Project",
                "plural": false,
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
        "name": "projects",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ProjectConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "Project",
                "plural": false,
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
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectTypeahead_ProjectsQuery",
    "id": null,
    "text": "query ProjectTypeahead_ProjectsQuery(\n  $limit: Int\n  $filters: [ProjectFilterInput!]!\n) {\n  projects(first: $limit, filterBy: $filters) {\n    edges {\n      node {\n        id\n        name\n        type {\n          name\n          id\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4f9669091bda026ea27aab7f07fd1b0b';
module.exports = node;
