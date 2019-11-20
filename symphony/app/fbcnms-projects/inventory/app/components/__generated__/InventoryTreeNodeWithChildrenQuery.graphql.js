/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 06dc9eb32f3b1b07af4ecc9112bc4675
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type InventoryTreeNodeWithChildrenQueryVariables = {|
  id: string
|};
export type InventoryTreeNodeWithChildrenQueryResponse = {|
  +location: ?{|
    +id?: string,
    +externalId?: ?string,
    +name?: string,
    +locationType?: {|
      +id: string,
      +name: string,
    |},
    +numChildren?: number,
    +siteSurveyNeeded?: boolean,
    +children?: $ReadOnlyArray<?{|
      +id: string,
      +externalId: ?string,
      +name: string,
      +locationType: {|
        +id: string,
        +name: string,
      |},
      +numChildren: number,
      +siteSurveyNeeded: boolean,
    |}>,
  |}
|};
export type InventoryTreeNodeWithChildrenQuery = {|
  variables: InventoryTreeNodeWithChildrenQueryVariables,
  response: InventoryTreeNodeWithChildrenQueryResponse,
|};
*/


/*
query InventoryTreeNodeWithChildrenQuery(
  $id: ID!
) {
  location: node(id: $id) {
    __typename
    ... on Location {
      id
      externalId
      name
      locationType {
        id
        name
      }
      numChildren
      siteSurveyNeeded
      children {
        id
        externalId
        name
        locationType {
          id
          name
        }
        numChildren
        siteSurveyNeeded
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
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
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
  "name": "externalId",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v4/*: any*/)
  ]
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "numChildren",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "siteSurveyNeeded",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "children",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    (v4/*: any*/),
    (v5/*: any*/),
    (v6/*: any*/),
    (v7/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "InventoryTreeNodeWithChildrenQuery",
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
              (v7/*: any*/),
              (v8/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "InventoryTreeNodeWithChildrenQuery",
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
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "InventoryTreeNodeWithChildrenQuery",
    "id": null,
    "text": "query InventoryTreeNodeWithChildrenQuery(\n  $id: ID!\n) {\n  location: node(id: $id) {\n    __typename\n    ... on Location {\n      id\n      externalId\n      name\n      locationType {\n        id\n        name\n      }\n      numChildren\n      siteSurveyNeeded\n      children {\n        id\n        externalId\n        name\n        locationType {\n          id\n          name\n        }\n        numChildren\n        siteSurveyNeeded\n      }\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '76a5d4ca2dafaf7f713fd3916fc4420a';
module.exports = node;
