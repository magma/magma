/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 605bfb0dce4a13b5c64fcc4e18828c5a
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type FilterEntity = "EQUIPMENT" | "LINK" | "LOCATION" | "PORT" | "SERVICE" | "WORK_ORDER" | "%future added value";
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ReportFilterInput = {|
  name: string,
  entity: FilterEntity,
  filters?: ?$ReadOnlyArray<?GeneralFilterInput>,
|};
export type GeneralFilterInput = {|
  filterType: string,
  operator: FilterOperator,
  stringValue?: ?string,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  boolValue?: ?boolean,
  propertyValue?: ?PropertyTypeInput,
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
export type AddReportFilterMutationVariables = {|
  input: ReportFilterInput
|};
export type AddReportFilterMutationResponse = {|
  +addReportFilter: {|
    +id: string,
    +name: string,
    +entity: FilterEntity,
    +filters: $ReadOnlyArray<{|
      +filterType: string,
      +operator: FilterOperator,
      +stringValue: ?string,
      +idSet: ?$ReadOnlyArray<string>,
      +stringSet: ?$ReadOnlyArray<string>,
      +boolValue: ?boolean,
      +propertyValue: ?{|
        +id: string,
        +name: string,
        +type: PropertyKind,
        +isEditable: ?boolean,
        +isInstanceProperty: ?boolean,
        +stringValue: ?string,
        +intValue: ?number,
        +floatValue: ?number,
        +booleanValue: ?boolean,
        +latitudeValue: ?number,
        +longitudeValue: ?number,
        +rangeFromValue: ?number,
        +rangeToValue: ?number,
      |},
    |}>,
  |}
|};
export type AddReportFilterMutation = {|
  variables: AddReportFilterMutationVariables,
  response: AddReportFilterMutationResponse,
|};
*/


/*
mutation AddReportFilterMutation(
  $input: ReportFilterInput!
) {
  addReportFilter(input: $input) {
    id
    name
    entity
    filters {
      filterType
      operator
      stringValue
      idSet
      stringSet
      boolValue
      propertyValue {
        id
        name
        type
        isEditable
        isInstanceProperty
        stringValue
        intValue
        floatValue
        booleanValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "ReportFilterInput!",
    "defaultValue": null
  }
],
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
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
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v4 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addReportFilter",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "ReportFilter",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      (v2/*: any*/),
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "entity",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "filters",
        "storageKey": null,
        "args": null,
        "concreteType": "GeneralFilter",
        "plural": true,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "filterType",
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
          (v3/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "idSet",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "stringSet",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "boolValue",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "propertyValue",
            "storageKey": null,
            "args": null,
            "concreteType": "PropertyType",
            "plural": false,
            "selections": [
              (v1/*: any*/),
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
    "name": "AddReportFilterMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddReportFilterMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddReportFilterMutation",
    "id": null,
    "text": "mutation AddReportFilterMutation(\n  $input: ReportFilterInput!\n) {\n  addReportFilter(input: $input) {\n    id\n    name\n    entity\n    filters {\n      filterType\n      operator\n      stringValue\n      idSet\n      stringSet\n      boolValue\n      propertyValue {\n        id\n        name\n        type\n        isEditable\n        isInstanceProperty\n        stringValue\n        intValue\n        floatValue\n        booleanValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '990637d691b7e35198b5c445629f8990';
module.exports = node;
