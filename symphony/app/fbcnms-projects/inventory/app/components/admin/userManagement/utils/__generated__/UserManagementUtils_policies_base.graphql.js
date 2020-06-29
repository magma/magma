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
import type { FragmentReference } from "relay-runtime";
declare export opaque type UserManagementUtils_policies_base$ref: FragmentReference;
declare export opaque type UserManagementUtils_policies_base$fragmentType: UserManagementUtils_policies_base$ref;
export type UserManagementUtils_policies_base = {|
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
  +$refType: UserManagementUtils_policies_base$ref,
|};
export type UserManagementUtils_policies_base$data = UserManagementUtils_policies_base;
export type UserManagementUtils_policies_base$key = {
  +$data?: UserManagementUtils_policies_base$data,
  +$fragmentRefs: UserManagementUtils_policies_base$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isAllowed",
  "args": null,
  "storageKey": null
},
v1 = [
  (v0/*: any*/)
],
v2 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  }
];
return {
  "kind": "Fragment",
  "name": "UserManagementUtils_policies_base",
  "type": "PermissionsPolicy",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "id",
      "args": null,
      "storageKey": null
    },
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
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "read",
              "storageKey": null,
              "args": null,
              "concreteType": "BasicPermissionRule",
              "plural": false,
              "selections": (v1/*: any*/)
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
                  "selections": (v1/*: any*/)
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
                    (v0/*: any*/),
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
                  "selections": (v1/*: any*/)
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
              "selections": (v2/*: any*/)
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "equipmentType",
              "storageKey": null,
              "args": null,
              "concreteType": "CUD",
              "plural": false,
              "selections": (v2/*: any*/)
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "locationType",
              "storageKey": null,
              "args": null,
              "concreteType": "CUD",
              "plural": false,
              "selections": (v2/*: any*/)
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "portType",
              "storageKey": null,
              "args": null,
              "concreteType": "CUD",
              "plural": false,
              "selections": (v2/*: any*/)
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "serviceType",
              "storageKey": null,
              "args": null,
              "concreteType": "CUD",
              "plural": false,
              "selections": (v2/*: any*/)
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
                (v0/*: any*/),
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
              "selections": (v2/*: any*/)
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
                  "selections": (v1/*: any*/)
                },
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "update",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "WorkforcePermissionRule",
                  "plural": false,
                  "selections": (v1/*: any*/)
                },
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "delete",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "WorkforcePermissionRule",
                  "plural": false,
                  "selections": (v1/*: any*/)
                },
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "assign",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "WorkforcePermissionRule",
                  "plural": false,
                  "selections": (v1/*: any*/)
                },
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "transferOwnership",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "WorkforcePermissionRule",
                  "plural": false,
                  "selections": (v1/*: any*/)
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
(node/*: any*/).hash = '7a02ff1e3fdafba9f4043d5321fb0ff4';
module.exports = node;
