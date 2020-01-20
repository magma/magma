/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 036b7f9a705890393b751664d71352b2
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type PropertyTypeFormField_propertyType$ref = any;
export type LocationsQueryVariables = {||};
export type LocationsQueryResponse = {|
  +locationTypes: ?{|
    +edges: ?$ReadOnlyArray<?{|
      +node: ?{|
        +id: string,
        +name: string,
        +propertyTypes: $ReadOnlyArray<?{|
          +$fragmentRefs: PropertyTypeFormField_propertyType$ref
        |}>,
        +numberOfLocations: number,
      |}
    |}>
  |}
|};
export type LocationsQuery = {|
  variables: LocationsQueryVariables,
  response: LocationsQueryResponse,
|};
*/


/*
query LocationsQuery {
  locationTypes(first: 50) {
    edges {
      node {
        id
        name
        propertyTypes {
          ...PropertyTypeFormField_propertyType
          id
        }
        numberOfLocations
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}

fragment PropertyTypeFormField_propertyType on PropertyType {
  id
  name
  type
  stringValue
  intValue
  booleanValue
  floatValue
  latitudeValue
  longitudeValue
  isEditable
  isInstanceProperty
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "numberOfLocations",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "pageInfo",
  "storageKey": null,
  "args": null,
  "concreteType": "PageInfo",
  "plural": false,
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "endCursor",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "hasNextPage",
      "args": null,
      "storageKey": null
    }
  ]
},
v6 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50,
    "type": "Int"
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "locationTypes",
        "name": "__Catalog_locationTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "LocationTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "LocationTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
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
                        "name": "PropertyTypeFormField_propertyType",
                        "args": null
                      }
                    ]
                  },
                  (v2/*: any*/),
                  (v3/*: any*/)
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationsQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "locationTypes",
        "storageKey": "locationTypes(first:50)",
        "args": (v6/*: any*/),
        "concreteType": "LocationTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "LocationTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
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
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "booleanValue",
                        "args": null,
                        "storageKey": null
                      },
                      (v0/*: any*/),
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
                      (v1/*: any*/),
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
                      }
                    ]
                  },
                  (v2/*: any*/),
                  (v3/*: any*/)
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "locationTypes",
        "args": (v6/*: any*/),
        "handle": "connection",
        "key": "Catalog_locationTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "LocationsQuery",
    "id": null,
    "text": "query LocationsQuery {\n  locationTypes(first: 50) {\n    edges {\n      node {\n        id\n        name\n        propertyTypes {\n          ...PropertyTypeFormField_propertyType\n          id\n        }\n        numberOfLocations\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  isEditable\n  isInstanceProperty\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "locationTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '60c38374e2f2a92250f9910897429aef';
module.exports = node;
