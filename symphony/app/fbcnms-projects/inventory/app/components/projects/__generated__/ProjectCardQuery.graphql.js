/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 541594022b88bbdde7811183fc28f5cc
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ProjectDetails_project$ref = any;
type ProjectMoreActionsButton_project$ref = any;
export type ProjectCardQueryVariables = {|
  projectId: string
|};
export type ProjectCardQueryResponse = {|
  +project: ?{|
    +$fragmentRefs: ProjectMoreActionsButton_project$ref & ProjectDetails_project$ref
  |}
|};
export type ProjectCardQuery = {|
  variables: ProjectCardQueryVariables,
  response: ProjectCardQueryResponse,
|};
*/


/*
query ProjectCardQuery(
  $projectId: ID!
) {
  project: node(id: $projectId) {
    __typename
    ... on Project {
      ...ProjectMoreActionsButton_project
      ...ProjectDetails_project
    }
    id
  }
}

fragment CommentsBox_comments on Comment {
  ...CommentsLog_comments
}

fragment CommentsLog_comments on Comment {
  id
  ...TextCommentPost_comment
}

fragment LocationBreadcrumbsTitle_locationDetails on Location {
  id
  name
  locationType {
    name
    id
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
}

fragment ProjectDetails_project on Project {
  id
  name
  description
  createdBy {
    id
    email
  }
  type {
    name
    id
  }
  location {
    id
    name
    latitude
    longitude
    locationType {
      mapType
      mapZoomLevel
      id
    }
    ...LocationBreadcrumbsTitle_locationDetails
  }
  properties {
    id
    stringValue
    intValue
    floatValue
    booleanValue
    latitudeValue
    longitudeValue
    rangeFromValue
    rangeToValue
    nodeValue {
      __typename
      id
      name
    }
    propertyType {
      id
      name
      type
      nodeType
      isEditable
      isMandatory
      isInstanceProperty
      stringValue
      intValue
      floatValue
      booleanValue
      latitudeValue
      longitudeValue
      rangeFromValue
      rangeToValue
    }
  }
  workOrders {
    ...ProjectWorkOrdersList_workOrders
    id
  }
  comments {
    ...CommentsBox_comments
    id
  }
}

fragment ProjectMoreActionsButton_project on Project {
  id
  name
  numberOfWorkOrders
}

fragment ProjectWorkOrdersList_workOrders on WorkOrder {
  id
  workOrderType {
    name
    id
  }
  name
  description
  owner {
    id
    email
  }
  creationDate
  installDate
  status
  priority
}

fragment TextCommentPost_comment on Comment {
  id
  author {
    email
    id
  }
  text
  createTime
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
  "name": "__typename",
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
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "email",
  "args": null,
  "storageKey": null
},
v7 = [
  (v3/*: any*/),
  (v6/*: any*/)
],
v8 = [
  (v4/*: any*/),
  (v3/*: any*/)
],
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ProjectCardQuery",
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
              {
                "kind": "FragmentSpread",
                "name": "ProjectMoreActionsButton_project",
                "args": null
              },
              {
                "kind": "FragmentSpread",
                "name": "ProjectDetails_project",
                "args": null
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ProjectCardQuery",
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
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Project",
            "selections": [
              (v4/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "numberOfWorkOrders",
                "args": null,
                "storageKey": null
              },
              (v5/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "createdBy",
                "storageKey": null,
                "args": null,
                "concreteType": "User",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "type",
                "storageKey": null,
                "args": null,
                "concreteType": "ProjectType",
                "plural": false,
                "selections": (v8/*: any*/)
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
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationType",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "mapType",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "mapZoomLevel",
                        "args": null,
                        "storageKey": null
                      },
                      (v3/*: any*/),
                      (v4/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationHierarchy",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Location",
                    "plural": true,
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "locationType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationType",
                        "plural": false,
                        "selections": (v8/*: any*/)
                      }
                    ]
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "properties",
                "storageKey": null,
                "args": null,
                "concreteType": "Property",
                "plural": true,
                "selections": [
                  (v3/*: any*/),
                  (v9/*: any*/),
                  (v10/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "nodeValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": [
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v4/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": false,
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
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
                        "name": "nodeType",
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
                      (v9/*: any*/),
                      (v10/*: any*/),
                      (v11/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/),
                      (v14/*: any*/),
                      (v15/*: any*/),
                      (v16/*: any*/)
                    ]
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "workOrders",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrder",
                "plural": true,
                "selections": [
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "workOrderType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkOrderType",
                    "plural": false,
                    "selections": (v8/*: any*/)
                  },
                  (v4/*: any*/),
                  (v5/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "owner",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "User",
                    "plural": false,
                    "selections": (v7/*: any*/)
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
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "priority",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "comments",
                "storageKey": null,
                "args": null,
                "concreteType": "Comment",
                "plural": true,
                "selections": [
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "author",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "User",
                    "plural": false,
                    "selections": [
                      (v6/*: any*/),
                      (v3/*: any*/)
                    ]
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "text",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "createTime",
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
    "name": "ProjectCardQuery",
    "id": null,
    "text": "query ProjectCardQuery(\n  $projectId: ID!\n) {\n  project: node(id: $projectId) {\n    __typename\n    ... on Project {\n      ...ProjectMoreActionsButton_project\n      ...ProjectDetails_project\n    }\n    id\n  }\n}\n\nfragment CommentsBox_comments on Comment {\n  ...CommentsLog_comments\n}\n\nfragment CommentsLog_comments on Comment {\n  id\n  ...TextCommentPost_comment\n}\n\nfragment LocationBreadcrumbsTitle_locationDetails on Location {\n  id\n  name\n  locationType {\n    name\n    id\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n}\n\nfragment ProjectDetails_project on Project {\n  id\n  name\n  description\n  createdBy {\n    id\n    email\n  }\n  type {\n    name\n    id\n  }\n  location {\n    id\n    name\n    latitude\n    longitude\n    locationType {\n      mapType\n      mapZoomLevel\n      id\n    }\n    ...LocationBreadcrumbsTitle_locationDetails\n  }\n  properties {\n    id\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    nodeValue {\n      __typename\n      id\n      name\n    }\n    propertyType {\n      id\n      name\n      type\n      nodeType\n      isEditable\n      isMandatory\n      isInstanceProperty\n      stringValue\n      intValue\n      floatValue\n      booleanValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n    }\n  }\n  workOrders {\n    ...ProjectWorkOrdersList_workOrders\n    id\n  }\n  comments {\n    ...CommentsBox_comments\n    id\n  }\n}\n\nfragment ProjectMoreActionsButton_project on Project {\n  id\n  name\n  numberOfWorkOrders\n}\n\nfragment ProjectWorkOrdersList_workOrders on WorkOrder {\n  id\n  workOrderType {\n    name\n    id\n  }\n  name\n  description\n  owner {\n    id\n    email\n  }\n  creationDate\n  installDate\n  status\n  priority\n}\n\nfragment TextCommentPost_comment on Comment {\n  id\n  author {\n    email\n    id\n  }\n  text\n  createTime\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4e882cb1745fd90983186a7b66242bfa';
module.exports = node;
