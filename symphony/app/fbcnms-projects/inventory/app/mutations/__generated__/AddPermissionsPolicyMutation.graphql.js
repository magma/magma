/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 28a9ae057345ab05e3095a37134dc2ae
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type AddPermissionsPolicyInput = {|
  name: string,
  description?: ?string,
  isGlobal?: ?boolean,
  inventoryInput?: ?InventoryPolicyInput,
  workforceInput?: ?WorkforcePolicyInput,
|};
export type InventoryPolicyInput = {|
  read?: ?BasicPermissionRuleInput,
  location?: ?BasicCUDInput,
  equipment?: ?BasicCUDInput,
  equipmentType?: ?BasicCUDInput,
  locationType?: ?BasicCUDInput,
  portType?: ?BasicCUDInput,
  serviceType?: ?BasicCUDInput,
|};
export type BasicPermissionRuleInput = {|
  isAllowed: PermissionValue
|};
export type BasicCUDInput = {|
  create?: ?BasicPermissionRuleInput,
  update?: ?BasicPermissionRuleInput,
  delete?: ?BasicPermissionRuleInput,
|};
export type WorkforcePolicyInput = {|
  read?: ?BasicPermissionRuleInput,
  data?: ?BasicWorkforceCUDInput,
  templates?: ?BasicCUDInput,
|};
export type BasicWorkforceCUDInput = {|
  create?: ?BasicPermissionRuleInput,
  update?: ?BasicPermissionRuleInput,
  delete?: ?BasicPermissionRuleInput,
  assign?: ?BasicPermissionRuleInput,
  transferOwnership?: ?BasicPermissionRuleInput,
|};
export type AddPermissionsPolicyMutationVariables = {|
  input: AddPermissionsPolicyInput
|};
export type AddPermissionsPolicyMutationResponse = {|
  +addPermissionsPolicy: {|
    +id: string,
    +name: string,
    +description: ?string,
    +isGlobal: boolean,
    +policy: {|
      +__typename: "InventoryPolicy",
      +read: {|
        +isAllowed: PermissionValue
      |},
    |} | {|
      +__typename: "WorkforcePolicy",
      +read: {|
        +isAllowed: PermissionValue
      |},
    |} | {|
      // This will never be '%other', but we need some
      // value in case none of the concrete values match.
      +__typename: "%other"
    |},
    +groups: $ReadOnlyArray<{|
      +id: string
    |}>,
  |}
|};
export type AddPermissionsPolicyMutation = {|
  variables: AddPermissionsPolicyMutationVariables,
  response: AddPermissionsPolicyMutationResponse,
|};
*/


/*
mutation AddPermissionsPolicyMutation(
  $input: AddPermissionsPolicyInput!
) {
  addPermissionsPolicy(input: $input) {
    id
    name
    description
    isGlobal
    policy {
      __typename
      ... on InventoryPolicy {
        __typename
        read {
          isAllowed
        }
      }
      ... on WorkforcePolicy {
        __typename
        read {
          isAllowed
        }
      }
    }
    groups {
      id
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddPermissionsPolicyInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isGlobal",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v7 = [
  (v6/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "read",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "isAllowed",
        "args": null,
        "storageKey": null
      }
    ]
  }
],
v8 = {
  "kind": "InlineFragment",
  "type": "InventoryPolicy",
  "selections": (v7/*: any*/)
},
v9 = {
  "kind": "InlineFragment",
  "type": "WorkforcePolicy",
  "selections": (v7/*: any*/)
},
v10 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "groups",
  "storageKey": null,
  "args": null,
  "concreteType": "UsersGroup",
  "plural": true,
  "selections": [
    (v2/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddPermissionsPolicyMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addPermissionsPolicy",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "PermissionsPolicy",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "policy",
            "storageKey": null,
            "args": null,
            "concreteType": null,
            "plural": false,
            "selections": [
              (v8/*: any*/),
              (v9/*: any*/)
            ]
          },
          (v10/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddPermissionsPolicyMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addPermissionsPolicy",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "PermissionsPolicy",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "policy",
            "storageKey": null,
            "args": null,
            "concreteType": null,
            "plural": false,
            "selections": [
              (v6/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/)
            ]
          },
          (v10/*: any*/)
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddPermissionsPolicyMutation",
    "id": null,
    "text": "mutation AddPermissionsPolicyMutation(\n  $input: AddPermissionsPolicyInput!\n) {\n  addPermissionsPolicy(input: $input) {\n    id\n    name\n    description\n    isGlobal\n    policy {\n      __typename\n      ... on InventoryPolicy {\n        __typename\n        read {\n          isAllowed\n        }\n      }\n      ... on WorkforcePolicy {\n        __typename\n        read {\n          isAllowed\n        }\n      }\n    }\n    groups {\n      id\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fe64ae13190a12f82069e6eceeba7980';
module.exports = node;
