/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 62b9adb12f43b92d5e2d97cb3a40a7ba
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type DynamicPropertiesGrid_properties$ref = any;
type DynamicPropertiesGrid_propertyTypes$ref = any;
type EquipmentTable_equipment$ref = any;
type LocationBreadcrumbsTitle_locationDetails$ref = any;
type LocationDocumentsCard_location$ref = any;
type LocationFloorPlansTab_location$ref = any;
type LocationMoreActionsButton_location$ref = any;
type LocationSiteSurveyTab_location$ref = any;
type PropertyFormField_property$ref = any;
type PropertyTypeFormField_propertyType$ref = any;
export type LocationPropertiesCardQueryVariables = {|
  locationId: string
|};
export type LocationPropertiesCardQueryResponse = {|
  +location: ?{|
    +id?: string,
    +name?: string,
    +latitude?: number,
    +longitude?: number,
    +externalId?: ?string,
    +locationType?: {|
      +id: string,
      +name: string,
      +mapType: ?string,
      +mapZoomLevel: ?number,
      +propertyTypes: $ReadOnlyArray<?{|
        +$fragmentRefs: PropertyTypeFormField_propertyType$ref & DynamicPropertiesGrid_propertyTypes$ref
      |}>,
    |},
    +parentLocation?: ?{|
      +id: string
    |},
    +children?: $ReadOnlyArray<?{|
      +id: string
    |}>,
    +equipments?: $ReadOnlyArray<?{|
      +$fragmentRefs: EquipmentTable_equipment$ref
    |}>,
    +properties?: $ReadOnlyArray<?{|
      +$fragmentRefs: PropertyFormField_property$ref & DynamicPropertiesGrid_properties$ref
    |}>,
    +images?: $ReadOnlyArray<?{|
      +id: string
    |}>,
    +files?: $ReadOnlyArray<?{|
      +id: string
    |}>,
    +hyperlinks?: $ReadOnlyArray<{|
      +id: string
    |}>,
    +surveys?: $ReadOnlyArray<?{|
      +id: string
    |}>,
    +$fragmentRefs: LocationBreadcrumbsTitle_locationDetails$ref & LocationSiteSurveyTab_location$ref & LocationDocumentsCard_location$ref & LocationFloorPlansTab_location$ref & LocationMoreActionsButton_location$ref,
  |}
|};
export type LocationPropertiesCardQuery = {|
  variables: LocationPropertiesCardQueryVariables,
  response: LocationPropertiesCardQueryResponse,
|};
*/


