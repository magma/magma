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
  federation_gateway,
  gateway_id,
} from '../../../generated/MagmaAPIBindings';

import ActionTable from '../../components/ActionTable';
// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import CheckIcon from '@material-ui/icons/Check';
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import Grid from '@material-ui/core/Grid';
import Link from '@material-ui/core/Link';
import React, {useContext, useEffect, useState} from 'react';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';

import {GatewayTypeEnum, HEALTHY_STATUS} from '../../components/GatewayUtils';
import {
  REFRESH_INTERVAL,
  RefreshTypeEnum,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

type EquipmentFegGatewayRowType = {
  name: string,
  id: gateway_id,
  is_primary: boolean,
  health: string,
  gx: string,
  gy: string,
  swx: string,
  s6a: string,
  csfb: string,
};

/**
 * Displays the number of federation gateways alonside with a table showing
 * each federation gateway.
 */
export default function GatewayTable() {
  const ctx = useContext(FEGGatewayContext);
  const [refresh, setRefresh] = useState(true);

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
            <Grid item>
              <AutorefreshCheckbox
                autorefreshEnabled={refresh}
                onToggle={() => setRefresh(current => !current)}
              />
            </Grid>
          </Grid>
        )}
      />
      <StatusTable refresh={refresh} />
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
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const networkId: string = nullthrows(params.networkId);
  const gwCtx = useContext(FEGGatewayContext);
  const [lastRefreshTime, setLastRefreshTime] = useState(
    new Date().toLocaleString(),
  );
  // Auto refresh gateways every 30 seconds
  const state = useRefreshingContext({
    context: FEGGatewayContext,
    networkId: networkId,
    type: RefreshTypeEnum.FEG_GATEWAY,
    interval: REFRESH_INTERVAL,
    enqueueSnackbar,
    refresh: props.refresh,
    lastRefreshTime: lastRefreshTime,
  });
  const fegGateways = state?.fegGateways || {};
  const health = state?.health || {};
  const activeFegGatewayId = state?.activeFegGatewayId || '';
  const ctxValues = [...Object.values(gwCtx.state)];
  useEffect(() => {
    setLastRefreshTime(new Date().toLocaleString());
  }, [ctxValues.length]);

  const [currRow, setCurrRow] = useState<EquipmentFegGatewayRowType>({});
  const fegGatewayRows: Array<EquipmentFegGatewayRowType> = [];
  Object?.keys(fegGateways)
    .map((gwId: string) => fegGateways[gwId])
    .filter((g: federation_gateway) => g.federation && g.id)
    .map((gateway: federation_gateway) => {
      fegGatewayRows.push({
        name: gateway.name,
        id: gateway.id,
        is_primary: activeFegGatewayId === gateway.id,
        health: health[gateway.id].status
          ? health[gateway.id]?.status === HEALTHY_STATUS
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
                onClick={() => navigate(currRow.id)}>
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
