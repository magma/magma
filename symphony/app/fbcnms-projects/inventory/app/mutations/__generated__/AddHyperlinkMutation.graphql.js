/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9af15badd8a1718dc07da7d410b5bfe9
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ImageEntity = "EQUIPMENT" | "LOCATION" | "SITE_SURVEY" | "WORK_ORDER" | "%future added value";
export type AddHyperlinkInput = {|
  entityType: ImageEntity,
  entityId: string,
  url: string,
  displayName?: ?string,
  category?: ?string,
|};
export type AddHyperlinkMutationVariables = {|
  input: AddHyperlinkInput
|};
export type AddHyperlinkMutationResponse = {|
  +addHyperlink: {|
    +id: string,
    +url: string,
    +displayName: ?string,
    +category: ?string,
    +createTime: any,
  |}
|};
export type AddHyperlinkMutation = {|
  variables: AddHyperlinkMutationVariables,
  response: AddHyperlinkMutationResponse,
|};
*/


/*
mutation AddHyperlinkMutation(
  $input: AddHyperlinkInput!
) {
  addHyperlink(input: $input) {
    id
    url
    displayName
    category
    createTime
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddHyperlinkInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addHyperlink",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "Hyperlink",
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
        "name": "url",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "displayName",
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
        "name": "createTime",
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
    "name": "AddHyperlinkMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddHyperlinkMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddHyperlinkMutation",
    "id": null,
    "text": "mutation AddHyperlinkMutation(\n  $input: AddHyperlinkInput!\n) {\n  addHyperlink(input: $input) {\n    id\n    url\n    displayName\n    category\n    createTime\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8aeb4cde118c6a2a5365fe50a82c419e';
module.exports = node;
