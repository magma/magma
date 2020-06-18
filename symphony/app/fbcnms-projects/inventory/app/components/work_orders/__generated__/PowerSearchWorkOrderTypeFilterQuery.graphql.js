/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash c27160cec3f78222bdba8e76696031c8
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PowerSearchWorkOrderTypeFilterQueryVariables = {||};
export type PowerSearchWorkOrderTypeFilterQueryResponse = {|
  +workOrderTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
      |}
    |}>
  |}
|};
export type PowerSearchWorkOrderTypeFilterQuery = {|
  variables: PowerSearchWorkOrderTypeFilterQueryVariables,
  response: PowerSearchWorkOrderTypeFilterQueryResponse,
|};
*/


/*
query PowerSearchWorkOrderTypeFilterQuery {
  workOrderTypes {
    edges {
      node {
        id
        name
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "workOrderTypes",
    "storageKey": null,
    "args": null,
    "concreteType": "WorkOrderTypeConnection",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "edges",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderTypeEdge",
        "plural": true,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "node",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderType",
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
                "name": "name",
                "args": null,
                "storageKey": null
              }
            ]
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
    "name": "PowerSearchWorkOrderTypeFilterQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "PowerSearchWorkOrderTypeFilterQuery",
    "argumentDefinitions": [],
    "selections": (v0/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "PowerSearchWorkOrderTypeFilterQuery",
    "id": null,
    "text": "query PowerSearchWorkOrderTypeFilterQuery {\n  workOrderTypes {\n    edges {\n      node {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd50e2d512e76bfdc0b0e7bb397d4ec3f';
module.exports = node;
