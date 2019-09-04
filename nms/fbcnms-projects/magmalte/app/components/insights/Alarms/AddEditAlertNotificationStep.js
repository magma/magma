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
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import alertsTheme from '@fbcnms/ui/theme/alerts';

import {makeStyles} from '@material-ui/styles';

type Props = {
  alertConfig: AlertConfig,
  setAlertConfig: AlertConfig => void,
  onSave: () => void,
  onPrevious: () => void,
};

const useStyles = makeStyles(() => ({
  body: alertsTheme.formBody,
  buttonGroup: alertsTheme.buttonGroup,
}));

export default function AddEditAlertNotificationStep(props: Props) {
  const classes = useStyles();

  return (
    <>
      <Typography variant="h6">SET YOUR NOTIFICATIONS</Typography>
      <div className={classes.body}>
        <div>
          <Typography variant="subtitle1">Notification Time</Typography>
          <Grid container spacing={1} alignItems="flex-end">
            <Grid item>
              <TextField
                required
                placeholder="Ex: 5m"
                label="Required"
                value={props.alertConfig.for}
                onChange={event =>
                  props.setAlertConfig({
                    ...props.alertConfig,
                    for: event.target.value,
                  })
                }
              />
            </Grid>
            <Grid item>
              <Tooltip
                title={
                  'Enter the amount of time the alert expression needs to be ' +
                  'true for before the alert fires. Click on the help icon to ' +
                  'open the prometheus time duration string formatting guide.'
                }
                placement="right">
                <IconButton
                  className={classes.iconButton}
                  href="https://prometheus.io/docs/prometheus/latest/querying/basics/#range-vector-selectors"
                  target="_blank">
                  <HelpIcon />
                </IconButton>
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
