/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 7bc415ec59f72e400314cfc284133ffd
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ImageEntity = "EQUIPMENT" | "LOCATION" | "SITE_SURVEY" | "WORK_ORDER" | "%future added value";
export type AddFloorPlanInput = {|
  name: string,
  locationID: string,
  image: AddImageInput,
  referenceX: number,
  referenceY: number,
  latitude: number,
  longitude: number,
  referencePoint1X: number,
  referencePoint1Y: number,
  referencePoint2X: number,
  referencePoint2Y: number,
  scaleInMeters: number,
|};
export type AddImageInput = {|
  entityType: ImageEntity,
  entityId: string,
  imgKey: string,
  fileName: string,
  fileSize: number,
  modified: any,
  contentType: string,
  category?: ?string,
|};
export type AddFloorPlanMutationVariables = {|
  input: AddFloorPlanInput
|};
export type AddFloorPlanMutationResponse = {|
  +addFloorPlan: ?{|
    +id: string,
    +name: string,
  |}
|};
export type AddFloorPlanMutation = {|
  variables: AddFloorPlanMutationVariables,
  response: AddFloorPlanMutationResponse,
|};
*/


/*
mutation AddFloorPlanMutation(
  $input: AddFloorPlanInput!
) {
  addFloorPlan(input: $input) {
    id
    name
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddFloorPlanInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addFloorPlan",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "FloorPlan",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "id",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "name",
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
    "name": "AddFloorPlanMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddFloorPlanMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddFloorPlanMutation",
    "id": null,
    "text": "mutation AddFloorPlanMutation(\n  $input: AddFloorPlanInput!\n) {\n  addFloorPlan(input: $input) {\n    id\n    name\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '6d933e125c353c1d23ed2b8f2d9ba362';
module.exports = node;
