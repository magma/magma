/**
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

import type {DataRows} from '../../components/DataGrid';

import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import GroupWorkIcon from '@material-ui/icons/GroupWork';
import MagmaV1API from '../../../generated/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Typography from '@material-ui/core/Typography';
import moment from 'moment';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

import {GatewayTypeEnum, HEALTHY_STATUS} from '../../components/GatewayUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(_ => ({
  paperRoot: {
    display: 'flex',
    padding: '10px 20px',
  },
  paperRow: {
    padding: '0px 20px',
    fontWeight: 'light',
  },
}));

/**
 * Displays the last fallover time, health of primary gateway & secondary gateway,
 * and the name of the primary/active gateway.
 */
export default function FEGClusterStatus() {
  const classes = useStyles();
  const params = useParams();
  const ctx = useContext(FEGGatewayContext);
  const networkId: string = nullthrows(params.networkId);
  const timeRange = '7d';
  const {response: lastFalloverResponse} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusQuery,
    {
      networkId: networkId,
      query: `max_over_time(timestamp(changes(active_gateway_changed_total[30s]) > 0)[${timeRange}:10s])`,
    },
  );
  const getGatewayHealthStatus = (fegGatewaysHealthStatus, gatewayId) => {
    const gatewayHealthStatus = fegGatewaysHealthStatus[gatewayId]?.status;
    if (gatewayId && gatewayHealthStatus) {
      // gateway exists and health status was fetched without error
      return gatewayHealthStatus == HEALTHY_STATUS
        ? GatewayTypeEnum.HEALTHY_GATEWAY
        : GatewayTypeEnum.UNHEALTHY_GATEWAY;
    }
    return 'N/A';
  };
  const getLastFalloverStatus = lastFalloverResponse => {
    let lastFalloverStatus = '-';
    let lastFalloverTime = 0;
    const lastFalloverResult = lastFalloverResponse?.data?.result || [];

    lastFalloverResult.map(res => {
      const curUpdate = parseFloat(res?.value?.[1]) || 0;
      // save the latest update
      lastFalloverTime = Math.max(lastFalloverTime, curUpdate);
    });
    lastFalloverTime &&
      (lastFalloverStatus = moment.unix(lastFalloverTime).calendar());
    return lastFalloverStatus;
  };
  const getSecondaryFegGatewayId = (fegGateways, activeFegGatewayId) => {
    const fegGatewaysId = Object.keys(fegGateways);
    if (fegGatewaysId.length > 1) {
      // has secondary gateway
      return fegGatewaysId[0] == activeFegGatewayId
        ? fegGatewaysId[1]
        : fegGatewaysId[0];
    }
    return '';
  };
  const isGatewayHealthStatusInactive = (
    fegGatewayId,
    fegGatewayHealthStatus,
  ) => {
    // is inactive if gateway doesn't exits or have no health status
    return !(fegGatewayId && fegGatewayHealthStatus?.status);
  };
  const fegGateways = ctx.state || {};
  const fegGatewaysHealthStatus = ctx.health || {};
  const activeFegGatewayId = ctx.activeFegGatewayId || '';
  const secondaryFegGatewayId = getSecondaryFegGatewayId(
    fegGateways,
    activeFegGatewayId,
  );
  const activeFegGatewayHealthStatus = getGatewayHealthStatus(
    fegGatewaysHealthStatus,
    activeFegGatewayId,
  );
  const secondaryFegGatewayHealthStatus = getGatewayHealthStatus(
    fegGatewaysHealthStatus,
    secondaryFegGatewayId,
  );
  const lastFalloverStatus = getLastFalloverStatus(lastFalloverResponse);

  const kpiData: DataRows[] = [
    [
      {
        category: 'Last Fallover Time',
        value: lastFalloverStatus,
        tooltip: 'The last time the active gateway of the network was changed',
      },
      {
        category: 'Primary Health',
        value: activeFegGatewayHealthStatus,
        statusCircle: true,
        // make kpi inactive if gateway doesn't exist or had no health status
        statusInactive: isGatewayHealthStatusInactive(
          activeFegGatewayId,
          fegGatewaysHealthStatus[activeFegGatewayId],
        ),
        status: activeFegGatewayHealthStatus === HEALTHY_STATUS,
        tooltip: 'Health of primary federation gateway',
      },
      {
        category: 'Secondary Health',
        value: secondaryFegGatewayHealthStatus,
        statusCircle: true,
        // make kpi inactive if secondary gateway doesn't exist or had no health status
        statusInactive: isGatewayHealthStatusInactive(
          secondaryFegGatewayId,
          fegGatewaysHealthStatus[secondaryFegGatewayId],
        ),
        status: secondaryFegGatewayHealthStatus === HEALTHY_STATUS,
        tooltip: 'Health of secondary federation gateway',
      },
    ],
  ];

  return (
    <>
      <CardTitleRow icon={GroupWorkIcon} label="Cluster Status" />
      <DataGrid data={kpiData} />
      <Paper className={classes.paperRoot}>
        <Typography color="textSecondary" className={classes.paperRow}>
          Primary
        </Typography>
        <Typography
          color="textSecondary"
          className={classes.paperRow}
          data-testid="Primary Gateway Name">
          {fegGateways[activeFegGatewayId]
            ? fegGateways[activeFegGatewayId].name
            : 'N/A'}
        </Typography>
      </Paper>
    </>
  );
}
