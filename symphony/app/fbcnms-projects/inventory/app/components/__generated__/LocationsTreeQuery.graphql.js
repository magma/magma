/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9b72c9304d67049927dffc3f57e315ae
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type LocationsTreeQueryVariables = {||};
export type LocationsTreeQueryResponse = {|
  +locations: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +externalId: ?string,
        +name: string,
        +locationType: {|
          +id: string,
          +name: string,
        |},
        +numChildren: number,
        +siteSurveyNeeded: boolean,
      |}
    |}>
  |}
|};
export type LocationsTreeQuery = {|
  variables: LocationsTreeQueryVariables,
  response: LocationsTreeQueryResponse,
|};
*/


/*
query LocationsTreeQuery {
  locations(first: 50, onlyTopLevel: true) {
    edges {
      node {
        id
        externalId
        name
        locationType {
          id
          name
        }
        numChildren
        siteSurveyNeeded
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
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "Literal",
  "name": "onlyTopLevel",
  "value": true
},
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
          (v1/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "externalId",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "locationType",
            "storageKey": null,
            "args": null,
            "concreteType": "LocationType",
            "plural": false,
            "selections": [
              (v1/*: any*/),
              (v2/*: any*/)
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numChildren",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "siteSurveyNeeded",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          }
        ]
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "cursor",
        "args": null,
        "storageKey": null
      }
    ]
  },
  {
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
  }
],
v4 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50
  },
  (v0/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationsTreeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "locations",
        "name": "__LocationsTree_locations_connection",
        "storageKey": "__LocationsTree_locations_connection(onlyTopLevel:true)",
        "args": [
          (v0/*: any*/)
        ],
        "concreteType": "LocationConnection",
        "plural": false,
        "selections": (v3/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationsTreeQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "locations",
        "storageKey": "locations(first:50,onlyTopLevel:true)",
        "args": (v4/*: any*/),
        "concreteType": "LocationConnection",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "locations",
        "args": (v4/*: any*/),
        "handle": "connection",
        "key": "LocationsTree_locations",
        "filters": [
          "onlyTopLevel"
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "LocationsTreeQuery",
    "id": null,
    "text": "query LocationsTreeQuery {\n  locations(first: 50, onlyTopLevel: true) {\n    edges {\n      node {\n        id\n        externalId\n        name\n        locationType {\n          id\n          name\n        }\n        numChildren\n        siteSurveyNeeded\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "locations"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7fe3739f9c1d4a0d510669710d0dd3e3';
module.exports = node;
