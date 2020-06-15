/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Card from '@material-ui/core/Card';
import EnodebKPIs from './EnodebKPIs';
import GatewayKPIs from './GatewayKPIs';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Text from '../theme/design-system/Text';

import {colors} from '../theme/default';
import {GpsFixed} from '@material-ui/icons';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitle: {
    marginBottom: theme.spacing(1),
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  eventsTable: {
    marginTop: theme.spacing(4),
  },
}));

export default function() {
  const classes = useStyles();

  return (
    <>
      {/* TODO: Can come back and make this a reusable component for other cards */}
      <Grid container xs={12} className={classes.cardTitle}>
        <GpsFixed className={classes.cardTitleIcon} />
        <Text variant="body1">Events (388)</Text>
      </Grid>
      <Grid zeroMinWidth container alignItems="center" spacing={4}>
        <Grid item xs={12} md={6} alignItems="center">
          <Paper elevation={0}>
            <GatewayKPIs />
          </Paper>
        </Grid>
        <Grid item xs={12} md={6} alignItems="center">
          <Paper elevation={0}>
            <EnodebKPIs />
          </Paper>
        </Grid>
      </Grid>
      <Card elevation={0} className={classes.eventsTable}>
        <Text variant="h4">Events Table Goes Here</Text>
      </Card>
    </>
  );
}
