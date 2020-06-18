/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash f1055cf578ddef96d81f013206edc0cb
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type NetworkServiceInput = {|
  ExternalID?: ?number,
  Customer: string,
  Name: string,
  Model: NetworkServiceModelInput,
  DeviceSites: $ReadOnlyArray<SiteInput>,
  AdditionalParams?: ?any,
|};
export type NetworkServiceModelInput = {|
  Name: string
|};
export type SiteInput = {|
  siteNumber: number,
  siteModelName: string,
  deviceName: string,
  deviceId: number,
  parameters: any,
  userPort: UserPortInput,
  accessMethod: string,
|};
export type UserPortInput = {|
  id: number,
  name: string,
|};
export type CreateServiceMutationVariables = {|
  nsi: NetworkServiceInput
|};
export type CreateServiceMutationResponse = {|
  +createService: string
|};
export type CreateServiceMutation = {|
  variables: CreateServiceMutationVariables,
  response: CreateServiceMutationResponse,
|};
*/


/*
mutation CreateServiceMutation(
  $nsi: NetworkServiceInput!
) {
  createService(input: $nsi)
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "nsi",
    "type": "NetworkServiceInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "createService",
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "nsi"
      }
    ],
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "CreateServiceMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "CreateServiceMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "CreateServiceMutation",
    "id": null,
    "text": "mutation CreateServiceMutation(\n  $nsi: NetworkServiceInput!\n) {\n  createService(input: $nsi)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e2ad8acf01814141fa4911f024690c1a';
module.exports = node;
