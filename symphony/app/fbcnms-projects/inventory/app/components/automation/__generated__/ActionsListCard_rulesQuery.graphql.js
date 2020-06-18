/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1c80c064fc304a4994166ab76fd7bc5c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ActionsListCard_actionsRule$ref = any;
export type ActionsListCard_rulesQueryVariables = {||};
export type ActionsListCard_rulesQueryResponse = {|
  +actionsRules: ?{|
    +results: $ReadOnlyArray<?{|
      +$fragmentRefs: ActionsListCard_actionsRule$ref
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
      ...ActionsListCard_actionsRule
      id
    }
  }
}

fragment ActionRow_data on ActionsTrigger {
  triggerID
  supportedActions {
    actionID
    dataType
    description
  }
}

fragment ActionsAddDialog_triggerData on ActionsTrigger {
  triggerID
  description
  ...ActionRow_data
  ...TriggerFilterRow_data
}

fragment ActionsListCard_actionsRule on ActionsRule {
  id
  name
  trigger {
    description
    ...ActionsAddDialog_triggerData
    id
  }
  ruleActions {
    actionID
    data
  }
  ruleFilters {
    filterID
    operatorID
    data
  }
}

fragment TriggerFilterOperator_data on ActionsFilter {
  supportedOperators {
    operatorID
    description
    dataType
  }
}

fragment TriggerFilterRow_data on ActionsTrigger {
  triggerID
  supportedFilters {
    filterID
    description
    supportedOperators {
      operatorID
    }
    ...TriggerFilterOperator_data
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
  "name": "description",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "actionID",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "dataType",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "filterID",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "operatorID",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "data",
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
              {
                "kind": "FragmentSpread",
                "name": "ActionsListCard_actionsRule",
                "args": null
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
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "name",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "trigger",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsTrigger",
                "plural": false,
                "selections": [
                  (v1/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "triggerID",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "supportedActions",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ActionsAction",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v1/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "supportedFilters",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ActionsFilter",
                    "plural": true,
                    "selections": [
                      (v4/*: any*/),
                      (v1/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "supportedOperators",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "ActionsOperator",
                        "plural": true,
                        "selections": [
                          (v5/*: any*/),
                          (v1/*: any*/),
                          (v3/*: any*/)
                        ]
                      }
                    ]
                  },
                  (v0/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "ruleActions",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsRuleAction",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v6/*: any*/)
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "ruleFilters",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsRuleFilter",
                "plural": true,
                "selections": [
                  (v4/*: any*/),
                  (v5/*: any*/),
                  (v6/*: any*/)
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
    "text": "query ActionsListCard_rulesQuery {\n  actionsRules {\n    results {\n      ...ActionsListCard_actionsRule\n      id\n    }\n  }\n}\n\nfragment ActionRow_data on ActionsTrigger {\n  triggerID\n  supportedActions {\n    actionID\n    dataType\n    description\n  }\n}\n\nfragment ActionsAddDialog_triggerData on ActionsTrigger {\n  triggerID\n  description\n  ...ActionRow_data\n  ...TriggerFilterRow_data\n}\n\nfragment ActionsListCard_actionsRule on ActionsRule {\n  id\n  name\n  trigger {\n    description\n    ...ActionsAddDialog_triggerData\n    id\n  }\n  ruleActions {\n    actionID\n    data\n  }\n  ruleFilters {\n    filterID\n    operatorID\n    data\n  }\n}\n\nfragment TriggerFilterOperator_data on ActionsFilter {\n  supportedOperators {\n    operatorID\n    description\n    dataType\n  }\n}\n\nfragment TriggerFilterRow_data on ActionsTrigger {\n  triggerID\n  supportedFilters {\n    filterID\n    description\n    supportedOperators {\n      operatorID\n    }\n    ...TriggerFilterOperator_data\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ea80e4c4f3ab927ae77127e709c1ca77';
module.exports = node;
