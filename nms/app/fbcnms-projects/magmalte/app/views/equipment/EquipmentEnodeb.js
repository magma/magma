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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebContext from '../../components/context/EnodebContext';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const CHART_TITLE = 'Total Throughput';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function Enodeb() {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          <DateTimeMetricChart
            title={CHART_TITLE}
            queries={[
              `sum(pdcp_user_plane_bytes_dl{service="enodebd"} + pdcp_user_plane_bytes_ul{service="enodebd"})/1000`,
            ]}
            legendLabels={['mbps']}
          />
        </Grid>
        <Grid item xs={12}>
          <EnodebTable />
        </Grid>
      </Grid>
    </div>
  );
}

type EnodebRowType = {
  name: string,
  id: string,
  sessionName: string,
  health: string,
  reportedTime: Date,
};

function EnodebTableRaw(props: WithAlert) {
  const {history, relativeUrl} = useRouter();
  const ctx = useContext(EnodebContext);
  const [currRow, setCurrRow] = useState<EnodebRowType>({});
  const enbInfo = ctx.state.enbInfo;
  const enqueueSnackbar = useEnqueueSnackbar();
  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        name: enbInf.enb.name,
        id: serialNum,
        sessionName: enbInf.enb_state?.fsm_state ?? 'not available',
        health: isEnodebHealthy(enbInf) ? 'Good' : 'Bad',
        reportedTime: new Date(enbInf.enb_state.time_reported ?? 0),
      };
    },
  );

  return (
    <ActionTable
      titleIcon={SettingsInputAntennaIcon}
      title="EnodeBs"
      data={enbRows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Serial Number', field: 'id'},
        {title: 'Session State Name', field: 'sessionName'},
        {title: 'Health', field: 'health'},
        {title: 'Reported Time', field: 'reportedTime', type: 'datetime'},
      ]}
      handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
      menuItems={[
        {
          name: 'View',
          handleFunc: () => {
            history.push(relativeUrl('/' + currRow.id));
          },
        },
        {
          name: 'Edit',
          handleFunc: () => {
            history.push(relativeUrl('/' + currRow.id + '/config'));
          },
        },
        {
          name: 'Remove',
          handleFunc: () => {
            props
              .confirm(`Are you sure you want to delete ${currRow.id}?`)
              .then(async confirmed => {
                if (!confirmed) {
                  return;
                }

                try {
                  await ctx.setState(currRow.id);
                } catch (e) {
                  enqueueSnackbar('failed deleting enodeb ' + currRow.id, {
                    variant: 'error',
                  });
                }
              });
          },
        },
        {name: 'Deactivate'},
        {name: 'Reboot'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
      }}
    />
  );
}

const EnodebTable = withAlert(EnodebTableRaw);
