/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 8b80ea99865d1434e55d2d9de7fdf619
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type EquipmentAddEditCardQuery__equipmentTypeQueryVariables = {|
  equipmentTypeId: string
|};
export type EquipmentAddEditCardQuery__equipmentTypeQueryResponse = {|
  +equipmentType: ?{|
    +id?: string,
    +name?: string,
    +propertyTypes?: $ReadOnlyArray<?{|
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
    |}>,
  |}
|};
export type EquipmentAddEditCardQuery__equipmentTypeQuery = {|
  variables: EquipmentAddEditCardQuery__equipmentTypeQueryVariables,
  response: EquipmentAddEditCardQuery__equipmentTypeQueryResponse,
|};
*/


/*
query EquipmentAddEditCardQuery__equipmentTypeQuery(
  $equipmentTypeId: ID!
) {
  equipmentType: node(id: $equipmentTypeId) {
    __typename
    ... on EquipmentType {
      id
      name
      propertyTypes {
        id
        name
        type
        nodeType
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
    "name": "equipmentTypeId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "equipmentTypeId"
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
    }
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EquipmentAddEditCardQuery__equipmentTypeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentType",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "EquipmentType",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EquipmentAddEditCardQuery__equipmentTypeQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentType",
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
            "type": "EquipmentType",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentAddEditCardQuery__equipmentTypeQuery",
    "id": null,
    "text": "query EquipmentAddEditCardQuery__equipmentTypeQuery(\n  $equipmentTypeId: ID!\n) {\n  equipmentType: node(id: $equipmentTypeId) {\n    __typename\n    ... on EquipmentType {\n      id\n      name\n      propertyTypes {\n        id\n        name\n        type\n        nodeType\n        index\n        stringValue\n        intValue\n        booleanValue\n        floatValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        isEditable\n        isMandatory\n        isInstanceProperty\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd0bad0b552a71b56cb8f59b3283de631';
module.exports = node;
