/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e2a6da25973b9e3712aebdfcff64f826
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type LocationTypeahead_LocationsQueryVariables = {|
  name: string
|};
export type LocationTypeahead_LocationsQueryResponse = {|
  +searchForNode: {|
    +edges: ?$ReadOnlyArray<{|
      +node: ?({|
        +__typename: "Location",
        +id: string,
        +externalId: ?string,
        +name: string,
        +locationType: {|
          +name: string
        |},
        +locationHierarchy: $ReadOnlyArray<{|
          +id: string,
          +name: string,
          +locationType: {|
            +name: string
          |},
        |}>,
      |} | {|
        // This will never be '%other', but we need some
        // value in case none of the concrete values match.
        +__typename: "%other"
      |})
    |}>
  |}
|};
export type LocationTypeahead_LocationsQuery = {|
  variables: LocationTypeahead_LocationsQueryVariables,
  response: LocationTypeahead_LocationsQueryResponse,
|};
*/


/*
query LocationTypeahead_LocationsQuery(
  $name: String!
) {
  searchForNode(name: $name, first: 10) {
    edges {
      node {
        __typename
        ... on Location {
          id
          externalId
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
        id
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "name",
    "type": "String!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 10
  },
  {
    "kind": "Variable",
    "name": "name",
    "variableName": "name"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "externalId",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": [
    (v5/*: any*/)
  ]
},
v7 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": [
    (v5/*: any*/),
    (v3/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationTypeahead_LocationsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "searchForNode",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "SearchNodesConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "SearchNodeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": null,
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "InlineFragment",
                    "type": "Location",
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v6/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "locationHierarchy",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Location",
                        "plural": true,
                        "selections": [
                          (v3/*: any*/),
                          (v5/*: any*/),
                          (v6/*: any*/)
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
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationTypeahead_LocationsQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "searchForNode",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "SearchNodesConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "SearchNodeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": null,
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "InlineFragment",
                    "type": "Location",
                    "selections": [
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v7/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "locationHierarchy",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Location",
                        "plural": true,
                        "selections": [
                          (v3/*: any*/),
                          (v5/*: any*/),
                          (v7/*: any*/)
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
  },
  "params": {
    "operationKind": "query",
    "name": "LocationTypeahead_LocationsQuery",
    "id": null,
    "text": "query LocationTypeahead_LocationsQuery(\n  $name: String!\n) {\n  searchForNode(name: $name, first: 10) {\n    edges {\n      node {\n        __typename\n        ... on Location {\n          id\n          externalId\n          name\n          locationType {\n            name\n            id\n          }\n          locationHierarchy {\n            id\n            name\n            locationType {\n              name\n              id\n            }\n          }\n        }\n        id\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9e4ab0ce3b932ff34a72c8db6f428190';
module.exports = node;
