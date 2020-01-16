/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0ff1a35e93f3210893e943e00117feb0
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditProjectTypeCard_editingProjectType$ref = any;
type ProjectTypeCard_projectType$ref = any;
type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref = any;
export type WorkOrderProjectTypesQueryVariables = {||};
export type WorkOrderProjectTypesQueryResponse = {|
  +projectTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +$fragmentRefs: ProjectTypeCard_projectType$ref & AddEditProjectTypeCard_editingProjectType$ref,
      |}
    |}>
  |},
  +workOrderTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +$fragmentRefs: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref
      |}
    |}>
  |},
|};
export type WorkOrderProjectTypesQuery = {|
  variables: WorkOrderProjectTypesQueryVariables,
  response: WorkOrderProjectTypesQueryResponse,
|};
*/


/*
query WorkOrderProjectTypesQuery {
  projectTypes(first: 50) {
    edges {
      node {
        id
        ...ProjectTypeCard_projectType
        ...AddEditProjectTypeCard_editingProjectType
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
  workOrderTypes(first: 50) {
    edges {
      node {
        ...ProjectTypeWorkOrderTemplatesPanel_workOrderTypes
        id
      }
    }
  }
}

fragment AddEditProjectTypeCard_editingProjectType on ProjectType {
  id
  name
  description
  workOrders {
    id
    type {
      id
      name
    }
  }
  properties {
    id
    name
    type
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
  }
}

fragment ProjectTypeCard_projectType on ProjectType {
  id
  name
  description
  numberOfProjects
  workOrders {
    id
  }
}

fragment ProjectTypeWorkOrderTemplatesPanel_workOrderTypes on WorkOrderType {
  id
  name
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
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v3 = {
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
v4 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50
  }
],
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v6 = [
  (v0/*: any*/),
  (v5/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderProjectTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "projectTypes",
        "name": "__WorkOrderProjectTypesQuery_projectTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "ProjectTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "ProjectType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "ProjectTypeCard_projectType",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "AddEditProjectTypeCard_editingProjectType",
                    "args": null
                  }
                ]
              },
              (v2/*: any*/)
            ]
          },
          (v3/*: any*/)
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderTypes",
        "storageKey": "workOrderTypes(first:50)",
        "args": (v4/*: any*/),
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
                  {
                    "kind": "FragmentSpread",
                    "name": "ProjectTypeWorkOrderTemplatesPanel_workOrderTypes",
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
    "name": "WorkOrderProjectTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "projectTypes",
        "storageKey": "projectTypes(first:50)",
        "args": (v4/*: any*/),
        "concreteType": "ProjectTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "ProjectType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v5/*: any*/),
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
                    "name": "numberOfProjects",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "workOrders",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkOrderDefinition",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "type",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkOrderType",
                        "plural": false,
                        "selections": (v6/*: any*/)
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "properties",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      (v5/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "type",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "index",
                        "args": null,
                        "storageKey": null
                      },
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
                      }
                    ]
                  },
                  (v1/*: any*/)
                ]
              },
              (v2/*: any*/)
            ]
          },
          (v3/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "projectTypes",
        "args": (v4/*: any*/),
        "handle": "connection",
        "key": "WorkOrderProjectTypesQuery_projectTypes",
        "filters": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderTypes",
        "storageKey": "workOrderTypes(first:50)",
        "args": (v4/*: any*/),
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
                "selections": (v6/*: any*/)
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrderProjectTypesQuery",
    "id": null,
    "text": "query WorkOrderProjectTypesQuery {\n  projectTypes(first: 50) {\n    edges {\n      node {\n        id\n        ...ProjectTypeCard_projectType\n        ...AddEditProjectTypeCard_editingProjectType\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n  workOrderTypes(first: 50) {\n    edges {\n      node {\n        ...ProjectTypeWorkOrderTemplatesPanel_workOrderTypes\n        id\n      }\n    }\n  }\n}\n\nfragment AddEditProjectTypeCard_editingProjectType on ProjectType {\n  id\n  name\n  description\n  workOrders {\n    id\n    type {\n      id\n      name\n    }\n  }\n  properties {\n    id\n    name\n    type\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n    isDeleted\n  }\n}\n\nfragment ProjectTypeCard_projectType on ProjectType {\n  id\n  name\n  description\n  numberOfProjects\n  workOrders {\n    id\n  }\n}\n\nfragment ProjectTypeWorkOrderTemplatesPanel_workOrderTypes on WorkOrderType {\n  id\n  name\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "projectTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '510f622e6e5b237f445a90a4a3e2531d';
module.exports = node;
