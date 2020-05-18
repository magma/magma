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
export type CheckListItemEnumSelectionMode = "multiple" | "single" | "%future added value";
export type CheckListItemType = "cell_scan" | "enum" | "files" | "simple" | "string" | "wifi_scan" | "yes_no" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type AddEditWorkOrderTypeCard_workOrderType$ref: FragmentReference;
declare export opaque type AddEditWorkOrderTypeCard_workOrderType$fragmentType: AddEditWorkOrderTypeCard_workOrderType$ref;
export type AddEditWorkOrderTypeCard_workOrderType = {|
  +id: string,
  +name: string,
  +description: ?string,
  +numberOfWorkOrders: number,
  +propertyTypes: $ReadOnlyArray<?{|
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
    +isMandatory: ?boolean,
    +isInstanceProperty: ?boolean,
    +isDeleted: ?boolean,
    +category: ?string,
  |}>,
  +checkListCategoryDefinitions: $ReadOnlyArray<{|
    +id: string,
    +title: string,
    +description: ?string,
    +checklistItemDefinitions: $ReadOnlyArray<{|
      +id: string,
      +title: string,
      +type: CheckListItemType,
      +index: ?number,
      +enumValues: ?string,
      +enumSelectionMode: ?CheckListItemEnumSelectionMode,
      +helpText: ?string,
    |}>,
  |}>,
  +$refType: AddEditWorkOrderTypeCard_workOrderType$ref,
|};
export type AddEditWorkOrderTypeCard_workOrderType$data = AddEditWorkOrderTypeCard_workOrderType;
export type AddEditWorkOrderTypeCard_workOrderType$key = {
  +$data?: AddEditWorkOrderTypeCard_workOrderType$data,
  +$fragmentRefs: AddEditWorkOrderTypeCard_workOrderType$ref,
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
  "name": "type",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "title",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "AddEditWorkOrderTypeCard_workOrderType",
  "type": "WorkOrderType",
  "metadata": null,
  "argumentDefinitions": [],
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
        (v3/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "nodeType",
          "args": null,
          "storageKey": null
        },
        (v4/*: any*/),
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
        (v5/*: any*/),
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
            (v5/*: any*/),
            (v3/*: any*/),
            (v4/*: any*/),
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
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd41a35b8680d6ea20ea8a79786fba884';
module.exports = node;
