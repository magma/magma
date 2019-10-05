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
import React from 'react';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import alertsTheme from '@fbcnms/ui/theme/alerts';

import {SEVERITY_STYLE} from './AlarmsTable';
import {makeStyles} from '@material-ui/styles';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  alertConfig: AlertConfig,
  setAlertConfig: ((AlertConfig => AlertConfig) | AlertConfig) => void,
  onNext: () => void,
  onPrevious: () => void,
};

type AlertInfoFields = 'name' | 'description' | 'severity' | 'team';

const useStyles = makeStyles(() => ({
  body: alertsTheme.formBody,
  buttonGroup: alertsTheme.buttonGroup,
  inputArea: {
    width: '100%',
    display: 'flex',
  },
  inputColumn: {
    flexGrow: 1,
  },
}));

const AlertInput = withStyles({
  root: {
    marginBottom: '20px',
  },
})(TextField);

export default function AddEditAlertInfoStep(props: Props) {
  const classes = useStyles();

  const fieldChangedHandler = (field: AlertInfoFields) => event => {
    switch (field) {
      case 'name':
        props.setAlertConfig({...props.alertConfig, alert: event.target.value});
        break;
      case 'description':
        props.setAlertConfig({
          ...props.alertConfig,
          annotations: {
            ...props.alertConfig.annotations,
            description: event.target.value,
          },
        });
        break;
      case 'severity':
        props.setAlertConfig({
          ...props.alertConfig,
          labels: {
            ...props.alertConfig.labels,
            severity: event.target.value,
          },
        });
        break;
    }
  };

  const {alert, labels, annotations} = props.alertConfig;

  return (
    <>
      <Typography variant="h6">ALERT INFO</Typography>
      <div className={classes.body}>
        <div className={classes.inputArea}>
          <div className={classes.inputColumn}>
            <Typography variant="subtitle1">Name your new alert</Typography>
            <AlertInput
              required
              label="Required"
              placeholder="Ex: Service Down"
              value={alert}
              onChange={fieldChangedHandler('name')}
            />
            <Typography variant="subtitle1">Describe your new alert</Typography>
            <AlertInput
              placeholder="Ex: The service is down at xyz."
              value={annotations?.description || ''}
              onChange={fieldChangedHandler('description')}
            />
            <Typography variant="subtitle1">
              What would be the severity of this alert?
            </Typography>
            <AlertInput
              select
              placeholder="Ex: Critical"
              value={labels?.severity || Object.keys(SEVERITY_STYLE)[0]}
              onChange={fieldChangedHandler('severity')}>
              {Object.keys(SEVERITY_STYLE).map(severity => (
                <MenuItem key={severity} value={severity}>
                  {severity.toUpperCase()}
                </MenuItem>
              ))}
            </AlertInput>
          </div>
          <div className={classes.inputColumn}>
            <Typography variant="subtitle1">What team?</Typography>
            <Grid container spacing={1} alignItems="flex-end">
              <Grid item>
                <TextField disabled value="operations" />
              </Grid>
              <Grid item>
                <Tooltip
                  title={
                    'Teams are used for routing your alert to a receiver. ' +
                    'Since setting routes and receivers is not available yet ' +
                    "(we're working on it!), we've made this field default " +
                    'to operations.'
                  }
                  placement="right">
                  <HelpIcon />
                </Tooltip>
              </Grid>
            </Grid>
          </div>
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
            className={classes.button}
            variant="contained"
            color="primary"
            onClick={() => props.onNext()}>
            Next
          </Button>
        </div>
      </div>
    </>
  );
}
