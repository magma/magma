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
type SiteSurveyQuestionReply_question$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type SiteSurveyPane_survey$ref: FragmentReference;
declare export opaque type SiteSurveyPane_survey$fragmentType: SiteSurveyPane_survey$ref;
export type SiteSurveyPane_survey = {|
  +name: string,
  +completionTimestamp: number,
  +surveyResponses: $ReadOnlyArray<?{|
    +id: string,
    +questionText: string,
    +formName: ?string,
    +formIndex: number,
    +questionIndex: number,
    +$fragmentRefs: SiteSurveyQuestionReply_question$ref,
  |}>,
  +$refType: SiteSurveyPane_survey$ref,
|};
export type SiteSurveyPane_survey$data = SiteSurveyPane_survey;
export type SiteSurveyPane_survey$key = {
  +$data?: SiteSurveyPane_survey$data,
  +$fragmentRefs: SiteSurveyPane_survey$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "SiteSurveyPane_survey",
  "type": "Survey",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
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
      "name": "completionTimestamp",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "surveyResponses",
      "storageKey": null,
      "args": null,
      "concreteType": "SurveyQuestion",
      "plural": true,
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
          "name": "questionText",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "formName",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "formIndex",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "questionIndex",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "FragmentSpread",
          "name": "SiteSurveyQuestionReply_question",
          "args": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'abb3682eadcff198afba66898e28a2a0';
module.exports = node;
