/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4275cd5876fcf908f1cdee2ff5181e99
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterOperator = "CONTAINS" | "DATE_GREATER_OR_EQUAL_THAN" | "DATE_GREATER_THAN" | "DATE_LESS_OR_EQUAL_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type ServiceFilterType = "EQUIPMENT_IN_SERVICE" | "LOCATION_INST" | "LOCATION_INST_EXTERNAL_ID" | "PROPERTY" | "SERVICE_DISCOVERY_METHOD" | "SERVICE_INST_CUSTOMER_NAME" | "SERVICE_INST_EXTERNAL_ID" | "SERVICE_INST_NAME" | "SERVICE_STATUS" | "SERVICE_TYPE" | "%future added value";
export type ServiceFilterInput = {|
  filterType: ServiceFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  externalId?: ?string,
  name: string,
  type: PropertyKind,
  nodeType?: ?string,
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
  +services: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
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
  services(filterBy: $filters, first: $limit) {
    edges {
      node {
        id
        name
      }
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
    "name": "services",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "filterBy",
        "variableName": "filters"
      },
      {
        "kind": "Variable",
        "name": "first",
        "variableName": "limit"
      }
    ],
    "concreteType": "ServiceConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": false,
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
    "text": "query ServiceTypeahead_ServicesQuery(\n  $filters: [ServiceFilterInput!]!\n  $limit: Int\n) {\n  services(filterBy: $filters, first: $limit) {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '76ab89e627f48adb6ac1601ad66f151f';
module.exports = node;
