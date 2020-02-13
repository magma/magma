/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1a3e933bcb237418193525a8665479d1
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type EquipmentFilterType = "EQUIPMENT_TYPE" | "EQUIP_INST_NAME" | "LOCATION_INST" | "PROPERTY" | "%future added value";
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type EquipmentFilterInput = {|
  filterType: EquipmentFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  name: string,
  type: PropertyKind,
  index?: ?number,
  category?: ?string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};
export type EquipmentTypeahead_equipmentQueryVariables = {|
  filters: $ReadOnlyArray<EquipmentFilterInput>
|};
export type EquipmentTypeahead_equipmentQueryResponse = {|
  +equipmentSearch: {|
    +equipment: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
    |}>
  |}
|};
export type EquipmentTypeahead_equipmentQuery = {|
  variables: EquipmentTypeahead_equipmentQueryVariables,
  response: EquipmentTypeahead_equipmentQueryResponse,
|};
*/


/*
query EquipmentTypeahead_equipmentQuery(
  $filters: [EquipmentFilterInput!]!
) {
  equipmentSearch(limit: 10, filters: $filters) {
    equipment {
      id
      name
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "filters",
    "type": "[EquipmentFilterInput!]!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "equipmentSearch",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filters",
        "variableName": "filters"
      },
      {
        "kind": "Literal",
        "name": "limit",
        "value": 10
      }
    ],
    "concreteType": "EquipmentSearchResult",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipment",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": true,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "id",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "name",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EquipmentTypeahead_equipmentQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EquipmentTypeahead_equipmentQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentTypeahead_equipmentQuery",
    "id": null,
    "text": "query EquipmentTypeahead_equipmentQuery(\n  $filters: [EquipmentFilterInput!]!\n) {\n  equipmentSearch(limit: 10, filters: $filters) {\n    equipment {\n      id\n      name\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '432c2659bc5f9c4f3650a27a9e6690d6';
module.exports = node;
