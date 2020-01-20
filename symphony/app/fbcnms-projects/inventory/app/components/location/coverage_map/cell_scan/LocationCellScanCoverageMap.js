/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

'use strict';

import type {Location} from '../../../../common/Location.js';
import type {LocationCellScanCoverageMap_cellData} from './__generated__/LocationCellScanCoverageMap_cellData.graphql.js';
import type {MapLayerStyles} from '../../../map/MapView';

import * as React from 'react';
import CellScanPopout from '../../../map/CellScanPopout';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import LocationCellScansPane from './LocationCellScansPane';
import MapView from '../../../map/MapView';
import PowerSearchBar from '../../../power_search/PowerSearchBar';
import shortid from 'shortid';
import {
  LocationCellScanSearchConfig,
  buildCellScanCellPropertiesFilterConfigs,
} from './LocationCellScanSearchConfig';
import {useMemo, useState} from 'react';

import {
  addCellScanSpread,
  aggregateCellScan,
  cellScanIndexToLayer,
  cleanCellData,
} from './CellScanUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {getInitialFilterValue} from '../../../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

export type AggregatedCellScan = {
  latitude: number,
  longitude: number,
  signalStrength: number,
  cells: LocationCellScanCoverageMap_cellData,
};

export type CellScanCollection = {
  data: LocationCellScanCoverageMap_cellData,
  networkTypes: Array<string>,
};

export type CellScanIndex = {[string]: AggregatedCellScan};

type Props = {
  location: Location,
  cellData: LocationCellScanCoverageMap_cellData,
  circleLayerStyles: MapLayerStyles,
  heatmapLayerStyles: MapLayerStyles,
  selector?: React.Node,
  legend?: React.Node,
};

// size of grid where cell scan data is aggregated for heatmap
const HEATMAP_LAT_LNG_GRID_SIZE = 0.01;

const LocationCellScanCoverageMap = (props: Props) => {
  const classes = useStyles();
  const {
    location,
    circleLayerStyles,
    heatmapLayerStyles,
    selector,
    legend,
  } = props;
  const [selectedCellScans, setSelectedCellScans] = useState(null);
  const [filters, setFilters] = useState([]);

  const cellScanDataCollection = useMemo(
    () => cleanCellData(props.cellData, filters),
    [filters, props.cellData],
  );
  const cellData = useMemo(() => cellScanDataCollection.data, [
    cellScanDataCollection,
  ]);
  const filterConfigs = useMemo(
    () =>
      buildCellScanCellPropertiesFilterConfigs(
        cellScanDataCollection.networkTypes,
      ),
    [cellScanDataCollection],
  );

  const aggregatedCellDataIndex = useMemo(
    () => aggregateCellScan(cellData, false),
    [cellData],
  );
  const gridCellDataIndex = useMemo(
    () =>
      addCellScanSpread(
        aggregateCellScan(cellData, true, HEATMAP_LAT_LNG_GRID_SIZE),
        HEATMAP_LAT_LNG_GRID_SIZE,
      ),
    [cellData],
  );

  const aggregatedCircleLayer = useMemo(
    () =>
      cellScanIndexToLayer(
        `${location.id}_circle_${shortid.generate()}`,
        aggregatedCellDataIndex,
        circleLayerStyles,
      ),
    [aggregatedCellDataIndex, circleLayerStyles, location.id],
  );
  const aggregatedHeatmapLayer = useMemo(
    () =>
      cellScanIndexToLayer(
        `${location.id}_heatmap_${shortid.generate()}`,
        gridCellDataIndex,
        heatmapLayerStyles,
      ),
    [gridCellDataIndex, heatmapLayerStyles, location.id],
  );
  const layers = useMemo(
    () => [aggregatedCircleLayer, aggregatedHeatmapLayer],
    [aggregatedCircleLayer, aggregatedHeatmapLayer],
  );

  const center = {
    lat: location.latitude,
    lng: location.longitude,
  };

  const getCellScanPopoutContent = feature => {
    const featureProps: ?{[string]: string} = feature.properties;
    const aggregatedCellData: ?AggregatedCellScan =
      featureProps == null || featureProps.id == null
        ? null
        : aggregatedCellDataIndex[featureProps.id];
    return (
      <CellScanPopout
        aggregatedCellData={aggregatedCellData}
        renderCellScansDialog={setSelectedCellScans}
      />
    );
  };

  return (
    <div className={classes.root}>
      {legend != null && <div className={classes.overlay}>{legend}</div>}
      <div className={classes.powerSearchContainer}>
        <PowerSearchBar
          placeholder="Filter cell scans"
          filterConfigs={filterConfigs}
          searchConfig={LocationCellScanSearchConfig}
          onFiltersChanged={setFilters}
          getSelectedFilter={filterConfig =>
            getInitialFilterValue(
              filterConfig.key,
              filterConfig.name,
              filterConfig.defaultOperator,
            )
          }
          header={selector}
          footer={`Sample Count: ${cellData.length}`}
        />
      </div>
      <div className={classes.mapContainer}>
        <MapView
          id="cellCoverageMap"
          mode={
            location.locationType.mapType === 'satellite'
              ? 'satellite'
              : 'streets'
          }
          center={center}
          zoomLevel={location.locationType.mapZoomLevel}
          layers={layers}
          getFeaturePopoutContent={getCellScanPopoutContent}
        />
      </div>
      {selectedCellScans && (
        <Dialog
          maxWidth="md"
          fullWidth={true}
          onClose={() => setSelectedCellScans(null)}
          open={true}>
          <DialogContent>
            <LocationCellScansPane
              latitude={selectedCellScans.latitude}
              longitude={selectedCellScans.longitude}
              cellData={selectedCellScans.cells}
            />
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
};

const useStyles = makeStyles(theme => ({
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  overlay: {
    position: 'absolute',
    bottom: '30px',
    right: '20px',
    zIndex: 2,
  },
  powerSearchContainer: {
    margin: '10px',
    backgroundColor: theme.palette.background.paper,
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
  mapContainer: {
    flexGrow: 1,
  },
}));

export default createFragmentContainer(LocationCellScanCoverageMap, {
  cellData: graphql`
    fragment LocationCellScanCoverageMap_cellData on SurveyCellScan
      @relay(plural: true) {
      id
      latitude
      longitude
      networkType
      signalStrength
      mobileCountryCode
      mobileNetworkCode
      operator
    }
  `,
});
