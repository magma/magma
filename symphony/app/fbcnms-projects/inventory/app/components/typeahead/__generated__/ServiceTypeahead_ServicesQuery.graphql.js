/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash b17ed1a6bb6f3102fa19753cffa502b0
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ServiceFilterType = "EQUIPMENT_IN_SERVICE" | "LOCATION_INST" | "PROPERTY" | "SERVICE_INST_CUSTOMER_NAME" | "SERVICE_INST_EXTERNAL_ID" | "SERVICE_INST_NAME" | "SERVICE_STATUS" | "SERVICE_TYPE" | "%future added value";
export type ServiceFilterInput = {|
  filterType: ServiceFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
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
export type ServiceTypeahead_ServicesQueryVariables = {|
  filters: $ReadOnlyArray<ServiceFilterInput>,
  limit?: ?number,
|};
export type ServiceTypeahead_ServicesQueryResponse = {|
  +serviceSearch: {|
    +services: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
    |}>
  |}
|};
export type ServiceTypeahead_ServicesQuery = {|
  variables: ServiceTypeahead_ServicesQueryVariables,
  response: ServiceTypeahead_ServicesQueryResponse,
|};
*/


/*
query ServiceTypeahead_ServicesQuery(
  $filters: [ServiceFilterInput!]!
  $limit: Int
) {
  serviceSearch(filters: $filters, limit: $limit) {
    services {
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
    "type": "[ServiceFilterInput!]!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "limit",
    "type": "Int",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "serviceSearch",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filters",
        "variableName": "filters"
      },
      {
        "kind": "Variable",
        "name": "limit",
        "variableName": "limit"
      }
    ],
    "concreteType": "ServiceSearchResult",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "services",
        "storageKey": null,
        "args": null,
        "concreteType": "Service",
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
    "name": "ServiceTypeahead_ServicesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceTypeahead_ServicesQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceTypeahead_ServicesQuery",
    "id": null,
    "text": "query ServiceTypeahead_ServicesQuery(\n  $filters: [ServiceFilterInput!]!\n  $limit: Int\n) {\n  serviceSearch(filters: $filters, limit: $limit) {\n    services {\n      id\n      name\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e7811185c06ef9a240fecd3d4fb3dd00';
module.exports = node;
