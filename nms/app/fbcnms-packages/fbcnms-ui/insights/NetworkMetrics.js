/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';
import type {TimeRange} from '../insights/AsyncMetric';

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '../insights/AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from '../insights/TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';
import {resolveQuery} from '../insights/Metrics';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
  appBar: {
    display: 'inline-block',
  },
}));

export default function (props: {configs: Array<MetricGraphConfig>}) {
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('3_hours');

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
      </AppBar>
      <GridList cols={2} cellHeight={300}>
        {props.configs.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Text component="h6" variant="h6">
                  {config.label}
                </Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(config, '', '')}
                    timeRange={timeRange}
                    legendLabels={config.legendLabels}
                  />
                </div>
              </CardContent>
            </Card>
          </GridListTile>
        ))}
      </GridList>
    </>
  );
}
