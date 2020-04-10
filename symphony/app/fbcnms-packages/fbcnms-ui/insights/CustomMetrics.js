/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AppBar from '@material-ui/core/AppBar';
import AsyncMetric from '../insights/AsyncMetric';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import Close from '@material-ui/icons/Close';
import FormControl from '@material-ui/core/FormControl';
import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import IconButton from '@material-ui/core/IconButton';
import InputLabel from '@material-ui/core/InputLabel';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage';
import React from 'react';
import Select from '@material-ui/core/Select';
import TimeRangeSelector from '../insights/TimeRangeSelector';
import Typography from '@material-ui/core/Typography';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

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
  addButton: {marginTop: '20px'},
  removeButton: {float: 'right', padding: '3px'},
}));

export default function() {
  const classes = useStyles();
  const {history, match} = useRouter();
  const [selectedMetric, setSelectedMetric] = useState('');
  const [timeRange, setTimeRange] = useState('12_hours');
  const [deviceID, setDeviceID] = useState('');
  const [allMetrics, setAllMetrics] = useState();
  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusSeries,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => {
      const metricsByDeviceID = {};
      response.forEach(item => {
        if (item.deviceID) {
          metricsByDeviceID[item.deviceID] =
            metricsByDeviceID[item.deviceID] || new Set();
          metricsByDeviceID[item.deviceID].add(item.__name__);
        }
      });
      setDeviceID(Object.keys(metricsByDeviceID)[0]);
      setAllMetrics(metricsByDeviceID);
    }, []),
  );

  if (error || isLoading || !allMetrics) {
    return <LoadingFiller />;
  }

  if (Object.keys(allMetrics).length === 0) {
    return (
      <NoNetworksMessage>
        There are currently no metrics available for display
      </NoNetworksMessage>
    );
  }

  const queryParams = new URLSearchParams(history.location.search);
  const queryMetrics = JSON.parse(queryParams.get('metrics') || '[]');
  const configs = queryMetrics.map(m => ({
    label: '',
    basicQueryConfigs: [
      {metric: m.metric, filters: [`deviceID=${m.deviceID}`]},
    ],
    timeRange: m.timeRange,
    deviceID: m.deviceID,
    metric: m.metric,
  }));

  const onAdd = () => {
    const newMetrics = [
      ...queryMetrics,
      {metric: selectedMetric, deviceID: deviceID, timeRange},
    ];
    queryParams.set('metrics', JSON.stringify(newMetrics));
    history.push({search: queryParams.toString()});
  };

  const onRemove = index => {
    const newMetrics = [...queryMetrics];
    newMetrics.splice(index, 1);
    queryParams.set('metrics', JSON.stringify(newMetrics));
    history.push({search: queryParams.toString()});
  };

  return (
    <>
      <AppBar className={classes.appBar} position="static" color="default">
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="devices">Device</InputLabel>
          <Select
            inputProps={{id: 'devices'}}
            value={deviceID}
            onChange={event => setDeviceID(event.target.value)}>
            {Object.keys(allMetrics).map(device => (
              <MenuItem value={device} key={device}>
                {device}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControl variant="filled" className={classes.formControl}>
          <InputLabel htmlFor="metrics">Metric</InputLabel>
          <Select
            inputProps={{id: 'metrics'}}
            value={selectedMetric}
            onChange={event => setSelectedMetric(event.target.value)}>
            {allMetrics[deviceID].map(metric => (
              <MenuItem value={metric} key={metric}>
                {metric}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <TimeRangeSelector
          className={classes.formControl}
          value={timeRange}
          onChange={setTimeRange}
        />
        <Button onClick={onAdd} className={classes.addButton}>
          Add Chart
        </Button>
      </AppBar>
      <GridList cols={2} cellHeight={300}>
        {configs.map((config, i) => (
          <GridListTile key={i} cols={1}>
            <Card>
              <CardContent>
                <Typography component="h6" variant="h6">
                  {config.label}
                  <IconButton
                    className={classes.removeButton}
                    onClick={() => onRemove(i)}>
                    <Close />
                  </IconButton>
                </Typography>
                <div style={{height: 250}}>
                  <AsyncMetric
                    label={config.label}
                    unit=""
                    queries={[
                      `${config.metric}{deviceID="${config.deviceID}"}`,
                    ]}
                    timeRange={config.timeRange}
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
