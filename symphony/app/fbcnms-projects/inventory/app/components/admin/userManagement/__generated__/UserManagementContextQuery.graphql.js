/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 6e71c418ddf16e72c00d0021f5a551d0
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UserManagementContextQueryVariables = {||};
export type UserManagementContextQueryResponse = {|
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
        |}>,
        +profilePhoto: ?{|
          +id: string,
          +fileName: string,
          +storeKey: ?string,
        |},
      |}
    |}>
  |},
  +usersGroups: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +description: ?string,
        +status: UsersGroupStatus,
        +members: $ReadOnlyArray<{|
          +id: string,
          +authID: string,
        |}>,
      |}
    |}>
  |},
  +permissionsPolicies: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
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
    |}>
  |},
|};
export type UserManagementContextQuery = {|
  variables: UserManagementContextQueryVariables,
  response: UserManagementContextQueryResponse,
|};
*/


/*
query UserManagementContextQuery {
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
        }
        profilePhoto {
          id
          fileName
          storeKey
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
  usersGroups(first: 500) {
    edges {
      node {
        id
        name
        description
        status
        members {
          id
          authID
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
  permissionsPolicies(first: 500) {
    edges {
      node {
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
  "name": "status",
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
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v6 = {
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
},
v7 = [
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
          (v2/*: any*/),
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
            "name": "groups",
            "storageKey": null,
            "args": null,
            "concreteType": "UsersGroup",
            "plural": true,
            "selections": [
              (v0/*: any*/),
              (v3/*: any*/)
            ]
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
              (v0/*: any*/),
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
          },
          (v4/*: any*/)
        ]
      },
      (v5/*: any*/)
    ]
  },
  (v6/*: any*/)
],
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "description",
  "args": null,
  "storageKey": null
},
v9 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "edges",
    "storageKey": null,
    "args": null,
    "concreteType": "UsersGroupEdge",
    "plural": true,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "node",
        "storageKey": null,
        "args": null,
        "concreteType": "UsersGroup",
        "plural": false,
        "selections": [
          (v0/*: any*/),
          (v3/*: any*/),
          (v8/*: any*/),
          (v2/*: any*/),
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
              (v1/*: any*/)
            ]
          },
          (v4/*: any*/)
        ]
      },
      (v5/*: any*/)
    ]
  },
  (v6/*: any*/)
],
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isGlobal",
  "args": null,
  "storageKey": null
},
v11 = [
  (v4/*: any*/),
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
v12 = {
  "kind": "InlineFragment",
  "type": "InventoryPolicy",
  "selections": (v11/*: any*/)
},
v13 = {
  "kind": "InlineFragment",
  "type": "WorkforcePolicy",
  "selections": (v11/*: any*/)
},
v14 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "groups",
  "storageKey": null,
  "args": null,
  "concreteType": "UsersGroup",
  "plural": true,
  "selections": [
    (v0/*: any*/)
  ]
},
v15 = [
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
    "name": "UserManagementContextQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "users",
        "name": "__UserManagementContext_users_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "UserConnection",
        "plural": false,
        "selections": (v7/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": "usersGroups",
        "name": "__UserManagementContext_usersGroups_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "UsersGroupConnection",
        "plural": false,
        "selections": (v9/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": "permissionsPolicies",
        "name": "__UserManagementContext_permissionsPolicies_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "PermissionsPolicyConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "PermissionsPolicyEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "PermissionsPolicy",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v3/*: any*/),
                  (v8/*: any*/),
                  (v10/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "policy",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": [
                      (v12/*: any*/),
                      (v13/*: any*/)
                    ]
                  },
                  (v14/*: any*/),
                  (v4/*: any*/)
                ]
              },
              (v5/*: any*/)
            ]
          },
          (v6/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "UserManagementContextQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "users",
        "storageKey": "users(first:500)",
        "args": (v15/*: any*/),
        "concreteType": "UserConnection",
        "plural": false,
        "selections": (v7/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "users",
        "args": (v15/*: any*/),
        "handle": "connection",
        "key": "UserManagementContext_users",
        "filters": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "usersGroups",
        "storageKey": "usersGroups(first:500)",
        "args": (v15/*: any*/),
        "concreteType": "UsersGroupConnection",
        "plural": false,
        "selections": (v9/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "usersGroups",
        "args": (v15/*: any*/),
        "handle": "connection",
        "key": "UserManagementContext_usersGroups",
        "filters": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "permissionsPolicies",
        "storageKey": "permissionsPolicies(first:500)",
        "args": (v15/*: any*/),
        "concreteType": "PermissionsPolicyConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "PermissionsPolicyEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "PermissionsPolicy",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v3/*: any*/),
                  (v8/*: any*/),
                  (v10/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "policy",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": [
                      (v4/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/)
                    ]
                  },
                  (v14/*: any*/),
                  (v4/*: any*/)
                ]
              },
              (v5/*: any*/)
            ]
          },
          (v6/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "permissionsPolicies",
        "args": (v15/*: any*/),
        "handle": "connection",
        "key": "UserManagementContext_permissionsPolicies",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "UserManagementContextQuery",
    "id": null,
    "text": "query UserManagementContextQuery {\n  users(first: 500) {\n    edges {\n      node {\n        id\n        authID\n        firstName\n        lastName\n        email\n        status\n        role\n        groups {\n          id\n          name\n        }\n        profilePhoto {\n          id\n          fileName\n          storeKey\n        }\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n  usersGroups(first: 500) {\n    edges {\n      node {\n        id\n        name\n        description\n        status\n        members {\n          id\n          authID\n        }\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n  permissionsPolicies(first: 500) {\n    edges {\n      node {\n        id\n        name\n        description\n        isGlobal\n        policy {\n          __typename\n          ... on InventoryPolicy {\n            __typename\n            read {\n              isAllowed\n            }\n          }\n          ... on WorkforcePolicy {\n            __typename\n            read {\n              isAllowed\n            }\n          }\n        }\n        groups {\n          id\n        }\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "users"
          ]
        },
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "usersGroups"
          ]
        },
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "permissionsPolicies"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '065678b56d20e6737426839b9bc05040';
module.exports = node;
