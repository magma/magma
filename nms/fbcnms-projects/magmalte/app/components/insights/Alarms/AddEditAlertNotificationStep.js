/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {AlertConfig} from './AlarmAPIType';

import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import HelpIcon from '@material-ui/icons/Help';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useEffect, useState} from 'react';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import alertsTheme from '@fbcnms/ui/theme/alerts';

import {makeStyles} from '@material-ui/styles';

type Props = {
  alertConfig: AlertConfig,
  setAlertConfig: ((AlertConfig => AlertConfig) | AlertConfig) => void,
  onSave: () => void,
  onPrevious: () => void,
};

const useStyles = makeStyles(() => ({
  body: alertsTheme.formBody,
  buttonGroup: alertsTheme.buttonGroup,
}));

const timeUnits = [
  {
    value: '',
    label: '',
  },
  {
    value: 's',
    label: 'seconds',
  },
  {
    value: 'm',
    label: 'minutes',
  },
  {
    value: 'h',
    label: 'hours',
  },
];

export default function AddEditAlertNotificationStep(props: Props) {
  const classes = useStyles();
  const {alertConfig, setAlertConfig} = props;
  const duration = alertConfig.for ?? '';
  const [timeNumber, setTimeNumber] = useState<string>(duration.slice(0, -1));
  const [timeUnit, setTimeUnit] = useState<string>(
    duration[duration.length - 1],
  );

  // Since "for" is a reserved key word, need to create this function to handle
  // deleting the key in an immutable way. Eventually this function will be
  // useful for deleting other optional fields, like annotations
  const deleteAlertConfigKey = (obj: AlertConfig, prop: string) => {
    const {[prop]: _, ...res} = obj;
    return res;
  };

  useEffect(() => {
    if (timeNumber && timeUnit) {
      setAlertConfig(prevConfig => {
        return {
          ...prevConfig,
          for: timeNumber + timeUnit,
        };
      });
    } else {
      setAlertConfig(prevConfig => {
        return {
          ...deleteAlertConfigKey(prevConfig, 'for'),
        };
      });
    }
  }, [setAlertConfig, timeNumber, timeUnit]);

  return (
    <>
      <Typography variant="h6">SET YOUR NOTIFICATIONS</Typography>
      <div className={classes.body}>
        <div>
          <Typography variant="subtitle1">Grace Period</Typography>
          <Grid container spacing={1} alignItems="flex-end">
            <Grid item>
              <TextField
                type="number"
                value={timeNumber || ''}
                onChange={event => setTimeNumber(event.target.value)}
              />
            </Grid>
            <Grid item>
              <TextField
                select
                value={timeUnit || timeUnits[0].value}
                onChange={event => setTimeUnit(event.target.value)}>
                {timeUnits.map(option => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item>
              <Tooltip
                title={
                  'Enter the amount of time the alert expression needs to be ' +
                  'true for before the alert fires.'
                }
                placement="right">
                <HelpIcon />
              </Tooltip>
            </Grid>
          </Grid>
        </div>
        <div className={classes.buttonGroup}>
          <Button
            style={{marginRight: '10px'}}
            variant="contained"
            color="secondary"
            onClick={() => props.onPrevious()}>
            Previous
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={() => props.onSave()}>
            Save
          </Button>
        </div>
      </div>
    </>
  );
}
