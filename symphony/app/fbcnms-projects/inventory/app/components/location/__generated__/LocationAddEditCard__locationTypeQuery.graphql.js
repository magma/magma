/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9212c67f1516f904f7f2a5c1135e4a34
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type LocationAddEditCard__locationTypeQueryVariables = {|
  locationTypeId: string
|};
export type LocationAddEditCard__locationTypeQueryResponse = {|
  +locationType: ?{|
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
      +isInstanceProperty: ?boolean,
      +isMandatory: ?boolean,
    |}>,
  |}
|};
export type LocationAddEditCard__locationTypeQuery = {|
  variables: LocationAddEditCard__locationTypeQueryVariables,
  response: LocationAddEditCard__locationTypeQueryResponse,
|};
*/


/*
query LocationAddEditCard__locationTypeQuery(
  $locationTypeId: ID!
) {
  locationType: node(id: $locationTypeId) {
    __typename
    ... on LocationType {
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
        isInstanceProperty
        isMandatory
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
    "name": "locationTypeId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "locationTypeId"
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
      "name": "isInstanceProperty",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "isMandatory",
      "args": null,
      "storageKey": null
    }
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationAddEditCard__locationTypeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "locationType",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "LocationType",
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
    "name": "LocationAddEditCard__locationTypeQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "locationType",
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
            "type": "LocationType",
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
    "name": "LocationAddEditCard__locationTypeQuery",
    "id": null,
    "text": "query LocationAddEditCard__locationTypeQuery(\n  $locationTypeId: ID!\n) {\n  locationType: node(id: $locationTypeId) {\n    __typename\n    ... on LocationType {\n      id\n      name\n      propertyTypes {\n        id\n        name\n        type\n        nodeType\n        index\n        stringValue\n        intValue\n        booleanValue\n        floatValue\n        latitudeValue\n        longitudeValue\n        rangeFromValue\n        rangeToValue\n        isEditable\n        isInstanceProperty\n        isMandatory\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f48e051d918fb0d2fd6ff5ca4a7d441c';
module.exports = node;
