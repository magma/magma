/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4f9cab2ea6d7b3e09b3e9f40ac456f98
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type EditPermissionsPolicyInput = {|
  id: string,
  name?: ?string,
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
export type EditPermissionsPolicyMutationVariables = {|
  input: EditPermissionsPolicyInput
|};
export type EditPermissionsPolicyMutationResponse = {|
  +editPermissionsPolicy: {|
    +id: string,
    +name: string,
    +description: ?string,
    +isGlobal: boolean,
    +policy: {|
      +__typename: "InventoryPolicy",
      +read: {|
        +isAllowed: PermissionValue
      |},
      +location: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +equipment: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +equipmentType: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +locationType: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +portType: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +serviceType: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
    |} | {|
      +__typename: "WorkforcePolicy",
      +read: {|
        +isAllowed: PermissionValue
      |},
      +templates: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
      |},
      +data: {|
        +create: {|
          +isAllowed: PermissionValue
        |},
        +update: {|
          +isAllowed: PermissionValue
        |},
        +delete: {|
          +isAllowed: PermissionValue
        |},
        +assign: {|
          +isAllowed: PermissionValue
        |},
        +transferOwnership: {|
          +isAllowed: PermissionValue
        |},
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
export type EditPermissionsPolicyMutation = {|
  variables: EditPermissionsPolicyMutationVariables,
  response: EditPermissionsPolicyMutationResponse,
|};
*/


/*
mutation EditPermissionsPolicyMutation(
  $input: EditPermissionsPolicyInput!
) {
  editPermissionsPolicy(input: $input) {
    id
    name
    description
    isGlobal
    policy {
      __typename
      ... on InventoryPolicy {
        read {
          isAllowed
        }
        location {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        equipment {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        equipmentType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        locationType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        portType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        serviceType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
      }
      ... on WorkforcePolicy {
        read {
          isAllowed
        }
        templates {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        data {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
          assign {
            isAllowed
          }
          transferOwnership {
            isAllowed
          }
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
    "type": "EditPermissionsPolicyInput!",
    "defaultValue": null
  }
],
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v2 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isAllowed",
    "args": null,
    "storageKey": null
  }
],
v3 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "read",
  "storageKey": null,
  "args": null,
  "concreteType": "BasicPermissionRule",
  "plural": false,
  "selections": (v2/*: any*/)
},
v4 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "create",
  "storageKey": null,
  "args": null,
  "concreteType": "BasicPermissionRule",
  "plural": false,
  "selections": (v2/*: any*/)
},
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "update",
  "storageKey": null,
  "args": null,
  "concreteType": "BasicPermissionRule",
  "plural": false,
  "selections": (v2/*: any*/)
},
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "delete",
  "storageKey": null,
  "args": null,
  "concreteType": "BasicPermissionRule",
  "plural": false,
  "selections": (v2/*: any*/)
},
v7 = [
  (v4/*: any*/),
  (v5/*: any*/),
  (v6/*: any*/)
],
v8 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "editPermissionsPolicy",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "PermissionsPolicy",
    "plural": false,
    "selections": [
      (v1/*: any*/),
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
        "name": "description",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "isGlobal",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "policy",
        "storageKey": null,
        "args": null,
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "InlineFragment",
            "type": "InventoryPolicy",
            "selections": [
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "location",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipment",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "portType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "serviceType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              }
            ]
          },
          {
            "kind": "InlineFragment",
            "type": "WorkforcePolicy",
            "selections": [
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "templates",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "data",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkforceCUD",
                "plural": false,
                "selections": [
                  (v4/*: any*/),
                  (v5/*: any*/),
                  (v6/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "assign",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "BasicPermissionRule",
                    "plural": false,
                    "selections": (v2/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "transferOwnership",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "BasicPermissionRule",
                    "plural": false,
                    "selections": (v2/*: any*/)
                  }
                ]
              }
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "groups",
        "storageKey": null,
        "args": null,
        "concreteType": "UsersGroup",
        "plural": true,
        "selections": [
          (v1/*: any*/)
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditPermissionsPolicyMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v8/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EditPermissionsPolicyMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v8/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditPermissionsPolicyMutation",
    "id": null,
    "text": "mutation EditPermissionsPolicyMutation(\n  $input: EditPermissionsPolicyInput!\n) {\n  editPermissionsPolicy(input: $input) {\n    id\n    name\n    description\n    isGlobal\n    policy {\n      __typename\n      ... on InventoryPolicy {\n        read {\n          isAllowed\n        }\n        location {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        equipment {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        equipmentType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        locationType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        portType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        serviceType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n      }\n      ... on WorkforcePolicy {\n        read {\n          isAllowed\n        }\n        templates {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        data {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n          assign {\n            isAllowed\n          }\n          transferOwnership {\n            isAllowed\n          }\n        }\n      }\n    }\n    groups {\n      id\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '91fc92917a36ff108e90243e830927f6';
module.exports = node;
