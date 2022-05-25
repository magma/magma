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
 */

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from './DataGrid';
import LoadingFiller from './LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import nullthrows from '../../shared/util/nullthrows';
import {FetchGateways} from '../state/lte/EquipmentState';
import {NetworkId} from '../../shared/types/network';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {DataRows} from './DataGrid';
import type {FegLteNetwork} from '../../generated-ts';
import type {OptionsObject} from 'notistack';

/**
 * Returns the list of federated lte networks serviced by the federation
 *  network with the id: federationNetworkId.
 * @param {NetworkId} federationNetworkId id of the federation network
 * @param {function} enqueueSnackbar snackbar used to display information
 */
export async function getServicedAccessNetworks(
  federationNetworkId: NetworkId,
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => (string | number) | null | undefined,
): Promise<Array<FegLteNetwork>> {
  const servicedAccessNetworks: Array<FegLteNetwork> = [];
  const fegLteNetworkIdList = (await MagmaAPI.federatedLTENetworks.fegLteGet())
    .data;
  const requests = fegLteNetworkIdList.map(async fegLteNetworkId => {
    try {
      return (
        await MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet({
          networkId: fegLteNetworkId,
        })
      ).data;
    } catch (e) {
      enqueueSnackbar?.(
        'failed fetching network information for ' + fegLteNetworkId,
        {
          variant: 'error',
        },
      );
      return;
    }
  });
  const fegLteNetworks = await Promise.all(requests);
  fegLteNetworks.filter(Boolean).forEach(fegLteNetwork => {
    if (fegLteNetwork?.federation?.feg_network_id === federationNetworkId) {
      servicedAccessNetworks.push(fegLteNetwork);
    }
  });
  return servicedAccessNetworks;
}

/**
 * Returns the total count of access gateways serviced by the
 * federation network.
 */

export default function ServicingAccessGatewayKPIs() {
  const params = useParams();
  const networkId = nullthrows(params.networkId);
  const [isLoading, setIsLoading] = useState(true);
  const [
    servicedAccessGatewaysCount,
    setServicedAccessGatewaysCount,
  ] = useState(0);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const getServicedAccessGatewaysCount = async (
      servicedAccessNetworks: Array<FegLteNetwork>,
    ): Promise<number> => {
      let totalServicedAccessGateways = 0;

      for (const servicedAccessNetwork of servicedAccessNetworks) {
        const servicedAccessGateways = await FetchGateways({
          networkId: servicedAccessNetwork.id,
          enqueueSnackbar,
        });
        if (servicedAccessGateways) {
          totalServicedAccessGateways += Object.keys(
            servicedAccessGateways,
          ).filter(Boolean).length;
        }
      }

      return totalServicedAccessGateways;
    };

    const fetchServicedAccessGateways = async () => {
      try {
        const servicedAccessNetworks = await getServicedAccessNetworks(
          networkId,
          enqueueSnackbar,
        );
        const totalServicedAccessGateways = await getServicedAccessGatewaysCount(
          servicedAccessNetworks,
        );
        setServicedAccessGatewaysCount(totalServicedAccessGateways);
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
  const data: Array<DataRows> = [
    [
      {
        icon: CellWifiIcon,
        value: 'Serviced Access Gateways',
      },
      {
        category: 'Gateway Count',
        value: servicedAccessGatewaysCount,
        tooltip: 'Number of gateways checked in recently',
      },
    ],
  ];

  if (isLoading) {
    return <LoadingFiller />;
  }

  return <DataGrid data={data} />;
}
