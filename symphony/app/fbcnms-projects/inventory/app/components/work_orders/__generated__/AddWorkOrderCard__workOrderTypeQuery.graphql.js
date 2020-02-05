/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash f18e22f5500c22197744d066c0f8483a
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type AddWorkOrderCard__workOrderTypeQueryVariables = {|
  workOrderTypeId: string
|};
export type AddWorkOrderCard__workOrderTypeQueryResponse = {|
  +workOrderType: ?{|
    +id?: string,
    +name?: string,
    +description?: ?string,
    +propertyTypes?: $ReadOnlyArray<?{|
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
      +isMandatory: ?boolean,
      +isInstanceProperty: ?boolean,
      +isDeleted: ?boolean,
    |}>,
    +checkListDefinitions?: $ReadOnlyArray<?{|
      +id: string,
      +title: string,
      +type: CheckListItemType,
      +index: ?number,
      +helpText: ?string,
      +enumValues: ?string,
    |}>,
  |}
|};
export type AddWorkOrderCard__workOrderTypeQuery = {|
  variables: AddWorkOrderCard__workOrderTypeQueryVariables,
  response: AddWorkOrderCard__workOrderTypeQueryResponse,
|};
*/


/*
query AddWorkOrderCard__workOrderTypeQuery(
  $workOrderTypeId: ID!
) {
  workOrderType: node(id: $workOrderTypeId) {
    __typename
    ... on WorkOrderType {
      id
      name
      description
      propertyTypes {
        id
        name
        type
        index
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        isEditable
        isMandatory
        isInstanceProperty
        isDeleted
      }
      checkListDefinitions {
        id
        title
        type
        index
        helpText
        enumValues
      }
    }
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "workOrderTypeId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "workOrderTypeId"
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
  "name": "description",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "propertyTypes",
  "storageKey": null,
  "args": null,
  "concreteType": "PropertyType",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    (v5/*: any*/),
    (v6/*: any*/),
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
    }
  ]
},
v8 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "checkListDefinitions",
  "storageKey": null,
  "args": null,
  "concreteType": "CheckListItemDefinition",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "title",
      "args": null,
      "storageKey": null
    },
    (v5/*: any*/),
    (v6/*: any*/),
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
      "name": "enumValues",
      "args": null,
      "storageKey": null
    }
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddWorkOrderCard__workOrderTypeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrderType",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "WorkOrderType",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddWorkOrderCard__workOrderTypeQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrderType",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "WorkOrderType",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "AddWorkOrderCard__workOrderTypeQuery",
    "id": null,
    "text": "query AddWorkOrderCard__workOrderTypeQuery(\n  $workOrderTypeId: ID!\n) {\n  workOrderType: node(id: $workOrderTypeId) {\n    __typename\n    ... on WorkOrderType {\n      id\n      name\n      description\n      propertyTypes {\n        id\n        name\n        type\n        index\n        stringValue\n        intValue\n        booleanValue\n        floatValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        isEditable\n        isMandatory\n        isInstanceProperty\n        isDeleted\n      }\n      checkListDefinitions {\n        id\n        title\n        type\n        index\n        helpText\n        enumValues\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9b6728b6413abd3164ccd872c6e16f39';
module.exports = node;
