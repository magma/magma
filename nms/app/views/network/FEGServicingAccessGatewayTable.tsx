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
import Link from '@material-ui/core/Link';
import LoadingFiller from '../../components/LoadingFiller';
import React, {useEffect, useState} from 'react';
import nullthrows from '../../../shared/util/nullthrows';
import {FetchGateways} from '../../state/lte/EquipmentState';
import {
  GatewayId,
  GatewayName,
  NetworkId,
  NetworkName,
} from '../../../shared/types/network';
import {GatewayTypeEnum} from '../../components/GatewayUtils';
import {getServicedAccessNetworks} from '../../components/FEGServicingAccessGatewayKPIs';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {FegLteNetwork, LteGateway} from '../../../generated-ts';
import type {OptionsObject} from 'notistack';

type ServicingAccessGatewayRowType = {
  networkId: NetworkId;
  networkName: NetworkName;
  gatewayId: GatewayId;
  gatewayName: GatewayName;
  gatewayHealth: string;
};

/**
 * Returns an array which holds information about each serviced
 * access gateways servicied by this federation network.
 * The information about the serviced gateways includes the gateway
 * id, name, network name & network id under which the gateway exists.
 *
 * @param {Array<FegLteNetwork>} servicedAccessNetworks List of federated LTE networks serviced by this federation network.
 * @param {(msg, cfg) => ?(string | number),} enqueueSnackbar A snackbar to display errors.
 * @returns An Array of the serviced access gateways with information about each.
 */
async function getServicedAccessGatewaysInfo(
  servicedAccessNetworks: Array<FegLteNetwork>,
  enqueueSnackbar: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined,
): Promise<Array<ServicingAccessGatewayRowType>> {
  const newServicedAccessGatewaysInfo: Array<ServicingAccessGatewayRowType> = [];

  for (const servicedAccessNetwork of servicedAccessNetworks) {
    const servicedAccessGateways = await FetchGateways({
      networkId: servicedAccessNetwork.id,
      enqueueSnackbar,
    });
    //Add the gateways of the serviced network
    if (servicedAccessGateways) {
      Object.keys(servicedAccessGateways).map(servicedAccessGatewayId => {
        const newServicedAccessGatewayInfo: ServicingAccessGatewayRowType = {
          networkId: servicedAccessNetwork.id,
          networkName: servicedAccessNetwork.name,
          gatewayId: servicedAccessGatewayId,
          gatewayName:
            servicedAccessGateways[servicedAccessGatewayId]?.name || '',
          gatewayHealth: servicedAccessGateways[servicedAccessGatewayId]
            ?.checked_in_recently
            ? GatewayTypeEnum.HEALTHY_GATEWAY
            : GatewayTypeEnum.UNHEALTHY_GATEWAY,
        };
        newServicedAccessGatewaysInfo.push(newServicedAccessGatewayInfo);
      });
    }
  }

  return newServicedAccessGatewaysInfo;
}

/**
 * Returns a table consisting of the serviced access gateways alongside
 * the serviced network in which they are under.
 */

export default function ServicingAccessGatewayInfo() {
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();
  const networkId: string = nullthrows(params.networkId);
  const [servicedAccessGatewaysInfo, setServicedAccessGatewaysInfo] = useState<
    Array<ServicingAccessGatewayRowType>
  >();
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

    void fetchServicedAccessGateways();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <div>
      <ActionTable
        data={servicedAccessGatewaysInfo!}
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
