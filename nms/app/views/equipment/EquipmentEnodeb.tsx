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
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EmptyState from '../../components/EmptyState';
import EnodebContext from '../../context/EnodebContext';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import React from 'react';
import SettingsInputAntennaIcon from '@mui/icons-material/SettingsInputAntenna';
import withAlert from '../../components/Alert/withAlert';
import {EnodeEditDialog} from './EnodebDetailConfigEdit';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {Theme} from '@mui/material/styles';
import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@mui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useInterval} from '../../hooks';
import {useNavigate} from 'react-router-dom';
import type {WithAlert} from '../../components/Alert/withAlert';

const CHART_TITLE = 'Total Throughput';
const EMPTY_STATE_INSTRUCTIONS =
  'eNodeBs can be either managed (via TR-069) or unmanaged.  Managed eNodeBs can be configured directly' +
  ' from the Access Gateway via the enodebd service. Unmanaged eNodeBs are configured externally on their own device management portal.';
const EMPTY_STATE_OVERVIEW =
  'The eNodeB is the Radio Access Network that connects the user devices to the Packet Core. ' +
  'eNodeBs (eNBs) can be either managed (via TR-069) or unmanaged.';

const useStyles = makeStyles<Theme>(theme => ({
  cardContent: {
    padding: '0 16px',
  },
  cardHeaderTitle: {
    fontSize: '14px',
    fontWeight: 'bold',
  },
  customIntructions: {
    marginTop: '24px',
    width: '100%',
  },
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
  instructions: {
    backgroundColor: colors.primary.concrete,
    height: '100%',
  },
  bulletList: {
    padding: '8px',
    listStyleType: 'disc',
    fontSize: '14px',
  },
}));
function AddEnodebInstructions(props: {setOpen: () => void}) {
  const classes = useStyles();

  return (
    <Grid
      className={classes.customIntructions}
      spacing={3}
      container
      justifyContent="space-between">
      <Grid item xs={6}>
        <Card className={classes.instructions}>
          <CardHeader
            classes={{title: classes.cardHeaderTitle}}
            title={"If you're provisioning a managed eNodeB"}
          />
          <CardContent className={classes.cardContent}>
            <ul className={classes.bulletList}>
              <li>Enter name and serial number</li>
              <li>
                Configure the RAN parameters (Note that fields left blank will
                be inherited from either the network or gateway LTE parameters)
              </li>
            </ul>
          </CardContent>
          <CardActions disableSpacing={true}>
            <Grid container direction="column" spacing={1}>
              <Grid item xs={5}>
                <Button
                  variant="contained"
                  color="primary"
                  onClick={() => props.setOpen()}>
                  Add enodeb
                </Button>
              </Grid>
              <Grid item>
                <Link
                  href="https://docs.magmacore.org/docs/next/lte/deploy_config_enodebd#configure-enodeb"
                  target="_blank"
                  underline="hover">
                  Learn more about the supported eNodeB and protocols
                </Link>
              </Grid>
            </Grid>
          </CardActions>
        </Card>
      </Grid>
      <Grid item xs={6}>
        <Card className={classes.instructions}>
          <CardHeader
            classes={{title: classes.cardHeaderTitle}}
            title={"If you're configuring an unmanaged eNodeB"}
          />
          <CardContent className={classes.cardContent}>
            <ul className={classes.bulletList}>
              <li>The unmanaged eNodeB is configured manually on the eNodeB</li>
              <li>
                Optionally, you can add the eNodeB to the NMS for tracking
                purposes
              </li>
            </ul>
          </CardContent>
          <CardActions disableSpacing={true}>
            <Link
              href="https://docs.magmacore.org/docs/next/lte/deploy_config_enodebd#manual-configuration"
              target="_blank"
              underline="hover">
              How to configure an unmanaged eNodeB
            </Link>
          </CardActions>
        </Card>
      </Grid>
    </Grid>
  );
}
export default function Enodeb() {
  const classes = useStyles();
  const ctx = useContext(EnodebContext);
  const [open, setOpen] = useState(false);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justifyContent="space-between" spacing={3}>
        <EnodeEditDialog
          open={open}
          onClose={() => {
            setOpen(false);
          }}
        />
        {Object.keys(ctx.state?.enbInfo || {}).length > 0 ? (
          <>
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
          </>
        ) : (
          <>
            <EmptyState
              title={'Set up an eNodeB'}
              instructions={EMPTY_STATE_INSTRUCTIONS}
              overviewTitle={'eNodeB Overview'}
              overviewDescription={EMPTY_STATE_OVERVIEW}
              customIntructions={
                <AddEnodebInstructions setOpen={() => setOpen(true)} />
              }
            />
          </>
        )}
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
  const enqueueSnackbar = useEnqueueSnackbar();
  const enodebContext = useContext(EnodebContext);
  const [refresh, setRefresh] = useState(true);
  const [currRow, setCurrRow] = useState<EnodebRowType>({} as EnodebRowType);

  useInterval(() => enodebContext.refetch(), refresh ? REFRESH_INTERVAL : null);

  const enbRows: Array<EnodebRowType> = enodebContext.state?.enbInfo
    ? Object.keys(enodebContext.state?.enbInfo).map((serialNum: string) => {
        const enbInf = enodebContext.state?.enbInfo[serialNum];
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
        label={`Enodebs (${
          Object.keys(enodebContext.state?.enbInfo || {}).length
        })`}
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
                onClick={() => navigate(currRow.id)}
                underline="hover">
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
                    await enodebContext.setState(currRow.id);
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
