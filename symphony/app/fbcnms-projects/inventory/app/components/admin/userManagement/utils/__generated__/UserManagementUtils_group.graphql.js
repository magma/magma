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
};
return {
  "kind": "Fragment",
  "name": "UserManagementUtils_group",
  "type": "UsersGroup",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
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
      "name": "status",
      "args": null,
      "storageKey": null
    },
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
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e6d2e87628742e930ae12e88ae7d4566';
module.exports = node;
