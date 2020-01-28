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
declare export opaque type LocationCellScanCoverageMap_cellData$ref: FragmentReference;
declare export opaque type LocationCellScanCoverageMap_cellData$fragmentType: LocationCellScanCoverageMap_cellData$ref;
export type LocationCellScanCoverageMap_cellData = $ReadOnlyArray<{|
  +id: string,
  +latitude: ?number,
  +longitude: ?number,
  +networkType: CellularNetworkType,
  +signalStrength: number,
  +mobileCountryCode: ?string,
  +mobileNetworkCode: ?string,
  +operator: ?string,
  +$refType: LocationCellScanCoverageMap_cellData$ref,
|}>;
export type LocationCellScanCoverageMap_cellData$data = LocationCellScanCoverageMap_cellData;
export type LocationCellScanCoverageMap_cellData$key = $ReadOnlyArray<{
  +$data?: LocationCellScanCoverageMap_cellData$data,
  +$fragmentRefs: LocationCellScanCoverageMap_cellData$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "LocationCellScanCoverageMap_cellData",
  "type": "SurveyCellScan",
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
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "operator",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'd94c9d40f7baef9bcd50963bcb149e0e';
module.exports = node;
