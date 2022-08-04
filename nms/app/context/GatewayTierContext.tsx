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
import LoadingFiller from '../components/LoadingFiller';
import MagmaAPI from '../api/MagmaAPI';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {NetworkId} from '../../shared/types/network';
import {Tier} from '../../generated';
import {useEffect, useState} from 'react';

type GatewayTierState = {
  tiers: Record<string, Tier>;
  supportedVersions: Array<string>;
};
type GatewayTierContextType = {
  state: GatewayTierState;
  setState: (key: string, val?: Tier) => Promise<void>;
};

const GatewayTierContext = React.createContext<GatewayTierContextType>(
  {} as GatewayTierContextType,
);

async function initTierState(params: {
  networkId: NetworkId;
  setTiers: (tiers: Record<string, Tier>) => void;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {networkId, setTiers, enqueueSnackbar} = params;
  let tierIdList: Array<string> = [];
  try {
    tierIdList = (
      await MagmaAPI.upgrades.networksNetworkIdTiersGet({
        networkId,
      })
    ).data;
  } catch (e) {
    enqueueSnackbar?.('failed fetching tier information', {
      variant: 'error',
    });
  }

  const requests = tierIdList.map(tierId => {
    try {
      return MagmaAPI.upgrades.networksNetworkIdTiersTierIdGet({
        networkId,
        tierId,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching tier information for ' + tierId, {
        variant: 'error',
      });
      return;
    }
  });

  const tierResponse = await Promise.all(requests);
  const tiers: Record<string, Tier> = {};
  tierResponse
    .filter(Boolean)
    .map(res => res!.data)
    .forEach(item => {
      tiers[item.id] = item;
    });
  setTiers(tiers);
}

async function setTierState(params: {
  networkId: NetworkId;
  tiers: Record<string, Tier>;
  setTiers: (tiers: Record<string, Tier>) => void;
  key: string;
  value?: Tier;
}) {
  const {networkId, tiers, setTiers, key, value} = params;

  if (value != null) {
    if (!(key in tiers)) {
      await MagmaAPI.upgrades.networksNetworkIdTiersPost({
        networkId: networkId,
        tier: value,
      });
    } else {
      await MagmaAPI.upgrades.networksNetworkIdTiersTierIdPut({
        networkId: networkId,
        tierId: key,
        tier: value,
      });
    }
    setTiers({...tiers, [key]: value});
  } else {
    await MagmaAPI.upgrades.networksNetworkIdTiersTierIdDelete({
      networkId: networkId,
      tierId: key,
    });
    const newTiers = {...tiers};
    delete newTiers[key];
    setTiers(newTiers);
  }
}

export function GatewayTierContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
  const {networkId} = props;
  const [tiers, setTiers] = useState<Record<string, Tier>>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [supportedVersions, setSupportedVersions] = useState<Array<string>>([]);

  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        const [stableChannelResp] = await Promise.allSettled([
          MagmaAPI.upgrades.channelsChannelIdGet({channelId: 'stable'}),
          initTierState({networkId, setTiers, enqueueSnackbar}),
        ]);
        if (stableChannelResp.status === 'fulfilled') {
          setSupportedVersions(
            stableChannelResp.value.data.supported_versions.reverse(),
          );
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching tier information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <GatewayTierContext.Provider
      value={{
        state: {supportedVersions, tiers},
        setState: (key, value?) =>
          setTierState({tiers, setTiers, networkId, key, value}),
      }}>
      {props.children}
    </GatewayTierContext.Provider>
  );
}

export default GatewayTierContext;
