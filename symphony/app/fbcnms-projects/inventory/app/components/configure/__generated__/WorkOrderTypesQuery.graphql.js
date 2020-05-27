/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash f01e19bdbea4d4c0875e494f217aaf7d
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditWorkOrderTypeCard_workOrderType$ref = any;
export type WorkOrderTypesQueryVariables = {||};
export type WorkOrderTypesQueryResponse = {|
  +workOrderTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +description: ?string,
        +$fragmentRefs: AddEditWorkOrderTypeCard_workOrderType$ref,
      |}
    |}>
  |}
|};
export type WorkOrderTypesQuery = {|
  variables: WorkOrderTypesQueryVariables,
  response: WorkOrderTypesQueryResponse,
|};
*/


/*
query WorkOrderTypesQuery {
  workOrderTypes(first: 500) {
    edges {
      node {
        id
        name
        description
        ...AddEditWorkOrderTypeCard_workOrderType
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}

fragment AddEditWorkOrderTypeCard_workOrderType on WorkOrderType {
  id
  name
  description
  numberOfWorkOrders
  propertyTypes {
    id
    name
    type
    nodeType
    index
    stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue
    rangeFromValue
    rangeToValue
    isEditable
    isMandatory
    isInstanceProperty
    isDeleted
    category
  }
  checkListCategoryDefinitions {
    id
    title
    description
    checklistItemDefinitions {
      id
      title
      type
      index
      enumValues
      enumSelectionMode
      helpText
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "pageInfo",
  "storageKey": null,
  "args": null,
  "concreteType": "PageInfo",
  "plural": false,
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "endCursor",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "hasNextPage",
      "args": null,
      "storageKey": null
    }
  ]
},
v6 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 500
  }
],
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrderTypes",
        "name": "__Configure_workOrderTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrderType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "AddEditWorkOrderTypeCard_workOrderType",
                    "args": null
                  }
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "WorkOrderTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderTypes",
        "storageKey": "workOrderTypes(first:500)",
        "args": (v6/*: any*/),
        "concreteType": "WorkOrderTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrderType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "numberOfWorkOrders",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      (v1/*: any*/),
                      (v7/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "nodeType",
                        "args": null,
                        "storageKey": null
                      },
                      (v8/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "stringValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "intValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "booleanValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "floatValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "latitudeValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "longitudeValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "rangeFromValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "rangeToValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isEditable",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isMandatory",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isInstanceProperty",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isDeleted",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "category",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "checkListCategoryDefinitions",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CheckListCategoryDefinition",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      (v9/*: any*/),
                      (v2/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "checklistItemDefinitions",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "CheckListItemDefinition",
                        "plural": true,
                        "selections": [
                          (v0/*: any*/),
                          (v9/*: any*/),
                          (v7/*: any*/),
                          (v8/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "enumValues",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "enumSelectionMode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "helpText",
                            "args": null,
                            "storageKey": null
                          }
                        ]
                      }
                    ]
                  },
                  (v3/*: any*/)
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "workOrderTypes",
        "args": (v6/*: any*/),
        "handle": "connection",
        "key": "Configure_workOrderTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrderTypesQuery",
    "id": null,
    "text": "query WorkOrderTypesQuery {\n  workOrderTypes(first: 500) {\n    edges {\n      node {\n        id\n        name\n        description\n        ...AddEditWorkOrderTypeCard_workOrderType\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment AddEditWorkOrderTypeCard_workOrderType on WorkOrderType {\n  id\n  name\n  description\n  numberOfWorkOrders\n  propertyTypes {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n    isDeleted\n    category\n  }\n  checkListCategoryDefinitions {\n    id\n    title\n    description\n    checklistItemDefinitions {\n      id\n      title\n      type\n      index\n      enumValues\n      enumSelectionMode\n      helpText\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "workOrderTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '3ccc96c93793d16a4cd794405847c2f5';
module.exports = node;
