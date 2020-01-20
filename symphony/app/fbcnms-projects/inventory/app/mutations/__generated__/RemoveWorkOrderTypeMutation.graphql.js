/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 477ffd930d739a116f9c92fd75e9df48
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveWorkOrderTypeMutationVariables = {|
  id: string
|};
export type RemoveWorkOrderTypeMutationResponse = {|
  +removeWorkOrderType: string
|};
export type RemoveWorkOrderTypeMutation = {|
  variables: RemoveWorkOrderTypeMutationVariables,
  response: RemoveWorkOrderTypeMutationResponse,
|};
*/


/*
mutation RemoveWorkOrderTypeMutation(
  $id: ID!
) {
  removeWorkOrderType(id: $id)
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
    "name": "removeWorkOrderType",
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
    "name": "RemoveWorkOrderTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveWorkOrderTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveWorkOrderTypeMutation",
    "id": null,
    "text": "mutation RemoveWorkOrderTypeMutation(\n  $id: ID!\n) {\n  removeWorkOrderType(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '6fec7c186557570f206537c12c607614';
module.exports = node;
