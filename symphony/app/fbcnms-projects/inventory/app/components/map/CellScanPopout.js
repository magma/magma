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

import type {AggregatedCellScan} from '../location/coverage_map/cell_scan/LocationCellScanCoverageMap';

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Text from '@fbcnms/ui/components/design-system/Text';

import {makeStyles} from '@material-ui/styles';

type Props = {
  aggregatedCellData: ?AggregatedCellScan,
  renderCellScansDialog?: (cell: AggregatedCellScan) => void,
};

const CellScanPopout = (props: Props) => {
  const classes = useStyles();
  const {aggregatedCellData, renderCellScansDialog} = props;
  if (aggregatedCellData == null) {
    return null;
  }
  let cellSections = [
    {
      title: 'latitude',
      value: aggregatedCellData.latitude.toFixed(3),
    },
    {
      title: 'longitude',
      value: aggregatedCellData.longitude.toFixed(3),
    },
    {
      title: aggregatedCellData.cells.length === 1 ? 'dBm' : 'max dBm',
      value:
        aggregatedCellData.signalStrength == null
          ? 'unknown'
          : aggregatedCellData.signalStrength,
    },
  ];
  if (aggregatedCellData.cells.length === 1) {
    const cell = aggregatedCellData.cells[0];
    cellSections = cellSections.concat([
      {
        title: 'network',
        value: cell.networkType,
      },
      {
        title: 'operator',
        value: cell.operator == null ? 'unknown' : cell.operator,
      },
      {
        title: 'mcc',
        value:
          cell.mobileCountryCode == null ? 'unknown' : cell.mobileCountryCode,
      },
      {
        title: 'mnc',
        value:
          cell.mobileNetworkCode == null ? 'unknown' : cell.mobileNetworkCode,
      },
    ]);
  }

  const onClickAllCellScans = () => {
    renderCellScansDialog &&
      aggregatedCellData &&
      renderCellScansDialog(aggregatedCellData);
  };

  return (
    <div className={classes.root}>
      {cellSections.map((section, idx) => (
        <div className={classes.section} key={idx}>
          <Text variant="subtitle2" className={classes.title}>
            {`${section.title}:`}
          </Text>
          <Text variant="body2" className={classes.body}>
            {section.value}
          </Text>
        </div>
      ))}
      {aggregatedCellData.cells.length > 1 && (
        <Button onClick={onClickAllCellScans}>
          {`View ${aggregatedCellData.cells.length} Cell Scans`}
        </Button>
      )}
    </div>
  );
};

const useStyles = makeStyles({
  root: {
    marginTop: '8px',
    maxWidth: '600px',
    minWidth: '100px',
  },
  section: {
    marginBottom: '5px',
  },
  title: {
    display: 'inline',
    marginRight: '3px',
  },
  body: {
    display: 'inline',
  },
});

export default CellScanPopout;
