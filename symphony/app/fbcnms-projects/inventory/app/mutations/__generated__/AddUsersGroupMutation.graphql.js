/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 3ba15c32f451155ec8102fcbfaf90382
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type AddUsersGroupInput = {|
  name: string,
  description?: ?string,
  members?: ?$ReadOnlyArray<string>,
  policies?: ?$ReadOnlyArray<string>,
|};
export type AddUsersGroupMutationVariables = {|
  input: AddUsersGroupInput
|};
export type AddUsersGroupMutationResponse = {|
  +addUsersGroup: {|
    +id: string,
    +name: string,
    +description: ?string,
    +status: UsersGroupStatus,
    +members: $ReadOnlyArray<{|
      +id: string,
      +authID: string,
      +firstName: string,
      +lastName: string,
      +email: string,
      +status: UserStatus,
      +role: UserRole,
      +profilePhoto: ?{|
        +id: string,
        +fileName: string,
        +storeKey: ?string,
      |},
    |}>,
    +policies: $ReadOnlyArray<{|
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
    |}>,
  |}
|};
export type AddUsersGroupMutation = {|
  variables: AddUsersGroupMutationVariables,
  response: AddUsersGroupMutationResponse,
|};
*/


/*
mutation AddUsersGroupMutation(
  $input: AddUsersGroupInput!
) {
  addUsersGroup(input: $input) {
    id
    name
    description
    status
    members {
      id
      authID
      firstName
      lastName
      email
      status
      role
      profilePhoto {
        id
        fileName
        storeKey
      }
    }
    policies {
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
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddUsersGroupInput!",
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
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "status",
  "args": null,
  "storageKey": null
},
v5 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isAllowed",
    "args": null,
    "storageKey": null
  }
],
v6 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v5/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v5/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v5/*: any*/)
  }
],
v7 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addUsersGroup",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "UsersGroup",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      (v2/*: any*/),
      (v3/*: any*/),
      (v4/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "members",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "authID",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "firstName",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "lastName",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "email",
            "args": null,
            "storageKey": null
          },
          (v4/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "role",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "profilePhoto",
            "storageKey": null,
            "args": null,
            "concreteType": "File",
            "plural": false,
            "selections": [
              (v1/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "fileName",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "storeKey",
                "args": null,
                "storageKey": null
              }
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "policies",
        "storageKey": null,
        "args": null,
        "concreteType": "PermissionsPolicy",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
          (v3/*: any*/),
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
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "read",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "BasicPermissionRule",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "location",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationCUD",
                    "plural": false,
                    "selections": [
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "create",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationPermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "update",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationPermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "delete",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationPermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipment",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "serviceType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
                  }
                ]
              },
              {
                "kind": "InlineFragment",
                "type": "WorkforcePolicy",
                "selections": [
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "read",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "templates",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v6/*: any*/)
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
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "create",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "update",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "delete",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "assign",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "transferOwnership",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v5/*: any*/)
                      }
                    ]
                  }
                ]
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
    "name": "AddUsersGroupMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v7/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddUsersGroupMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v7/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddUsersGroupMutation",
    "id": null,
    "text": "mutation AddUsersGroupMutation(\n  $input: AddUsersGroupInput!\n) {\n  addUsersGroup(input: $input) {\n    id\n    name\n    description\n    status\n    members {\n      id\n      authID\n      firstName\n      lastName\n      email\n      status\n      role\n      profilePhoto {\n        id\n        fileName\n        storeKey\n      }\n    }\n    policies {\n      id\n      name\n      description\n      isGlobal\n      policy {\n        __typename\n        ... on InventoryPolicy {\n          read {\n            isAllowed\n          }\n          location {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          equipment {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          equipmentType {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          locationType {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          portType {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          serviceType {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n        }\n        ... on WorkforcePolicy {\n          read {\n            isAllowed\n          }\n          templates {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n          }\n          data {\n            create {\n              isAllowed\n            }\n            update {\n              isAllowed\n            }\n            delete {\n              isAllowed\n            }\n            assign {\n              isAllowed\n            }\n            transferOwnership {\n              isAllowed\n            }\n          }\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a7b68e738f78bffc91dc722538502f62';
module.exports = node;
