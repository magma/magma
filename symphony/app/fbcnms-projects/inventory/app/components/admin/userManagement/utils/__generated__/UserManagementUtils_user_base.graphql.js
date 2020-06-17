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
export type UserRole = "ADMIN" | "OWNER" | "USER" | "%future added value";
export type UserStatus = "ACTIVE" | "DEACTIVATED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type UserManagementUtils_user_base$ref: FragmentReference;
declare export opaque type UserManagementUtils_user_base$fragmentType: UserManagementUtils_user_base$ref;
export type UserManagementUtils_user_base = {|
  +id: string,
  +authID: string,
  +firstName: string,
  +lastName: string,
  +email: string,
  +status: UserStatus,
  +role: UserRole,
  +$refType: UserManagementUtils_user_base$ref,
|};
export type UserManagementUtils_user_base$data = UserManagementUtils_user_base;
export type UserManagementUtils_user_base$key = {
  +$data?: UserManagementUtils_user_base$data,
  +$fragmentRefs: UserManagementUtils_user_base$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "UserManagementUtils_user_base",
  "type": "User",
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
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "role",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '37042ebf2f46af736aac063cef60994e';
module.exports = node;
