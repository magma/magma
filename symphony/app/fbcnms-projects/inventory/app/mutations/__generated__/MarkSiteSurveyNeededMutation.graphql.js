/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e2f73cec2f03c2d1ce797125a4cf903d
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MarkSiteSurveyNeededMutationVariables = {|
  locationId: string,
  needed: boolean,
|};
export type MarkSiteSurveyNeededMutationResponse = {|
  +markSiteSurveyNeeded: {|
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
export type MarkSiteSurveyNeededMutation = {|
  variables: MarkSiteSurveyNeededMutationVariables,
  response: MarkSiteSurveyNeededMutationResponse,
|};
*/


/*
mutation MarkSiteSurveyNeededMutation(
  $locationId: ID!
  $needed: Boolean!
) {
  markSiteSurveyNeeded(locationId: $locationId, needed: $needed) {
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
    "name": "locationId",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "needed",
    "type": "Boolean!",
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
    "name": "markSiteSurveyNeeded",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "locationId",
        "variableName": "locationId"
      },
      {
        "kind": "Variable",
        "name": "needed",
        "variableName": "needed"
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
    "name": "MarkSiteSurveyNeededMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "MarkSiteSurveyNeededMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "MarkSiteSurveyNeededMutation",
    "id": null,
    "text": "mutation MarkSiteSurveyNeededMutation(\n  $locationId: ID!\n  $needed: Boolean!\n) {\n  markSiteSurveyNeeded(locationId: $locationId, needed: $needed) {\n    id\n    externalId\n    name\n    locationType {\n      id\n      name\n    }\n    numChildren\n    siteSurveyNeeded\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a5c70fb90bb9de415d3c763e0d8617a8';
module.exports = node;
