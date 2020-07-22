/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import EnodebKPIs from './EnodebKPIs';
import GatewayKPIs from './GatewayKPIs';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import React from 'react';

import {CardTitleRow} from './layout/CardTitleRow';
import {GpsFixed} from '@material-ui/icons';

export default function () {
  return (
    <>
      <CardTitleRow icon={GpsFixed} label="Events" />
      <Grid container item zeroMinWidth alignItems="center" spacing={4}>
        <Grid item xs={12} md={6}>
          <Paper elevation={0}>
            <GatewayKPIs />
          </Paper>
        </Grid>
        <Grid item xs={12} md={6}>
          <Paper elevation={0}>
            <EnodebKPIs />
          </Paper>
        </Grid>
      </Grid>
    </>
  );
}
