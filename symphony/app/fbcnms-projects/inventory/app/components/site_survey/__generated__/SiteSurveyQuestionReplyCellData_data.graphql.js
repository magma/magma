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
export type CellularNetworkType = "CDMA" | "GSM" | "LTE" | "WCDMA" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type SiteSurveyQuestionReplyCellData_data$ref: FragmentReference;
declare export opaque type SiteSurveyQuestionReplyCellData_data$fragmentType: SiteSurveyQuestionReplyCellData_data$ref;
export type SiteSurveyQuestionReplyCellData_data = {|
  +cellData: ?$ReadOnlyArray<?{|
    +networkType: CellularNetworkType,
    +signalStrength: number,
    +baseStationID: ?string,
    +cellID: ?string,
    +locationAreaCode: ?string,
    +mobileCountryCode: ?string,
    +mobileNetworkCode: ?string,
  |}>,
  +$refType: SiteSurveyQuestionReplyCellData_data$ref,
|};
export type SiteSurveyQuestionReplyCellData_data$data = SiteSurveyQuestionReplyCellData_data;
export type SiteSurveyQuestionReplyCellData_data$key = {
  +$data?: SiteSurveyQuestionReplyCellData_data$data,
  +$fragmentRefs: SiteSurveyQuestionReplyCellData_data$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "SiteSurveyQuestionReplyCellData_data",
  "type": "SurveyQuestion",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "cellData",
      "storageKey": null,
      "args": null,
      "concreteType": "SurveyCellScan",
      "plural": true,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "networkType",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "signalStrength",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "baseStationID",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "cellID",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "locationAreaCode",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "mobileCountryCode",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "mobileNetworkCode",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '586ccb85b3fc3c5e37ce3014c2a26c0b';
module.exports = node;
