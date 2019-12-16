/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash c10665f2bcab1fa8b00ea7c90081f6ac
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type GraphVertexDetailsQueryVariables = {|
  id: string
|};
export type GraphVertexDetailsQueryResponse = {|
  +vertex: ?{|
    +id: string,
    +type: string,
    +fields: $ReadOnlyArray<{|
      +name: string,
      +value: string,
    |}>,
  |}
|};
export type GraphVertexDetailsQuery = {|
  variables: GraphVertexDetailsQueryVariables,
  response: GraphVertexDetailsQueryResponse,
|};
*/


/*
query GraphVertexDetailsQuery(
  $id: ID!
) {
  vertex(id: $id) {
    id
    type
    fields {
      name
      value
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
v1 = [
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
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "type",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "fields",
        "storageKey": null,
        "args": null,
        "concreteType": "Field",
        "plural": true,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "name",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "value",
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
    "name": "GraphVertexDetailsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "GraphVertexDetailsQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "GraphVertexDetailsQuery",
    "id": null,
    "text": "query GraphVertexDetailsQuery(\n  $id: ID!\n) {\n  vertex(id: $id) {\n    id\n    type\n    fields {\n      name\n      value\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b810a4963a64c2695e2cb81fbd02be9e';
module.exports = node;
