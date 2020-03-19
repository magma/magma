/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 40f5938b7c62ab6435d905cbd586e396
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type DeleteReportFilterMutationVariables = {|
  id: string
|};
export type DeleteReportFilterMutationResponse = {|
  +deleteReportFilter: boolean
|};
export type DeleteReportFilterMutation = {|
  variables: DeleteReportFilterMutationVariables,
  response: DeleteReportFilterMutationResponse,
|};
*/


/*
mutation DeleteReportFilterMutation(
  $id: ID!
) {
  deleteReportFilter(id: $id)
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
    "name": "deleteReportFilter",
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
    "name": "DeleteReportFilterMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "DeleteReportFilterMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "DeleteReportFilterMutation",
    "id": null,
    "text": "mutation DeleteReportFilterMutation(\n  $id: ID!\n) {\n  deleteReportFilter(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7963d86165414f401b01981da338f8de';
module.exports = node;
