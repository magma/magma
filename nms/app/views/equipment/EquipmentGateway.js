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
import type {WithAlert} from '../../components/Alert/withAlert';
import type {
  gateway_id,
  lte_gateway,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import EquipmentGatewayKPIs from './EquipmentGatewayKPIs';
import GatewayCheckinChart from './GatewayCheckinChart';
// $FlowFixMe migrated to typescript
import GatewayContext from '../../components/context/GatewayContext';
import GatewayTierContext from '../../components/context/GatewayTierContext';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React, {useContext, useEffect, useState} from 'react';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../components/context/SubscriberContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TypedSelect from '../../components/TypedSelect';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {SelectEditComponent} from '../../components/ActionTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

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
  viewLabelText: {
    color: colors.primary.comet,
  },
}));

const UPGRADE_VIEW = 'UPGRADE';

export default function Gateway() {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justifyContent="space-between" spacing={3}>
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
      </Grid>
    </div>
  );
}

type EquipmentGatewayRowType = {
  name: string,
  id: gateway_id,
  num_enodeb: number,
  num_subscribers: number,
  health: string,
  checkInTime: Date | string,
};

type EquipmentGatewayUpgradeType = {
  name: string,
  id: gateway_id,
  hardwareId: string,
  tier: string,
  currentVersion: string,
};

const ViewTypes = {
  STATUS: 'Status',
  UPGRADE: 'Upgrade',
};

function GatewayTable() {
  const classes = useStyles();
  const [currentView, setCurrentView] = useState<$Keys<typeof ViewTypes>>(
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
  const lteGateways = gwCtx.state;
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();

  const lteGatewayRows: Array<EquipmentGatewayUpgradeType> = [];
  Object.keys(lteGateways)
    .map((gwId: string) => lteGateways[gwId])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map((gateway: lte_gateway) => {
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
              onClick={() => navigate(currRow.id)}>
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
              defaultValue={props.value}
              value={props.value}
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
        onRowUpdate: async (newData, oldData) =>
          new Promise(async (resolve, reject) => {
            try {
              await gwCtx.updateGateway({
                gatewayId: newData.id,
                tierId: newData.tier,
              });
              const dataUpdate = [...lteGatewayUpgradeRows];
              const index = oldData.tableData.id;
              dataUpdate[index] = newData;
              setLteGatewayUpgradeRows([...dataUpdate]);
              resolve();
            } catch (e) {
              enqueueSnackbar('failed saving gateway tier information', {
                variant: 'error',
              });
              reject();
            }
          }),
      }}
    />
  );
}

function GatewayStatusTable(props: WithAlert & {refresh: boolean}) {
  const navigate = useNavigate();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const gwCtx = useContext(GatewayContext);
  const [lastRefreshTime, setLastRefreshTime] = useState(
    new Date().toLocaleString(),
  );
  // Auto refresh gateways every 30 seconds
  const state = useRefreshingContext({
    context: GatewayContext,
    networkId: networkId,
    type: 'gateway',
    interval: REFRESH_INTERVAL,
    enqueueSnackbar,
    refresh: props.refresh,
    lastRefreshTime: lastRefreshTime,
  });
  const ctxValues = [...Object.values(gwCtx.state)];
  useEffect(() => {
    setLastRefreshTime(new Date().toLocaleString());
  }, [ctxValues.length]);

  const subscriberCtx = useContext(SubscriberContext);
  const gwSubscriberMap = subscriberCtx.gwSubscriberMap;

  const lteGateways = state;
  const [currRow, setCurrRow] = useState<EquipmentGatewayRowType>({});
  const lteGatewayRows: Array<EquipmentGatewayRowType> = [];

  Object.keys(lteGateways)
    .map((gwId: string) => lteGateways[gwId])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map((gateway: lte_gateway) => {
      let numEnodeBs = 0;
      if (gateway.connected_enodeb_serials) {
        numEnodeBs = gateway.connected_enodeb_serials.length;
      }

      let checkInTime = '-';
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
          // $FlowIgnore: gateway.device should be present
          gwSubscriberMap?.[gateway.device.hardware_id]?.length ?? 0,
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
                onClick={() => navigate(currRow.id)}>
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
              props
                .confirm(`Are you sure you want to delete ${currRow.id}?`)
                .then(async confirmed => {
                  if (!confirmed) {
                    return;
                  }

                  try {
                    await gwCtx.setState(currRow.id);
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
