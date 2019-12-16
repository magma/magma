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

import type {LocationCellScanCoverageMap_cellData} from './__generated__/LocationCellScanCoverageMap_cellData.graphql.js';

import * as React from 'react';
import PerfectScrollbar from 'react-perfect-scrollbar';
import Text from '@fbcnms/ui/components/design-system/Text';

import {makeStyles} from '@material-ui/styles';

type Props = {
  latitude: number,
  longitude: number,
  cellData: LocationCellScanCoverageMap_cellData,
};

const LocationCellScansPane = (props: Props) => {
  const classes = useStyles();
  const {latitude, longitude, cellData} = props;
  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <Text variant="h6">Cell Scans</Text>
        <Text variant="body2">
          {`latitude: ${latitude.toFixed(5)}, longitude: ${longitude.toFixed(
            5,
          )}`}
        </Text>
      </div>
      <div className={classes.body}>
        <div className={classes.row}>
          <Text variant="subtitle2" className={classes.col}>
            network
          </Text>
          <Text variant="subtitle2" className={classes.col}>
            strength(dBm)
          </Text>
          <Text variant="subtitle2" className={classes.col}>
            MCC
          </Text>
          <Text variant="subtitle2" className={classes.col}>
            MNC
          </Text>
        </div>
        <PerfectScrollbar>
          {cellData.map((cell, index) => (
            <div key={index} className={classes.row}>
              <Text variant="body2" className={classes.col}>
                {cell.networkType}
              </Text>
              <Text variant="body2" className={classes.col}>
                {cell.signalStrength == null ? 'unknown' : cell.signalStrength}
              </Text>
              <Text variant="body2" className={classes.col}>
                {cell.mobileCountryCode || 'unknown'}
              </Text>
              <Text variant="body2" className={classes.col}>
                {cell.mobileNetworkCode || 'unknown'}
              </Text>
            </div>
          ))}
        </PerfectScrollbar>
      </div>
    </div>
  );
};

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
  },
  header: {
    padding: '20px 0px',
    borderBottom: `1px solid ${theme.palette.grey[200]}`,
  },
  body: {
    padding: '20px 0px',
  },
  row: {
    display: 'flex',
    flexDirection: 'row',
    marginBottom: '10px',
  },
  col: {
    flex: 1,
  },
}));

export default LocationCellScansPane;
