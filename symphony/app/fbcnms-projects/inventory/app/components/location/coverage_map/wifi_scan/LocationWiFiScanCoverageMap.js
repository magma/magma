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
import type {LocationWiFiScanCoverageMap_wifiData} from './__generated__/LocationWiFiScanCoverageMap_wifiData.graphql.js';
import type {MapLayerStyles} from '../../../map/MapView';

import * as React from 'react';
import MapView from '../../../map/MapView';
import PowerSearchBar from '../../../power_search/PowerSearchBar';
import shortid from 'shortid';
import {
  LocationWiFiScanSearchConfig,
  buildWiFiScanWiFiPropertiesFilterConfigs,
} from './LocationWiFiScanSearchConfig';
import {useMemo, useState} from 'react';

import {
  aggregateWiFiScan,
  cleanWiFiData,
  wifiScanIndexToLayer,
} from './WiFiScanUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {getInitialFilterValue} from '../../../comparison_view/FilterUtils';
import {makeStyles} from '@material-ui/styles';

export type WiFiScanCollection = {
  data: LocationWiFiScanCoverageMap_wifiData,
  bands: Array<string>,
};

type AggregatedWiFiScan = {
  latitude: number,
  longitude: number,
  strength: number,
  wifis: LocationWiFiScanCoverageMap_wifiData,
};

export type WiFiScanIndex = {[string]: AggregatedWiFiScan};

type Props = {
  location: Location,
  wifiData: LocationWiFiScanCoverageMap_wifiData,
  circleLayerStyles: MapLayerStyles,
  selector?: React.Node,
  legend?: React.Node,
};

const LocationWiFiScanCoverageMap = (props: Props) => {
  const {location, circleLayerStyles, selector, legend} = props;
  const classes = useStyles();

  const [filters, setFilters] = useState([]);

  const wifiDataCollection = useMemo(
    () => cleanWiFiData(props.wifiData, filters),
    [filters, props.wifiData],
  );
  const wifiData = useMemo(() => wifiDataCollection.data, [wifiDataCollection]);
  const filterConfigs = useMemo(
    () => buildWiFiScanWiFiPropertiesFilterConfigs(wifiDataCollection.bands),
    [wifiDataCollection],
  );

  const aggregatedWiFiDataIndex = useMemo(() => aggregateWiFiScan(wifiData), [
    wifiData,
  ]);

  const layers = useMemo(
    () => [
      wifiScanIndexToLayer(
        `${location.id}_circle_${shortid.generate()}`,
        aggregatedWiFiDataIndex,
        circleLayerStyles,
      ),
    ],
    [aggregatedWiFiDataIndex, circleLayerStyles, location.id],
  );

  const center = {
    lat: location.latitude,
    lng: location.longitude,
  };

  return (
    <div className={classes.root}>
      {legend != null && <div className={classes.overlay}>{legend}</div>}
      <div className={classes.powerSearchContainer}>
        <PowerSearchBar
          placeholder="Filter wifi scans"
          filterConfigs={filterConfigs}
          searchConfig={LocationWiFiScanSearchConfig}
          onFiltersChanged={setFilters}
          getSelectedFilter={filterConfig =>
            getInitialFilterValue(
              filterConfig.key,
              filterConfig.name,
              filterConfig.defaultOperator,
            )
          }
          header={selector}
          footer={`Sample Count: ${wifiData.length}`}
        />
      </div>
      <div className={classes.mapContainer}>
        <MapView
          id="wifiCoverageMap"
          mode={
            location.locationType.mapType === 'satellite'
              ? 'satellite'
              : 'streets'
          }
          center={center}
          zoomLevel={location.locationType.mapZoomLevel}
          layers={layers}
        />
      </div>
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

export default createFragmentContainer(LocationWiFiScanCoverageMap, {
  wifiData: graphql`
    fragment LocationWiFiScanCoverageMap_wifiData on SurveyWiFiScan
      @relay(plural: true) {
      id
      latitude
      longitude
      frequency
      channel
      bssid
      ssid
      strength
      band
    }
  `,
});