/*
query LocationPropertiesCardQuery(
  $locationId: ID!
) {
  location: node(id: $locationId) {
    __typename
    ... on Location {
      id
      name
      latitude
      longitude
      externalId
      locationType {
        id
        name
        mapType
        mapZoomLevel
        propertyTypes {
          ...PropertyTypeFormField_propertyType
          ...DynamicPropertiesGrid_propertyTypes
          id
        }
      }
      ...LocationBreadcrumbsTitle_locationDetails
      parentLocation {
        id
      }
      children {
        id
      }
      equipments {
        ...EquipmentTable_equipment
        id
      }
      properties {
        ...PropertyFormField_property
        ...DynamicPropertiesGrid_properties
        id
      }
      images {
        id
      }
      files {
        id
      }
      hyperlinks {
        id
      }
      surveys {
        id
      }
      ...LocationSiteSurveyTab_location
      ...LocationDocumentsCard_location
      ...LocationFloorPlansTab_location
      ...LocationMoreActionsButton_location
    }
    id
  }
}

fragment DocumentMenu_document on File {
  id
  fileName
  storeKey
  fileType
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

fragment DynamicPropertiesGrid_properties on Property {
  ...PropertyFormField_property
  propertyType {
    id
    index
  }
}

fragment DynamicPropertiesGrid_propertyTypes on PropertyType {
  id
  name
  index
  isInstanceProperty
  type
  stringValue
  intValue
  booleanValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  floatValue
}

fragment EntityDocumentsTable_files on File {
  ...DocumentTable_files
}

fragment EntityDocumentsTable_hyperlinks on Hyperlink {
  ...DocumentTable_hyperlinks
}

fragment EquipmentTable_equipment on Equipment {
  id
  name
  futureState
  equipmentType {
    id
    name
  }
  workOrder {
    id
    status
  }
  device {
    up
  }
  services {
    id
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
  ...DocumentMenu_document
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

fragment LocationDocumentsCard_location on Location {
  id
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
}

fragment LocationFloorPlansTab_location on Location {
  id
  floorPlans {
    id
    name
    image {
      ...FileAttachment_file
      id
    }
  }
}

fragment LocationMoreActionsButton_location on Location {
  id
  parentLocation {
    id
  }
  children {
    id
  }
  equipments {
    id
  }
  images {
    id
  }
  files {
    id
  }
  surveys {
    id
  }
}

fragment LocationSiteSurveyTab_location on Location {
  id
  siteSurveyNeeded
  surveys {
    id
    completionTimestamp
    name
    ownerName
    sourceFile {
      id
      fileName
      storeKey
    }
    ...SiteSurveyPane_survey
  }
}

fragment PropertyFormField_property on Property {
  id
  propertyType {
    id
    name
    type
    isEditable
    isMandatory
    isInstanceProperty
    stringValue
  }
  stringValue
  intValue
  floatValue
  booleanValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  equipmentValue {
    id
    name
  }
  locationValue {
    id
    name
  }
  serviceValue {
    id
    name
  }
}

fragment PropertyTypeFormField_propertyType on PropertyType {
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
  isInstanceProperty
  isMandatory
}

fragment SiteSurveyPane_survey on Survey {
  name
  completionTimestamp
  surveyResponses {
    id
    questionText
    formName
    formIndex
    questionIndex
    ...SiteSurveyQuestionReply_question
  }
}

fragment SiteSurveyQuestionReplyCellData_data on SurveyQuestion {
  cellData {
    networkType
    signalStrength
    baseStationID
    cellID
    locationAreaCode
    mobileCountryCode
    mobileNetworkCode
    id
  }
}

fragment SiteSurveyQuestionReplyWifiData_data on SurveyQuestion {
  wifiData {
    band
    bssid
    channel
    frequency
    strength
    ssid
    id
  }
}

fragment SiteSurveyQuestionReply_question on SurveyQuestion {
  questionFormat
  longitude
  latitude
  boolData
  textData
  emailData
  phoneData
  floatData
  intData
  dateData
  photoData {
    storeKey
    id
  }
  ...SiteSurveyQuestionReplyWifiData_data
  ...SiteSurveyQuestionReplyCellData_data
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "locationId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "locationId"
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
  "name": "latitude",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitude",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "externalId",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "mapType",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "mapZoomLevel",
  "args": null,
  "storageKey": null
},
v9 = [
  (v2/*: any*/)
],
v10 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "parentLocation",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": false,
  "selections": (v9/*: any*/)
},
v11 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "children",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": true,
  "selections": (v9/*: any*/)
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v22 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v23 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v24 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v25 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v26 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileName",
  "args": null,
  "storageKey": null
},
v27 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "category",
  "args": null,
  "storageKey": null
},
v28 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "sizeInBytes",
  "args": null,
  "storageKey": null
},
v29 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "uploaded",
  "args": null,
  "storageKey": null
},
v30 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "fileType",
  "args": null,
  "storageKey": null
},
v31 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "storeKey",
  "args": null,
  "storageKey": null
},
v32 = [
  (v2/*: any*/),
  (v26/*: any*/),
  (v27/*: any*/),
  (v28/*: any*/),
  (v29/*: any*/),
  (v30/*: any*/),
  (v31/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationPropertiesCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "FragmentSpread",
                        "name": "PropertyTypeFormField_propertyType",
                        "args": null
                      },
                      {
                        "kind": "FragmentSpread",
                        "name": "DynamicPropertiesGrid_propertyTypes",
                        "args": null
                      }
                    ]
                  }
                ]
              },
              (v10/*: any*/),
              (v11/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipments",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "EquipmentTable_equipment",
                    "args": null
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
                  {
                    "kind": "FragmentSpread",
                    "name": "PropertyFormField_property",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "DynamicPropertiesGrid_properties",
                    "args": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "images",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v9/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "files",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v9/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "hyperlinks",
                "storageKey": null,
                "args": null,
                "concreteType": "Hyperlink",
                "plural": true,
                "selections": (v9/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "surveys",
                "storageKey": null,
                "args": null,
                "concreteType": "Survey",
                "plural": true,
                "selections": (v9/*: any*/)
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationBreadcrumbsTitle_locationDetails",
                "args": null
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationSiteSurveyTab_location",
                "args": null
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationDocumentsCard_location",
                "args": null
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationFloorPlansTab_location",
                "args": null
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationMoreActionsButton_location",
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
    "name": "LocationPropertiesCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
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
            "type": "Location",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/),
                      (v14/*: any*/),
                      (v15/*: any*/),
                      (v16/*: any*/),
                      (v17/*: any*/),
                      (v18/*: any*/),
                      (v19/*: any*/),
                      (v20/*: any*/),
                      (v21/*: any*/),
                      (v22/*: any*/),
                      (v23/*: any*/),
                      (v24/*: any*/)
                    ]
                  }
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
                    "selections": [
                      (v3/*: any*/),
                      (v2/*: any*/)
                    ]
                  }
                ]
              },
              (v10/*: any*/),
              (v11/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipments",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "futureState",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentType",
                    "plural": false,
                    "selections": (v25/*: any*/)
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
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "status",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "device",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Device",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "up",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "services",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Service",
                    "plural": true,
                    "selections": (v9/*: any*/)
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
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": false,
                    "selections": [
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v12/*: any*/),
                      (v22/*: any*/),
                      (v24/*: any*/),
                      (v23/*: any*/),
                      (v14/*: any*/),
                      (v13/*: any*/)
                    ]
                  },
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v17/*: any*/),
                  (v16/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  (v20/*: any*/),
                  (v21/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Equipment",
                    "plural": false,
                    "selections": (v25/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Location",
                    "plural": false,
                    "selections": (v25/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "serviceValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Service",
                    "plural": false,
                    "selections": (v25/*: any*/)
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "images",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v32/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "files",
                "storageKey": null,
                "args": null,
                "concreteType": "File",
                "plural": true,
                "selections": (v32/*: any*/)
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
                  (v27/*: any*/),
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
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "createTime",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "surveys",
                "storageKey": null,
                "args": null,
                "concreteType": "Survey",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "completionTimestamp",
                    "args": null,
                    "storageKey": null
                  },
                  (v3/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "ownerName",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "sourceFile",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "File",
                    "plural": false,
                    "selections": [
                      (v2/*: any*/),
                      (v26/*: any*/),
                      (v31/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "surveyResponses",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "SurveyQuestion",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "questionText",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "formName",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "formIndex",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "questionIndex",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "questionFormat",
                        "args": null,
                        "storageKey": null
                      },
                      (v5/*: any*/),
                      (v4/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "boolData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "textData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "emailData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "phoneData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "floatData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "intData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "dateData",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "photoData",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "File",
                        "plural": false,
                        "selections": [
                          (v31/*: any*/),
                          (v2/*: any*/)
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
                            "name": "bssid",
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
                            "name": "frequency",
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
                          (v2/*: any*/)
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
                          (v2/*: any*/)
                        ]
                      }
                    ]
                  }
                ]
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "siteSurveyNeeded",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "floorPlans",
                "storageKey": null,
                "args": null,
                "concreteType": "FloorPlan",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "image",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "File",
                    "plural": false,
                    "selections": [
                      (v2/*: any*/),
                      (v26/*: any*/),
                      (v28/*: any*/),
                      (v29/*: any*/),
                      (v30/*: any*/),
                      (v31/*: any*/),
                      (v27/*: any*/)
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
    "name": "LocationPropertiesCardQuery",
    "id": null,
    "text": "query LocationPropertiesCardQuery(\n  $locationId: ID!\n) {\n  location: node(id: $locationId) {\n    __typename\n    ... on Location {\n      id\n      name\n      latitude\n      longitude\n      externalId\n      locationType {\n        id\n        name\n        mapType\n        mapZoomLevel\n        propertyTypes {\n          ...PropertyTypeFormField_propertyType\n          ...DynamicPropertiesGrid_propertyTypes\n          id\n        }\n      }\n      ...LocationBreadcrumbsTitle_locationDetails\n      parentLocation {\n        id\n      }\n      children {\n        id\n      }\n      equipments {\n        ...EquipmentTable_equipment\n        id\n      }\n      properties {\n        ...PropertyFormField_property\n        ...DynamicPropertiesGrid_properties\n        id\n      }\n      images {\n        id\n      }\n      files {\n        id\n      }\n      hyperlinks {\n        id\n      }\n      surveys {\n        id\n      }\n      ...LocationSiteSurveyTab_location\n      ...LocationDocumentsCard_location\n      ...LocationFloorPlansTab_location\n      ...LocationMoreActionsButton_location\n    }\n    id\n  }\n}\n\nfragment DocumentMenu_document on File {\n  id\n  fileName\n  storeKey\n  fileType\n}\n\nfragment DocumentTable_files on File {\n  id\n  fileName\n  category\n  ...FileAttachment_file\n}\n\nfragment DocumentTable_hyperlinks on Hyperlink {\n  id\n  category\n  url\n  displayName\n  ...HyperlinkTableRow_hyperlink\n}\n\nfragment DynamicPropertiesGrid_properties on Property {\n  ...PropertyFormField_property\n  propertyType {\n    id\n    index\n  }\n}\n\nfragment DynamicPropertiesGrid_propertyTypes on PropertyType {\n  id\n  name\n  index\n  isInstanceProperty\n  type\n  stringValue\n  intValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  floatValue\n}\n\nfragment EntityDocumentsTable_files on File {\n  ...DocumentTable_files\n}\n\nfragment EntityDocumentsTable_hyperlinks on Hyperlink {\n  ...DocumentTable_hyperlinks\n}\n\nfragment EquipmentTable_equipment on Equipment {\n  id\n  name\n  futureState\n  equipmentType {\n    id\n    name\n  }\n  workOrder {\n    id\n    status\n  }\n  device {\n    up\n  }\n  services {\n    id\n  }\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...DocumentMenu_document\n  ...ImageDialog_img\n}\n\nfragment HyperlinkTableRow_hyperlink on Hyperlink {\n  id\n  category\n  url\n  displayName\n  createTime\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n\nfragment LocationBreadcrumbsTitle_locationDetails on Location {\n  id\n  name\n  locationType {\n    name\n    id\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n}\n\nfragment LocationDocumentsCard_location on Location {\n  id\n  images {\n    ...EntityDocumentsTable_files\n    id\n  }\n  files {\n    ...EntityDocumentsTable_files\n    id\n  }\n  hyperlinks {\n    ...EntityDocumentsTable_hyperlinks\n    id\n  }\n}\n\nfragment LocationFloorPlansTab_location on Location {\n  id\n  floorPlans {\n    id\n    name\n    image {\n      ...FileAttachment_file\n      id\n    }\n  }\n}\n\nfragment LocationMoreActionsButton_location on Location {\n  id\n  parentLocation {\n    id\n  }\n  children {\n    id\n  }\n  equipments {\n    id\n  }\n  images {\n    id\n  }\n  files {\n    id\n  }\n  surveys {\n    id\n  }\n}\n\nfragment LocationSiteSurveyTab_location on Location {\n  id\n  siteSurveyNeeded\n  surveys {\n    id\n    completionTimestamp\n    name\n    ownerName\n    sourceFile {\n      id\n      fileName\n      storeKey\n    }\n    ...SiteSurveyPane_survey\n  }\n}\n\nfragment PropertyFormField_property on Property {\n  id\n  propertyType {\n    id\n    name\n    type\n    isEditable\n    isMandatory\n    isInstanceProperty\n    stringValue\n  }\n  stringValue\n  intValue\n  floatValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  equipmentValue {\n    id\n    name\n  }\n  locationValue {\n    id\n    name\n  }\n  serviceValue {\n    id\n    name\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n}\n\nfragment SiteSurveyPane_survey on Survey {\n  name\n  completionTimestamp\n  surveyResponses {\n    id\n    questionText\n    formName\n    formIndex\n    questionIndex\n    ...SiteSurveyQuestionReply_question\n  }\n}\n\nfragment SiteSurveyQuestionReplyCellData_data on SurveyQuestion {\n  cellData {\n    networkType\n    signalStrength\n    baseStationID\n    cellID\n    locationAreaCode\n    mobileCountryCode\n    mobileNetworkCode\n    id\n  }\n}\n\nfragment SiteSurveyQuestionReplyWifiData_data on SurveyQuestion {\n  wifiData {\n    band\n    bssid\n    channel\n    frequency\n    strength\n    ssid\n    id\n  }\n}\n\nfragment SiteSurveyQuestionReply_question on SurveyQuestion {\n  questionFormat\n  longitude\n  latitude\n  boolData\n  textData\n  emailData\n  phoneData\n  floatData\n  intData\n  dateData\n  photoData {\n    storeKey\n    id\n  }\n  ...SiteSurveyQuestionReplyWifiData_data\n  ...SiteSurveyQuestionReplyCellData_data\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8d74fb9c2960bf14d0a475e53878a4b0';
module.exports = node;
