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
  CellScanCollection,
  CellScanIndex,
} from './LocationCellScanCoverageMap';
import type {CoordsWithProps} from '../../../map/MapUtil';
import type {FilterValue} from '../../../comparison_view/ComparisonViewTypes';
import type {LocationCellScanCoverageMap_cellData} from './__generated__/LocationCellScanCoverageMap_cellData.graphql.js';
import type {MapLayer, MapLayerStyles} from '../../../map/MapView';

import * as colors from '@fbcnms/ui/theme/colors';
import {
  CELL_SCAN_MCC_FILTER,
  CELL_SCAN_MNC_FILTER,
  CELL_SCAN_NETWORK_TYPE_FILTER,
} from './LocationCellScanSearchConfig';

import {coordsToGeoJSONSource} from '../../../map/MapUtil';
import {
  generateLatLngKey,
  getGridCoordinate,
} from '../../coverage_map/CoverageMapUtils';
import {groupBy, isEqual} from 'lodash';

// attenuation of cell signal strength in dBm
const CELL_SIGNAL_ATTENUATION = 20;

// Filter cell query response based on applied filters.
function cleanCellData(
  cellData: LocationCellScanCoverageMap_cellData,
  filters: Array<FilterValue>,
): CellScanCollection {
  const filteredCellData = cellData
    .filter(cell => cell.latitude != null && cell.longitude != null)
    .filter(cell => cellScanPassesFilters(cell, filters));
  return {
    data: filteredCellData,
    networkTypes: Object.keys(
      groupBy(filteredCellData, cell => cell.networkType),
    ),
  };
}

/* Generate cell scan data points around collected data points to balance out
 * the effect of number of samples on heatmap */
function addCellScanSpread(
  aggregatedIndex: CellScanIndex,
  granularity: number,
): CellScanIndex {
  const indexWithSpread = Object.assign({}, aggregatedIndex);
  Object.keys(indexWithSpread).forEach(key => {
    const {latitude, longitude, signalStrength} = indexWithSpread[key];
    const deltas = [-granularity, 0, granularity];
    deltas.forEach(latDelta => {
      deltas.forEach(lngDelta => {
        if (latDelta === 0 && lngDelta === 0) {
          return;
        }
        const newLatitude = latitude + latDelta;
        const newLongitude = longitude + lngDelta;
        const newKey = generateLatLngKey(newLatitude, newLongitude);
        if (newKey in indexWithSpread) {
          return;
        }
        indexWithSpread[newKey] = {
          latitude: newLatitude,
          longitude: newLongitude,
          signalStrength:
            signalStrength == null
              ? null
              : signalStrength - CELL_SIGNAL_ATTENUATION,
          cells: [],
        };
      });
    });
  });

  return indexWithSpread;
}

// Aggregate cell scan at same locations, indexed by latitude and longitude.
function aggregateCellScan(
  cellData: LocationCellScanCoverageMap_cellData,
  collapseToGrid: boolean,
  /* default precision is 5 decimal places
   * (individual trees, door entrances, ...)
   * https://en.wikipedia.org/wiki/Decimal_degrees */
  granularity: number = 0.00001,
): CellScanIndex {
  const groupedCellData = groupBy(cellData, cell => {
    const latitude = getGridCoordinate(
      cell.latitude || 0,
      collapseToGrid,
      granularity,
    );
    const longitude = getGridCoordinate(
      cell.longitude || 0,
      collapseToGrid,
      granularity,
    );
    return generateLatLngKey(latitude, longitude);
  });

  const cellDataIndex: CellScanIndex = {};
  Object.keys(groupedCellData).forEach(key => {
    const cells: LocationCellScanCoverageMap_cellData = groupedCellData[key];
    const cell = cells.reduce((aggregatedCell, currentCell) => {
      return aggregatedCell.signalStrength == null ||
        (currentCell.signalStrength != null &&
          currentCell.signalStrength > aggregatedCell.signalStrength)
        ? currentCell
        : aggregatedCell;
    });
    if (cell.signalStrength != null) {
      cellDataIndex[key] = {
        latitude: getGridCoordinate(
          cell.latitude || 0,
          collapseToGrid,
          granularity,
        ),
        longitude: getGridCoordinate(
          cell.longitude || 0,
          collapseToGrid,
          granularity,
        ),
        signalStrength: cell.signalStrength,
        cells: cells,
      };
    }
  });

  return cellDataIndex;
}

// Build map layer from aggregated cell scans.
function cellScanIndexToLayer(
  sourceKey: string,
  aggregatedCellDataIndex: CellScanIndex,
  styles: MapLayerStyles,
): MapLayer {
  const coordsWithProps: Array<CoordsWithProps> = Object.keys(
    aggregatedCellDataIndex,
  ).map(id => {
    const cells = aggregatedCellDataIndex[id];
    return {
      latitude: cells.latitude,
      longitude: cells.longitude,
      properties: {
        color: colors.black,
        id: id,
        signalStrength: cells.signalStrength,
      },
    };
  });

  return {
    source: coordsToGeoJSONSource(sourceKey, coordsWithProps),
    styles: styles,
  };
}

function cellScanPassesFilters(cell, filters: Array<FilterValue>): boolean {
  return (
    filters.filter(filter => cellScanPassesFilter(cell, filter)).length ===
    filters.length
  );
}

function cellScanPassesFilter(cell, filter: FilterValue): boolean {
  if (filter.key === CELL_SCAN_NETWORK_TYPE_FILTER) {
    if (filter.idSet == null) {
      return true;
    }
    return (
      filter.idSet.filter(networkType => networkType === cell.networkType)
        .length > 0
    );
  }
  if (filter.key === CELL_SCAN_MCC_FILTER) {
    if (filter.stringValue == null) {
      return true;
    }
    return isEqual(cell.mobileCountryCode, filter.stringValue);
  }
  if (filter.key === CELL_SCAN_MNC_FILTER) {
    if (filter.stringValue == null) {
      return true;
    }
    return isEqual(cell.mobileNetworkCode, filter.stringValue);
  }
  return true;
}

export {
  cleanCellData,
  aggregateCellScan,
  cellScanIndexToLayer,
  addCellScanSpread,
};
