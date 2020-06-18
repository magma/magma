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

import type {ColorStop, MapLayerStyles} from '../map/MapView';
import type {Location} from '../../common/Location.js';
import type {LocationCoverageMapTabQueryResponse} from './__generated__/LocationCoverageMapTabQuery.graphql.js';

import * as colors from '@fbcnms/ui/theme/colors';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import FormControl from '@material-ui/core/FormControl';
import Input from '@material-ui/core/Input';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import LocationCellScanCoverageMap from './coverage_map/cell_scan/LocationCellScanCoverageMap';
import LocationWiFiScanCoverageMap from './coverage_map/wifi_scan/LocationWiFiScanCoverageMap';
import MapColorSchemeLegend from '../map/MapColorSchemeLegend';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useState} from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import amber from '@material-ui/core/colors/amber';
import green from '@material-ui/core/colors/green';
import lime from '@material-ui/core/colors/lime';

import {graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

type Props = {
  location: Location,
};

// signal strength(dBm) range: https://powerfulsignal.com/cell-signal-strength
const CELL_SIGNAL_STRENGTH_THRESHOLD_LOW = -120;
const CELL_SIGNAL_STRENGTH_THRESHOLD_HIGH = -60;

const WIFI_SIGNAL_STRENGTH_THRESHOLD_LOW = 1;
const WIFI_SIGNAL_STRENGTH_THRESHOLD_HIGH = 5;

const HEATMAP_DENSITY_THRESHOLD_LOW = 0.1;
const HEATMAP_DENSITY_THRESHOLD_HIGH = 0.9;

const MAP_CELL = 'cell';
const MAP_WIFI = 'wifi';

const signalStrengthColors = [
  colors.red,
  amber[400],
  lime[200],
  green[200],
  green[600],
];

const cellCircleLayerColorStops = generateColorStops(
  CELL_SIGNAL_STRENGTH_THRESHOLD_LOW,
  CELL_SIGNAL_STRENGTH_THRESHOLD_HIGH,
  signalStrengthColors,
);

const cellCircleLayerStyles: MapLayerStyles = {
  circle: {
    colorInterpolation: {
      property: 'signalStrength',
      type: 'exponential',
      stops: cellCircleLayerColorStops,
    },
    fadeInZoomLevel: 11,
  },
};

const heatmapLayerColorStops = [
  {
    threshold: 0,
    color: 'transparent',
  },
].concat(
  generateColorStops(
    HEATMAP_DENSITY_THRESHOLD_LOW,
    HEATMAP_DENSITY_THRESHOLD_HIGH,
    signalStrengthColors,
  ),
);

const cellHeatmapLayerStyles: MapLayerStyles = {
  heatmap: {
    weight: {
      property: 'signalStrength',
      type: 'exponential',
      stops: [
        {
          threshold: CELL_SIGNAL_STRENGTH_THRESHOLD_LOW,
          weight: 0,
        },
        {
          threshold: CELL_SIGNAL_STRENGTH_THRESHOLD_HIGH,
          weight: 1,
        },
      ],
    },
    colorStops: heatmapLayerColorStops,
    fadeOutZoomLevel: 11,
  },
};

const wifiCircleLayerColorStops = generateColorStops(
  WIFI_SIGNAL_STRENGTH_THRESHOLD_LOW,
  WIFI_SIGNAL_STRENGTH_THRESHOLD_HIGH,
  signalStrengthColors,
);

const wifiCircleLayerStyles: MapLayerStyles = {
  circle: {
    colorInterpolation: {
      property: 'strength',
      type: 'exponential',
      stops: wifiCircleLayerColorStops,
    },
  },
};

const locationScansQuery = graphql`
  query LocationCoverageMapTabQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        cellData {
          ...LocationCellScanCoverageMap_cellData
        }
        wifiData {
          ...LocationWiFiScanCoverageMap_wifiData
        }
      }
    }
  }
`;

const LocationCoverageMapTab = (props: Props) => {
  const {location} = props;
  const classes = useStyles();
  const [selectedMap, setSelectedMap] = useState(MAP_CELL);

  function getSelector() {
    return (
      <FormControl className={classes.selector}>
        <Select
          value={selectedMap}
          onChange={event => {
            setSelectedMap(event.target.value);
          }}
          input={
            <Input className={classes.selectorInput} disableUnderline={true} />
          }>
          <MenuItem value={MAP_CELL}>Cell</MenuItem>
          <MenuItem value={MAP_WIFI}>Wifi</MenuItem>
        </Select>
      </FormControl>
    );
  }

  return (
    <InventoryQueryRenderer
      query={locationScansQuery}
      variables={{locationId: location.id}}
      render={(props: LocationCoverageMapTabQueryResponse) => {
        if (props.location == null) {
          return (
            <Text variant="body2" className={classes.note}>
              Location not found
            </Text>
          );
        }
        return (
          <Card className={classes.root}>
            {selectedMap === MAP_WIFI ? (
              <LocationWiFiScanCoverageMap
                location={location}
                wifiData={props.location.wifiData}
                circleLayerStyles={wifiCircleLayerStyles}
                selector={getSelector()}
                legend={
                  <MapColorSchemeLegend
                    title="Signal Strength Level"
                    colorScheme={{
                      lowLabel: WIFI_SIGNAL_STRENGTH_THRESHOLD_LOW.toString(),
                      highLabel: WIFI_SIGNAL_STRENGTH_THRESHOLD_HIGH.toString(),
                      colors: signalStrengthColors,
                    }}
                  />
                }
              />
            ) : (
              <LocationCellScanCoverageMap
                location={location}
                cellData={props.location.cellData}
                circleLayerStyles={cellCircleLayerStyles}
                heatmapLayerStyles={cellHeatmapLayerStyles}
                selector={getSelector()}
                legend={
                  <MapColorSchemeLegend
                    title="Signal Strength (dBm)"
                    colorScheme={{
                      lowLabel: CELL_SIGNAL_STRENGTH_THRESHOLD_LOW.toString(),
                      highLabel: CELL_SIGNAL_STRENGTH_THRESHOLD_HIGH.toString(),
                      colors: signalStrengthColors,
                    }}
                  />
                }
              />
            )}
          </Card>
        );
      }}
    />
  );
};

function generateColorStops(
  thresholdLow: number,
  thresholdHigh: number,
  colorScale: Array<string>,
): Array<ColorStop> {
  const step = (thresholdHigh - thresholdLow) / (colorScale.length - 1);
  return colorScale.map((color, index) => ({
    threshold: thresholdLow + step * index,
    color: color,
  }));
}

const useStyles = makeStyles(() => ({
  root: {
    height: '100%',
  },
  selector: {
    marginRight: '20px',
  },
  selectorInput: {
    fontWeight: 'bold',
  },
  note: {
    width: '100%',
    textAlign: 'center',
  },
}));

export default LocationCoverageMapTab;
