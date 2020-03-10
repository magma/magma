/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 93b8b0ab231ec3948c9edbde543e35af
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type FileAttachment_file$ref = any;
export type ImageEntity = "EQUIPMENT" | "LOCATION" | "SITE_SURVEY" | "USER" | "WORK_ORDER" | "%future added value";
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
export type AddImageMutationVariables = {|
  input: AddImageInput
|};
export type AddImageMutationResponse = {|
  +addImage: {|
    +$fragmentRefs: FileAttachment_file$ref
  |}
|};
export type AddImageMutation = {|
  variables: AddImageMutationVariables,
  response: AddImageMutationResponse,
|};
*/


/*
mutation AddImageMutation(
  $input: AddImageInput!
) {
  addImage(input: $input) {
    ...FileAttachment_file
    id
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
    "type": "AddImageInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddImageMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addImage",
        "storageKey": null,
        "args": (v1/*: any*/),
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
  },
  "operation": {
    "kind": "Operation",
    "name": "AddImageMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addImage",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "File",
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
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddImageMutation",
    "id": null,
    "text": "mutation AddImageMutation(\n  $input: AddImageInput!\n) {\n  addImage(input: $input) {\n    ...FileAttachment_file\n    id\n  }\n}\n\nfragment DocumentMenu_document on File {\n  id\n  fileName\n  storeKey\n  fileType\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...DocumentMenu_document\n  ...ImageDialog_img\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'cfb8687a3be6e209c5d3c9a6f94c249e';
module.exports = node;
