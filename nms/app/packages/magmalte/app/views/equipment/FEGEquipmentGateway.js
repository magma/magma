/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import FEGEquipmentGatewayKPIs from './FEGEquipmentGatewayKPIs';
import GatewayCheckinChart from './GatewayCheckinChart';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
}));

/**
 * Returns the federation gateways check in chart, the gateways KPI,
 * cluster KPI, and a table of the federation gateways.
 */
export default function FEGGateway() {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          <GatewayCheckinChart />
        </Grid>
        <Grid item xs={12}>
          <Paper elevation={0}>
            <FEGEquipmentGatewayKPIs />
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}
