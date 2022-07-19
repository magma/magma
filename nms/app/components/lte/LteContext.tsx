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
import InitSubscriberState, {
  FetchSubscriberState,
  getGatewaySubscriberMap,
  setSubscriberState,
} from '../../state/lte/SubscriberState';
import LoadingFiller from '../LoadingFiller';
import SubscriberContext from '../context/SubscriberContext';
import {ApnContextProvider} from '../context/ApnContext';
import {CbsdContextProvider} from '../context/CbsdContext';
import {EnodebContextProvider} from '../context/EnodebContext';
import {
  FEG_LTE,
  LTE,
  NetworkId,
  SubscriberId,
} from '../../../shared/types/network';
import {GatewayContextProvider} from '../context/GatewayContext';
import {GatewayPoolsContextProvider} from '../context/GatewayPoolsContext';
import {GatewayTierContextProvider} from '../context/GatewayTierContext';
import {LteNetworkContextProvider} from '../context/LteNetworkContext';
import {PolicyProvider} from '../context/PolicyContext';
import {TraceContextProvider} from '../context/TraceContext';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import type {
  MutableSubscriber,
  Subscriber,
  SubscriberState,
} from '../../../generated';

type Props = {
  networkId: NetworkId;
  networkType: string;
  children: React.ReactNode;
};

export function SubscriberContextProvider(props: Props) {
  const {networkId} = props;
  const [subscriberMap, setSubscriberMap] = useState<
    Record<string, Subscriber>
  >({});
  const [sessionState, setSessionState] = useState<
    Record<string, SubscriberState>
  >({});
  const [subscriberMetrics, setSubscriberMetrics] = useState({});
  const [isLoading, setIsLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);
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
        setSessionState,
        setTotalCount,
        enqueueSnackbar,
      });
      setIsLoading(false);
    };
    void fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <SubscriberContext.Provider
      value={{
        forbiddenNetworkTypes: {},
        state: subscriberMap,
        metrics: subscriberMetrics,
        sessionState: sessionState,
        totalCount: totalCount,
        gwSubscriberMap: getGatewaySubscriberMap(sessionState),
        setState: (
          key: SubscriberId,
          value?: MutableSubscriber | Array<MutableSubscriber>,
          newState?,
        ) =>
          setSubscriberState({
            networkId,
            subscriberMap,
            setSubscriberMap,
            setSessionState,
            key,
            value,
            newState,
          }),
        refetchSessionState: (id?: SubscriberId) => {
          void FetchSubscriberState({networkId, id}).then(sessions => {
            if (sessions) {
              setSessionState(currentSessionState =>
                id
                  ? {
                      ...currentSessionState,
                      ...sessions,
                    }
                  : sessions,
              );
            }
          });
        },
      }}>
      {props.children}
    </SubscriberContext.Provider>
  );
}

export function LteContextProvider(props: Props) {
  const {networkId, networkType} = props;
  const lteNetwork = networkType === LTE || networkType === FEG_LTE;
  if (!lteNetwork) {
    return <>{props.children}</>;
  }

  return (
    <LteNetworkContextProvider {...{networkId, networkType}}>
      <PolicyProvider networkId={networkId}>
        <ApnContextProvider networkId={networkId}>
          <SubscriberContextProvider {...{networkId, networkType}}>
            <GatewayTierContextProvider {...{networkId, networkType}}>
              <EnodebContextProvider networkId={networkId}>
                <GatewayContextProvider networkId={networkId}>
                  <GatewayPoolsContextProvider networkId={networkId}>
                    <TraceContextProvider networkId={networkId}>
                      <CbsdContextProvider {...{networkId, networkType}}>
                        {props.children}
                      </CbsdContextProvider>
                    </TraceContextProvider>
                  </GatewayPoolsContextProvider>
                </GatewayContextProvider>
              </EnodebContextProvider>
            </GatewayTierContextProvider>
          </SubscriberContextProvider>
        </ApnContextProvider>
      </PolicyProvider>
    </LteNetworkContextProvider>
  );
}
