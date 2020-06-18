/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 8bf8652045bc85dff716d1b1774fa0f1
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type SurveyQuestionType = "BOOL" | "CELLULAR" | "COORDS" | "DATE" | "EMAIL" | "FLOAT" | "INTEGER" | "PHONE" | "PHOTO" | "TEXT" | "TEXTAREA" | "WIFI" | "%future added value";
export type SurveyTemplateCategoryInput = {|
  id?: ?string,
  categoryTitle: string,
  categoryDescription: string,
  surveyTemplateQuestions?: ?$ReadOnlyArray<?SurveyTemplateQuestionInput>,
|};
export type SurveyTemplateQuestionInput = {|
  id?: ?string,
  questionTitle: string,
  questionDescription: string,
  questionType: SurveyQuestionType,
  index: number,
|};
export type EditLocationTypeSurveyTemplateCategoriesMutationVariables = {|
  id: string,
  surveyTemplateCategories: $ReadOnlyArray<SurveyTemplateCategoryInput>,
|};
export type EditLocationTypeSurveyTemplateCategoriesMutationResponse = {|
  +editLocationTypeSurveyTemplateCategories: ?$ReadOnlyArray<{|
    +id: string
  |}>
|};
export type EditLocationTypeSurveyTemplateCategoriesMutation = {|
  variables: EditLocationTypeSurveyTemplateCategoriesMutationVariables,
  response: EditLocationTypeSurveyTemplateCategoriesMutationResponse,
|};
*/


/*
mutation EditLocationTypeSurveyTemplateCategoriesMutation(
  $id: ID!
  $surveyTemplateCategories: [SurveyTemplateCategoryInput!]!
) {
  editLocationTypeSurveyTemplateCategories(id: $id, surveyTemplateCategories: $surveyTemplateCategories) {
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "surveyTemplateCategories",
    "type": "[SurveyTemplateCategoryInput!]!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "editLocationTypeSurveyTemplateCategories",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      },
      {
        "kind": "Variable",
        "name": "surveyTemplateCategories",
        "variableName": "surveyTemplateCategories"
      }
    ],
    "concreteType": "SurveyTemplateCategory",
    "plural": true,
    "selections": [
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "id",
        "args": null,
        "storageKey": null
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditLocationTypeSurveyTemplateCategoriesMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "EditLocationTypeSurveyTemplateCategoriesMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v1/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditLocationTypeSurveyTemplateCategoriesMutation",
    "id": null,
    "text": "mutation EditLocationTypeSurveyTemplateCategoriesMutation(\n  $id: ID!\n  $surveyTemplateCategories: [SurveyTemplateCategoryInput!]!\n) {\n  editLocationTypeSurveyTemplateCategories(id: $id, surveyTemplateCategories: $surveyTemplateCategories) {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '85792e6b61c311d8ce3645703deb69b6';
module.exports = node;
