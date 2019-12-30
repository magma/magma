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
type CheckListItem_item$ref = any;
type CommentsBox_comments$ref = any;
type EntityDocumentsTable_files$ref = any;
type LocationBreadcrumbsTitle_locationDetails$ref = any;
type WorkOrderDetailsPane_workOrder$ref = any;
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
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
  +ownerName: string,
  +assignee: ?string,
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
      +isEditable: ?boolean,
      +isMandatory: ?boolean,
      +isInstanceProperty: ?boolean,
      +stringValue: ?string,
    |},
    +stringValue: ?string,
    +intValue: ?number,
    +floatValue: ?number,
    +booleanValue: ?boolean,
    +latitudeValue: ?number,
    +longitudeValue: ?number,
    +rangeFromValue: ?number,
    +rangeToValue: ?number,
    +equipmentValue: ?{|
      +id: string,
      +name: string,
    |},
    +locationValue: ?{|
      +id: string,
      +name: string,
    |},
    +serviceValue: ?{|
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
  +comments: $ReadOnlyArray<?{|
    +$fragmentRefs: CommentsBox_comments$ref
  |}>,
  +project: ?{|
    +name: string,
    +id: string,
    +type: {|
      +id: string,
      +name: string,
    |},
  |},
  +checkList: $ReadOnlyArray<?{|
    +id: string,
    +index: ?number,
    +type: CheckListItemType,
    +title: string,
    +checked: ?boolean,
    +$fragmentRefs: CheckListItem_item$ref,
  |}>,
  +$fragmentRefs: WorkOrderDetailsPane_workOrder$ref,
  +$refType: WorkOrderDetails_workOrder$ref,
|};
export type WorkOrderDetails_workOrder$data = WorkOrderDetails_workOrder;
export type WorkOrderDetails_workOrder$key = {
  +$data?: WorkOrderDetails_workOrder$data,
  +$fragmentRefs: WorkOrderDetails_workOrder$ref,
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
  "name": "type",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v4 = [
  (v0/*: any*/),
  (v1/*: any*/)
],
v5 = [
  {
    "kind": "FragmentSpread",
    "name": "EntityDocumentsTable_files",
    "args": null
  }
];
return {
  "kind": "Fragment",
  "name": "WorkOrderDetails_workOrder",
  "type": "WorkOrder",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
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
      "kind": "ScalarField",
      "alias": null,
      "name": "ownerName",
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
            (v2/*: any*/),
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
            (v3/*: any*/)
          ]
        },
        (v3/*: any*/),
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
          "name": "floatValue",
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
          "kind": "LinkedField",
          "alias": null,
          "name": "equipmentValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Equipment",
          "plural": false,
          "selections": (v4/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "locationValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Location",
          "plural": false,
          "selections": (v4/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "serviceValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Service",
          "plural": false,
          "selections": (v4/*: any*/)
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
      "selections": (v5/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "files",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v5/*: any*/)
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
          "name": "CommentsBox_comments",
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
          "selections": (v4/*: any*/)
        }
      ]
    },
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
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "index",
          "args": null,
          "storageKey": null
        },
        (v2/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "title",
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
          "kind": "FragmentSpread",
          "name": "CheckListItem_item",
          "args": null
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
(node/*: any*/).hash = '6d5b85a397a803ac802745f16cdeeee2';
module.exports = node;
