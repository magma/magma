/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 47894127b28d44985070388a5ca8975d
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type FileAttachment_file$ref = any;
export type ImageEntity = "CHECKLIST_ITEM" | "EQUIPMENT" | "LOCATION" | "SITE_SURVEY" | "USER" | "WORK_ORDER" | "%future added value";
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
  +addFloorPlan: {|
    +id: string,
    +name: string,
    +image: {|
      +$fragmentRefs: FileAttachment_file$ref
    |},
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
    image {
      ...FileAttachment_file
      id
    }
  }
}

fragment DocumentMenu_document on File {
  id
  fileName
  storeKey
  fileType
}

fragment FileAttachment_file on File {
  id
  fileName
  sizeInBytes
  uploaded
  fileType
  storeKey
  category
  ...DocumentMenu_document
  ...ImageDialog_img
}

fragment ImageDialog_img on File {
  storeKey
  fileName
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
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
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
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddFloorPlanMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addFloorPlan",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "FloorPlan",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "image",
            "storageKey": null,
            "args": null,
            "concreteType": "File",
            "plural": false,
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "FileAttachment_file",
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
    "name": "AddFloorPlanMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addFloorPlan",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "FloorPlan",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "image",
            "storageKey": null,
            "args": null,
            "concreteType": "File",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "fileName",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "sizeInBytes",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "uploaded",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "fileType",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "storeKey",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "category",
                "args": null,
                "storageKey": null
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddFloorPlanMutation",
    "id": null,
    "text": "mutation AddFloorPlanMutation(\n  $input: AddFloorPlanInput!\n) {\n  addFloorPlan(input: $input) {\n    id\n    name\n    image {\n      ...FileAttachment_file\n      id\n    }\n  }\n}\n\nfragment DocumentMenu_document on File {\n  id\n  fileName\n  storeKey\n  fileType\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...DocumentMenu_document\n  ...ImageDialog_img\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd318199b52e4e0c7bbc6d5467669fe28';
module.exports = node;
