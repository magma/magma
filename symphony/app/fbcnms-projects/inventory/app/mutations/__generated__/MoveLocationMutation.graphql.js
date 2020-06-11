/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 77db1948ec186d25858fdaf4227fd999
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MoveLocationMutationVariables = {|
  locationID: string,
  parentLocationID?: ?string,
|};
export type MoveLocationMutationResponse = {|
  +moveLocation: {|
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
|};
export type MoveLocationMutation = {|
  variables: MoveLocationMutationVariables,
  response: MoveLocationMutationResponse,
|};
*/


/*
mutation MoveLocationMutation(
  $locationID: ID!
  $parentLocationID: ID
) {
  moveLocation(locationID: $locationID, parentLocationID: $parentLocationID) {
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
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "locationID",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "parentLocationID",
    "type": "ID",
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
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "moveLocation",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "locationID",
        "variableName": "locationID"
      },
      {
        "kind": "Variable",
        "name": "parentLocationID",
        "variableName": "parentLocationID"
      }
    ],
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
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "MoveLocationMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "MoveLocationMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "MoveLocationMutation",
    "id": null,
    "text": "mutation MoveLocationMutation(\n  $locationID: ID!\n  $parentLocationID: ID\n) {\n  moveLocation(locationID: $locationID, parentLocationID: $parentLocationID) {\n    id\n    externalId\n    name\n    locationType {\n      id\n      name\n    }\n    numChildren\n    siteSurveyNeeded\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '98502526513b758cda745d28f8b0f7ee';
module.exports = node;
