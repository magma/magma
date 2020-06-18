/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash d7fd158f964542694e7512b8d5bfba1c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DeletePermissionsPolicyMutationVariables = {|
  id: string
|};
export type DeletePermissionsPolicyMutationResponse = {|
  +deletePermissionsPolicy: boolean
|};
export type DeletePermissionsPolicyMutation = {|
  variables: DeletePermissionsPolicyMutationVariables,
  response: DeletePermissionsPolicyMutationResponse,
|};
*/


/*
mutation DeletePermissionsPolicyMutation(
  $id: ID!
) {
  deletePermissionsPolicy(id: $id)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "deletePermissionsPolicy",
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      }
    ],
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "DeletePermissionsPolicyMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "DeletePermissionsPolicyMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeletePermissionsPolicyMutation",
    "id": null,
    "text": "mutation DeletePermissionsPolicyMutation(\n  $id: ID!\n) {\n  deletePermissionsPolicy(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '14f797ed5eccde6161e1583e51596f09';
module.exports = node;
