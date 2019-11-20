/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash cf22b9996609b238dc365a5dafe9554c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type EntDetailsQueryVariables = {|
  id: string
|};
export type EntDetailsQueryResponse = {|
  +vertex: ?{|
    +id: string,
    +type: string,
    +fields: $ReadOnlyArray<{|
      +name: string,
      +value: string,
      +type: string,
    |}>,
    +edges: $ReadOnlyArray<{|
      +name: string,
      +type: string,
      +ids: $ReadOnlyArray<string>,
    |}>,
  |}
|};
export type EntDetailsQuery = {|
  variables: EntDetailsQueryVariables,
  response: EntDetailsQueryResponse,
|};
*/


/*
query EntDetailsQuery(
  $id: ID!
) {
  vertex(id: $id) {
    id
    type
    fields {
      name
      value
      type
    }
    edges {
      name
      type
      ids
    }
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
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
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
    "name": "vertex",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      }
    ],
    "concreteType": "Vertex",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "id",
        "args": null,
        "storageKey": null
      },
      (v1/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "fields",
        "storageKey": null,
        "args": null,
        "concreteType": "Field",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "value",
            "args": null,
            "storageKey": null
          },
          (v1/*: any*/)
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "Edge",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v1/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "ids",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EntDetailsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EntDetailsQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v3/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "EntDetailsQuery",
    "id": null,
    "text": "query EntDetailsQuery(\n  $id: ID!\n) {\n  vertex(id: $id) {\n    id\n    type\n    fields {\n      name\n      value\n      type\n    }\n    edges {\n      name\n      type\n      ids\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4be33521bcf990ae9df86c758a47fcde';
module.exports = node;
