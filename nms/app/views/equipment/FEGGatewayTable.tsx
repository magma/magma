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

import type {FederationGateway} from '../../../generated';
import type {WithAlert} from '../../components/Alert/withAlert';

import ActionTable from '../../components/ActionTable';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@mui/icons-material/CellWifi';
import CheckIcon from '@mui/icons-material/Check';
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import FEGGatewayContext from '../../context/FEGGatewayContext';
import GatewayTierContext from '../../context/GatewayTierContext';
import Grid from '@mui/material/Grid';
import Link from '@mui/material/Link';
import MenuItem from '@mui/material/MenuItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import React, {useContext, useState} from 'react';
import Select from '@mui/material/Select';
import Text from '../../theme/design-system/Text';

import withAlert from '../../components/Alert/withAlert';
import {GatewayId} from '../../../shared/types/network';
import {GatewayTypeEnum, HEALTHY_STATUS} from '../../components/GatewayUtils';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {SelectEditComponent} from '../../components/ActionTable';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useInterval} from '../../hooks';
import {useNavigate} from 'react-router-dom';

type EquipmentFegGatewayRowType = {
  name: string;
  id: GatewayId;
  is_primary: boolean;
  health: string;
  gx: string;
  gy: string;
  swx: string;
  s6a: string;
  s8: string;
  csfb: string;
};

const ViewTypes = {
  STATUS: 'Status',
  UPGRADE: 'Upgrade',
};
/**
 * Displays the number of federation gateways alonside with a table showing
 * each federation gateway.
 */
export default function GatewayTable() {
  const ctx = useContext(FEGGatewayContext);
  const [refresh, setRefresh] = useState(true);
  const [currentView, setCurrentView] = useState<
    typeof ViewTypes[keyof typeof ViewTypes]
  >('Status');

  return (
    <>
      <CardTitleRow
        key="title"
        icon={CellWifiIcon}
        label={`Federation Gateways (${Object.keys(ctx.state).length})`}
        filter={() => (
          <Grid
            container
            justifyContent="flex-end"
            alignItems="center"
            spacing={1}>
            {currentView !== ViewTypes.UPGRADE && (
              <Grid item>
                <AutorefreshCheckbox
                  autorefreshEnabled={refresh}
                  onToggle={() => setRefresh(current => !current)}
                />
              </Grid>
            )}
            <Grid item>
              <Text variant="body3">View</Text>
            </Grid>
            <Grid item>
              <Select
                value={currentView}
                input={<OutlinedInput />}
                onChange={({target}) => setCurrentView(target.value)}>
                <MenuItem key={ViewTypes.STATUS} value={ViewTypes.STATUS}>
                  Status
                </MenuItem>
                <MenuItem key={ViewTypes.UPGRADE} value={ViewTypes.UPGRADE}>
                  Upgrade
                </MenuItem>
              </Select>
            </Grid>
          </Grid>
        )}
      />
      {currentView === ViewTypes.UPGRADE ? (
        <UpgradeTable />
      ) : (
        <StatusTable refresh={refresh} />
      )}
    </>
  );
}

/**
 * Returns a table containing the federation gateways which shows their
 * basic configurations.
 * @param {boolean} refresh Tells to autorefresh after 30 seconds or not.
 */
