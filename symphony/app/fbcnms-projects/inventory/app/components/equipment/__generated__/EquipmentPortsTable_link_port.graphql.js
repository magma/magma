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
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPortsTable_link_port$ref: FragmentReference;
declare export opaque type EquipmentPortsTable_link_port$fragmentType: EquipmentPortsTable_link_port$ref;
export type EquipmentPortsTable_link_port = {|
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
    +definition: {|
      +role: ?string
    |},
    +service: {|
      +name: string
    |},
  |}>,
  +$refType: EquipmentPortsTable_link_port$ref,
|};
export type EquipmentPortsTable_link_port$data = EquipmentPortsTable_link_port;
export type EquipmentPortsTable_link_port$key = {
  +$data?: EquipmentPortsTable_link_port$data,
  +$fragmentRefs: EquipmentPortsTable_link_port$ref,
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
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "EquipmentPortsTable_link_port",
  "type": "EquipmentPort",
  "metadata": null,
  "argumentDefinitions": [],
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
        (v1/*: any*/),
        (v2/*: any*/),
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
              "selections": [
                (v0/*: any*/),
                (v1/*: any*/),
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
              ]
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
        (v1/*: any*/),
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
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
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
                (v1/*: any*/),
                (v2/*: any*/),
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
                  "selections": [
                    (v0/*: any*/),
                    (v1/*: any*/)
                  ]
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
            }
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
          "selections": [
            (v1/*: any*/)
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a52288c25e58426be38cca7b0a1d61ce';
module.exports = node;
