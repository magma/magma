/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e7e0218b9ab1d962f71d0e79ffc34c4b
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type WorkOrdersMap_workOrders$ref = any;
type WorkOrdersView_workOrder$ref = any;
export type FilterOperator = "CONTAINS" | "DATE_GREATER_OR_EQUAL_THAN" | "DATE_GREATER_THAN" | "DATE_LESS_OR_EQUAL_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type WorkOrderFilterType = "LOCATION_INST" | "LOCATION_INST_EXTERNAL_ID" | "WORK_ORDER_ASSIGNED_TO" | "WORK_ORDER_CLOSE_DATE" | "WORK_ORDER_CREATION_DATE" | "WORK_ORDER_LOCATION_INST" | "WORK_ORDER_NAME" | "WORK_ORDER_OWNED_BY" | "WORK_ORDER_PRIORITY" | "WORK_ORDER_STATUS" | "WORK_ORDER_TYPE" | "%future added value";
export type WorkOrderFilterInput = {|
  filterType: WorkOrderFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  propertyValue?: ?PropertyTypeInput,
  timeValue?: ?any,
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
export type WorkOrderComparisonViewQueryRendererSearchQueryVariables = {|
  limit?: ?number,
  filters: $ReadOnlyArray<WorkOrderFilterInput>,
|};
export type WorkOrderComparisonViewQueryRendererSearchQueryResponse = {|
  +workOrders: {|
    +totalCount: number,
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +$fragmentRefs: WorkOrdersView_workOrder$ref & WorkOrdersMap_workOrders$ref
      |}
    |}>,
  |}
|};
export type WorkOrderComparisonViewQueryRendererSearchQuery = {|
  variables: WorkOrderComparisonViewQueryRendererSearchQueryVariables,
  response: WorkOrderComparisonViewQueryRendererSearchQueryResponse,
|};
*/


/*
query WorkOrderComparisonViewQueryRendererSearchQuery(
  $limit: Int
  $filters: [WorkOrderFilterInput!]!
) {
  workOrders(first: $limit, filterBy: $filters) {
    totalCount
    edges {
      node {
        ...WorkOrdersView_workOrder
        ...WorkOrdersMap_workOrders
        id
      }
    }
  }
}

fragment WorkOrdersMap_workOrders on WorkOrder {
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

fragment WorkOrdersView_workOrder on WorkOrder {
  id
  name
  description
  owner {
    id
    email
  }
  creationDate
  installDate
  status
  assignedTo {
    id
    email
  }
  location {
    id
    name
  }
  workOrderType {
    id
    name
  }
  project {
    id
    name
  }
  closeDate
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
    "type": "[WorkOrderFilterInput!]!",
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
  "name": "totalCount",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v5 = [
  (v3/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "email",
    "args": null,
    "storageKey": null
  }
],
v6 = [
  (v3/*: any*/),
  (v4/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderComparisonViewQueryRendererSearchQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrders",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrder",
                "plural": false,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "WorkOrdersView_workOrder",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "WorkOrdersMap_workOrders",
                    "args": null
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
    "name": "WorkOrderComparisonViewQueryRendererSearchQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrders",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderConnection",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrder",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v4/*: any*/),
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
                    "name": "creationDate",
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
                    "name": "assignedTo",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "User",
                    "plural": false,
                    "selections": (v5/*: any*/)
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
                      (v3/*: any*/),
                      (v4/*: any*/),
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
                    "name": "workOrderType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkOrderType",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "project",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Project",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "closeDate",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "priority",
                    "args": null,
                    "storageKey": null
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
    "name": "WorkOrderComparisonViewQueryRendererSearchQuery",
    "id": null,
    "text": "query WorkOrderComparisonViewQueryRendererSearchQuery(\n  $limit: Int\n  $filters: [WorkOrderFilterInput!]!\n) {\n  workOrders(first: $limit, filterBy: $filters) {\n    totalCount\n    edges {\n      node {\n        ...WorkOrdersView_workOrder\n        ...WorkOrdersMap_workOrders\n        id\n      }\n    }\n  }\n}\n\nfragment WorkOrdersMap_workOrders on WorkOrder {\n  id\n  name\n  description\n  owner {\n    id\n    email\n  }\n  status\n  priority\n  assignedTo {\n    id\n    email\n  }\n  installDate\n  location {\n    id\n    name\n    latitude\n    longitude\n  }\n}\n\nfragment WorkOrdersView_workOrder on WorkOrder {\n  id\n  name\n  description\n  owner {\n    id\n    email\n  }\n  creationDate\n  installDate\n  status\n  assignedTo {\n    id\n    email\n  }\n  location {\n    id\n    name\n  }\n  workOrderType {\n    id\n    name\n  }\n  project {\n    id\n    name\n  }\n  closeDate\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a2134c21f3a4ecdab7e0de3d8f8066e9';
module.exports = node;