function GatewayStatusTable(props: WithAlert & {refresh: boolean}) {
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();
  const gwCtx = useContext(FEGGatewayContext);

  // Auto refresh gateways every 30 seconds
  useInterval(() => gwCtx.refetch(), props.refresh ? REFRESH_INTERVAL : null);

  const [currRow, setCurrRow] = useState<EquipmentFegGatewayRowType>(
    {} as EquipmentFegGatewayRowType,
  );
  const fegGatewayRows: Array<EquipmentFegGatewayRowType> = [];
  Object?.keys(gwCtx.state)
    .map((gwId: string) => gwCtx.state[gwId])
    .filter((g: FederationGateway) => g.federation && g.id)
    .map((gateway: FederationGateway) => {
      fegGatewayRows.push({
        name: gateway.name,
        id: gateway.id,
        is_primary: gwCtx.activeFegGatewayId === gateway.id,
        health: gwCtx.health[gateway.id]?.status
          ? gwCtx.health[gateway.id]?.status === HEALTHY_STATUS
            ? GatewayTypeEnum.HEALTHY_GATEWAY
            : GatewayTypeEnum.UNHEALTHY_GATEWAY
          : GatewayTypeEnum.UNKNOWN,
        gx: gateway.federation?.gx?.server?.address || '-',
        gy: gateway.federation?.gy?.server?.address || '-',
        swx: gateway.federation?.swx?.server?.address || '-',
        s6a: gateway.federation?.s6a?.server?.address || '-',
        s8: gateway.federation?.s8?.local_address || '-',
        csfb: gateway.federation?.csfb?.client?.server_address || '-',
      });
    });

  return (
    <>
      <ActionTable
        data={fegGatewayRows}
        columns={[
          {
            title: 'Name',
            field: 'name',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => navigate(currRow.id)}
                underline="hover">
                {currRow.name}
              </Link>
            ),
          },
          {
            title: 'Primary',
            field: 'is_primary',
            render: currRow =>
              currRow.is_primary && (
                <CheckIcon data-testid={`${currRow.id} is primary`} />
              ),
          },
          {
            title: 'Health',
            field: 'health',
            width: 100,
            render: currRow => (
              <>
                <DeviceStatusCircle
                  isActive={currRow.health === GatewayTypeEnum.HEALTHY_GATEWAY}
                  // grey out status if gateway had no health status
                  isGrey={currRow.health === GatewayTypeEnum.UNKNOWN}
                />
                {currRow.health}
              </>
            ),
          },
          {title: 'Gx', field: 'gx'},
          {title: 'Gy', field: 'gy'},
          {title: 'SWx', field: 'swx'},
          {title: 'S6a', field: 's6a'},
          {title: 'S8', field: 's8'},
          {title: 'CSFB', field: 'csfb'},
        ]}
        handleCurrRow={(row: EquipmentFegGatewayRowType) => setCurrRow(row)}
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

type EquipmentGatewayUpgradeType = {
  name: string;
  id: GatewayId;
  hardwareId: string;
  tier: string;
  currentVersion: string;
};
function UpgradeTable() {
  const tierCtx = useContext(GatewayTierContext);
  const gwCtx = useContext(FEGGatewayContext);
  const navigate = useNavigate();
  const enqueueSnackbar = useEnqueueSnackbar();

  const fegGatewayRows: Array<EquipmentGatewayUpgradeType> = [];
  Object.keys(gwCtx.state)
    .map((gwId: string) => gwCtx.state[gwId])
    .map((gateway: FederationGateway) => {
      const packages = gateway.status?.platform_info?.packages || [];
      fegGatewayRows.push({
        name: gateway.name,
        id: gateway.id,
        hardwareId: gateway.device?.hardware_id || '-',
        tier: gateway.tier,
        currentVersion:
          packages.find(p => p.name === 'magma')?.version || 'Not Reported',
      });
    });
  const [fegGatewayUpgradeRows, setFegGatewayUpgradeRows] = useState(
    fegGatewayRows,
  );
  return (
    <ActionTable
      data={fegGatewayUpgradeRows}
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
              defaultValue={props.value as string}
              value={props.value as string}
              content={Object.keys(tierCtx.state.tiers)}
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
            if (newData.tier) {
              // Update tier id
              await gwCtx.updateGateway({
                gatewayId: newData.id,
                tierId: newData.tier,
              });
              const dataUpdate = [...fegGatewayUpgradeRows];
              // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
              const index = (oldData as any).tableData.index as number;
              dataUpdate[index] = newData;
              setFegGatewayUpgradeRows([...dataUpdate]);
            } else {
              throw Error('Invalid tier');
            }
          } catch (e) {
            enqueueSnackbar(
              `failed saving gateway tier information: ${getErrorMessage(e)}`,
              {
                variant: 'error',
              },
            );
            throw e;
          }
        },
      }}
    />
  );
}
