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

import type {LteGateway} from '../../../generated';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@mui/icons-material/CellWifi';
import EmptyState from '../../components/EmptyState';
import EquipmentGatewayKPIs from './EquipmentGatewayKPIs';
import GatewayCheckinChart from './GatewayCheckinChart';
import GatewayContext from '../../context/GatewayContext';
import GatewayTierContext from '../../context/GatewayTierContext';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemText from '@mui/material/ListItemText';
import OutlinedInput from '@mui/material/OutlinedInput';
import Paper from '@mui/material/Paper';
import React, {useContext, useState} from 'react';
import SubscriberContext from '../../context/SubscriberContext';
import Text from '../../theme/design-system/Text';
import TypedSelect from '../../components/TypedSelect';
import withAlert from '../../components/Alert/withAlert';
import {GatewayEditDialog} from './GatewayDetailConfigEdit';
import {GatewayId} from '../../../shared/types/network';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {SelectEditComponent} from '../../components/ActionTable';
import {Theme} from '@mui/material/styles';
import {colors} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useInterval} from '../../hooks';
import {useNavigate} from 'react-router-dom';

const useStyles = makeStyles<Theme>(theme => ({
  avatar: {
    backgroundColor: colors.primary.comet,
    color: colors.primary.white,
    width: '24px',
    height: '24px',
    fontSize: '12px',
    fontWeight: 'bold',
  },
  bulletList: {
    padding: '8px',
    listStyleType: 'disc',
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
  emptyState: {
    margin: '0',
  },
  emptyStateLink: {
    marginLeft: '0px',
    fontSize: '14px',
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
  viewLabelText: {
    color: colors.primary.comet,
  },
  listItemText: {
    marginTop: '0px',
  },
  listItemPrimary: {
    fontSize: '16px',
    color: colors.primary.mirage,
    marginBottom: '18px',
  },
  listItemSecondary: {
    fontSize: '14px',
    color: colors.primary.mirage,
  },
  listItem: {
    alignItems: 'flex-start',
  },
}));

const UPGRADE_VIEW = 'UPGRADE';
const EMPTY_STATE_OVERVIEW =
  'The Access Gateway (AGW) provides network services and policy enforcement. In an LTE network,' +
  ' the AGW implements an evolved packet core (EPC). It works with existing, unmodified commercial radio hardware.';
const EMPTY_STATE_INSTRUCTIONS_STEP_2 =
  'The Access Gateway (AGW) is configured and managed via the orchestrator and is part of a specific organization.' +
  'This configuration is made through the NMS as part of adding a new gateway to the system. ';
const INSTALL_AGW_LINK = 'https://docs.magmacore.org/docs/lte/deploy_install';
const CONFIGURE_AGW_LINK =
  'https://docs.magmacore.org/docs/lte/deploy_config_agw#access-gateway-configuration';
const LEARN_MORE_LINK = 'https://docs.magmacore.org/docs/lte/deploy_config_agw';

export default function Gateway() {
  const classes = useStyles();
  const ctx = useContext(GatewayContext);
  const [open, setOpen] = useState(false);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justifyContent="space-between" spacing={3}>
        <GatewayEditDialog open={open} onClose={() => setOpen(false)} />
        {Object.keys(ctx.state).length === 0 ? (
          <EmptyState
            title={'Set up a Gateway'}
            customIntructions={
              <GatewayEmptyStateList setOpen={() => setOpen(true)} />
            }
            instructions={''}
            overviewTitle={'Access Gateway Overview'}
            overviewDescription={EMPTY_STATE_OVERVIEW}
          />
        ) : (
          <>
            <Grid item xs={12}>
              <GatewayCheckinChart />
            </Grid>
            <Grid item xs={12}>
              <Paper elevation={0}>
                <EquipmentGatewayKPIs />
              </Paper>
            </Grid>
            <Grid item xs={12}>
              <GatewayTable />
            </Grid>
          </>
        )}
      </Grid>
    </div>
  );
}

function InstallGatewayList() {
  const classes = useStyles();
  return (
    <ul className={classes.bulletList}>
      <li>
        Create bootable USB with OS (Ubuntu).{' '}
        <Link href={INSTALL_AGW_LINK} target="_blank" underline="hover">
          View documentation
        </Link>
      </li>
      <li>Install Magma service</li>
      <li>
        Install <code>rootca.pem</code> and <code>control_proxy.yml</code>.{' '}
        <Link href={CONFIGURE_AGW_LINK} target="_blank" underline="hover">
          View documentation
        </Link>
      </li>
      <li>Restart Magma services</li>
      <li>
        Gather hardware ID and challenge key to add to the NMS
        <code>(show_gateway_info.py</code>)
      </li>
    </ul>
  );
}

type GatewayInstructionsProps = {
  setOpen: () => void;
};
function AddGatewayInstructions(props: GatewayInstructionsProps) {
  const classes = useStyles();
  return (
    <Grid container direction="column" spacing={2}>
      <Grid item>{EMPTY_STATE_INSTRUCTIONS_STEP_2}</Grid>
      <Grid item>
        <Grid container direction="column" spacing={3}>
          <Grid item xs={3}>
            <Button
              variant="contained"
              color="primary"
              onClick={() => props.setOpen()}>
              Add Gateway
            </Button>
          </Grid>
          <Grid item>
            <Link
              className={classes.emptyStateLink}
              href={LEARN_MORE_LINK}
              target="_blank"
              underline="hover">
              Learn more about Access Gateway Configuration
            </Link>
          </Grid>
        </Grid>
      </Grid>
    </Grid>
  );
}

function GatewayEmptyStateList(props: GatewayInstructionsProps) {
  const classes = useStyles();
  return (
    <List dense disablePadding>
      <ListItem classes={{root: classes.listItem}} disableGutters>
        <ListItemAvatar>
          <Avatar classes={{root: classes.avatar}}>1</Avatar>
        </ListItemAvatar>
        <ListItemText
          classes={{
            root: classes.listItemText,
            secondary: classes.listItemSecondary,
            primary: classes.listItemPrimary,
          }}
          primary="Install and configure an Access Gateway"
          secondary={<InstallGatewayList />}
          secondaryTypographyProps={{component: 'div'}}
        />
      </ListItem>
      <ListItem classes={{root: classes.listItem}} disableGutters>
        <ListItemAvatar>
          <Avatar classes={{root: classes.avatar}}>2</Avatar>
        </ListItemAvatar>
        <ListItemText
          classes={{
            root: classes.listItemText,
            secondary: classes.listItemSecondary,
            primary: classes.listItemPrimary,
          }}
          primary="Add an Access Gateway"
          secondary={<AddGatewayInstructions setOpen={() => props.setOpen()} />}
          secondaryTypographyProps={{component: 'div'}}
        />
      </ListItem>
    </List>
  );
}

type EquipmentGatewayRowType = {
  name: string;
  id: GatewayId;
  num_enodeb: number;
  num_subscribers: number;
  health: string;
  checkInTime: Date | string;
};

type EquipmentGatewayUpgradeType = {
  name: string;
  id: GatewayId;
  hardwareId: string;
  tier: string;
  currentVersion: string;
};

const ViewTypes = {
  STATUS: 'Status',
  UPGRADE: 'Upgrade',
};

function GatewayTable() {
  const classes = useStyles();
  const [currentView, setCurrentView] = useState<keyof typeof ViewTypes>(
    'STATUS',
  );
  const ctx = useContext(GatewayContext);
  const [refresh, setRefresh] = useState(true);

  return (
    <>
      <CardTitleRow
        key="title"
        icon={CellWifiIcon}
        label={`Gateways (${Object.keys(ctx.state).length})`}
        filter={() => (
          <Grid
            container
            justifyContent="flex-end"
            alignItems="center"
            spacing={1}>
            {currentView !== UPGRADE_VIEW && (
              <Grid item>
                <AutorefreshCheckbox
                  autorefreshEnabled={refresh}
                  onToggle={() => setRefresh(current => !current)}
                />
              </Grid>
            )}
            <Grid item>
              <Text variant="body3" className={classes.viewLabelText}>
                View
              </Text>
            </Grid>
            <Grid item>
              <TypedSelect
                input={<OutlinedInput />}
                value={currentView}
                items={{
                  STATUS: 'Status',
                  UPGRADE: 'Upgrade',
                }}
                onChange={setCurrentView}
              />
            </Grid>
          </Grid>
        )}
      />
      {currentView === 'UPGRADE' ? (
        <UpgradeTable />
      ) : (
        <StatusTable refresh={refresh} />
      )}
    </>
  );
}

function UpgradeTable() {
  const ctx = useContext(GatewayTierContext);
  const gwCtx = useContext(GatewayContext);
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();

  const lteGatewayRows: Array<EquipmentGatewayUpgradeType> = [];
  Object.keys(gwCtx.state)
    .map((gwId: string) => gwCtx.state[gwId])
    .filter((g: LteGateway) => g.cellular && g.id)
    .map((gateway: LteGateway) => {
      const packages = gateway.status?.platform_info?.packages || [];
      lteGatewayRows.push({
        name: gateway.name,
        id: gateway.id,
        hardwareId: gateway.device?.hardware_id || '-',
        tier: gateway.tier,
        currentVersion:
          packages.find(p => p.name === 'magma')?.version || 'Not Reported',
      });
    });
  const [lteGatewayUpgradeRows, setLteGatewayUpgradeRows] = useState(
    lteGatewayRows,
  );
  return (
    <ActionTable
      data={lteGatewayUpgradeRows}
      columns={[
        {title: 'Name', field: 'name', editable: 'never'},
        {
          title: 'ID',
          field: 'id',
          editable: 'never',
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
        {
          title: 'Hardware ID',
          field: 'hardwareId',
          editable: 'never',
        },
        {
          title: 'Current Version',
          field: 'currentVersion',
          editable: 'never',
          width: 250,
        },
        {
          title: 'Tier',
          field: 'tier',
          width: 100,
          editComponent: props => (
            <SelectEditComponent
              {...props}
              defaultValue={props.value as string}
              value={props.value as string}
              content={Object.keys(ctx.state.tiers)}
              onChange={value => props.onChange(value)}
            />
          ),
        },
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
      }}
      editable={{
        onRowUpdate: async (newData, oldData) => {
          try {
            await gwCtx.updateGateway({
              gatewayId: newData.id,
              tierId: newData.tier,
            });
            const dataUpdate = [...lteGatewayUpgradeRows];
            // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
            const index = (oldData as any).tableData.index as number;
            dataUpdate[index] = newData;
            setLteGatewayUpgradeRows([...dataUpdate]);
          } catch (e) {
            enqueueSnackbar('failed saving gateway tier information', {
              variant: 'error',
            });
            throw e;
          }
        },
      }}
    />
  );
}

function GatewayStatusTable(props: WithAlert & {refresh: boolean}) {
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();
  const gatewayContext = useContext(GatewayContext);
  const subscriberCtx = useContext(SubscriberContext);
  const gwSubscriberMap = subscriberCtx.gwSubscriberMap;
  const [currRow, setCurrRow] = useState<EquipmentGatewayRowType>(
    {} as EquipmentGatewayRowType,
  );
  const lteGatewayRows: Array<EquipmentGatewayRowType> = [];

  useInterval(
    () => gatewayContext.refetch(),
    props.refresh ? REFRESH_INTERVAL : null,
  );

  Object.keys(gatewayContext.state)
    .map((gwId: string) => gatewayContext.state[gwId])
    .filter((g: LteGateway) => g.cellular && g.id)
    .map((gateway: LteGateway) => {
      let numEnodeBs = 0;
      if (gateway.connected_enodeb_serials) {
        numEnodeBs = gateway.connected_enodeb_serials.length;
      }

      let checkInTime: string | Date = '-';
      if (
        gateway.status &&
        gateway.status.checkin_time != null &&
        gateway.status.checkin_time > 0
      ) {
        checkInTime = new Date(gateway.status.checkin_time);
      }

      lteGatewayRows.push({
        name: gateway.name,
        id: gateway.id,
        num_enodeb: numEnodeBs,
        num_subscribers:
          gwSubscriberMap?.[gateway.device!.hardware_id]?.length ?? 0,
        health: gateway.checked_in_recently ? 'Good' : 'Bad',
        checkInTime: checkInTime,
      });
    });
  return (
    <>
      <ActionTable
        data={lteGatewayRows}
        columns={[
          {title: 'Name', field: 'name'},
          {
            title: 'ID',
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
          {
            title: 'enodeBs',
            field: 'num_enodeb',
            width: 100,
          },
          {title: 'Subscribers', field: 'num_subscribers', width: 100},
          {title: 'Health', field: 'health', width: 100},
          {title: 'Check In Time', field: 'checkInTime', type: 'datetime'},
        ]}
        handleCurrRow={(row: EquipmentGatewayRowType) => setCurrRow(row)}
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
                    await gatewayContext.setState(currRow.id);
                  } catch (e) {
                    enqueueSnackbar('failed deleting gateway ' + currRow.id, {
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
const StatusTable = withAlert(GatewayStatusTable);
