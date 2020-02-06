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
import type { FragmentReference } from "relay-runtime";
declare export opaque type SiteSurveyQuestionReplyWifiData_data$ref: FragmentReference;
declare export opaque type SiteSurveyQuestionReplyWifiData_data$fragmentType: SiteSurveyQuestionReplyWifiData_data$ref;
export type SiteSurveyQuestionReplyWifiData_data = {|
  +wifiData: ?$ReadOnlyArray<?{|
    +band: ?string,
    +bssid: string,
    +channel: number,
    +frequency: number,
    +strength: number,
    +ssid: ?string,
  |}>,
  +$refType: SiteSurveyQuestionReplyWifiData_data$ref,
|};
export type SiteSurveyQuestionReplyWifiData_data$data = SiteSurveyQuestionReplyWifiData_data;
export type SiteSurveyQuestionReplyWifiData_data$key = {
  +$data?: SiteSurveyQuestionReplyWifiData_data$data,
  +$fragmentRefs: SiteSurveyQuestionReplyWifiData_data$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "SiteSurveyQuestionReplyWifiData_data",
  "type": "SurveyQuestion",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "wifiData",
      "storageKey": null,
      "args": null,
      "concreteType": "SurveyWiFiScan",
      "plural": true,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "band",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "bssid",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "channel",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "frequency",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "strength",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "ssid",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '2d96d30d887d78992884a1de2ceca46c';
module.exports = node;
