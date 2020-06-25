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
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isAllowed",
  "args": null,
  "storageKey": null
},
v1 = [
  (v0/*: any*/)
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
      "selections": [
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
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd110f4f516c9551f09b4cd7f70c6849b';
module.exports = node;
