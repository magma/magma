/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4743a0f85e8e11a48147322daeabff79
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type EquipmentAddEditCardQueryVariables = {|
  equipmentId: string
|};
export type EquipmentAddEditCardQueryResponse = {|
  +equipment: ?{|
    +id?: string,
    +name?: string,
    +parentLocation?: ?{|
      +id: string
    |},
    +parentPosition?: ?{|
      +id: string
    |},
    +device?: ?{|
      +id: string
    |},
    +equipmentType?: {|
      +id: string,
      +name: string,
      +propertyTypes: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
        +index: ?number,
        +isInstanceProperty: ?boolean,
        +type: PropertyKind,
        +isMandatory: ?boolean,
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
    +properties?: $ReadOnlyArray<?{|
      +propertyType: {|
        +id: string,
        +name: string,
        +index: ?number,
        +isInstanceProperty: ?boolean,
        +type: PropertyKind,
        +stringValue: ?string,
        +isMandatory: ?boolean,
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
      +serviceValue: ?{|
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
  equipment: node(id: $equipmentId) {
    __typename
    ... on Equipment {
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
          isMandatory
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
          isMandatory
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
        serviceValue {
          id
          name
        }
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
    "name": "equipmentId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "equipmentId"
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
v4 = [
  (v2/*: any*/)
],
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "parentLocation",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": false,
  "selections": (v4/*: any*/)
},
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "parentPosition",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPosition",
  "plural": false,
  "selections": (v4/*: any*/)
},
v7 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "device",
  "storageKey": null,
  "args": null,
  "concreteType": "Device",
  "plural": false,
  "selections": (v4/*: any*/)
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
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
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
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
        (v8/*: any*/),
        (v9/*: any*/),
        (v10/*: any*/),
        (v11/*: any*/),
        (v12/*: any*/),
        (v13/*: any*/),
        (v14/*: any*/),
        (v15/*: any*/),
        (v16/*: any*/),
        (v17/*: any*/),
        (v18/*: any*/),
        (v19/*: any*/)
      ]
    }
  ]
},
v21 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v22 = {
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
        (v2/*: any*/),
        (v3/*: any*/),
        (v8/*: any*/),
        (v9/*: any*/),
        (v10/*: any*/),
        (v12/*: any*/),
        (v11/*: any*/)
      ]
    },
    (v2/*: any*/),
    (v12/*: any*/),
    (v13/*: any*/),
    (v14/*: any*/),
    (v15/*: any*/),
    (v16/*: any*/),
    (v17/*: any*/),
    (v18/*: any*/),
    (v19/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v21/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v21/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": false,
      "selections": (v21/*: any*/)
    }
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EquipmentAddEditCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipment",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Equipment",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v20/*: any*/),
              (v22/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EquipmentAddEditCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipment",
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
            "type": "Equipment",
            "selections": [
              (v3/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v20/*: any*/),
              (v22/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentAddEditCardQuery",
    "id": null,
    "text": "query EquipmentAddEditCardQuery(\n  $equipmentId: ID!\n) {\n  equipment: node(id: $equipmentId) {\n    __typename\n    ... on Equipment {\n      id\n      name\n      parentLocation {\n        id\n      }\n      parentPosition {\n        id\n      }\n      device {\n        id\n      }\n      equipmentType {\n        id\n        name\n        propertyTypes {\n          id\n          name\n          index\n          isInstanceProperty\n          type\n          isMandatory\n          stringValue\n          intValue\n          floatValue\n          booleanValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n        }\n      }\n      properties {\n        propertyType {\n          id\n          name\n          index\n          isInstanceProperty\n          type\n          stringValue\n          isMandatory\n        }\n        id\n        stringValue\n        intValue\n        floatValue\n        booleanValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        equipmentValue {\n          id\n          name\n        }\n        locationValue {\n          id\n          name\n        }\n        serviceValue {\n          id\n          name\n        }\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9c48efc187ea06723b66c965ab24c6dc';
module.exports = node;
