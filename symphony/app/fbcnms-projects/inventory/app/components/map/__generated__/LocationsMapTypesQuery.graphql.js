/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 52284b6d63204f84e17e43b314827f8d
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type LocationsMapTypesQueryVariables = {||};
export type LocationsMapTypesQueryResponse = {|
  +locationTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +locations: ?{|
          +edges: $ReadOnlyArray<{|
            +node: ?{|
              +id: string,
              +name: string,
              +latitude: number,
              +longitude: number,
            |}
          |}>
        |},
      |}
    |}>
  |}
|};
export type LocationsMapTypesQuery = {|
  variables: LocationsMapTypesQueryVariables,
  response: LocationsMapTypesQueryResponse,
|};
*/


/*
query LocationsMapTypesQuery {
  locationTypes {
    edges {
      node {
        id
        name
        locations(enforceHasLatLong: true) {
          edges {
            node {
              id
              name
              latitude
              longitude
            }
          }
        }
      }
    }
  }
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
v2 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "locationTypes",
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
                "name": "locations",
                "storageKey": "locations(enforceHasLatLong:true)",
                "args": [
                  {
                    "kind": "Literal",
                    "name": "enforceHasLatLong",
                    "value": true
                  }
                ],
                "concreteType": "LocationConnection",
                "plural": false,
                "selections": [
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "edges",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationEdge",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "node",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Location",
                        "plural": false,
                        "selections": [
                          (v0/*: any*/),
                          (v1/*: any*/),
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "latitude",
                            "args": null,
                            "storageKey": null
                          },
                          {
                            "kind": "ScalarField",
                            "alias": null,
                            "name": "longitude",
                            "args": null,
                            "storageKey": null
                          }
                        ]
                      }
                    ]
                  }
                ]
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
    "name": "LocationsMapTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v2/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationsMapTypesQuery",
    "argumentDefinitions": [],
    "selections": (v2/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "LocationsMapTypesQuery",
    "id": null,
    "text": "query LocationsMapTypesQuery {\n  locationTypes {\n    edges {\n      node {\n        id\n        name\n        locations(enforceHasLatLong: true) {\n          edges {\n            node {\n              id\n              name\n              latitude\n              longitude\n            }\n          }\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fc689b6b48f62c12d1c96d10d3f80c7c';
module.exports = node;
