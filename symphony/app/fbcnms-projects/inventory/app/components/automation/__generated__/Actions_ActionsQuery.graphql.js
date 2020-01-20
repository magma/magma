/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 89cc945f16da8636cfcd4b9b2f85ba2f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ActionsAddDialog_triggerData$ref = any;
export type TriggerID = "magma_alert" | "%future added value";
export type Actions_ActionsQueryVariables = {||};
export type Actions_ActionsQueryResponse = {|
  +actionsTriggers: ?{|
    +results: $ReadOnlyArray<?{|
      +triggerID: TriggerID,
      +description: string,
      +$fragmentRefs: ActionsAddDialog_triggerData$ref,
    |}>,
    +count: number,
  |}
|};
export type Actions_ActionsQuery = {|
  variables: Actions_ActionsQueryVariables,
  response: Actions_ActionsQueryResponse,
|};
*/


/*
query Actions_ActionsQuery {
  actionsTriggers {
    results {
      triggerID
      description
      ...ActionsAddDialog_triggerData
      id
    }
    count
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
  "name": "triggerID",
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
  "name": "count",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "dataType",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "Actions_ActionsQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "actionsTriggers",
        "storageKey": null,
        "args": null,
        "concreteType": "ActionsTriggersSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "results",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsTrigger",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v1/*: any*/),
              {
                "kind": "FragmentSpread",
                "name": "ActionsAddDialog_triggerData",
                "args": null
              }
            ]
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "Actions_ActionsQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "actionsTriggers",
        "storageKey": null,
        "args": null,
        "concreteType": "ActionsTriggersSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "results",
            "storageKey": null,
            "args": null,
            "concreteType": "ActionsTrigger",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v1/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "supportedActions",
                "storageKey": null,
                "args": null,
                "concreteType": "ActionsAction",
                "plural": true,
                "selections": [
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "actionID",
                    "args": null,
                    "storageKey": null
                  },
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
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "filterID",
                    "args": null,
                    "storageKey": null
                  },
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
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "operatorID",
                        "args": null,
                        "storageKey": null
                      },
                      (v1/*: any*/),
                      (v3/*: any*/)
                    ]
                  }
                ]
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "id",
                "args": null,
                "storageKey": null
              }
            ]
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "Actions_ActionsQuery",
    "id": null,
    "text": "query Actions_ActionsQuery {\n  actionsTriggers {\n    results {\n      triggerID\n      description\n      ...ActionsAddDialog_triggerData\n      id\n    }\n    count\n  }\n}\n\nfragment ActionRow_data on ActionsTrigger {\n  triggerID\n  supportedActions {\n    actionID\n    dataType\n    description\n  }\n}\n\nfragment ActionsAddDialog_triggerData on ActionsTrigger {\n  triggerID\n  description\n  ...ActionRow_data\n  ...TriggerFilterRow_data\n}\n\nfragment TriggerFilterOperator_data on ActionsFilter {\n  supportedOperators {\n    operatorID\n    description\n    dataType\n  }\n}\n\nfragment TriggerFilterRow_data on ActionsTrigger {\n  triggerID\n  supportedFilters {\n    filterID\n    description\n    supportedOperators {\n      operatorID\n    }\n    ...TriggerFilterOperator_data\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '3d7dfd46b65c97833ed524cbe563e058';
module.exports = node;
