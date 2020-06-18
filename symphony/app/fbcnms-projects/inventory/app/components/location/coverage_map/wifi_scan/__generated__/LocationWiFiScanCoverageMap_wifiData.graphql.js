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
declare export opaque type LocationWiFiScanCoverageMap_wifiData$ref: FragmentReference;
declare export opaque type LocationWiFiScanCoverageMap_wifiData$fragmentType: LocationWiFiScanCoverageMap_wifiData$ref;
export type LocationWiFiScanCoverageMap_wifiData = $ReadOnlyArray<{|
  +id: string,
  +latitude: ?number,
  +longitude: ?number,
  +frequency: number,
  +channel: number,
  +bssid: string,
  +ssid: ?string,
  +strength: number,
  +band: ?string,
  +$refType: LocationWiFiScanCoverageMap_wifiData$ref,
|}>;
export type LocationWiFiScanCoverageMap_wifiData$data = LocationWiFiScanCoverageMap_wifiData;
export type LocationWiFiScanCoverageMap_wifiData$key = $ReadOnlyArray<{
  +$data?: LocationWiFiScanCoverageMap_wifiData$data,
  +$fragmentRefs: LocationWiFiScanCoverageMap_wifiData$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "LocationWiFiScanCoverageMap_wifiData",
  "type": "SurveyWiFiScan",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
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
      "name": "latitude",
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
      "name": "frequency",
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
      "name": "bssid",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "ssid",
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
      "name": "band",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'c19f1446b6147fde7dff384c07c86b60';
module.exports = node;
