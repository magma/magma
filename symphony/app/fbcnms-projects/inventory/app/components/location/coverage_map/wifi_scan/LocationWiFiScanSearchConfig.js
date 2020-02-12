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

import type {
  EntityConfig,
  FilterConfig,
} from '../../../comparison_view/ComparisonViewTypes';

import LocationWiFiScanBandFilter from './LocationWiFiScanBandFilter';

const WIFI_SCAN_BAND_FILTER = 'wifi_scan_band';

const LocationWiFiScanSearchConfig: Array<EntityConfig> = [
  {
    type: 'wifi_scan',
    label: 'WiFi Properties',
    filters: [],
  },
];

function buildWiFiScanWiFiPropertiesFilterConfigs(
  bands: Array<string>,
): Array<FilterConfig> {
  return [
    {
      key: WIFI_SCAN_BAND_FILTER,
      name: WIFI_SCAN_BAND_FILTER,
      entityType: 'wifi_scan',
      label: 'Band',
      component: LocationWiFiScanBandFilter,
      defaultOperator: 'is',
      extraData: {
        bands,
      },
    },
  ];
}

export {
  LocationWiFiScanSearchConfig,
  buildWiFiScanWiFiPropertiesFilterConfigs,
  WIFI_SCAN_BAND_FILTER,
};
