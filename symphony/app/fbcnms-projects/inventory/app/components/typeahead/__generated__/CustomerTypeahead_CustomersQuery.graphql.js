/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 2beb626c12ee24bd5f4c4199a825e03f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type CustomerTypeahead_CustomersQueryVariables = {|
  limit?: ?number
|};
export type CustomerTypeahead_CustomersQueryResponse = {|
  +customerSearch: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +externalId: ?string,
  |}>
|};
export type CustomerTypeahead_CustomersQuery = {|
  variables: CustomerTypeahead_CustomersQueryVariables,
  response: CustomerTypeahead_CustomersQueryResponse,
|};
*/


/*
query CustomerTypeahead_CustomersQuery(
  $limit: Int
) {
  customerSearch(limit: $limit) {
    id
    name
    externalId
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "limit",
    "type": "Int",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "customerSearch",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "limit",
        "variableName": "limit"
      }
    ],
    "concreteType": "Customer",
    "plural": true,
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
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "externalId",
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
    "name": "CustomerTypeahead_CustomersQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "CustomerTypeahead_CustomersQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "CustomerTypeahead_CustomersQuery",
    "id": null,
    "text": "query CustomerTypeahead_CustomersQuery(\n  $limit: Int\n) {\n  customerSearch(limit: $limit) {\n    id\n    name\n    externalId\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ce52a287db805de0357e25ea1c021505';
module.exports = node;
