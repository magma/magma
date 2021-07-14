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

import type {EnqueueSnackbarOptions} from 'notistack';
import type {
  feg_lte_network,
  gateway_id,
  gateway_name,
  lte_gateway,
  network_id,
  network_name,
} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import Link from '@material-ui/core/Link';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import React, {useEffect, useState} from 'react';
import isGatewayHealthy from '../../components/GatewayUtils';
import nullthrows from '@fbcnms/util/nullthrows';

import {FetchGateways} from '../../state/lte/EquipmentState';
import {
  HEALTHY_GATEWAY,
  UNHEALTHY_GATEWAY,
} from '../../components/GatewayUtils';
import {getServicedAccessNetworks} from '../../components/FEGServicingAccessGatewayKPIs';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

type ServicingAccessGatewayRowType = {
  networkId: network_id,
  networkName: network_name,
  gatewayId: gateway_id,
  gatewayName: gateway_name,
  gatewayHealth: string,
};

/**
 * Returns an array which holds information about each serviced
 * access gateways servicied by this federation network.
 * The information about the serviced gateways includes the gateway
 * id, name, network name & network id under which the gateway exists.
 *
 * @param {Array<feg_lte_network>} servicedAccessNetworks List of federated LTE networks serviced by this federation network.
 * @param {(msg, cfg) => ?(string | number),} enqueueSnackbar A snackbar to display errors.
 * @returns An Array of the serviced access gateways with information about each.
 */
async function getServicedAccessGatewaysInfo(
  servicedAccessNetworks: Array<feg_lte_network>,
  enqueueSnackbar: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
): Promise<Array<ServicingAccessGatewayRowType>> {
  const newServicedAccessGatewaysInfo = [];
  for (const servicedAccessNetwork of servicedAccessNetworks) {
    const servicedAccessGateways: {
      [string]: lte_gateway,
    } = await FetchGateways({
      networkId: servicedAccessNetwork.id,
      undefined,
      enqueueSnackbar,
    });
    //Add the gateways of the serviced network
    Object.keys(servicedAccessGateways).map(servicedAccessGatewayId => {
      const newServicedAccessGatewayInfo: ServicingAccessGatewayRowType = {
        networkId: servicedAccessNetwork.id,
        networkName: servicedAccessNetwork.name,
        gatewayId: servicedAccessGatewayId,
        gatewayName:
          servicedAccessGateways[servicedAccessGatewayId]?.name || '',
        gatewayHealth: isGatewayHealthy(
          servicedAccessGateways[servicedAccessGatewayId] || {},
        )
          ? HEALTHY_GATEWAY
          : UNHEALTHY_GATEWAY,
      };
      newServicedAccessGatewaysInfo.push(newServicedAccessGatewayInfo);
    });
  }
  return newServicedAccessGatewaysInfo;
}

/**
 * Returns a table consisting of the serviced access gateways alongside
 * the serviced network in which they are under.
 */
export default function ServicingAccessGatewayInfo() {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const networkId: string = nullthrows(match.params.networkId);
  const [servicedAccessGatewaysInfo, setServicedAccessGatewaysInfo] = useState<
    Array<ServicingAccessGatewayRowType>,
  >([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchServicedAccessGateways = async () => {
      try {
        const servicedAccessNetworks = await getServicedAccessNetworks(
          networkId,
          enqueueSnackbar,
        );
        const newServicedAccessGatewaysInfo = await getServicedAccessGatewaysInfo(
          servicedAccessNetworks,
          enqueueSnackbar,
        );
        setServicedAccessGatewaysInfo(newServicedAccessGatewaysInfo);
        setIsLoading(false);
      } catch (e) {
        enqueueSnackbar?.(
          'failed fetching servicing access gateway information',
          {
            variant: 'error',
          },
        );
      }
    };
    fetchServicedAccessGateways();
  }, [networkId, enqueueSnackbar]);
  if (isLoading) {
    return <LoadingFiller />;
  }
  return (
    <div>
      <ActionTable
        data={servicedAccessGatewaysInfo}
        columns={[
          {title: 'Access Network', field: 'networkName'},
          {title: 'Access Gateway Id', field: 'gatewayId'},
          {
            title: 'Access Gateway Name',
            field: 'gatewayName',
            render: currRow => (
              <Link
                variant="body2"
                component="button"
                onClick={() => {
                  window.open(
                    `${window.location.origin}/nms/${currRow.networkId}/equipment/overview/gateway/${currRow.gatewayId}`,
                  );
                }}>
                {currRow.gatewayName}
              </Link>
            ),
          },
          {title: 'Access Gateway Health', field: 'gatewayHealth'},
        ]}
        options={{
          actionsColumnIndex: -1,
          pageSizeOptions: [5, 10],
        }}
      />
    </div>
  );
}
