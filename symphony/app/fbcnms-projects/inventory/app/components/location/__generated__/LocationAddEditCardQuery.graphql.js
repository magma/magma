/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash dbacd2259cca71b92644d3c977105d89
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type LocationAddEditCardQueryVariables = {|
  locationId: string
|};
export type LocationAddEditCardQueryResponse = {|
  +location: ?{|
    +id?: string,
    +name?: string,
    +latitude?: number,
    +longitude?: number,
    +externalId?: ?string,
    +locationType?: {|
      +id: string,
      +name: string,
      +mapType: ?string,
      +mapZoomLevel: ?number,
      +propertyTypes: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
        +index: ?number,
        +isInstanceProperty: ?boolean,
        +type: PropertyKind,
        +nodeType: ?string,
        +stringValue: ?string,
        +intValue: ?number,
        +floatValue: ?number,
        +booleanValue: ?boolean,
        +latitudeValue: ?number,
        +longitudeValue: ?number,
        +rangeFromValue: ?number,
        +rangeToValue: ?number,
        +isMandatory: ?boolean,
      |}>,
    |},
    +equipments?: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +equipmentType: {|
        +id: string,
        +name: string,
      |},
    |}>,
    +properties?: $ReadOnlyArray<?{|
      +id: string,
      +propertyType: {|
        +id: string,
        +name: string,
        +type: PropertyKind,
        +nodeType: ?string,
        +index: ?number,
        +isEditable: ?boolean,
        +isInstanceProperty: ?boolean,
        +stringValue: ?string,
        +isMandatory: ?boolean,
      |},
      +stringValue: ?string,
      +intValue: ?number,
      +booleanValue: ?boolean,
      +floatValue: ?number,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +nodeValue: ?{|
        +id: string,
        +name: string,
      |},
    |}>,
  |}
|};
export type LocationAddEditCardQuery = {|
  variables: LocationAddEditCardQueryVariables,
  response: LocationAddEditCardQueryResponse,
|};
*/


/*
query LocationAddEditCardQuery(
  $locationId: ID!
) {
  location: node(id: $locationId) {
    __typename
    ... on Location {
      id
      name
      latitude
      longitude
      externalId
      locationType {
        id
        name
        mapType
        mapZoomLevel
        propertyTypes {
          id
          name
          index
          isInstanceProperty
          type
          nodeType
          stringValue
          intValue
          floatValue
          booleanValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          isMandatory
        }
      }
      equipments {
        id
        name
        equipmentType {
          id
          name
        }
      }
      properties {
        id
        propertyType {
          id
          name
          type
          nodeType
          index
          isEditable
          isInstanceProperty
          stringValue
          isMandatory
        }
        stringValue
        intValue
        booleanValue
        floatValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        nodeValue {
          __typename
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
    "name": "locationId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "locationId"
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
  "name": "latitude",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitude",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "externalId",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "nodeType",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "mapType",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "mapZoomLevel",
      "args": null,
      "storageKey": null
    },
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
        (v7/*: any*/),
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
  "name": "equipments",
  "storageKey": null,
  "args": null,
  "concreteType": "Equipment",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": (v21/*: any*/)
    }
  ]
},
v23 = {
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
    (v9/*: any*/),
    (v10/*: any*/),
    (v7/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "isEditable",
      "args": null,
      "storageKey": null
    },
    (v8/*: any*/),
    (v11/*: any*/),
    (v19/*: any*/)
  ]
},
v24 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationAddEditCardQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v20/*: any*/),
              (v22/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "properties",
                "storageKey": null,
                "args": null,
                "concreteType": "Property",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v23/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v14/*: any*/),
                  (v13/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "nodeValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": (v21/*: any*/)
                  }
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
    "name": "LocationAddEditCardQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          (v24/*: any*/),
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v20/*: any*/),
              (v22/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "properties",
                "storageKey": null,
                "args": null,
                "concreteType": "Property",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v23/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v14/*: any*/),
                  (v13/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "nodeValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": [
                      (v24/*: any*/),
                      (v2/*: any*/),
                      (v3/*: any*/)
                    ]
                  }
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
    "name": "LocationAddEditCardQuery",
    "id": null,
    "text": "query LocationAddEditCardQuery(\n  $locationId: ID!\n) {\n  location: node(id: $locationId) {\n    __typename\n    ... on Location {\n      id\n      name\n      latitude\n      longitude\n      externalId\n      locationType {\n        id\n        name\n        mapType\n        mapZoomLevel\n        propertyTypes {\n          id\n          name\n          index\n          isInstanceProperty\n          type\n          nodeType\n          stringValue\n          intValue\n          floatValue\n          booleanValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          isMandatory\n        }\n      }\n      equipments {\n        id\n        name\n        equipmentType {\n          id\n          name\n        }\n      }\n      properties {\n        id\n        propertyType {\n          id\n          name\n          type\n          nodeType\n          index\n          isEditable\n          isInstanceProperty\n          stringValue\n          isMandatory\n        }\n        stringValue\n        intValue\n        booleanValue\n        floatValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        nodeValue {\n          __typename\n          id\n          name\n        }\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'abfd7091b5c16d02b59c1efa02e234f2';
module.exports = node;
