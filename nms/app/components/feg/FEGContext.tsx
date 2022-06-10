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
import FEGNetworkContext from '../context/FEGNetworkContext';
import FEGSubscriberContext from '../context/FEGSubscriberContext';
import LoadingFiller from '../LoadingFiller';
import useMagmaAPI from '../../../api/useMagmaAPI';

import type {
  FederationGateway,
  FegNetwork,
  SubscriberState,
} from '../../../generated-ts';
import type {FederationGatewayHealthStatus} from '../GatewayUtils';

import MagmaAPI from '../../../api/MagmaAPI';
import {FetchFegSubscriberState} from '../../state/feg/SubscriberState';
import {GatewayId, NetworkId, NetworkType} from '../../../shared/types/network';
import {
  InitGatewayState,
  SetGatewayState,
} from '../../state/feg/EquipmentState';
import {UpdateNetworkState as UpdateFegNetworkState} from '../../state/feg/NetworkState';
import {useCallback, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

type Props = {
  networkId: NetworkId;
  networkType: NetworkType;
  children: React.ReactNode;
};

/**
 * Fetches and saves the subscriber session states of networks
 * serviced by this federation network and whose subscriber
 * information is not managed by the HSS.
 *
 * @param {network_id} networkId Id of the network
 * @param {network_type} networkType Type of the network
 */
export function FEGSubscriberContextProvider(props: Props) {
  const {networkId} = props;
  const [sessionState, setSessionState] = useState<
    Record<NetworkId, Record<string, SubscriberState>>
  >({});
  const [isLoading, setIsLoading] = useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const fetchFegState = async () => {
      if (networkId == null) {
        return;
      }
      const sessionState = await FetchFegSubscriberState({
        networkId,
        enqueueSnackbar,
      });
      setSessionState(sessionState);
      setIsLoading(false);
    };
    void fetchFegState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <FEGSubscriberContext.Provider
      value={{
        sessionState: sessionState,
        setSessionState: newSessionState => {
          return setSessionState(newSessionState);
        },
      }}>
      {props.children}
    </FEGSubscriberContext.Provider>
  );
}

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
        setState: (key, value?, newState?) => {
          return SetGatewayState({
            networkId,
            fegGateways,
            fegGatewaysHealthStatus,
            setFegGateways,
            setFegGatewaysHealthStatus,
            setActiveFegGatewayId,
            key,
            value,
            newState,
            enqueueSnackbar,
          });
        },
        health: fegGatewaysHealthStatus,
        activeFegGatewayId,
      }}>
      {props.children}
    </FEGGatewayContext.Provider>
  );
}

/**
 * Fetches and returns information about the federation network inside
 * a context provider.
 * @param {object} props: contains the network id and its type
 */
export function FEGNetworkContextProvider(props: Props) {
  const {networkId} = props;
  const [fegNetwork, setFegNetwork] = useState<Partial<FegNetwork>>({});
  const enqueueSnackbar = useEnqueueSnackbar();
  const {error, isLoading} = useMagmaAPI(
    MagmaAPI.federationNetworks.fegNetworkIdGet,
    {networkId: networkId},
    useCallback((response: Partial<FegNetwork>) => setFegNetwork(response), []),
  );

  if (error) {
    enqueueSnackbar?.('failed fetching network information', {
      variant: 'error',
    });
  }

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <FEGNetworkContext.Provider
      value={{
        state: fegNetwork,
        updateNetworks: props => {
          let refreshState = true;
          if (networkId !== props.networkId) {
            refreshState = false;
          }
          return UpdateFegNetworkState({
            networkId,
            setFegNetwork,
            refreshState,
            ...props,
          });
        },
      }}>
      {props.children}
    </FEGNetworkContext.Provider>
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
    <FEGNetworkContextProvider {...{networkId, networkType}}>
      <FEGSubscriberContextProvider {...{networkId, networkType}}>
        <FEGGatewayContextProvider {...{networkId, networkType}}>
          {props.children}
        </FEGGatewayContextProvider>
      </FEGSubscriberContextProvider>
    </FEGNetworkContextProvider>
  );
}
