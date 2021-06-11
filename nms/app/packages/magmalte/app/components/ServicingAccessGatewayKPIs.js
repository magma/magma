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

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from './DataGrid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import type {DataRows} from './DataGrid';
import type {EnqueueSnackbarOptions} from 'notistack';
import type {feg_lte_network, network_id} from '@fbcnms/magma-api';

import {FetchGateways} from '../state/lte/EquipmentState';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

export async function getServicedAccessNetworks(
  networkId: network_id,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
): Promise<Array<feg_lte_network>> {
  const servicedAccessNetworks = [];
  const fegLteNetworkIdList = await MagmaV1API.getFegLte();
  const requests = fegLteNetworkIdList.map(async fegLteNetworkId => {
    try {
      return await MagmaV1API.getFegLteByNetworkId({
        networkId: fegLteNetworkId,
      });
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
    if (fegLteNetwork?.federation?.feg_network_id === networkId) {
      servicedAccessNetworks.push(fegLteNetwork);
    }
  });
  return servicedAccessNetworks;
}

export default function ServicingAccessGatewayKPIs() {
  const {match} = useRouter();
  const networkId = nullthrows(match.params.networkId);
  const [isLoading, setIsLoading] = useState(true);
  const [
    servicedAccessGatewaysCount,
    setServicedAccessGatewaysCount,
  ] = useState(0);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const getServicedAccessGatewaysCount = async (
      servicedAccessNetworks: Array<feg_lte_network>,
    ): Promise<number> => {
      let totalServicedAccessGateways = 0;
      for (const servicedAccessNetwork of servicedAccessNetworks) {
        const servicedAccessGateways = await FetchGateways({
          networkId: servicedAccessNetwork.id,
          undefined,
          enqueueSnackbar,
        });
        totalServicedAccessGateways += Object.keys(
          servicedAccessGateways,
        ).filter(Boolean).length;
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
    fetchServicedAccessGateways();
  }, [networkId, enqueueSnackbar]);
  const data: DataRows[] = [
    [
      {
        icon: CellWifiIcon,
        value: 'Serviced Access Gateways',
      },
      {
        category: 'Gateway Counts',
        value: servicedAccessGatewaysCount || 0,
        tooltip: 'Number of gateways checked in within last 5 minutes',
      },
    ],
  ];
  if (isLoading) {
    return <LoadingFiller />;
  }
  return <DataGrid data={data} />;
}
