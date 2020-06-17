/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
export type UsersGroupStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type UserManagementUtils_group$ref: FragmentReference;
declare export opaque type UserManagementUtils_group$fragmentType: UserManagementUtils_group$ref;
export type UserManagementUtils_group = {|
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
  +$refType: UserManagementUtils_group$ref,
|};
export type UserManagementUtils_group$data = UserManagementUtils_group;
export type UserManagementUtils_group$key = {
  +$data?: UserManagementUtils_group$data,
  +$fragmentRefs: UserManagementUtils_group$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
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
  "name": "status",
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
];
return {
  "kind": "Fragment",
  "name": "UserManagementUtils_group",
  "type": "UsersGroup",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    (v2/*: any*/),
    (v3/*: any*/),
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
        (v3/*: any*/),
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
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'aba8deb9bda0c2bceff11a5db4a576f8';
module.exports = node;
