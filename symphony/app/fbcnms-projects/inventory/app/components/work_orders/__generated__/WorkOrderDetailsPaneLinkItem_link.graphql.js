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
type EquipmentBreadcrumbs_equipment$ref = any;
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ServiceEndpointRole = "CONSUMER" | "PROVIDER" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type WorkOrderDetailsPaneLinkItem_link$ref: FragmentReference;
declare export opaque type WorkOrderDetailsPaneLinkItem_link$fragmentType: WorkOrderDetailsPaneLinkItem_link$ref;
export type WorkOrderDetailsPaneLinkItem_link = {|
  +id: string,
  +futureState: ?FutureState,
  +ports: $ReadOnlyArray<?{|
    +id: string,
    +definition: {|
      +id: string,
      +name: string,
      +visibleLabel: ?string,
      +portType: ?{|
        +linkPropertyTypes: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +type: PropertyKind,
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
        |}>
      |},
    |},
    +parentEquipment: {|
      +id: string,
      +name: string,
      +futureState: ?FutureState,
      +equipmentType: {|
        +id: string,
        +name: string,
        +portDefinitions: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +visibleLabel: ?string,
          +bandwidth: ?string,
          +portType: ?{|
            +id: string,
            +name: string,
          |},
        |}>,
      |},
      +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
    |},
    +serviceEndpoints: $ReadOnlyArray<{|
      +role: ServiceEndpointRole,
      +service: {|
        +name: string
      |},
    |}>,
  |}>,
  +workOrder: ?{|
    +id: string,
    +status: WorkOrderStatus,
  |},
  +properties: $ReadOnlyArray<?{|
    +id: string,
    +propertyType: {|
      +id: string,
      +name: string,
      +type: PropertyKind,
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
  +services: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
  |}>,
  +$refType: WorkOrderDetailsPaneLinkItem_link$ref,
|};
export type WorkOrderDetailsPaneLinkItem_link$data = WorkOrderDetailsPaneLinkItem_link;
export type WorkOrderDetailsPaneLinkItem_link$key = {
  +$data?: WorkOrderDetailsPaneLinkItem_link$data,
  +$fragmentRefs: WorkOrderDetailsPaneLinkItem_link$ref,
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
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v12 = [
  (v0/*: any*/),
  (v2/*: any*/),
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
  (v4/*: any*/),
  (v5/*: any*/),
  (v6/*: any*/),
  (v7/*: any*/),
  (v8/*: any*/),
  (v9/*: any*/),
  (v10/*: any*/),
  (v11/*: any*/),
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
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "category",
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
],
v13 = [
  (v0/*: any*/),
  (v2/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "WorkOrderDetailsPaneLinkItem_link",
  "type": "Link",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "ports",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPort",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "definition",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v2/*: any*/),
            (v3/*: any*/),
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
                  "selections": (v12/*: any*/)
                }
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
            (v0/*: any*/),
            (v2/*: any*/),
            (v1/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "equipmentType",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentType",
              "plural": false,
              "selections": [
                (v0/*: any*/),
                (v2/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "portDefinitions",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentPortDefinition",
                  "plural": true,
                  "selections": [
                    (v0/*: any*/),
                    (v2/*: any*/),
                    (v3/*: any*/),
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
            {
              "kind": "FragmentSpread",
              "name": "EquipmentBreadcrumbs_equipment",
              "args": null
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
              "kind": "ScalarField",
              "alias": null,
              "name": "role",
              "args": null,
              "storageKey": null
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "service",
              "storageKey": null,
              "args": null,
              "concreteType": "Service",
              "plural": false,
              "selections": [
                (v2/*: any*/)
              ]
            }
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
        (v0/*: any*/),
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
          "selections": (v12/*: any*/)
        },
        (v4/*: any*/),
        (v5/*: any*/),
        (v7/*: any*/),
        (v6/*: any*/),
        (v8/*: any*/),
        (v9/*: any*/),
        (v10/*: any*/),
        (v11/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "equipmentValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Equipment",
          "plural": false,
          "selections": (v13/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "locationValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Location",
          "plural": false,
          "selections": (v13/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "serviceValue",
          "storageKey": null,
          "args": null,
          "concreteType": "Service",
          "plural": false,
          "selections": (v13/*: any*/)
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
      "selections": (v13/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '2b79255bca55026291ae3b75ae46f9f3';
module.exports = node;
