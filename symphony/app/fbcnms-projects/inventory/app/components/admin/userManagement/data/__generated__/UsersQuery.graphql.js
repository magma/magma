/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4d71541396eef70465756062855b1cc3
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersQueryVariables = {||};
export type UsersQueryResponse = {|
  +users: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +authID: string,
        +firstName: string,
        +lastName: string,
        +email: string,
        +status: UserStatus,
        +role: UserRole,
        +groups: $ReadOnlyArray<?{|
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
                  +isAllowed: PermissionValue,
                  +locationTypeIds: ?$ReadOnlyArray<string>,
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
                +isAllowed: PermissionValue,
                +projectTypeIds: ?$ReadOnlyArray<string>,
                +workOrderTypeIds: ?$ReadOnlyArray<string>,
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
        |}>,
      |}
    |}>
  |}
|};
export type UsersQuery = {|
  variables: UsersQueryVariables,
  response: UsersQueryResponse,
|};
*/


/*
query UsersQuery {
  users(first: 500) {
    edges {
      node {
        id
        authID
        firstName
        lastName
        email
        status
        role
        groups {
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
                    locationTypeIds
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
                  projectTypeIds
                  workOrderTypeIds
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
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
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
  "name": "authID",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "firstName",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "lastName",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "email",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "status",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "role",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isAllowed",
  "args": null,
  "storageKey": null
},
v11 = [
  (v10/*: any*/)
],
v12 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v11/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v11/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v11/*: any*/)
  }
],
v13 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "edges",
    "storageKey": null,
    "args": null,
    "concreteType": "UserEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": false,
        "selections": [
          (v0/*: any*/),
          (v1/*: any*/),
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
          (v5/*: any*/),
          (v6/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "groups",
            "storageKey": null,
            "args": null,
            "concreteType": "UsersGroup",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              (v5/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "members",
                "storageKey": null,
                "args": null,
                "concreteType": "User",
                "plural": true,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v4/*: any*/),
                  (v5/*: any*/),
                  (v6/*: any*/)
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
                  (v0/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
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
                      (v9/*: any*/),
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
                            "selections": (v11/*: any*/)
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
                                "selections": (v11/*: any*/)
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "update",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "LocationPermissionRule",
                                "plural": false,
                                "selections": [
                                  (v10/*: any*/),
                                  {
                                    "kind": "ScalarField",
                                    "alias": null,
                                    "name": "locationTypeIds",
                                    "args": null,
                                    "storageKey": null
                                  }
                                ]
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "delete",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "LocationPermissionRule",
                                "plural": false,
                                "selections": (v11/*: any*/)
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
                            "selections": (v12/*: any*/)
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "equipmentType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "CUD",
                            "plural": false,
                            "selections": (v12/*: any*/)
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "locationType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "CUD",
                            "plural": false,
                            "selections": (v12/*: any*/)
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "portType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "CUD",
                            "plural": false,
                            "selections": (v12/*: any*/)
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "serviceType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "CUD",
                            "plural": false,
                            "selections": (v12/*: any*/)
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
                            "selections": [
                              (v10/*: any*/),
                              {
                                "kind": "ScalarField",
                                "alias": null,
                                "name": "projectTypeIds",
                                "args": null,
                                "storageKey": null
                              },
                              {
                                "kind": "ScalarField",
                                "alias": null,
                                "name": "workOrderTypeIds",
                                "args": null,
                                "storageKey": null
                              }
                            ]
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "templates",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "CUD",
                            "plural": false,
                            "selections": (v12/*: any*/)
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
                                "selections": (v11/*: any*/)
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "update",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "WorkforcePermissionRule",
                                "plural": false,
                                "selections": (v11/*: any*/)
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "delete",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "WorkforcePermissionRule",
                                "plural": false,
                                "selections": (v11/*: any*/)
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "assign",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "WorkforcePermissionRule",
                                "plural": false,
                                "selections": (v11/*: any*/)
                              },
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "transferOwnership",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "WorkforcePermissionRule",
                                "plural": false,
                                "selections": (v11/*: any*/)
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
          },
          (v9/*: any*/)
        ]
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "cursor",
        "args": null,
        "storageKey": null
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "pageInfo",
    "storageKey": null,
    "args": null,
    "concreteType": "PageInfo",
    "plural": false,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "endCursor",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "hasNextPage",
        "args": null,
        "storageKey": null
      }
    ]
  }
],
v14 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 500
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "UsersQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "users",
        "name": "__Users_users_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "UserConnection",
        "plural": false,
        "selections": (v13/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "UsersQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "users",
        "storageKey": "users(first:500)",
        "args": (v14/*: any*/),
        "concreteType": "UserConnection",
        "plural": false,
        "selections": (v13/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "users",
        "args": (v14/*: any*/),
        "handle": "connection",
        "key": "Users_users",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "UsersQuery",
    "id": null,
    "text": "query UsersQuery {\n  users(first: 500) {\n    edges {\n      node {\n        id\n        authID\n        firstName\n        lastName\n        email\n        status\n        role\n        groups {\n          id\n          name\n          description\n          status\n          members {\n            id\n            authID\n            firstName\n            lastName\n            email\n            status\n            role\n          }\n          policies {\n            id\n            name\n            description\n            isGlobal\n            policy {\n              __typename\n              ... on InventoryPolicy {\n                read {\n                  isAllowed\n                }\n                location {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                    locationTypeIds\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                equipment {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                equipmentType {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                locationType {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                portType {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                serviceType {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n              }\n              ... on WorkforcePolicy {\n                read {\n                  isAllowed\n                  projectTypeIds\n                  workOrderTypeIds\n                }\n                templates {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                }\n                data {\n                  create {\n                    isAllowed\n                  }\n                  update {\n                    isAllowed\n                  }\n                  delete {\n                    isAllowed\n                  }\n                  assign {\n                    isAllowed\n                  }\n                  transferOwnership {\n                    isAllowed\n                  }\n                }\n              }\n            }\n          }\n        }\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "users"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '54d6b28dfd92b7fe7120702e142658a0';
module.exports = node;
