/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 3724b87b4b4bda0499e6778d2dde1c7f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type DocumentTable_files$ref = any;
type FileAttachment_file$ref = any;
export type ImageEntity = "CHECKLIST_ITEM" | "EQUIPMENT" | "LOCATION" | "SITE_SURVEY" | "USER" | "WORK_ORDER" | "%future added value";
export type DeleteImageMutationVariables = {|
  entityType: ImageEntity,
  entityId: string,
  id: string,
|};
export type DeleteImageMutationResponse = {|
  +deleteImage: {|
    +$fragmentRefs: DocumentTable_files$ref & FileAttachment_file$ref
  |}
|};
export type DeleteImageMutation = {|
  variables: DeleteImageMutationVariables,
  response: DeleteImageMutationResponse,
|};
*/


/*
mutation DeleteImageMutation(
  $entityType: ImageEntity!
  $entityId: ID!
  $id: ID!
) {
  deleteImage(entityType: $entityType, entityId: $entityId, id: $id) {
    ...DocumentTable_files
    ...FileAttachment_file
    id
  }
}

fragment DocumentTable_files on File {
  id
  fileName
  category
  ...FileAttachment_file
}

fragment FileAttachment_file on File {
  id
  fileName
  sizeInBytes
  uploaded
  fileType
  storeKey
  category
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
    "name": "entityType",
    "type": "ImageEntity!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "entityId",
    "type": "ID!",
    "defaultValue": null
  },
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
    "name": "entityId",
    "variableName": "entityId"
  },
  {
    "kind": "Variable",
    "name": "entityType",
    "variableName": "entityType"
  },
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "DeleteImageMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "deleteImage",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "File",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "DocumentTable_files",
            "args": null
          },
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
    "name": "DeleteImageMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "deleteImage",
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
            "name": "category",
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
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeleteImageMutation",
    "id": null,
    "text": "mutation DeleteImageMutation(\n  $entityType: ImageEntity!\n  $entityId: ID!\n  $id: ID!\n) {\n  deleteImage(entityType: $entityType, entityId: $entityId, id: $id) {\n    ...DocumentTable_files\n    ...FileAttachment_file\n    id\n  }\n}\n\nfragment DocumentTable_files on File {\n  id\n  fileName\n  category\n  ...FileAttachment_file\n}\n\nfragment FileAttachment_file on File {\n  id\n  fileName\n  sizeInBytes\n  uploaded\n  fileType\n  storeKey\n  category\n  ...ImageDialog_img\n}\n\nfragment ImageDialog_img on File {\n  storeKey\n  fileName\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e35a3c256648c1f9c986d51a8d7b77bb';
module.exports = node;
