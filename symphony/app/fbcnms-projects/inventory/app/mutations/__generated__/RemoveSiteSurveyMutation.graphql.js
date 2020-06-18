/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 6f6ebf8f00eb0e675c1ca70ea6b7f4b5
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type RemoveSiteSurveyMutationVariables = {|
  id: string
|};
export type RemoveSiteSurveyMutationResponse = {|
  +removeSiteSurvey: string
|};
export type RemoveSiteSurveyMutation = {|
  variables: RemoveSiteSurveyMutationVariables,
  response: RemoveSiteSurveyMutationResponse,
|};
*/


/*
mutation RemoveSiteSurveyMutation(
  $id: ID!
) {
  removeSiteSurvey(id: $id)
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
    "name": "removeSiteSurvey",
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
    "name": "RemoveSiteSurveyMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveSiteSurveyMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveSiteSurveyMutation",
    "id": null,
    "text": "mutation RemoveSiteSurveyMutation(\n  $id: ID!\n) {\n  removeSiteSurvey(id: $id)\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fc6117923fa8eda4898294028568180b';
module.exports = node;
