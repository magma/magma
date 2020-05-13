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
declare export opaque type UserManagementUtils_workforcePolicy$ref: FragmentReference;
declare export opaque type UserManagementUtils_workforcePolicy$fragmentType: UserManagementUtils_workforcePolicy$ref;
export type UserManagementUtils_workforcePolicy = {|
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
  +$refType: UserManagementUtils_workforcePolicy$ref,
|};
export type UserManagementUtils_workforcePolicy$data = UserManagementUtils_workforcePolicy;
export type UserManagementUtils_workforcePolicy$key = {
  +$data?: UserManagementUtils_workforcePolicy$data,
  +$fragmentRefs: UserManagementUtils_workforcePolicy$ref,
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
];
return {
  "kind": "Fragment",
  "name": "UserManagementUtils_workforcePolicy",
  "type": "WorkforcePolicy",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "read",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkforcePermissionRule",
      "plural": false,
      "selections": (v0/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "templates",
      "storageKey": null,
      "args": null,
      "concreteType": "CUD",
      "plural": false,
      "selections": [
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
      ]
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
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "update",
          "storageKey": null,
          "args": null,
          "concreteType": "WorkforcePermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "delete",
          "storageKey": null,
          "args": null,
          "concreteType": "WorkforcePermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "assign",
          "storageKey": null,
          "args": null,
          "concreteType": "WorkforcePermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "transferOwnership",
          "storageKey": null,
          "args": null,
          "concreteType": "WorkforcePermissionRule",
          "plural": false,
          "selections": (v0/*: any*/)
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '1fe6ffd822ae78d075bfdeae4a46f59a';
module.exports = node;
