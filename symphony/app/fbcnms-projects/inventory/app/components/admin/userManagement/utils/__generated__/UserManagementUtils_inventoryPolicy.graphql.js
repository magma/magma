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
declare export opaque type UserManagementUtils_inventoryPolicy$ref: FragmentReference;
declare export opaque type UserManagementUtils_inventoryPolicy$fragmentType: UserManagementUtils_inventoryPolicy$ref;
export type UserManagementUtils_inventoryPolicy = {|
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
  +$refType: UserManagementUtils_inventoryPolicy$ref,
|};
export type UserManagementUtils_inventoryPolicy$data = UserManagementUtils_inventoryPolicy;
export type UserManagementUtils_inventoryPolicy$key = {
  +$data?: UserManagementUtils_inventoryPolicy$data,
  +$fragmentRefs: UserManagementUtils_inventoryPolicy$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isAllowed",
    "args": null,
    "storageKey": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v0/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v0/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v0/*: any*/)
  }
];
return {
  "kind": "Fragment",
  "name": "UserManagementUtils_inventoryPolicy",
  "type": "InventoryPolicy",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "read",
      "storageKey": null,
      "args": null,
      "concreteType": "BasicPermissionRule",
      "plural": false,
      "selections": (v0/*: any*/)
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
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "update",
          "storageKey": null,
          "args": null,
          "concreteType": "LocationPermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "delete",
          "storageKey": null,
          "args": null,
          "concreteType": "LocationPermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
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
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "CUD",
      "plural": false,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationType",
      "storageKey": null,
      "args": null,
      "concreteType": "CUD",
      "plural": false,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portType",
      "storageKey": null,
      "args": null,
      "concreteType": "CUD",
      "plural": false,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceType",
      "storageKey": null,
      "args": null,
      "concreteType": "CUD",
      "plural": false,
      "selections": (v1/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e50c986b2421ebe08a5be0715e650bc8';
module.exports = node;
