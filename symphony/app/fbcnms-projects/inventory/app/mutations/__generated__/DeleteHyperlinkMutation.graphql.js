/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4957ba9f2df8b640e3c3af4cb0d044eb
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type HyperlinkTableRow_hyperlink$ref = any;
export type DeleteHyperlinkMutationVariables = {|
  id: string
|};
export type DeleteHyperlinkMutationResponse = {|
  +deleteHyperlink: {|
    +$fragmentRefs: HyperlinkTableRow_hyperlink$ref
  |}
|};
export type DeleteHyperlinkMutation = {|
  variables: DeleteHyperlinkMutationVariables,
  response: DeleteHyperlinkMutationResponse,
|};
*/


/*
mutation DeleteHyperlinkMutation(
  $id: ID!
) {
  deleteHyperlink(id: $id) {
    ...HyperlinkTableRow_hyperlink
    id
  }
}

fragment HyperlinkTableRow_hyperlink on Hyperlink {
  id
  category
  url
  displayName
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
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "DeleteHyperlinkMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "deleteHyperlink",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Hyperlink",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "HyperlinkTableRow_hyperlink",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "DeleteHyperlinkMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "deleteHyperlink",
        "storageKey": null,
        "args": (v1/*: any*/),
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
            "name": "category",
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
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeleteHyperlinkMutation",
    "id": null,
    "text": "mutation DeleteHyperlinkMutation(\n  $id: ID!\n) {\n  deleteHyperlink(id: $id) {\n    ...HyperlinkTableRow_hyperlink\n    id\n  }\n}\n\nfragment HyperlinkTableRow_hyperlink on Hyperlink {\n  id\n  category\n  url\n  displayName\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'eea56538c03c8ce55b62e4e48331a303';
module.exports = node;
