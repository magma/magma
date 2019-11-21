/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 910a8e4823506abef49d9b6d98353192
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type MarkLocationTypeIsSiteMutationVariables = {|
  id: string,
  isSite: boolean,
|};
export type MarkLocationTypeIsSiteMutationResponse = {|
  +markLocationTypeIsSite: ?{|
    +id: string,
    +isSite: boolean,
  |}
|};
export type MarkLocationTypeIsSiteMutation = {|
  variables: MarkLocationTypeIsSiteMutationVariables,
  response: MarkLocationTypeIsSiteMutationResponse,
|};
*/


/*
mutation MarkLocationTypeIsSiteMutation(
  $id: ID!
  $isSite: Boolean!
) {
  markLocationTypeIsSite(id: $id, isSite: $isSite) {
    id
    isSite
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
  },
  {
    "kind": "LocalArgument",
    "name": "isSite",
    "type": "Boolean!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "markLocationTypeIsSite",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      },
      {
        "kind": "Variable",
        "name": "isSite",
        "variableName": "isSite"
      }
    ],
    "concreteType": "LocationType",
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
        "name": "isSite",
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
    "name": "MarkLocationTypeIsSiteMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "MarkLocationTypeIsSiteMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "MarkLocationTypeIsSiteMutation",
    "id": null,
    "text": "mutation MarkLocationTypeIsSiteMutation(\n  $id: ID!\n  $isSite: Boolean!\n) {\n  markLocationTypeIsSite(id: $id, isSite: $isSite) {\n    id\n    isSite\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fbd84093825f107c5e04d591c1e5bc0e';
module.exports = node;
