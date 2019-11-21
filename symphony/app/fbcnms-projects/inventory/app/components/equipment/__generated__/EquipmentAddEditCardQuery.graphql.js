/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash bcaa045594122c9c92846094b06927da
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "string" | "%future added value";
export type EquipmentAddEditCardQueryVariables = {|
  equipmentId: string
|};
export type EquipmentAddEditCardQueryResponse = {|
  +equipment: ?{|
    +id: string,
    +name: string,
    +parentLocation: ?{|
      +id: string
    |},
    +parentPosition: ?{|
      +id: string
    |},
    +device: ?{|
      +id: string
    |},
    +equipmentType: {|
      +id: string,
      +name: string,
      +propertyTypes: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
        +index: ?number,
        +isInstanceProperty: ?boolean,
        +type: PropertyKind,
        +stringValue: ?string,
        +intValue: ?number,
        +floatValue: ?number,
        +booleanValue: ?boolean,
        +latitudeValue: ?number,
        +longitudeValue: ?number,
        +rangeFromValue: ?number,
        +rangeToValue: ?number,
      |}>,
    |},
    +properties: $ReadOnlyArray<?{|
      +propertyType: {|
        +id: string,
        +name: string,
        +index: ?number,
        +isInstanceProperty: ?boolean,
        +type: PropertyKind,
        +stringValue: ?string,
      |},
      +id: string,
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
    |}>,
  |}
|};
export type EquipmentAddEditCardQuery = {|
  variables: EquipmentAddEditCardQueryVariables,
  response: EquipmentAddEditCardQueryResponse,
|};
*/


/*
query EquipmentAddEditCardQuery(
  $equipmentId: ID!
) {
  equipment(id: $equipmentId) {
    id
    name
    parentLocation {
      id
    }
    parentPosition {
      id
    }
    device {
      id
    }
    equipmentType {
      id
      name
      propertyTypes {
        id
        name
        index
        isInstanceProperty
        type
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
    properties {
      propertyType {
        id
        name
        index
        isInstanceProperty
        type
        stringValue
      }
      id
      stringValue
      intValue
      floatValue
      booleanValue
      latitudeValue
      longitudeValue
      rangeFromValue
      rangeToValue
      equipmentValue {
        id
        name
      }
      locationValue {
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
    "name": "equipmentId",
    "type": "ID!",
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
v3 = [
  (v1/*: any*/)
],
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
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v15 = [
  (v1/*: any*/),
  (v2/*: any*/)
],
v16 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "equipment",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "equipmentId"
      }
    ],
    "concreteType": "Equipment",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentLocation",
        "storageKey": null,
        "args": null,
        "concreteType": "Location",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "parentPosition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPosition",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "device",
        "storageKey": null,
        "args": null,
        "concreteType": "Device",
        "plural": false,
        "selections": (v3/*: any*/)
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
          (v1/*: any*/),
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "propertyTypes",
            "storageKey": null,
            "args": null,
            "concreteType": "PropertyType",
            "plural": true,
            "selections": [
              (v1/*: any*/),
              (v2/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/),
              (v12/*: any*/),
              (v13/*: any*/),
              (v14/*: any*/)
            ]
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
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "propertyType",
            "storageKey": null,
            "args": null,
            "concreteType": "PropertyType",
            "plural": false,
            "selections": [
              (v1/*: any*/),
              (v2/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/)
            ]
          },
          (v1/*: any*/),
          (v7/*: any*/),
          (v8/*: any*/),
          (v9/*: any*/),
          (v10/*: any*/),
          (v11/*: any*/),
          (v12/*: any*/),
          (v13/*: any*/),
          (v14/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentValue",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": false,
            "selections": (v15/*: any*/)
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "locationValue",
            "storageKey": null,
            "args": null,
            "concreteType": "Location",
            "plural": false,
            "selections": (v15/*: any*/)
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
    "name": "EquipmentAddEditCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v16/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EquipmentAddEditCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v16/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentAddEditCardQuery",
    "id": null,
    "text": "query EquipmentAddEditCardQuery(\n  $equipmentId: ID!\n) {\n  equipment(id: $equipmentId) {\n    id\n    name\n    parentLocation {\n      id\n    }\n    parentPosition {\n      id\n    }\n    device {\n      id\n    }\n    equipmentType {\n      id\n      name\n      propertyTypes {\n        id\n        name\n        index\n        isInstanceProperty\n        type\n        stringValue\n        intValue\n        floatValue\n        booleanValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n      }\n    }\n    properties {\n      propertyType {\n        id\n        name\n        index\n        isInstanceProperty\n        type\n        stringValue\n      }\n      id\n      stringValue\n      intValue\n      floatValue\n      booleanValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n      equipmentValue {\n        id\n        name\n      }\n      locationValue {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'dfd1f534e95c18ec87062b42e4c129b4';
module.exports = node;
