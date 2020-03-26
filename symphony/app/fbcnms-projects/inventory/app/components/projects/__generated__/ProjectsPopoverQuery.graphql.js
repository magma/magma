/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash a13dbacf531c1f1fc0cbde4597ea4182
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
    +id?: string,
    +name?: string,
    +location?: ?{|
      +id: string,
      +name: string,
      +latitude: number,
      +longitude: number,
    |},
    +workOrders?: $ReadOnlyArray<{|
      +id: string,
      +name: string,
      +description: ?string,
      +owner: {|
        +id: string,
        +email: string,
      |},
      +status: WorkOrderStatus,
      +priority: WorkOrderPriority,
      +assignedTo: ?{|
        +id: string,
        +email: string,
      |},
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
  project: node(id: $projectId) {
    __typename
    ... on Project {
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
        owner {
          id
          email
        }
        status
        priority
        assignedTo {
          id
          email
        }
        installDate
        location {
          id
          name
          latitude
          longitude
        }
      }
    }
    id
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
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "projectId"
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
},
v4 = {
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
v5 = [
  (v2/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "email",
    "args": null,
    "storageKey": null
  }
],
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "workOrders",
  "storageKey": null,
  "args": null,
  "concreteType": "WorkOrder",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "description",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "owner",
      "storageKey": null,
      "args": null,
      "concreteType": "User",
      "plural": false,
      "selections": (v5/*: any*/)
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
      "kind": "LinkedField",
      "alias": null,
      "name": "assignedTo",
      "storageKey": null,
      "args": null,
      "concreteType": "User",
      "plural": false,
      "selections": (v5/*: any*/)
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "installDate",
      "args": null,
      "storageKey": null
    },
    (v4/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ProjectsPopoverQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "project",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Project",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v6/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectsPopoverQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "project",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Project",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              (v6/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ProjectsPopoverQuery",
    "id": null,
    "text": "query ProjectsPopoverQuery(\n  $projectId: ID!\n) {\n  project: node(id: $projectId) {\n    __typename\n    ... on Project {\n      id\n      name\n      location {\n        id\n        name\n        latitude\n        longitude\n      }\n      workOrders {\n        id\n        name\n        description\n        owner {\n          id\n          email\n        }\n        status\n        priority\n        assignedTo {\n          id\n          email\n        }\n        installDate\n        location {\n          id\n          name\n          latitude\n          longitude\n        }\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '94e8e66d8e2608553c3a5bb7e742468e';
module.exports = node;
