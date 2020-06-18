/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash a6d0a187f66037c90225ff738c3f494f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type PermissionsPoliciesViewQueryVariables = {||};
export type PermissionsPoliciesViewQueryResponse = {|
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
          +id: string,
          +name: string,
          +description: ?string,
          +status: UsersGroupStatus,
        |}>,
      |}
    |}>
  |}
|};
export type PermissionsPoliciesViewQuery = {|
  variables: PermissionsPoliciesViewQueryVariables,
  response: PermissionsPoliciesViewQueryResponse,
|};
*/


/*
query PermissionsPoliciesViewQuery {
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
          name
          description
          status
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
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v4 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isAllowed",
    "args": null,
    "storageKey": null
  }
],
v5 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v4/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v4/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v4/*: any*/)
  }
],
v6 = [
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
          (v1/*: any*/),
          (v2/*: any*/),
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
              (v3/*: any*/),
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
                    "selections": (v4/*: any*/)
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
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "update",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationPermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "delete",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "LocationPermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
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
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "serviceType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v5/*: any*/)
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
                    "selections": (v4/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "templates",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "CUD",
                    "plural": false,
                    "selections": (v5/*: any*/)
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
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "update",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "delete",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "assign",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
                      },
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "transferOwnership",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "WorkforcePermissionRule",
                        "plural": false,
                        "selections": (v4/*: any*/)
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
              (v0/*: any*/),
              (v1/*: any*/),
              (v2/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "status",
                "args": null,
                "storageKey": null
              }
            ]
          },
          (v3/*: any*/)
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
v7 = [
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
    "name": "PermissionsPoliciesViewQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "permissionsPolicies",
        "name": "__PermissionsPoliciesView_permissionsPolicies_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "PermissionsPolicyConnection",
        "plural": false,
        "selections": (v6/*: any*/)
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "PermissionsPoliciesViewQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "permissionsPolicies",
        "storageKey": "permissionsPolicies(first:500)",
        "args": (v7/*: any*/),
        "concreteType": "PermissionsPolicyConnection",
        "plural": false,
        "selections": (v6/*: any*/)
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "permissionsPolicies",
        "args": (v7/*: any*/),
        "handle": "connection",
        "key": "PermissionsPoliciesView_permissionsPolicies",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "PermissionsPoliciesViewQuery",
    "id": null,
    "text": "query PermissionsPoliciesViewQuery {\n  permissionsPolicies(first: 500) {\n    edges {\n      node {\n        id\n        name\n        description\n        isGlobal\n        policy {\n          __typename\n          ... on InventoryPolicy {\n            read {\n              isAllowed\n            }\n            location {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            equipment {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            equipmentType {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            locationType {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            portType {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            serviceType {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n          }\n          ... on WorkforcePolicy {\n            read {\n              isAllowed\n            }\n            templates {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n            }\n            data {\n              create {\n                isAllowed\n              }\n              update {\n                isAllowed\n              }\n              delete {\n                isAllowed\n              }\n              assign {\n                isAllowed\n              }\n              transferOwnership {\n                isAllowed\n              }\n            }\n          }\n        }\n        groups {\n          id\n          name\n          description\n          status\n        }\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n",
    "metadata": {
      "connection": [
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
(node/*: any*/).hash = '8de464be5a9f9dd8d8dfa8b5d60e6e78';
module.exports = node;
