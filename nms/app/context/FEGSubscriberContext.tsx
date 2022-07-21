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
import {NetworkId, SubscriberId} from '../../shared/types/network';
import {SubscriberState} from '../../generated';
import {fetchFegSubscriberState} from '../util/SubscriberState';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

type FEGSubscriberContextType = {
  refetch: () => void;
  sessionState: Record<NetworkId, Record<SubscriberId, SubscriberState>>;
  setSessionState: (
    newSessionState: Record<NetworkId, Record<SubscriberId, SubscriberState>>,
  ) => void;
};

const FEGSubscriberContext = React.createContext<FEGSubscriberContextType>(
  {} as FEGSubscriberContextType,
);

/**
 * Fetches and saves the subscriber session states of networks
 * serviced by this federation network and whose subscriber
 * information is not managed by the HSS.
 *
 * @param {network_id} networkId Id of the network
 */
export function FEGSubscriberContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
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
      const sessionState = await fetchFegSubscriberState({
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
        refetch: () => {
          void fetchFegSubscriberState({networkId}).then(fegSubscriberState => {
            setSessionState(fegSubscriberState);
            if (fegSubscriberState) {
            }
          });
        },
        sessionState: sessionState,
        setSessionState: newSessionState => {
          return setSessionState(newSessionState);
        },
      }}>
      {props.children}
    </FEGSubscriberContext.Provider>
  );
}

export default FEGSubscriberContext;
