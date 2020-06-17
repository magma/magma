/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4e95ef44d950f97aaf5bb7c64b9a4fef
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type HubVersionQueryVariables = {||};
export type HubVersionQueryResponse = {|
  +version: {|
    +string: string
  |}
|};
export type HubVersionQuery = {|
  variables: HubVersionQueryVariables,
  response: HubVersionQueryResponse,
|};
*/


/*
query HubVersionQuery {
  version {
    string
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "version",
    "storageKey": null,
    "args": null,
    "concreteType": "Version",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "string",
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
    "name": "HubVersionQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "HubVersionQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "HubVersionQuery",
    "id": null,
    "text": "query HubVersionQuery {\n  version {\n    string\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '34ebb96a32faf7bffcb55abca095e681';
module.exports = node;
