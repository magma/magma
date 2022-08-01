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
import {
  FederatedNetworkConfigs,
  FegLteNetwork,
  LteNetwork,
  NetworkDnsConfig,
  NetworkEpcConfigs,
  NetworkRanConfigs,
  NetworkSubscriberConfig,
} from '../../generated';

import LoadingFiller from '../components/LoadingFiller';
import MagmaAPI from '../api/MagmaAPI';
import NetworkContext from './NetworkContext';
import React, {useContext, useEffect, useState} from 'react';
import {FEG_LTE, NetworkId} from '../../shared/types/network';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

// TODO[TS-migration] This should probably be something like Partial<UpdateLteNetworkParams> | Partial<UpdateFegLteNetworkParams>
export type UpdateNetworkContextParams = Partial<
  UpdateLteNetworkParams & UpdateFegLteNetworkParams
>;

export type LteNetworkContextType = {
  state: Partial<LteNetwork & FegLteNetwork>;
  updateNetworks: (props: UpdateNetworkContextParams) => Promise<void>;
};

type UpdateFegLteNetworkParams = {
  networkId: NetworkId;
  lteNetwork?: FegLteNetwork & {
    subscriber_config: NetworkSubscriberConfig;
  };
  federation?: FederatedNetworkConfigs;
  epcConfigs?: NetworkEpcConfigs;
  lteRanConfigs?: NetworkRanConfigs;
  subscriberConfig?: NetworkSubscriberConfig;
  setLteNetwork: (
    lteNetwork: FegLteNetwork & {
      subscriber_config: NetworkSubscriberConfig;
    },
  ) => void;
  refreshState: boolean;
};

type UpdateLteNetworkParams = {
  networkId: NetworkId;
  lteNetwork?: LteNetwork;
  epcConfigs?: NetworkEpcConfigs;
  lteRanConfigs?: NetworkRanConfigs;
  lteDnsConfig?: NetworkDnsConfig;
  subscriberConfig?: NetworkSubscriberConfig;
  setLteNetwork: (lteNetwork: LteNetwork) => void;
  refreshState: boolean;
};

const LteNetworkContext = React.createContext<LteNetworkContextType>(
  {} as LteNetworkContextType,
);

async function updateFegLteNetworkState(params: UpdateFegLteNetworkParams) {
  const {networkId, setLteNetwork} = params;
  const requests = [];
  if (params.lteNetwork) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdPut({
        networkId: networkId,
        lteNetwork: {...params.lteNetwork},
      }),
    );
  }
  if (params.federation) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdFederationPut({
        networkId: networkId,
        config: params.federation,
      }),
    );
  }
  if (params.subscriberConfig) {
    requests.push(
      await MagmaAPI.federatedLTENetworks.fegLteNetworkIdSubscriberConfigPut({
        networkId: params.networkId,
        record: params.subscriberConfig,
      }),
    );
  }
  if (params.epcConfigs != null || params.lteRanConfigs != null) {
    await updateLteNetworkState({
      networkId,
      setLteNetwork: () => {},
      epcConfigs: params.epcConfigs,
      lteRanConfigs: params.lteRanConfigs,
      refreshState: false,
    });
  }
  await Promise.all(requests);
  if (params.refreshState) {
    const [fegLteResp, fegLteSubscriberConfigResp] = await Promise.allSettled([
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet({
        networkId,
      }),
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdSubscriberConfigGet({
        networkId,
      }),
    ]);
    if (fegLteResp.status === 'fulfilled') {
      let subscriber_config = {};
      if (fegLteSubscriberConfigResp.status === 'fulfilled') {
        subscriber_config = fegLteSubscriberConfigResp.value.data;
      }
      setLteNetwork({...fegLteResp.value.data, subscriber_config});
    }
  }
}

async function updateLteNetworkState(params: UpdateLteNetworkParams) {
  const {networkId, setLteNetwork} = params;
  const requests = [];

  if (params.lteNetwork) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdPut({
        networkId: networkId,
        lteNetwork: {...params.lteNetwork},
      }),
    );
  }

  if (params.epcConfigs) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdCellularEpcPut({
        networkId: params.networkId,
        config: params.epcConfigs,
      }),
    );
  }
  if (params.lteRanConfigs) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdCellularRanPut({
        networkId: params.networkId,
        config: params.lteRanConfigs,
      }),
    );
  }
  if (params.lteDnsConfig) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdDnsPut({
        networkId: params.networkId,
        config: params.lteDnsConfig,
      }),
    );
  }
  if (params.subscriberConfig) {
    requests.push(
      MagmaAPI.lteNetworks.lteNetworkIdSubscriberConfigPut({
        networkId: params.networkId,
        record: params.subscriberConfig,
      }),
    );
  }
  // TODO(andreilee): Provide a way to handle errors here
  await Promise.all(requests);
  if (params.refreshState) {
    setLteNetwork(
      (
        await MagmaAPI.lteNetworks.lteNetworkIdGet({
          networkId,
        })
      ).data,
    );
  }
}

export function LteNetworkContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
  const {networkId} = props;
  const networkCtx = useContext(NetworkContext);
  const [lteNetwork, setLteNetwork] = useState<
    Partial<LteNetwork & FegLteNetwork>
  >({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkCtx.networkType === FEG_LTE) {
          const [
            fegLteResp,
            fegLteSubscriberConfigResp,
          ] = await Promise.allSettled([
            MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet({networkId}),
            MagmaAPI.federatedLTENetworks.fegLteNetworkIdSubscriberConfigGet({
              networkId,
            }),
          ]);
          if (fegLteResp.status === 'fulfilled') {
            let subscriber_config = {};
            if (fegLteSubscriberConfigResp.status === 'fulfilled') {
              subscriber_config = fegLteSubscriberConfigResp.value.data;
            }
            setLteNetwork({...fegLteResp.value.data, subscriber_config});
          }
        } else {
          setLteNetwork(
            (await MagmaAPI.lteNetworks.lteNetworkIdGet({networkId})).data,
          );
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching network information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, networkCtx, enqueueSnackbar]);

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
          if (networkCtx.networkType === FEG_LTE) {
            return updateFegLteNetworkState({
              networkId,
              setLteNetwork,
              refreshState,
              ...props,
            });
          } else {
            return updateLteNetworkState({
              networkId,
              setLteNetwork,
              refreshState,
              ...props,
            });
          }
        },
      }}>
      {props.children}
    </LteNetworkContext.Provider>
  );
}

export default LteNetworkContext;
