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
type SiteSurveyQuestionReplyCellData_data$ref = any;
type SiteSurveyQuestionReplyWifiData_data$ref = any;
export type SurveyQuestionType = "BOOL" | "CELLULAR" | "COORDS" | "DATE" | "EMAIL" | "FLOAT" | "INTEGER" | "PHONE" | "PHOTO" | "TEXT" | "TEXTAREA" | "WIFI" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type SiteSurveyQuestionReply_question$ref: FragmentReference;
declare export opaque type SiteSurveyQuestionReply_question$fragmentType: SiteSurveyQuestionReply_question$ref;
export type SiteSurveyQuestionReply_question = {|
  +questionFormat: ?SurveyQuestionType,
  +longitude: ?number,
  +latitude: ?number,
  +boolData: ?boolean,
  +textData: ?string,
  +emailData: ?string,
  +phoneData: ?string,
  +floatData: ?number,
  +intData: ?number,
  +dateData: ?number,
  +photoData: ?{|
    +storeKey: ?string
  |},
  +$fragmentRefs: SiteSurveyQuestionReplyWifiData_data$ref & SiteSurveyQuestionReplyCellData_data$ref,
  +$refType: SiteSurveyQuestionReply_question$ref,
|};
export type SiteSurveyQuestionReply_question$data = SiteSurveyQuestionReply_question;
export type SiteSurveyQuestionReply_question$key = {
  +$data?: SiteSurveyQuestionReply_question$data,
  +$fragmentRefs: SiteSurveyQuestionReply_question$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "SiteSurveyQuestionReply_question",
  "type": "SurveyQuestion",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "questionFormat",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "longitude",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "latitude",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "boolData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "textData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "emailData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "phoneData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "floatData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "intData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "dateData",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "photoData",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": false,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "storeKey",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "SiteSurveyQuestionReplyWifiData_data",
      "args": null
    },
    {
      "kind": "FragmentSpread",
      "name": "SiteSurveyQuestionReplyCellData_data",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '4d9d45828bc4b9296fdd9b8347017e18';
module.exports = node;
