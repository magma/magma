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

import * as React from 'react';
import FEGGatewayContext from '../context/FEGGatewayContext';
// $FlowFixMe migrated to typescript
import FEGNetworkContext from '../context/FEGNetworkContext';
import FEGSubscriberContext from '../context/FEGSubscriberContext';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

import type {FederationGatewayHealthStatus} from '../../components/GatewayUtils';
import type {
  federation_gateway,
  feg_network,
  gateway_id,
  network_id,
  network_type,
  subscriber_state,
} from '../../../generated/MagmaAPIBindings';

import {FetchFegSubscriberState} from '../../state/feg/SubscriberState';
import {
  InitGatewayState,
  SetGatewayState,
} from '../../state/feg/EquipmentState';
// $FlowFixMe migrated to typescript
import {UpdateNetworkState as UpdateFegNetworkState} from '../../state/feg/NetworkState';
import {useCallback, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

type Props = {
  networkId: network_id,
  networkType: network_type,
  children: React.Node,
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
  const [sessionState, setSessionState] = useState<{
    [networkId: network_id]: {[string]: subscriber_state},
  }>({});
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
    fetchFegState();
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
  const [fegGateways, setFegGateways] = useState<{
    [gateway_id]: federation_gateway,
  }>({});
  const [fegGatewaysHealthStatus, setFegGatewaysHealthStatus] = useState<{
    [gateway_id]: FederationGatewayHealthStatus,
  }>({});
  const [activeFegGatewayId, setActiveFegGatewayId] = useState<gateway_id>('');
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
    fetchState();
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
  const [fegNetwork, setFegNetwork] = useState<$Shape<feg_network>>({});
  const enqueueSnackbar = useEnqueueSnackbar();
  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getFegByNetworkId,
    {networkId: networkId},
    useCallback(response => setFegNetwork(response), []),
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
