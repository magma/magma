/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';

import type {
  EntityConfig,
  FilterConfig,
} from '../../../comparison_view/ComparisonViewTypes';

import LocationCellScanMccFilter from './LocationCellScanMccFilter';
import LocationCellScanMncFilter from './LocationCellScanMncFilter';
import LocationCellScanNetworkFilter from './LocationCellScanNetworkFilter';

const CELL_SCAN_NETWORK_TYPE_FILTER = 'cell_scan_network_type';
const CELL_SCAN_MCC_FILTER = 'cell_scan_mcc';
const CELL_SCAN_MNC_FILTER = 'cell_scan_mnc';

const LocationCellScanSearchConfig: Array<EntityConfig> = [
  {
    type: 'cell_scan',
    label: 'Cell Properties',
    filters: [],
  },
];

function buildCellScanCellPropertiesFilterConfigs(
  networkTypes: Array<string>,
): Array<FilterConfig> {
  return [
    {
      key: CELL_SCAN_NETWORK_TYPE_FILTER,
      name: CELL_SCAN_NETWORK_TYPE_FILTER,
      entityType: 'cell_scan',
      label: 'Network',
      component: LocationCellScanNetworkFilter,
      defaultOperator: 'is',
      extraData: {
        networkTypes,
      },
    },
    {
      key: CELL_SCAN_MCC_FILTER,
      name: CELL_SCAN_MCC_FILTER,
      entityType: 'cell_scan',
      label: 'MCC',
      component: LocationCellScanMccFilter,
      defaultOperator: 'is',
    },
    {
      key: CELL_SCAN_MNC_FILTER,
      name: CELL_SCAN_MNC_FILTER,
      entityType: 'cell_scan',
      label: 'MNC',
      component: LocationCellScanMncFilter,
      defaultOperator: 'is',
    },
  ];
}

export {
  LocationCellScanSearchConfig,
  buildCellScanCellPropertiesFilterConfigs,
  CELL_SCAN_MCC_FILTER,
  CELL_SCAN_MNC_FILTER,
  CELL_SCAN_NETWORK_TYPE_FILTER,
};
