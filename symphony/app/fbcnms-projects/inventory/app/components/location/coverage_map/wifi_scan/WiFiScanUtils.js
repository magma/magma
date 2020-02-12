/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local strict-local
 * @format
 */

'use strict';

import type {CoordsWithProps} from '../../../map/MapUtil';
import type {FilterValue} from '../../../comparison_view/ComparisonViewTypes';
import type {LocationWiFiScanCoverageMap_wifiData} from './__generated__/LocationWiFiScanCoverageMap_wifiData.graphql.js';
import type {MapLayer, MapLayerStyles} from '../../../map/MapView';
import type {
  WiFiScanCollection,
  WiFiScanIndex,
} from './LocationWiFiScanCoverageMap';

import * as colors from '@fbcnms/ui/theme/colors';
import {WIFI_SCAN_BAND_FILTER} from './LocationWiFiScanSearchConfig';

import {coordsToGeoJSONSource} from '../../../map/MapUtil';
import {generateLatLngKey} from '../../coverage_map/CoverageMapUtils';
import {groupBy} from 'lodash';

// Filter wifi query response based on applied filters.
function cleanWiFiData(
  wifiData: LocationWiFiScanCoverageMap_wifiData,
  filters: Array<FilterValue>,
): WiFiScanCollection {
  const filteredWiFiData = wifiData
    .filter(wifi => wifi.latitude != null && wifi.longitude != null)
    .filter(wifi => wifiScanPassesFilters(wifi, filters));
  return {
    data: filteredWiFiData,
    bands: Object.keys(
      groupBy(
        filteredWiFiData.filter(wifi => wifi.band != null),
        wifi => wifi.band || '',
      ),
    ),
  };
}

// Aggregate wifi scan at same locations, indexed by latitude and longitude.
const aggregateWiFiScan = (
  wifiData: LocationWiFiScanCoverageMap_wifiData,
): WiFiScanIndex => {
  const groupedWiFiData = groupBy(wifiData, wifi => {
    return generateLatLngKey(wifi.latitude || 0, wifi.longitude || 0);
  });

  const wifiDataIndex: WiFiScanIndex = {};
  Object.keys(groupedWiFiData).forEach(key => {
    const wifis: LocationWiFiScanCoverageMap_wifiData = groupedWiFiData[key];
    const wifi = wifis.reduce((aggregatedWiFi, currentWiFi) => {
      return aggregatedWiFi.strength == null ||
        (currentWiFi.strength != null &&
          currentWiFi.strength > aggregatedWiFi.strength)
        ? currentWiFi
        : aggregatedWiFi;
    });
    if (wifi.strength != null) {
      wifiDataIndex[key] = {
        latitude: wifi.latitude || 0,
        longitude: wifi.longitude || 0,
        strength: wifi.strength,
        wifis: wifis,
      };
    }
  });

  return wifiDataIndex;
};

// Build map layer from aggregated wifi scans.
const wifiScanIndexToLayer = (
  sourceKey: string,
  aggregatedWiFiDataIndex: WiFiScanIndex,
  styles: MapLayerStyles,
): MapLayer => {
  const coordsWithProps: Array<CoordsWithProps> = Object.keys(
    aggregatedWiFiDataIndex,
  ).map(id => {
    const wifis = aggregatedWiFiDataIndex[id];
    return {
      latitude: wifis.latitude,
      longitude: wifis.longitude,
      properties: {
        color: colors.black,
        id: id,
        strength: wifis.strength,
      },
    };
  });

  return {
    source: coordsToGeoJSONSource(sourceKey, coordsWithProps),
    styles: styles,
  };
};

function wifiScanPassesFilters(wifi, filters: Array<FilterValue>): boolean {
  return (
    filters.filter(filter => wifiScanPassesFilter(wifi, filter)).length ===
    filters.length
  );
}

function wifiScanPassesFilter(wifi, filter: FilterValue): boolean {
  if (filter.key === WIFI_SCAN_BAND_FILTER) {
    if (filter.idSet == null) {
      return true;
    }
    return filter.idSet.filter(band => band === wifi.band).length > 0;
  }
  return true;
}

export {cleanWiFiData, aggregateWiFiScan, wifiScanIndexToLayer};
