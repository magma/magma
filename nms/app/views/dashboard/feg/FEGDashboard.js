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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DashboardAlertTable from '../../../components/DashboardAlertTable';
import EventAlertChart from '../../../components/EventAlertChart';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import EventsTable from '../../events/EventsTable';
import FEGDashboardKPIs from '../../../components/FEGDashboardKPIs';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import TopBar from '../../../components/TopBar';
import moment from 'moment';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {EVENT_STREAM} from '../../events/EventsTable';

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

/**
 * Returns the full federation network dashboard.
 * It consists of a top bar which helps in adjusting filters such as date
 * and a network dashboard which provides information about the network.
 */
function FEGDashboard() {
  // datetime picker
  const [startDate, setStartDate] = useState(moment().subtract(3, 'days'));
  const [endDate, setEndDate] = useState(moment());

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
 * @param {Array<moment>} startEnd: An array of two elements holding the
 * start and end date.
 */
function FEGNetworkDashboard({startEnd}: {startEnd: [moment, moment]}) {
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
function FEGNetworkTab(props) {
  const {startDate, endDate, setStartDate, setEndDate} = props;
  const classes = useStyles();
  return (
    <Grid container justifyContent="flex-end" alignItems="center" spacing={2}>
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
  );
}
export default FEGDashboard;
