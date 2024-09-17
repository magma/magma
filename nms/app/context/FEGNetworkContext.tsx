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
import useMagmaAPI from '../api/useMagmaAPI';
import {FegNetwork, NetworkSubscriberConfig} from '../../generated';
import {NetworkId} from '../../shared/types/network';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

export type FEGNetworkContextType = {
  state: Partial<FegNetwork>;
  updateNetworks: (props: Partial<UpdateNetworkParams>) => Promise<void>;
};

const FEGNetworkContext = React.createContext<FEGNetworkContextType>(
  {} as FEGNetworkContextType,
);

export type UpdateNetworkParams = {
  networkId: NetworkId;
  fegNetwork?: FegNetwork;
  subscriberConfig?: NetworkSubscriberConfig;
  setFegNetwork: (fn: FegNetwork) => void;
  refreshState: boolean;
};

export async function updateFegNetworkState(params: UpdateNetworkParams) {
  const {networkId, setFegNetwork} = params;
  const requests = [];
  if (params.fegNetwork) {
    requests.push(
      await MagmaAPI.federationNetworks.fegNetworkIdPut({
        networkId: networkId,
        fegNetwork: {
          ...params.fegNetwork,
        },
      }),
    );
  }
  if (params.subscriberConfig) {
    requests.push(
      await MagmaAPI.federationNetworks.fegNetworkIdSubscriberConfigPut({
        networkId: params.networkId,
        record: params.subscriberConfig,
      }),
    );
  }
  await Promise.all(requests);
  if (params.refreshState) {
    setFegNetwork(
      (await MagmaAPI.federationNetworks.fegNetworkIdGet({networkId})).data,
    );
  }
}

/**
 * Fetches and returns information about the federation network inside
 * a context provider.
 * @param {object} props: contains the network id and its type
 */
export function FEGNetworkContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
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
          return updateFegNetworkState({
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

export default FEGNetworkContext;
