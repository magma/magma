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
import EventAlertChart from '../../../components/EventAlertChart';
import EventsTable from '../../events/EventsTable';
import FEGDashboardKPIs from '../../../components/FEGDashboardKPIs';
import Grid from '@mui/material/Grid';
import React, {useState} from 'react';
import Text from '../../../theme/design-system/Text';
import TextField from '@mui/material/TextField';
import TopBar from '../../../components/TopBar';
import {DateTimePicker} from '@mui/x-date-pickers/DateTimePicker';
import {EVENT_STREAM} from '../../events/EventsTable';
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

/**
 * Returns the full federation network dashboard.
 * It consists of a top bar which helps in adjusting filters such as date
 * and a network dashboard which provides information about the network.
 */
function FEGDashboard() {
  // datetime picker
  const [startDate, setStartDate] = useState(subDays(new Date(), 3));
  const [endDate, setEndDate] = useState(new Date());

  return (
    <>
      <TopBar
        key="dashboard"
        header="Federated Network Dashboard"
        tabs={[
          {
            label: 'Network',
            to: 'network',
            icon: NetworkCheck,
            filters: (
              <FEGNetworkTab
                startDate={startDate}
                endDate={endDate}
                setStartDate={setStartDate}
                setEndDate={setEndDate}
              />
            ),
          },
        ]}
      />

      <Routes>
        <Route
          path="/network/*"
          element={<FEGNetworkDashboard startEnd={[startDate, endDate]} />}
        />
        <Route index element={<Navigate to="network" replace />} />
      </Routes>
    </>
  );
}

/**
 * Returns the network dashboard of the federation network.
 * It consists of an event alert chart, an alert table, a kpi for the
 * federation network and events table which helps in describing the
 * current network state.
 * @param {Array<Date>} startEnd: An array of two elements holding the
 * start and end date.
 */
function FEGNetworkDashboard({startEnd}: {startEnd: [Date, Date]}) {
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
          <FEGDashboardKPIs />
        </Grid>
        <Grid item xs={12}>
          <EventsTable
            eventStream={EVENT_STREAM.NETWORK}
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

/**
 * Returns the topbar of the dashboard which is useful in filtering out the dates.
 * @param {object} props: props consists of the startDate and endDate selected
 * by the user. It also consists of functions(setStartDate and setEndDate)
 * needed to change those values.
 */

type Props = {
  startDate: Date;
  endDate: Date;
  setStartDate: (startDate: Date) => void;
  setEndDate: (endDate: Date) => void;
};

function FEGNetworkTab(props: Props) {
  const {startDate, endDate, setStartDate, setEndDate} = props;
  const classes = useStyles();
  return (
    <Grid container justifyContent="flex-end" alignItems="center" spacing={2}>
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
  );
}

export default FEGDashboard;
