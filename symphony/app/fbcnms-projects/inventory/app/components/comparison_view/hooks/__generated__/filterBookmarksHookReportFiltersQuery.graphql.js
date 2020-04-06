/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash d8eee17f56d5ae9a296160ec06a59e95
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterEntity = "EQUIPMENT" | "LINK" | "LOCATION" | "PORT" | "SERVICE" | "WORK_ORDER" | "%future added value";
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type filterBookmarksHookReportFiltersQueryVariables = {|
  entity: FilterEntity
|};
export type filterBookmarksHookReportFiltersQueryResponse = {|
  +reportFilters: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +entity: FilterEntity,
    +filters: $ReadOnlyArray<{|
      +operator: FilterOperator,
      +key: string,
      +filterType: string,
      +stringValue: ?string,
      +idSet: ?$ReadOnlyArray<string>,
      +stringSet: ?$ReadOnlyArray<string>,
      +boolValue: ?boolean,
      +propertyValue: ?{|
        +index: ?number,
        +name: string,
        +type: PropertyKind,
        +stringValue: ?string,
        +intValue: ?number,
        +booleanValue: ?boolean,
        +floatValue: ?number,
        +latitudeValue: ?number,
        +longitudeValue: ?number,
        +rangeFromValue: ?number,
        +rangeToValue: ?number,
        +isDeleted: ?boolean,
      |},
    |}>,
  |}>
|};
export type filterBookmarksHookReportFiltersQuery = {|
  variables: filterBookmarksHookReportFiltersQueryVariables,
  response: filterBookmarksHookReportFiltersQueryResponse,
|};
*/


/*
query filterBookmarksHookReportFiltersQuery(
  $entity: FilterEntity!
) {
  reportFilters(entity: $entity) {
    id
    name
    entity
    filters {
      operator
      key
      filterType
      stringValue
      idSet
      stringSet
      boolValue
      propertyValue {
        index
        name
        type
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        isDeleted
        id
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "entity",
    "type": "FilterEntity!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "entity",
    "variableName": "entity"
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
  "name": "entity",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "operator",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "key",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "filterType",
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
  "name": "idSet",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringSet",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "boolValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isDeleted",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "filterBookmarksHookReportFiltersQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "reportFilters",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ReportFilter",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "filters",
            "storageKey": null,
            "args": null,
            "concreteType": "GeneralFilter",
            "plural": true,
            "selections": [
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "propertyValue",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": false,
                "selections": [
                  (v12/*: any*/),
                  (v3/*: any*/),
                  (v13/*: any*/),
                  (v8/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  (v20/*: any*/),
                  (v21/*: any*/)
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "filterBookmarksHookReportFiltersQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "reportFilters",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ReportFilter",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "filters",
            "storageKey": null,
            "args": null,
            "concreteType": "GeneralFilter",
            "plural": true,
            "selections": [
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "propertyValue",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": false,
                "selections": [
                  (v12/*: any*/),
                  (v3/*: any*/),
                  (v13/*: any*/),
                  (v8/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  (v20/*: any*/),
                  (v21/*: any*/),
                  (v2/*: any*/)
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
    "name": "filterBookmarksHookReportFiltersQuery",
    "id": null,
    "text": "query filterBookmarksHookReportFiltersQuery(\n  $entity: FilterEntity!\n) {\n  reportFilters(entity: $entity) {\n    id\n    name\n    entity\n    filters {\n      operator\n      key\n      filterType\n      stringValue\n      idSet\n      stringSet\n      boolValue\n      propertyValue {\n        index\n        name\n        type\n        stringValue\n        intValue\n        booleanValue\n        floatValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        isDeleted\n        id\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '19ea160f56380b4a7d93b2b0df38b959';
module.exports = node;
