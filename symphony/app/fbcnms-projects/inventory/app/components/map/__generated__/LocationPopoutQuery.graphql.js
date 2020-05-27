/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash d8b2952b5834a69ac6720183e5be305e
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type DynamicPropertiesGrid_properties$ref = any;
type DynamicPropertiesGrid_propertyTypes$ref = any;
type LocationBreadcrumbsTitle_locationDetails$ref = any;
export type LocationPopoutQueryVariables = {|
  locationId: string
|};
export type LocationPopoutQueryResponse = {|
  +location: ?{|
    +id?: string,
    +name?: string,
    +locationType?: {|
      +name: string,
      +propertyTypes: $ReadOnlyArray<?{|
        +$fragmentRefs: DynamicPropertiesGrid_propertyTypes$ref
      |}>,
    |},
    +properties?: $ReadOnlyArray<?{|
      +$fragmentRefs: DynamicPropertiesGrid_properties$ref
    |}>,
    +$fragmentRefs: LocationBreadcrumbsTitle_locationDetails$ref,
  |}
|};
export type LocationPopoutQuery = {|
  variables: LocationPopoutQueryVariables,
  response: LocationPopoutQueryResponse,
|};
*/


/*
query LocationPopoutQuery(
  $locationId: ID!
) {
  location: node(id: $locationId) {
    __typename
    ... on Location {
      id
      name
      locationType {
        name
        id
        propertyTypes {
          ...DynamicPropertiesGrid_propertyTypes
          id
        }
      }
      ...LocationBreadcrumbsTitle_locationDetails
      properties {
        ...DynamicPropertiesGrid_properties
        id
      }
    }
    id
  }
}

fragment DynamicPropertiesGrid_properties on Property {
  ...PropertyFormField_property
  propertyType {
    id
    index
  }
}

fragment DynamicPropertiesGrid_propertyTypes on PropertyType {
  id
  name
  index
  isInstanceProperty
  type
  nodeType
  stringValue
  intValue
  booleanValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  floatValue
}

fragment LocationBreadcrumbsTitle_locationDetails on Location {
  id
  name
  locationType {
    name
    id
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
}

fragment PropertyFormField_property on Property {
  id
  propertyType {
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
    category
    isDeleted
  }
  stringValue
  intValue
  floatValue
  booleanValue
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
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "nodeType",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationPopoutQuery",
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
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
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
                      {
                        "kind": "FragmentSpread",
                        "name": "DynamicPropertiesGrid_propertyTypes",
                        "args": null
                      }
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
                    "kind": "FragmentSpread",
                    "name": "DynamicPropertiesGrid_properties",
                    "args": null
                  }
                ]
              },
              {
                "kind": "FragmentSpread",
                "name": "LocationBreadcrumbsTitle_locationDetails",
                "args": null
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationPopoutQuery",
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
          (v4/*: any*/),
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
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
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v5/*: any*/),
                      (v6/*: any*/),
                      (v7/*: any*/),
                      (v8/*: any*/),
                      (v9/*: any*/),
                      (v10/*: any*/),
                      (v11/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/),
                      (v14/*: any*/),
                      (v15/*: any*/),
                      (v16/*: any*/)
                    ]
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationHierarchy",
                "storageKey": null,
                "args": null,
                "concreteType": "Location",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationType",
                    "plural": false,
                    "selections": [
                      (v3/*: any*/),
                      (v2/*: any*/)
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
                  (v2/*: any*/),
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
                      (v7/*: any*/),
                      (v8/*: any*/),
                      (v5/*: any*/),
                      (v9/*: any*/),
                      (v10/*: any*/),
                      (v11/*: any*/),
                      (v16/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/),
                      (v14/*: any*/),
                      (v15/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isEditable",
                        "args": null,
                        "storageKey": null
                      },
                      (v6/*: any*/),
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
                        "name": "category",
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
                  (v9/*: any*/),
                  (v10/*: any*/),
                  (v16/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "nodeValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": [
                      (v4/*: any*/),
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
    "name": "LocationPopoutQuery",
    "id": null,
    "text": "query LocationPopoutQuery(\n  $locationId: ID!\n) {\n  location: node(id: $locationId) {\n    __typename\n    ... on Location {\n      id\n      name\n      locationType {\n        name\n        id\n        propertyTypes {\n          ...DynamicPropertiesGrid_propertyTypes\n          id\n        }\n      }\n      ...LocationBreadcrumbsTitle_locationDetails\n      properties {\n        ...DynamicPropertiesGrid_properties\n        id\n      }\n    }\n    id\n  }\n}\n\nfragment DynamicPropertiesGrid_properties on Property {\n  ...PropertyFormField_property\n  propertyType {\n    id\n    index\n  }\n}\n\nfragment DynamicPropertiesGrid_propertyTypes on PropertyType {\n  id\n  name\n  index\n  isInstanceProperty\n  type\n  nodeType\n  stringValue\n  intValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  floatValue\n}\n\nfragment LocationBreadcrumbsTitle_locationDetails on Location {\n  id\n  name\n  locationType {\n    name\n    id\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n}\n\nfragment PropertyFormField_property on Property {\n  id\n  propertyType {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isInstanceProperty\n    isMandatory\n    category\n    isDeleted\n  }\n  stringValue\n  intValue\n  floatValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  nodeValue {\n    __typename\n    id\n    name\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '80eec4f9084409c4f63a764b3b6e91a1';
module.exports = node;
