/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TimeRange} from './AsyncMetric';

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from './AsyncMetric';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TimeRangeSelector from './TimeRangeSelector';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  appBar: {
    display: 'inline-block',
  },
  chartRow: {
    display: 'flex',
  },
  formControl: {
    minWidth: '200px',
    padding: theme.spacing(),
  },
}));

export type MetricGraphConfig = {
  basicQueryConfigs: BasicQueryConfig[],
  customQueryConfigs?: CustomQuery[],
  label: string,
  unit?: string,
  legendLabels?: string[],
};

export type CustomQuery = {
  resolveQuery: string => string,
};

export type BasicQueryConfig = {
  filters: MetricLabel[],
  metric: string,
};

export type MetricLabel = {
  name: string,
  value: string,
};

export function resolveQuery(
  config: MetricGraphConfig,
  filterName: string,
  filterValue: string,
): string[] {
  if (config.customQueryConfigs) {
    return resolveCustomQuery(config.customQueryConfigs, filterValue);
  }
  return resolveBasicQuery(config.basicQueryConfigs, filterName, filterValue);
}

function resolveBasicQuery(
  configs: BasicQueryConfig[],
  filterName: string,
  filterValue: string,
): string[] {
  return configs.map(config => {
    const filterString = resolveFilters(
      config.filters,
      filterName,
      filterValue,
    );
    return `${config.metric}{${filterString}}`;
  });
}

function resolveFilters(
  filters: MetricLabel[],
  filterName: string,
  filterValue: string,
): string {
  const dbFilters: string[] = filters.map(
    filter => filter.name + '="' + filter.value + '"',
  );
  dbFilters.push(`${filterName}="${filterValue}"`);
  return dbFilters.join(',');
}

function resolveCustomQuery(
  configs: CustomQuery[],
  filterValue: string,
): string[] {
  return configs.map(config => config.resolveQuery(filterValue));
}

export default function(props: {
  selectors: Array<string>,
  defaultSelector: string,
  onSelectorChange: (SyntheticInputEvent<EventTarget>) => void,
  configs: MetricGraphConfig[],
  selectorName: string,
}) {
  const {match} = useRouter();
  const classes = useStyles();
  const [timeRange, setTimeRange] = useState<TimeRange>('24_hours');

  const selectedID = match.params.selectedID;
  const selectedOrDefault = selectedID || props.defaultSelector;

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="devices">Device</InputLabel>
          <Select
            inputProps={{id: 'devices'}}
            value={selectedOrDefault}
            onChange={props.onSelectorChange}>
            {props.selectors.map(device => (
              <MenuItem value={device} key={device}>
                {device}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
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
                <Text variant="h6">{config.label}</Text>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit={config.unit || ''}
                    queries={resolveQuery(
                      config,
                      props.selectorName,
                      selectedOrDefault,
                    )}
                    timeRange={timeRange}
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
