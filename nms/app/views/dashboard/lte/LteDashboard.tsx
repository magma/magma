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
 */

import DashboardAlertTable from '../../../components/DashboardAlertTable';
import DashboardKPIs from '../../../components/DashboardKPIs';
import EventAlertChart from '../../../components/EventAlertChart';
import EventsTable from '../../events/EventsTable';
import Grid from '@mui/material/Grid';
import React, {useState} from 'react';
import Text from '../../../theme/design-system/Text';
import TextField from '@mui/material/TextField';
import TopBar from '../../../components/TopBar';
import {DateTimePicker} from '@mui/x-date-pickers/DateTimePicker';
import {Navigate, Route, Routes} from 'react-router-dom';
import {NetworkCheck} from '@mui/icons-material';
import {Theme} from '@mui/material/styles';
import {colors} from '../../../theme/default';
import {makeStyles} from '@mui/styles';
import {subDays} from 'date-fns';

const useStyles = makeStyles<Theme>(theme => ({
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
  const [startDate, setStartDate] = useState(subDays(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());

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
                <Grid item>
                  <DateTimePicker
                    renderInput={props => <TextField {...props} />}
                    maxDate={endDate}
                    disableFuture
                    value={startDate}
                    onChange={date => setStartDate(date!)}
                  />
                </Grid>
                <Grid item>
                  <Text variant="body3" className={classes.dateTimeText}>
                    to
                  </Text>
                </Grid>
                <Grid item>
                  <DateTimePicker
                    renderInput={props => <TextField {...props} />}
                    disableFuture
                    value={endDate}
                    onChange={date => setEndDate(date!)}
                  />
                </Grid>
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

function LteNetworkDashboard({startEnd}: {startEnd: [Date, Date]}) {
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
