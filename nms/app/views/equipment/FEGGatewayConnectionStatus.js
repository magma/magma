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

import DataGrid from '../../components/DataGrid';
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import React from 'react';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {
  AVAILABLE_STATUS,
  HEALTHY_STATUS,
  ServiceTypeEnum,
} from '../../components/GatewayUtils';
import {useContext} from 'react';
import {useParams} from 'react-router-dom';

/**
 * Displays the Connection status of the Gx/Gy servers, S6a server
 * and SWx server.
 */
export default function FEGGatewayConnectionStatus() {
  const params = useParams();
  const ctx = useContext(FEGGatewayContext);
  const gatewayId: string = nullthrows(params.gatewayId);
  const gwHealthStatus = ctx.health[gatewayId] || {};
  const getServiceHealthStatus = serviceStatus => {
    if (serviceStatus) {
      if (!(serviceStatus.service_state === AVAILABLE_STATUS)) {
        return ServiceTypeEnum.UNAVAILABLE_SERVICE;
      }
      return serviceStatus?.health_status === HEALTHY_STATUS
        ? ServiceTypeEnum.HEALTHY_SERVICE
        : ServiceTypeEnum.UNHEALTHY_SERVICE;
    }
    return ServiceTypeEnum.UNENABLED_SERVICE;
  };
  const isServiceStatusInactive = serviceStatus =>
    serviceStatus === ServiceTypeEnum.UNAVAILABLE_SERVICE ||
    serviceStatus === ServiceTypeEnum.UNENABLED_SERVICE;
  const isServiceStatusActive = serviceStatus =>
    serviceStatus === ServiceTypeEnum.HEALTHY_SERVICE;
  const gxGyConnectionStatus = getServiceHealthStatus(
    gwHealthStatus?.service_status?.SESSION_PROXY,
  );
  const swxConnectionStatus = getServiceHealthStatus(
    gwHealthStatus?.service_status?.SWX_PROXY,
  );
  const s6aConnectionStatus = getServiceHealthStatus(
    gwHealthStatus?.service_status?.S6A_PROXY,
  );

  const kpiData: DataRows[] = [
    [
      {
        category: 'Gx/Gy Watchdog',
        value: gxGyConnectionStatus,
        statusCircle: true,
        statusInactive: isServiceStatusInactive(gxGyConnectionStatus),
        status: isServiceStatusActive(gxGyConnectionStatus),
        tooltip: 'Connection status of Gx & Gy servers',
      },
      {
        category: 'SWx Watchdog',
        value: swxConnectionStatus,
        statusCircle: true,
        statusInactive: isServiceStatusInactive(swxConnectionStatus),
        status: isServiceStatusActive(swxConnectionStatus),
        tooltip: 'Connection status of SWx server',
      },
      {
        category: 'S6a Watchdog',
        value: s6aConnectionStatus,
        statusCircle: true,
        statusInactive: isServiceStatusInactive(s6aConnectionStatus),
        status: isServiceStatusActive(s6aConnectionStatus),
        tooltip: 'Connection status of S6a server',
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}
