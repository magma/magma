/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash fed592c78e90ce375867fc0308a37f6c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveProjectMutationVariables = {|
  id: string
|};
export type RemoveProjectMutationResponse = {|
  +deleteProject: boolean
|};
export type RemoveProjectMutation = {|
  variables: RemoveProjectMutationVariables,
  response: RemoveProjectMutationResponse,
|};
*/


/*
mutation RemoveProjectMutation(
  $id: ID!
) {
  deleteProject(id: $id)
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
    "name": "deleteProject",
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
    "name": "RemoveProjectMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveProjectMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveProjectMutation",
    "id": null,
    "text": "mutation RemoveProjectMutation(\n  $id: ID!\n) {\n  deleteProject(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7424ca24008e15a923f89b3fbc822395';
module.exports = node;
