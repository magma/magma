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
import ApnContext from '../context/ApnContext';
import EnodebContext from '../context/EnodebContext';
import GatewayContext from '../context/GatewayContext';
import GatewayTierContext from '../context/GatewayTierContext';
import InitSubscriberState from '../../state/lte/SubscriberState';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import LteNetworkContext from '../context/LteNetworkContext';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PolicyContext from '../context/PolicyContext';
import SubscriberContext from '../context/SubscriberContext';

import {
  InitEnodeState,
  InitTierState,
  SetEnodebState,
  SetGatewayState,
  SetTierState,
  UpdateGateway,
} from '../../state/lte/EquipmentState';
import {SetApnState} from '../../state/lte/ApnState';
import {SetPolicyState} from '../../state/lte/PolicyState';
import {UpdateNetworkState} from '../../state/lte/NetworkState';
import {
  getSubscriberGatewayMap,
  setSubscriberState,
} from '../../state/lte/SubscriberState';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import type {EnodebInfo} from '../lte/EnodebUtils';
import type {
  apn,
  lte_gateway,
  lte_network,
  mutable_subscriber,
  network_id,
  network_ran_configs,
  policy_rule,
  subscriber_id,
  tier,
} from '@fbcnms/magma-api';

type Props = {
  networkId: network_id,
  children: React.Node,
};

export function GatewayContextProvider(props: Props) {
  const {networkId} = props;
  const [lteGateways, setLteGateways] = useState<{[string]: lte_gateway}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      const lteGateways = await MagmaV1API.getLteByNetworkIdGateways({
        networkId,
      });
      setLteGateways(lteGateways);
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <GatewayContext.Provider
      value={{
        state: lteGateways,
        setState: (key, value?) => {
          return SetGatewayState({
            lteGateways,
            setLteGateways,
            networkId,
            key,
            value,
          });
        },
        updateGateway: props =>
          UpdateGateway({networkId, setLteGateways, ...props}),
      }}>
      {props.children}
    </GatewayContext.Provider>
  );
}

export function EnodebContextProvider(props: Props) {
  const {networkId} = props;
  const [enbInfo, setEnbInfo] = useState<{[string]: EnodebInfo}>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      if (networkId == null) {
        return;
      }
      const [lteRanConfigsResp] = await Promise.allSettled([
        MagmaV1API.getLteByNetworkIdCellularRan({networkId}),
        InitEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
      ]);
      if (lteRanConfigsResp.value) {
        setLteRanConfigs(lteRanConfigsResp.value);
      }
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }
  return (
    <EnodebContext.Provider
      value={{
        state: {enbInfo},
        lteRanConfigs: lteRanConfigs,
        setState: (key, value?) =>
          SetEnodebState({enbInfo, setEnbInfo, networkId, key, value}),
        setLteRanConfigs: lteRanConfigs => setLteRanConfigs(lteRanConfigs),
      }}>
      {props.children}
    </EnodebContext.Provider>
  );
}

export function SubscriberContextProvider(props: Props) {
  const {networkId} = props;
  const [subscriberMap, setSubscriberMap] = useState({});
  const [subscriberMetrics, setSubscriberMetrics] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchLteState = async () => {
      if (networkId == null) {
        return;
      }
      await InitSubscriberState({
        networkId,
        setSubscriberMap,
        setSubscriberMetrics,
        enqueueSnackbar,
      }),
        setIsLoading(false);
    };
    fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <SubscriberContext.Provider
      value={{
        state: subscriberMap,
        metrics: subscriberMetrics,
        gwSubscriberMap: getSubscriberGatewayMap(subscriberMap),
        setState: (key: subscriber_id, value?: mutable_subscriber) =>
          setSubscriberState({
            networkId,
            subscriberMap,
            setSubscriberMap,
            key,
            value,
          }),
      }}>
      {props.children}
    </SubscriberContext.Provider>
  );
}

export function GatewayTierContextProvider(props: Props) {
  const {networkId} = props;
  const [tiers, setTiers] = useState<{[string]: tier}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>([]);

  useEffect(() => {
    const fetchState = async () => {
      if (networkId == null) {
        return;
      }
      const [stableChannelResp] = await Promise.allSettled([
        MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
        InitTierState({networkId, setTiers, enqueueSnackbar}),
      ]);
      if (stableChannelResp.value) {
        setSupportedVersions(
          stableChannelResp.value.supported_versions.reverse(),
        );
      }
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <GatewayTierContext.Provider
      value={{
        state: {supportedVersions, tiers},
        setState: (key, value?) =>
          SetTierState({tiers, setTiers, networkId, key, value}),
      }}>
      {props.children}
    </GatewayTierContext.Provider>
  );
}

export function PolicyProvider(props: Props) {
  const {networkId} = props;
  const [policies, setPolicies] = useState<{[string]: policy_rule}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      setPolicies(
        await MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull({
          networkId,
        }),
      );
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <PolicyContext.Provider
      value={{
        state: policies,
        setState: (key, value?) => {
          return SetPolicyState({
            policies,
            setPolicies,
            networkId,
            key,
            value,
          });
        },
      }}>
      {props.children}
    </PolicyContext.Provider>
  );
}

export function ApnProvider(props: Props) {
  const {networkId} = props;
  const [apns, setApns] = useState<{[string]: apn}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      setApns(
        await MagmaV1API.getLteByNetworkIdApns({
          networkId,
        }),
      );
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <ApnContext.Provider
      value={{
        state: apns,
        setState: (key, value?) => {
          return SetApnState({
            apns,
            setApns,
            networkId,
            key,
            value,
          });
        },
      }}>
      {props.children}
    </ApnContext.Provider>
  );
}

export function LteNetworkContextProvider(props: Props) {
  const {networkId} = props;
  const [lteNetwork, setLteNetwork] = useState<lte_network>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      const lteNetwork = await MagmaV1API.getLteByNetworkId({
        networkId,
      });
      setLteNetwork(lteNetwork);
      setIsLoading(false);
    };
    fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <LteNetworkContext.Provider
      value={{
        state: lteNetwork,
        updateNetworks: props => {
          let refreshState = true;
          if (networkId !== props.networkId) {
            refreshState = false;
          }
          return UpdateNetworkState({
            networkId,
            setLteNetwork,
            refreshState,
            ...props,
          });
        },
      }}>
      {props.children}
    </LteNetworkContext.Provider>
  );
}

export function LteContextProvider(props: Props) {
  const {networkId} = props;
  return (
    <LteNetworkContextProvider networkId={networkId}>
      <PolicyProvider networkId={networkId}>
        <ApnProvider networkId={networkId}>
          <SubscriberContextProvider networkId={networkId}>
            <GatewayTierContextProvider networkId={networkId}>
              <EnodebContextProvider networkId={networkId}>
                <GatewayContextProvider networkId={networkId}>
                  {props.children}
                </GatewayContextProvider>
              </EnodebContextProvider>
            </GatewayTierContextProvider>
          </SubscriberContextProvider>
        </ApnProvider>
      </PolicyProvider>
    </LteNetworkContextProvider>
  );
}
