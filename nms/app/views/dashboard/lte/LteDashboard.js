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
import DashboardAlertTable from '../../../components/DashboardAlertTable';
import DashboardKPIs from '../../../components/DashboardKPIs';
import EventAlertChart from '../../../components/EventAlertChart';
import EventsTable from '../../events/EventsTable';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../../theme/design-system/Text';
import TopBar from '../../../components/TopBar';
import moment from 'moment';

import {DateTimePicker} from '@material-ui/pickers';
import {Navigate, Route, Routes} from 'react-router-dom';
import {NetworkCheck} from '@material-ui/icons';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  dateTimeText: {
    color: colors.primary.selago,
  },
}));

function LteDashboard() {
  const classes = useStyles();

  // datetime picker
  const [startDate, setStartDate] = useState(moment().subtract(3, 'days'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <TopBar
        key="dashboard"
        header="Dashboard"
        tabs={[
          {
            label: 'Network',
            to: 'network',
            icon: NetworkCheck,
            filters: (
              <Grid
                container
                justifyContent="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <Text variant="body3" className={classes.dateTimeText}>
                    Filter By Date
                  </Text>
                </Grid>
                <DateTimePicker
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  maxDate={endDate}
                  disableFuture
                  value={startDate}
                  onChange={setStartDate}
                />
                <Grid item>
                  <Text variant="body3" className={classes.dateTimeText}>
                    to
                  </Text>
                </Grid>
                <DateTimePicker
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  disableFuture
                  value={endDate}
                  onChange={setEndDate}
                />
              </Grid>
            ),
          },
        ]}
      />

      <Routes>
        <Route
          path="/network/*"
          element={<LteNetworkDashboard startEnd={[startDate, endDate]} />}
        />
        <Route index element={<Navigate to="network" replace />} />
      </Routes>
    </>
  );
}

function LteNetworkDashboard({startEnd}: {startEnd: [moment, moment]}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <EventAlertChart startEnd={startEnd} />
        </Grid>

        <Grid item xs={12}>
          <DashboardAlertTable />
        </Grid>

        <Grid item xs={12}>
          <DashboardKPIs />
        </Grid>
        <Grid item xs={12}>
          <EventsTable
            eventStream="NETWORK"
            sz="md"
            inStartDate={startEnd[0]}
            inEndDate={startEnd[1]}
            isAutoRefreshing={true}
          />
        </Grid>
      </Grid>
    </div>
  );
}

export default LteDashboard;
