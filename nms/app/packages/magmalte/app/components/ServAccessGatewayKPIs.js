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

import type {DataRows} from './DataGrid';
import type {EnqueueSnackbarOptions,federation_gateway,feg_lte_network, network_id} from '@fbcnms/magma-api';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from './DataGrid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import isGatewayHealthy from './GatewayUtils';
import nullthrows from '@fbcnms/util/nullthrows';

import {useContext, useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

export async function getServAccessGateways(networkId:network_id, enqueueSnackbar?: (msg: string, cfg: EnqueueSnackbarOptions,) => ?(string | number),): Promise<Array<feg_lte_network>> {
    const servedAccessGateways = [];
    //console.log("First")
    const fegLteNetworkIdList =  await MagmaV1API.getFegLte();
    //console.log("Reached here: " + fegLteNetworkIdList);
    const requests = fegLteNetworkIdList.map(fegLteNetworkId => {
        try {
          return MagmaV1API.getFegLteByNetworkId({
            fegLteNetworkId,
          });
        } catch (e) {
            console.log("error: " + e);
          enqueueSnackbar?.('failed fetching tier information for ' + fegLteNetworkId, {
            variant: 'error',
          });
          return;
        }
      });

    const fegLteNetworks = await Promise.all(requests);
    fegLteNetworks.filter(Boolean).forEach(fegLteNetwork => {
        if (fegLteNetwork?.federation?.feg_network_id == networkId) {
            servedAccessGateways.push(fegLteNetwork);
        }
    });
    return servedAccessGateways;
}

export default function ServAccessGatewayKPIs() {
    const {match} = useRouter();
    const networkId = nullthrows(match.params.networkId);
    const [isLoading, setIsLoading] = useState(true);
    const [servedAccessGateways, setServedAccessGateways] = useState([]);
    const enqueueSnackbar = useEnqueueSnackbar();
    useEffect(() => {
        const fetchServicedAccessGateways = async () => {
            try {
                const servedAccessGateways = await getServAccessGateways(networkId);
                setServedAccessGateways(servedAccessGateways);
                console.log("Served access gateways: ");
                servedAccessGateways.map(s => {
                    console.log(s);
                });
                setIsLoading(false);
            }
            catch (e) {
                console.log("Error: " + e);
                enqueueSnackbar?.('failed fetching servicing access gateway information', {
                variant: 'error',
                });
            }
        }
        fetchServicedAccessGateways();
    }, [networkId]);
    const data: DataRows[] = [
        [
          {
            icon: CellWifiIcon,
            value: 'Federation Gateway',
          },
          {
            category: 'Severe Events',
            value: 0,
            tooltip: 'Severe Events reported by the gateway',
          },
          {
            category: 'Connected',
            value: servedAccessGateways.length || 0,
            tooltip: 'Number of gateways checked in within last 5 minutes',
          },
          {
            category: 'Disconnected',
            value: 0,
            tooltip: 'Number of gateways not checked in within last 5 minutes',
          },
        ],
      ];
    if (isLoading) {
        return <LoadingFiller />
    }
    return <DataGrid data={data} />;

}
