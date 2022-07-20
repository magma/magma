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

import * as React from 'react';
import FEGGatewayContext from '../context/FEGGatewayContext';
import LoadingFiller from '../LoadingFiller';
import {FEGNetworkContextProvider} from '../context/FEGNetworkContext';
import {FEGSubscriberContextProvider} from '../context/FEGSubscriberContext';
import {useEffect, useState} from 'react';

import {
  FetchFegGateway,
  FetchFegGateways,
  InitGatewayState,
  SetGatewayState,
} from '../../state/feg/EquipmentState';
import {GatewayId, NetworkId, NetworkType} from '../../../shared/types/network';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import type {FederationGateway} from '../../../generated';
import type {FederationGatewayHealthStatus} from '../GatewayUtils';

type Props = {
  networkId: NetworkId;
  networkType: NetworkType;
  children: React.ReactNode;
};

/**
 * Fetches and returns the federation gateways, their health status and
 * the active federation gateway id.
 * @param {network_id} networkId Id of the network
 * @param {network_type} networkType Type of the network
 */
export function FEGGatewayContextProvider(props: Props) {
  const {networkId} = props;
  const [fegGateways, setFegGateways] = useState<
    Record<GatewayId, FederationGateway>
  >({});
  const [fegGatewaysHealthStatus, setFegGatewaysHealthStatus] = useState<
    Record<GatewayId, FederationGatewayHealthStatus>
  >({});
  const [activeFegGatewayId, setActiveFegGatewayId] = useState<GatewayId>('');
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      await InitGatewayState({
        networkId,
        setFegGateways,
        setFegGatewaysHealthStatus,
        setActiveFegGatewayId,
        enqueueSnackbar,
      });
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <FEGGatewayContext.Provider
      value={{
        state: fegGateways,
        setState: (key, value?) => {
          return SetGatewayState({
            networkId,
            fegGateways,
            fegGatewaysHealthStatus,
            setFegGateways,
            setFegGatewaysHealthStatus,
            setActiveFegGatewayId,
            key,
            value,
            enqueueSnackbar,
          });
        },
        refetch: (id?: GatewayId) => {
          if (id) {
            void FetchFegGateway({id, networkId, enqueueSnackbar}).then(
              response => {
                if (response) {
                  setFegGateways(gateways => ({
                    ...gateways,
                    [id]: response.fegGateway,
                  }));
                  setFegGatewaysHealthStatus(healthStatus => ({
                    ...healthStatus,
                    [id]: response.healthStatus,
                  }));
                }
              },
            );
          } else {
            void FetchFegGateways({networkId, enqueueSnackbar}).then(
              response => {
                if (response) {
                  setFegGateways(response.fegGateways);
                  setFegGatewaysHealthStatus(response.fegGatewaysHealthStatus);
                  setActiveFegGatewayId(response.activeFegGatewayId);
                }
              },
            );
          }
        },
        health: fegGatewaysHealthStatus,
        activeFegGatewayId,
      }}>
      {props.children}
    </FEGGatewayContext.Provider>
  );
}

/**
 * A context provider for federation networks. It is used in sharing
 * information like the network information or the gateways information.
 * @param {object} props contains the network id and its type
 */
export function FEGContextProvider(props: Props) {
  const {networkId, networkType} = props;

  return (
    <FEGNetworkContextProvider networkId={networkId}>
      <FEGSubscriberContextProvider {...{networkId, networkType}}>
        <FEGGatewayContextProvider {...{networkId, networkType}}>
          {props.children}
        </FEGGatewayContextProvider>
      </FEGSubscriberContextProvider>
    </FEGNetworkContextProvider>
  );
}
