/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 415f081227e6c8d2a67a16ced96d08f1
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type ActionsListCard_rulesQueryVariables = {||};
export type ActionsListCard_rulesQueryResponse = {|
  +actionsRules: ?{|
    +results: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +trigger: {|
        +description: string
      |},
    |}>
  |}
|};
export type ActionsListCard_rulesQuery = {|
  variables: ActionsListCard_rulesQueryVariables,
  response: ActionsListCard_rulesQueryResponse,
|};
*/


/*
query ActionsListCard_rulesQuery {
  actionsRules {
    results {
      id
      name
      trigger {
        description
        id
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ActionsListCard_rulesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "actionsRules",
        "storageKey": null,
        "args": null,
        "concreteType": "ActionsRulesSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "results",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsRule",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v1/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "trigger",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsTrigger",
                "plural": false,
                "selections": [
                  (v2/*: any*/)
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ActionsListCard_rulesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "actionsRules",
        "storageKey": null,
        "args": null,
        "concreteType": "ActionsRulesSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "results",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsRule",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v1/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "trigger",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsTrigger",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v0/*: any*/)
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ActionsListCard_rulesQuery",
    "id": null,
    "text": "query ActionsListCard_rulesQuery {\n  actionsRules {\n    results {\n      id\n      name\n      trigger {\n        description\n        id\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'aaa447f9e9053947abb83fa377bf9cec';
module.exports = node;
