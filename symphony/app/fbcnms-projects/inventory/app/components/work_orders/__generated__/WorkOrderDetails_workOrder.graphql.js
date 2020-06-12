/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
type CommentsActivitiesBox_activities$ref = any;
type CommentsActivitiesBox_comments$ref = any;
type EntityDocumentsTable_files$ref = any;
type EntityDocumentsTable_hyperlinks$ref = any;
type LocationBreadcrumbsTitle_locationDetails$ref = any;
type WorkOrderDetailsPane_workOrder$ref = any;
export type CellularNetworkType = "CDMA" | "GSM" | "LTE" | "WCDMA" | "%future added value";
export type CheckListItemEnumSelectionMode = "multiple" | "single" | "%future added value";
export type CheckListItemType = "cell_scan" | "enum" | "files" | "simple" | "string" | "wifi_scan" | "yes_no" | "%future added value";
export type FileType = "FILE" | "IMAGE" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type YesNoResponse = "NO" | "YES" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type WorkOrderDetails_workOrder$ref: FragmentReference;
declare export opaque type WorkOrderDetails_workOrder$fragmentType: WorkOrderDetails_workOrder$ref;
export type WorkOrderDetails_workOrder = {|
  +id: string,
  +name: string,
  +description: ?string,
  +workOrderType: {|
    +name: string,
    +id: string,
  |},
  +location: ?{|
    +name: string,
    +id: string,
    +latitude: number,
    +longitude: number,
    +locationType: {|
      +mapType: ?string,
      +mapZoomLevel: ?number,
    |},
    +$fragmentRefs: LocationBreadcrumbsTitle_locationDetails$ref,
  |},
  +owner: {|
    +id: string,
    +email: string,
  |},
  +assignedTo: ?{|
    +id: string,
    +email: string,
  |},
  +creationDate: any,
  +installDate: ?any,
  +status: WorkOrderStatus,
  +priority: WorkOrderPriority,
  +properties: $ReadOnlyArray<?{|
    +id: string,
    +propertyType: {|
      +id: string,
      +name: string,
      +type: PropertyKind,
      +nodeType: ?string,
      +index: ?number,
      +stringValue: ?string,
      +intValue: ?number,
      +booleanValue: ?boolean,
      +floatValue: ?number,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +isEditable: ?boolean,
      +isInstanceProperty: ?boolean,
      +isMandatory: ?boolean,
      +category: ?string,
      +isDeleted: ?boolean,
    |},
    +stringValue: ?string,
    +intValue: ?number,
    +floatValue: ?number,
    +booleanValue: ?boolean,
    +latitudeValue: ?number,
    +longitudeValue: ?number,
    +rangeFromValue: ?number,
    +rangeToValue: ?number,
    +nodeValue: ?{|
      +id: string,
      +name: string,
    |},
  |}>,
  +images: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +files: $ReadOnlyArray<?{|
    +$fragmentRefs: EntityDocumentsTable_files$ref
  |}>,
  +hyperlinks: $ReadOnlyArray<{|
    +$fragmentRefs: EntityDocumentsTable_hyperlinks$ref
  |}>,
  +comments: $ReadOnlyArray<?{|
    +$fragmentRefs: CommentsActivitiesBox_comments$ref
  |}>,
  +activities: $ReadOnlyArray<{|
    +$fragmentRefs: CommentsActivitiesBox_activities$ref
  |}>,
  +project: ?{|
    +name: string,
    +id: string,
    +type: {|
      +id: string,
      +name: string,
    |},
  |},
  +checkListCategories: $ReadOnlyArray<{|
    +id: string,
    +title: string,
    +description: ?string,
    +checkList: $ReadOnlyArray<{|
      +id: string,
      +index: ?number,
      +type: CheckListItemType,
      +title: string,
      +helpText: ?string,
      +checked: ?boolean,
      +enumValues: ?string,
      +stringValue: ?string,
      +enumSelectionMode: ?CheckListItemEnumSelectionMode,
      +selectedEnumValues: ?string,
      +yesNoResponse: ?YesNoResponse,
      +files: ?$ReadOnlyArray<{|
        +id: string,
        +fileName: string,
        +sizeInBytes: ?number,
        +modified: ?any,
        +uploaded: ?any,
        +fileType: ?FileType,
        +storeKey: ?string,
        +category: ?string,
      |}>,
      +cellData: ?$ReadOnlyArray<{|
        +id: string,
        +networkType: CellularNetworkType,
        +signalStrength: number,
        +timestamp: ?number,
        +baseStationID: ?string,
        +networkID: ?string,
        +systemID: ?string,
        +cellID: ?string,
        +locationAreaCode: ?string,
        +mobileCountryCode: ?string,
        +mobileNetworkCode: ?string,
        +primaryScramblingCode: ?string,
        +operator: ?string,
        +arfcn: ?number,
        +physicalCellID: ?string,
        +trackingAreaCode: ?string,
        +timingAdvance: ?number,
        +earfcn: ?number,
        +uarfcn: ?number,
        +latitude: ?number,
        +longitude: ?number,
      |}>,
      +wifiData: ?$ReadOnlyArray<{|
        +id: string,
        +timestamp: number,
        +frequency: number,
        +channel: number,
        +bssid: string,
        +strength: number,
        +ssid: ?string,
        +band: ?string,
        +channelWidth: ?number,
        +capabilities: ?string,
        +latitude: ?number,
        +longitude: ?number,
      |}>,
    |}>,
  |}>,
  +$fragmentRefs: WorkOrderDetailsPane_workOrder$ref,
  +$refType: WorkOrderDetails_workOrder$ref,
|};
export type WorkOrderDetails_workOrder$data = WorkOrderDetails_workOrder;
export type WorkOrderDetails_workOrder$key = {
  +$data?: WorkOrderDetails_workOrder$data,
  +$fragmentRefs: WorkOrderDetails_workOrder$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
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
  "name": "latitude",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitude",
  "args": null,
  "storageKey": null
},
v5 = [
  (v0/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "email",
    "args": null,
    "storageKey": null
  }
],
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
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
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "category",
  "args": null,
  "storageKey": null
},
v17 = [
  (v0/*: any*/),
  (v1/*: any*/)
],
v18 = [
  {
    "kind": "FragmentSpread",
    "name": "EntityDocumentsTable_files",
    "args": null
  }
],
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "timestamp",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "WorkOrderDetails_workOrder",
  "type": "WorkOrder",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "workOrderType",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkOrderType",
      "plural": false,
      "selections": [
        (v1/*: any*/),
        (v0/*: any*/)
      ]
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
        (v1/*: any*/),
        (v0/*: any*/),
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
            }
          ]
        },
        {
          "kind": "FragmentSpread",
          "name": "LocationBreadcrumbsTitle_locationDetails",
          "args": null
        }
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
      "selections": (v5/*: any*/)
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
        (v0/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "propertyType",
          "storageKey": null,
          "args": null,
          "concreteType": "PropertyType",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
            (v6/*: any*/),
            {
              "kind": "ScalarField",
              "alias": null,
              "name": "nodeType",
              "args": null,
              "storageKey": null
            },
            (v7/*: any*/),
            (v8/*: any*/),
            (v9/*: any*/),
            (v10/*: any*/),
            (v11/*: any*/),
            (v12/*: any*/),
            (v13/*: any*/),
            (v14/*: any*/),
            (v15/*: any*/),
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
            (v16/*: any*/),
            {
              "kind": "ScalarField",
              "alias": null,
              "name": "isDeleted",
              "args": null,
              "storageKey": null
            }
          ]
        },
        (v8/*: any*/),
        (v9/*: any*/),
        (v11/*: any*/),
        (v10/*: any*/),
        (v12/*: any*/),
        (v13/*: any*/),
        (v14/*: any*/),
        (v15/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "nodeValue",
          "storageKey": null,
          "args": null,
          "concreteType": null,
          "plural": false,
          "selections": (v17/*: any*/)
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
      "selections": (v18/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "files",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v18/*: any*/)
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
        {
          "kind": "FragmentSpread",
          "name": "EntityDocumentsTable_hyperlinks",
          "args": null
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
        {
          "kind": "FragmentSpread",
          "name": "CommentsActivitiesBox_comments",
          "args": null
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
        {
          "kind": "FragmentSpread",
          "name": "CommentsActivitiesBox_activities",
          "args": null
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
        (v1/*: any*/),
        (v0/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "type",
          "storageKey": null,
          "args": null,
          "concreteType": "ProjectType",
          "plural": false,
          "selections": (v17/*: any*/)
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
        (v0/*: any*/),
        (v19/*: any*/),
        (v2/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "checkList",
          "storageKey": null,
          "args": null,
          "concreteType": "CheckListItem",
          "plural": true,
          "selections": [
            (v0/*: any*/),
            (v7/*: any*/),
            (v6/*: any*/),
            (v19/*: any*/),
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
            (v8/*: any*/),
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
                (v0/*: any*/),
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "fileName",
                  "args": null,
                  "storageKey": null
                },
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "sizeInBytes",
                  "args": null,
                  "storageKey": null
                },
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "modified",
                  "args": null,
                  "storageKey": null
                },
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "uploaded",
                  "args": null,
                  "storageKey": null
                },
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "fileType",
                  "args": null,
                  "storageKey": null
                },
                {
                  "kind": "ScalarField",
                  "alias": null,
                  "name": "storeKey",
                  "args": null,
                  "storageKey": null
                },
                (v16/*: any*/)
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
                (v0/*: any*/),
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
                (v20/*: any*/),
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
                (v3/*: any*/),
                (v4/*: any*/)
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
                (v0/*: any*/),
                (v20/*: any*/),
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
                (v3/*: any*/),
                (v4/*: any*/)
              ]
            }
          ]
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "WorkOrderDetailsPane_workOrder",
      "args": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd11da7b4168f4e7a88e0e63a6de0467c';
module.exports = node;
