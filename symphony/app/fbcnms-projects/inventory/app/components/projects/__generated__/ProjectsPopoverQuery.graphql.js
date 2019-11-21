/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 846817c230867e91f1c7c73a2380e5c8
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type ProjectsPopoverQueryVariables = {|
  projectId: string
|};
export type ProjectsPopoverQueryResponse = {|
  +project: ?{|
    +id: string,
    +name: string,
    +location: ?{|
      +id: string,
      +name: string,
      +latitude: number,
      +longitude: number,
    |},
    +workOrders: $ReadOnlyArray<{|
      +id: string,
      +name: string,
      +description: ?string,
      +ownerName: string,
      +status: WorkOrderStatus,
      +priority: WorkOrderPriority,
      +assignee: ?string,
      +installDate: ?any,
      +location: ?{|
        +id: string,
        +name: string,
        +latitude: number,
        +longitude: number,
      |},
    |}>,
  |}
|};
export type ProjectsPopoverQuery = {|
  variables: ProjectsPopoverQueryVariables,
  response: ProjectsPopoverQueryResponse,
|};
*/


/*
query ProjectsPopoverQuery(
  $projectId: ID!
) {
  project(id: $projectId) {
    id
    name
    location {
      id
      name
      latitude
      longitude
    }
    workOrders {
      id
      name
      description
      ownerName
      status
      priority
      assignee
      installDate
      location {
        id
        name
        latitude
        longitude
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "projectId",
    "type": "ID!",
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
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "location",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": false,
  "selections": [
    (v1/*: any*/),
    (v2/*: any*/),
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
v4 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "project",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "projectId"
      }
    ],
    "concreteType": "Project",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      (v2/*: any*/),
      (v3/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrders",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrder",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
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
            "name": "ownerName",
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
            "name": "priority",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "assignee",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "installDate",
            "args": null,
            "storageKey": null
          },
          (v3/*: any*/)
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ProjectsPopoverQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectsPopoverQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectsPopoverQuery",
    "id": null,
    "text": "query ProjectsPopoverQuery(\n  $projectId: ID!\n) {\n  project(id: $projectId) {\n    id\n    name\n    location {\n      id\n      name\n      latitude\n      longitude\n    }\n    workOrders {\n      id\n      name\n      description\n      ownerName\n      status\n      priority\n      assignee\n      installDate\n      location {\n        id\n        name\n        latitude\n        longitude\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '499c4bc6e254946c1139480f67cb4eaf';
module.exports = node;
