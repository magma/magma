/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash fdfb23332182757e5c727d50b8736c70
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type WorkOrderDetails_workOrder$ref = any;
export type WorkOrderCardQueryVariables = {|
  workOrderId: string
|};
export type WorkOrderCardQueryResponse = {|
  +workOrder: ?{|
    +id?: string,
    +name?: string,
    +$fragmentRefs: WorkOrderDetails_workOrder$ref,
  |}
|};
export type WorkOrderCardQuery = {|
  variables: WorkOrderCardQueryVariables,
  response: WorkOrderCardQueryResponse,
|};
*/


/*
query WorkOrderCardQuery(
  $workOrderId: ID!
) {
  workOrder: node(id: $workOrderId) {
    __typename
    ... on WorkOrder {
      id
      name
      ...WorkOrderDetails_workOrder
    }
    id
  }
}

fragment ActivityPost_activity on Activity {
  id
  author {
    email
    id
  }
  isCreate
  changedField
  newRelatedNode {
    __typename
    ... on User {
      id
      email
    }
    id
  }
  oldRelatedNode {
    __typename
    ... on User {
      id
      email
    }
    id
  }
  oldValue
  newValue
  createTime
}

fragment CommentsActivitiesBox_activities on Activity {
  ...CommentsActivitiesLog_activities
}

fragment CommentsActivitiesBox_comments on Comment {
  ...CommentsActivitiesLog_comments
}

fragment CommentsActivitiesLog_activities on Activity {
  id
  createTime
  ...ActivityPost_activity
}

fragment CommentsActivitiesLog_comments on Comment {
  id
  createTime
  ...TextCommentPost_comment
}

fragment DocumentTable_files on File {
  id
  fileName
  category
  ...FileAttachment_file
}

fragment DocumentTable_hyperlinks on Hyperlink {
  id
  category
  url
  displayName
  ...HyperlinkTableRow_hyperlink
}

fragment EntityDocumentsTable_files on File {
  ...DocumentTable_files
}

fragment EntityDocumentsTable_hyperlinks on Hyperlink {
  ...DocumentTable_hyperlinks
}

fragment EquipmentBreadcrumbs_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
  positionHierarchy {
    id
    definition {
      id
      name
      visibleLabel
    }
    parentEquipment {
      id
      name
      equipmentType {
        id
        name
      }
    }
  }
}

fragment FileAttachment_file on File {
  id
  fileName
  sizeInBytes
  uploaded
  fileType
  storeKey
  category
  ...ImageDialog_img
}

fragment HyperlinkTableRow_hyperlink on Hyperlink {
  id
  category
  url
  displayName
  createTime
}

fragment ImageDialog_img on File {
  storeKey
  fileName
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

fragment TextCommentPost_comment on Comment {
  id
  author {
    email
    id
  }
  text
  createTime
}

fragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  parentLocation {
    id
    name
    locationType {
      id
      name
    }
  }
  parentPosition {
    id
    definition {
      name
      visibleLabel
      id
    }
    parentEquipment {
      id
      name
    }
  }
}

fragment WorkOrderDetailsPaneLinkItem_link on Link {
  id
  futureState
  ports {
    id
    definition {
      id
      name
      visibleLabel
      portType {
        linkPropertyTypes {
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
          isInstanceProperty
          isMandatory
          category
          isDeleted
        }
        id
      }
    }
    parentEquipment {
      id
      name
      futureState
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          bandwidth
          portType {
            id
            name
          }
        }
      }
      ...EquipmentBreadcrumbs_equipment
    }
    serviceEndpoints {
      definition {
        role
        id
      }
      service {
        name
        id
      }
      id
    }
  }
  workOrder {
    id
    status
  }
  properties {
    id
    propertyType {
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
      isInstanceProperty
      isMandatory
      category
      isDeleted
    }
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
  }
  services {
    id
    name
  }
}

fragment WorkOrderDetailsPane_workOrder on WorkOrder {
  id
  name
  equipmentToAdd {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  equipmentToRemove {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  linksToAdd {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
  linksToRemove {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
}

fragment WorkOrderDetails_workOrder on WorkOrder {
  id
  name
  description
  workOrderType {
    name
    id
  }
  location {
    name
    id
    latitude
    longitude
    locationType {
      mapType
      mapZoomLevel
      id
    }
    ...LocationBreadcrumbsTitle_locationDetails
  }
  owner {
    id
    email
  }
  assignedTo {
    id
    email
  }
  creationDate
  installDate
  status
  priority
  ...WorkOrderDetailsPane_workOrder
  properties {
    id
    propertyType {
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
      isInstanceProperty
      isMandatory
      category
      isDeleted
    }
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
  }
  images {
    ...EntityDocumentsTable_files
    id
  }
  files {
    ...EntityDocumentsTable_files
    id
  }
  hyperlinks {
    ...EntityDocumentsTable_hyperlinks
    id
  }
  comments {
    ...CommentsActivitiesBox_comments
    id
  }
  activities {
    ...CommentsActivitiesBox_activities
    id
  }
  project {
    name
    id
    type {
      id
      name
    }
  }
  checkListCategories {
    id
    title
    description
    checkList {
      id
      index
      type
      title
      helpText
      checked
      enumValues
      stringValue
      enumSelectionMode
      selectedEnumValues
      yesNoResponse
      files {
        id
        fileName
        sizeInBytes
        modified
        uploaded
        fileType
        storeKey
        category
      }
      cellData {
        id
        networkType
        signalStrength
        timestamp
        baseStationID
        networkID
        systemID
        cellID
        locationAreaCode
        mobileCountryCode
        mobileNetworkCode
        primaryScramblingCode
        operator
        arfcn
        physicalCellID
        trackingAreaCode
        timingAdvance
        earfcn
        uarfcn
        latitude
        longitude
      }
      wifiData {
        id
        timestamp
        frequency
        channel
        bssid
        strength
        ssid
        band
        channelWidth
        capabilities
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
    "name": "workOrderId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "workOrderId"
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
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
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
v6 = [
  (v3/*: any*/),
  (v2/*: any*/)
],
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitude",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitude",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationHierarchy",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationType",
      "storageKey": null,
      "args": null,
      "concreteType": "LocationType",
      "plural": false,
      "selections": (v6/*: any*/)
    }
  ]
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "email",
  "args": null,
  "storageKey": null
},
v11 = [
  (v2/*: any*/),
  (v10/*: any*/)
],
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "status",
  "args": null,
  "storageKey": null
},
v13 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v14 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v13/*: any*/)
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v16 = [
  (v2/*: any*/),
  (v3/*: any*/),
  (v14/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentLocation",
    "storageKey": null,
    "args": null,
    "concreteType": "Location",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      (v3/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "locationType",
        "storageKey": null,
        "args": null,
        "concreteType": "LocationType",
        "plural": false,
        "selections": (v13/*: any*/)
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentPosition",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPosition",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "definition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPositionDefinition",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          (v15/*: any*/),
          (v2/*: any*/)
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentEquipment",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": false,
        "selections": (v13/*: any*/)
      }
    ]
  }
],
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v22 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v23 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v24 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v25 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v26 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v27 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v28 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "category",
  "args": null,
  "storageKey": null
},
v29 = [
  (v2/*: any*/),
  (v3/*: any*/),
  (v18/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "nodeType",
    "args": null,
    "storageKey": null
  },
  (v19/*: any*/),
  (v20/*: any*/),
  (v21/*: any*/),
  (v22/*: any*/),
  (v23/*: any*/),
  (v24/*: any*/),
  (v25/*: any*/),
  (v26/*: any*/),
  (v27/*: any*/),
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
    "name": "isInstanceProperty",
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
  (v28/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isDeleted",
    "args": null,
    "storageKey": null
  }
],
v30 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "properties",
  "storageKey": null,
  "args": null,
  "concreteType": "Property",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
      "selections": (v29/*: any*/)
    },
    (v20/*: any*/),
    (v21/*: any*/),
    (v23/*: any*/),
    (v22/*: any*/),
    (v24/*: any*/),
    (v25/*: any*/),
    (v26/*: any*/),
    (v27/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "nodeValue",
      "storageKey": null,
      "args": null,
      "concreteType": null,
      "plural": false,
      "selections": [
        (v4/*: any*/),
        (v2/*: any*/),
        (v3/*: any*/)
      ]
    }
  ]
},
v31 = [
  (v2/*: any*/),
  (v17/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "ports",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPort",
    "plural": true,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "definition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortDefinition",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v15/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "portType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortType",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linkPropertyTypes",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": true,
                "selections": (v29/*: any*/)
              },
              (v2/*: any*/)
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentEquipment",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v17/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentType",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "portDefinitions",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPortDefinition",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v15/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "bandwidth",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortType",
                    "plural": false,
                    "selections": (v13/*: any*/)
                  }
                ]
              }
            ]
          },
          (v9/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "positionHierarchy",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPosition",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "definition",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPositionDefinition",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v15/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "parentEquipment",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v14/*: any*/)
                ]
              }
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceEndpoints",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceEndpoint",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "definition",
            "storageKey": null,
            "args": null,
            "concreteType": "ServiceEndpointDefinition",
            "plural": false,
            "selections": [
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "role",
                "args": null,
                "storageKey": null
              },
              (v2/*: any*/)
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "service",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": false,
            "selections": (v6/*: any*/)
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "workOrder",
    "storageKey": null,
    "args": null,
    "concreteType": "WorkOrder",
    "plural": false,
    "selections": [
      (v2/*: any*/),
      (v12/*: any*/)
    ]
  },
  (v30/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "services",
    "storageKey": null,
    "args": null,
    "concreteType": "Service",
    "plural": true,
    "selections": (v13/*: any*/)
  }
],
v32 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileName",
  "args": null,
  "storageKey": null
},
v33 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "sizeInBytes",
  "args": null,
  "storageKey": null
},
v34 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "uploaded",
  "args": null,
  "storageKey": null
},
v35 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileType",
  "args": null,
  "storageKey": null
},
v36 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "storeKey",
  "args": null,
  "storageKey": null
},
v37 = [
  (v2/*: any*/),
  (v32/*: any*/),
  (v28/*: any*/),
  (v33/*: any*/),
  (v34/*: any*/),
  (v35/*: any*/),
  (v36/*: any*/)
],
v38 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "createTime",
  "args": null,
  "storageKey": null
},
v39 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "author",
  "storageKey": null,
  "args": null,
  "concreteType": "User",
  "plural": false,
  "selections": [
    (v10/*: any*/),
    (v2/*: any*/)
  ]
},
v40 = [
  (v4/*: any*/),
  (v2/*: any*/),
  {
    "kind": "InlineFragment",
    "type": "User",
    "selections": [
      (v10/*: any*/)
    ]
  }
],
v41 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
},
v42 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "timestamp",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "WorkOrder",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "FragmentSpread",
                "name": "WorkOrderDetails_workOrder",
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
    "name": "WorkOrderCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          (v4/*: any*/),
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "WorkOrder",
            "selections": [
              (v3/*: any*/),
              (v5/*: any*/),
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
                "name": "location",
                "storageKey": null,
                "args": null,
                "concreteType": "Location",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v2/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
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
                      (v2/*: any*/),
                      (v3/*: any*/)
                    ]
                  },
                  (v9/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "owner",
                "storageKey": null,
                "args": null,
                "concreteType": "User",
                "plural": false,
                "selections": (v11/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "assignedTo",
                "storageKey": null,
                "args": null,
                "concreteType": "User",
                "plural": false,
                "selections": (v11/*: any*/)
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
              (v12/*: any*/),
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
                "name": "equipmentToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v16/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v16/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v31/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v31/*: any*/)
              },
              (v30/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "images",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v37/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "files",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v37/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "hyperlinks",
                "storageKey": null,
                "args": null,
                "concreteType": "Hyperlink",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v28/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "url",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "displayName",
                    "args": null,
                    "storageKey": null
                  },
                  (v38/*: any*/)
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
                  (v2/*: any*/),
                  (v38/*: any*/),
                  (v39/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "text",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "activities",
                "storageKey": null,
                "args": null,
                "concreteType": "Activity",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v38/*: any*/),
                  (v39/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "isCreate",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "changedField",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "newRelatedNode",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": (v40/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "oldRelatedNode",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": (v40/*: any*/)
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "oldValue",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "newValue",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "project",
                "storageKey": null,
                "args": null,
                "concreteType": "Project",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "type",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ProjectType",
                    "plural": false,
                    "selections": (v13/*: any*/)
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "checkListCategories",
                "storageKey": null,
                "args": null,
                "concreteType": "CheckListCategory",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v41/*: any*/),
                  (v5/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "checkList",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CheckListItem",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      (v19/*: any*/),
                      (v18/*: any*/),
                      (v41/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "helpText",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "checked",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "enumValues",
                        "args": null,
                        "storageKey": null
                      },
                      (v20/*: any*/),
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
                        "name": "selectedEnumValues",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "yesNoResponse",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "files",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "File",
                        "plural": true,
                        "selections": [
                          (v2/*: any*/),
                          (v32/*: any*/),
                          (v33/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "modified",
                            "args": null,
                            "storageKey": null
                          },
                          (v34/*: any*/),
                          (v35/*: any*/),
                          (v36/*: any*/),
                          (v28/*: any*/)
                        ]
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "cellData",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "SurveyCellScan",
                        "plural": true,
                        "selections": [
                          (v2/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "networkType",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "signalStrength",
                            "args": null,
                            "storageKey": null
                          },
                          (v42/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "baseStationID",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "networkID",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "systemID",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "cellID",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "locationAreaCode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "mobileCountryCode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "mobileNetworkCode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "primaryScramblingCode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "operator",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "arfcn",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "physicalCellID",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "trackingAreaCode",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "timingAdvance",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "earfcn",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "uarfcn",
                            "args": null,
                            "storageKey": null
                          },
                          (v7/*: any*/),
                          (v8/*: any*/)
                        ]
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "wifiData",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "SurveyWiFiScan",
                        "plural": true,
                        "selections": [
                          (v2/*: any*/),
                          (v42/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "frequency",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "channel",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "bssid",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "strength",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "ssid",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "band",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "channelWidth",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "capabilities",
                            "args": null,
                            "storageKey": null
                          },
                          (v7/*: any*/),
                          (v8/*: any*/)
                        ]
                      }
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
    "name": "WorkOrderCardQuery",
    "id": null,
    "text": "query WorkOrderCardQuery(\n  $workOrderId: ID!\n) {\n  workOrder: node(id: $workOrderId) {\n    __typename\n    ... on WorkOrder {\n      id\n      name\n      ...WorkOrderDetails_workOrder\n    }\n    id\n  }\n}\n\nfragment ActivityPost_activity on Activity {\n  id\n  author {\n    email\n    id\n  }\n  isCreate\n  changedField\n  newRelatedNode {\n    __typename\n    ... on User {\n      id\n      email\n    }\n    id\n  }\n  oldRelatedNode {\n    __typename\n    ... on User {\n      id\n      email\n    }\n    id\n  }\n  oldValue\n  newValue\n  createTime\n}\n\nfragment CommentsActivitiesBox_activities on Activity {\n  ...CommentsActivitiesLog_activities\n}\n\nfragment CommentsActivitiesBox_comments on Comment {\n  ...CommentsActivitiesLog_comments\n}\n\nfragment CommentsActivitiesLog_activities on Activity {\n  id\n  createTime\n  ...ActivityPost_activity\n}\n\nfragment CommentsActivitiesLog_comments on Comment {\n  id\n  createTime\n  ...TextCommentPost_comment\n}\n\nfragment DocumentTable_files on File {\n  id\n  fileName\n  category\n  ...FileAttachment_file\n}\n\nfragment DocumentTable_hyperlinks on Hyperlink {\n  id\n  category\n  url\n  displayName\n  ...HyperlinkTableRow_hyperlink\n}\n\nfragment EntityDocumentsTable_files on File {\n  ...DocumentTable_files\n}\n\nfragment EntityDocumentsTable_hyperlinks on Hyperlink {\n  ...DocumentTable_hyperlinks\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...ImageDialog_img\n}\n\nfragment HyperlinkTableRow_hyperlink on Hyperlink {\n  id\n  category\n  url\n  displayName\n  createTime\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n\nfragment LocationBreadcrumbsTitle_locationDetails on Location {\n  id\n  name\n  locationType {\n    name\n    id\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n}\n\nfragment TextCommentPost_comment on Comment {\n  id\n  author {\n    email\n    id\n  }\n  text\n  createTime\n}\n\nfragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  parentLocation {\n    id\n    name\n    locationType {\n      id\n      name\n    }\n  }\n  parentPosition {\n    id\n    definition {\n      name\n      visibleLabel\n      id\n    }\n    parentEquipment {\n      id\n      name\n    }\n  }\n}\n\nfragment WorkOrderDetailsPaneLinkItem_link on Link {\n  id\n  futureState\n  ports {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n      portType {\n        linkPropertyTypes {\n          id\n          name\n          type\n          nodeType\n          index\n          stringValue\n          intValue\n          booleanValue\n          floatValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          isEditable\n          isInstanceProperty\n          isMandatory\n          category\n          isDeleted\n        }\n        id\n      }\n    }\n    parentEquipment {\n      id\n      name\n      futureState\n      equipmentType {\n        id\n        name\n        portDefinitions {\n          id\n          name\n          visibleLabel\n          bandwidth\n          portType {\n            id\n            name\n          }\n        }\n      }\n      ...EquipmentBreadcrumbs_equipment\n    }\n    serviceEndpoints {\n      definition {\n        role\n        id\n      }\n      service {\n        name\n        id\n      }\n      id\n    }\n  }\n  workOrder {\n    id\n    status\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      nodeType\n      index\n      stringValue\n      intValue\n      booleanValue\n      floatValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n      isEditable\n      isInstanceProperty\n      isMandatory\n      category\n      isDeleted\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    nodeValue {\n      __typename\n      id\n      name\n    }\n  }\n  services {\n    id\n    name\n  }\n}\n\nfragment WorkOrderDetailsPane_workOrder on WorkOrder {\n  id\n  name\n  equipmentToAdd {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  equipmentToRemove {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  linksToAdd {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n  linksToRemove {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n}\n\nfragment WorkOrderDetails_workOrder on WorkOrder {\n  id\n  name\n  description\n  workOrderType {\n    name\n    id\n  }\n  location {\n    name\n    id\n    latitude\n    longitude\n    locationType {\n      mapType\n      mapZoomLevel\n      id\n    }\n    ...LocationBreadcrumbsTitle_locationDetails\n  }\n  owner {\n    id\n    email\n  }\n  assignedTo {\n    id\n    email\n  }\n  creationDate\n  installDate\n  status\n  priority\n  ...WorkOrderDetailsPane_workOrder\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      nodeType\n      index\n      stringValue\n      intValue\n      booleanValue\n      floatValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n      isEditable\n      isInstanceProperty\n      isMandatory\n      category\n      isDeleted\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    nodeValue {\n      __typename\n      id\n      name\n    }\n  }\n  images {\n    ...EntityDocumentsTable_files\n    id\n  }\n  files {\n    ...EntityDocumentsTable_files\n    id\n  }\n  hyperlinks {\n    ...EntityDocumentsTable_hyperlinks\n    id\n  }\n  comments {\n    ...CommentsActivitiesBox_comments\n    id\n  }\n  activities {\n    ...CommentsActivitiesBox_activities\n    id\n  }\n  project {\n    name\n    id\n    type {\n      id\n      name\n    }\n  }\n  checkListCategories {\n    id\n    title\n    description\n    checkList {\n      id\n      index\n      type\n      title\n      helpText\n      checked\n      enumValues\n      stringValue\n      enumSelectionMode\n      selectedEnumValues\n      yesNoResponse\n      files {\n        id\n        fileName\n        sizeInBytes\n        modified\n        uploaded\n        fileType\n        storeKey\n        category\n      }\n      cellData {\n        id\n        networkType\n        signalStrength\n        timestamp\n        baseStationID\n        networkID\n        systemID\n        cellID\n        locationAreaCode\n        mobileCountryCode\n        mobileNetworkCode\n        primaryScramblingCode\n        operator\n        arfcn\n        physicalCellID\n        trackingAreaCode\n        timingAdvance\n        earfcn\n        uarfcn\n        latitude\n        longitude\n      }\n      wifiData {\n        id\n        timestamp\n        frequency\n        channel\n        bssid\n        strength\n        ssid\n        band\n        channelWidth\n        capabilities\n        latitude\n        longitude\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ba8f1936a3ba4fd0f7aa7b7b90970432';
module.exports = node;
