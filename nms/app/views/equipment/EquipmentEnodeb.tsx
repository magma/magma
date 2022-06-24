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
import ActionTable from '../../components/ActionTable';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebContext from '../../components/context/EnodebContext';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';
import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
import {Theme} from '@material-ui/core/styles';
import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';
import type {WithAlert} from '../../components/Alert/withAlert';

const CHART_TITLE = 'Total Throughput';

const useStyles = makeStyles<Theme>(theme => ({
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
      <Grid container justifyContent="space-between" spacing={3}>
        <Grid item xs={12}>
          <DateTimeMetricChart
            unit={'Throughput(mb/s)'}
            title={CHART_TITLE}
            queries={[
              `sum(rate(gtp_port_user_plane_dl_bytes{service="pipelined"}[5m]) + rate(gtp_port_user_plane_ul_bytes{service="pipelined"}[5m]))/1000`,
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
  name: string;
  id: string;
  sessionName: string;
  mmeConnected: string;
  health: string;
  reportedTime: Date;
  numSubscribers: number;
};

function EnodebTableRaw(props: WithAlert) {
  const navigate = useNavigate();
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const networkId: string = nullthrows(params.networkId);
  const ctx = useContext(EnodebContext);
  const [refresh, setRefresh] = useState(true);
  const [lastRefreshTime, setLastRefreshTime] = useState(
    new Date().toLocaleString(),
  );

  // Auto refresh  every 30 seconds
  const state = useRefreshingContext({
    context: EnodebContext,
    networkId: networkId,
    type: 'enodeb',
    interval: REFRESH_INTERVAL,
    refresh: refresh,
    lastRefreshTime: lastRefreshTime,
  });
  const ctxValues = [...Object.values(ctx.state.enbInfo)];
  useEffect(() => {
    setLastRefreshTime(new Date().toLocaleString());
  }, [ctxValues.length]);

  const [currRow, setCurrRow] = useState<EnodebRowType>({} as EnodebRowType);
  const enbInfo = state?.enbInfo;
  const enbRows: Array<EnodebRowType> = enbInfo
    ? Object.keys(enbInfo).map((serialNum: string) => {
        const enbInf = enbInfo[serialNum];
        const isEnbManaged =
          enbInf.enb?.enodeb_config?.config_type === 'MANAGED';
        return {
          name: enbInf.enb.name,
          id: serialNum,
          numSubscribers: enbInf.enb_state?.ues_connected ?? 0,
          sessionName: enbInf.enb_state?.fsm_state ?? '-',
          ipAddress: enbInf.enb_state?.ip_address ?? '-',
          mmeConnected: enbInf.enb_state?.mme_connected
            ? 'Connected'
            : 'Disconnected',
          health: isEnbManaged
            ? isEnodebHealthy(enbInf)
              ? 'Good'
              : 'Bad'
            : '-',
          reportedTime: new Date(enbInf.enb_state.time_reported ?? 0),
        };
      })
    : [];

  return (
    <>
      <CardTitleRow
        key="title"
        icon={SettingsInputAntennaIcon}
        label={`Enodebs (${Object.keys(state?.enbInfo || {}).length})`}
        filter={() => (
          <Grid
            container
            justifyContent="flex-end"
            alignItems="center"
            spacing={1}>
            <Grid item>
              <AutorefreshCheckbox
                autorefreshEnabled={refresh}
                onToggle={() => setRefresh(current => !current)}
              />
            </Grid>
          </Grid>
        )}
      />
      <ActionTable
        title=""
        data={enbRows}
        columns={[
          {title: 'Name', field: 'name'},
          {
            title: 'Serial Number',
            field: 'id',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => navigate(currRow.id)}>
                {currRow.id}
              </Link>
            ),
          },
          {title: 'Session State Name', field: 'sessionName'},
          {title: 'IP Address', field: 'ipAddress'},
          {title: 'Subscribers', field: 'numSubscribers', width: 100},
          {title: 'MME', field: 'mmeConnected', width: 100},
          {title: 'Health', field: 'health', width: 100},
          {title: 'Reported Time', field: 'reportedTime', type: 'datetime'},
        ]}
        handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
        menuItems={[
          {
            name: 'View',
            handleFunc: () => {
              navigate(currRow.id);
            },
          },
          {
            name: 'Edit',
            handleFunc: () => {
              navigate(currRow.id + '/config');
            },
          },
          {
            name: 'Remove',
            handleFunc: () => {
              void props
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    await ctx.setState(currRow.id);
                    // setLastRefreshTime(new Date().toLocaleString());
                  } catch (e) {
                    enqueueSnackbar('failed deleting enodeb ' + currRow.id, {
                      variant: 'error',
                    });
                  }
                });
            },
          },
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </>
  );
}

const EnodebTable = withAlert(EnodebTableRaw);
