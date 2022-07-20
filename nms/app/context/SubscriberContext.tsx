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

import LoadingFiller from '../components/LoadingFiller';
import MagmaAPI from '../api/MagmaAPI';
import React, {useEffect, useState} from 'react';
import {EnqueueSnackbar, useEnqueueSnackbar} from '../hooks/useSnackbar';
import {GatewayId, NetworkId, SubscriberId} from '../../shared/types/network';
import {
  MutableSubscriber,
  PaginatedSubscribers,
  Subscriber,
  SubscriberForbiddenNetworkTypesEnum,
  SubscriberState,
} from '../../generated';
import {fetchSubscriberState, fetchSubscribers} from '../util/SubscriberState';
import {getLabelUnit} from '../views/subscriber/SubscriberUtils';

export type Metrics = {
  currentUsage: string;
  dailyAvg: string;
};

/** SubscriberContextType
 * state: paginated subscribers
 * sessionState: paginated subscribers session state
 * metrics: subscriber metrics
 * gwSubscriberMap: gateway subscriber map
 * totalCount: total count of subscribers
 * setState: POST, PUT, DELETE subscriber
 */
export type SubscriberContextType = {
  state: Record<string, Subscriber>;
  sessionState: Record<string, SubscriberState>;
  forbiddenNetworkTypes: Record<
    string,
    Array<SubscriberForbiddenNetworkTypesEnum>
  >;
  metrics?: Record<string, Metrics>;
  gwSubscriberMap: Record<GatewayId, Array<SubscriberId>>;
  totalCount: number;
  setState?: (
    key: string,
    val?: MutableSubscriber | Array<MutableSubscriber>,
    newState?: Record<string, Subscriber>,
  ) => Promise<void>;
  refetchSessionState: (subscriberId?: SubscriberId) => void;
};

const SubscriberContext = React.createContext<SubscriberContextType>(
  {} as SubscriberContextType,
);

async function initSubscriberState(params: {
  networkId: NetworkId;
  setSubscriberMap: (subscriberMap: Record<string, Subscriber>) => void;
  setSessionState: (sessionState: Record<string, SubscriberState>) => void;
  setSubscriberMetrics: (subscriberMetrics: Record<string, Metrics>) => void;
  setTotalCount: (count: number) => void;
  enqueueSnackbar?: EnqueueSnackbar;
}) {
  const {
    networkId,
    setSubscriberMap,
    setSubscriberMetrics,
    setSessionState,
    setTotalCount,
    enqueueSnackbar,
  } = params;
  const subscriberResponse = (await fetchSubscribers({
    networkId,
    enqueueSnackbar,
  })) as PaginatedSubscribers;
  if (subscriberResponse) {
    setSubscriberMap(subscriberResponse.subscribers);
    setTotalCount(subscriberResponse.total_count);
  }

  const state = await fetchSubscriberState({
    networkId,
    enqueueSnackbar,
  });
  if (state) {
    setSessionState(state);
  }

  if (setSubscriberMetrics) {
    const subscriberMetrics = await fetchSubscriberMetrics({
      networkId,
      enqueueSnackbar,
    });
    if (subscriberMetrics) {
      setSubscriberMetrics(subscriberMetrics);
    }
  }
}

export async function fetchSubscriberMetrics(params: {
  enqueueSnackbar?: EnqueueSnackbar;
  networkId: string;
}) {
  const {networkId, enqueueSnackbar} = params;
  const queries = {
    dailyAvg: 'avg (avg_over_time(ue_reported_usage[24h])) by (IMSI)',
    currentUsage: 'sum (ue_reported_usage) by (IMSI)',
  };
  const subscriberMetrics: Record<string, Metrics> = {};

  const requests = (Object.keys(queries) as Array<keyof typeof queries>).map(
    async queryType => {
      try {
        const resp = (
          await MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet({
            networkId,
            query: queries[queryType],
          })
        ).data;

        resp?.data?.result?.filter(Boolean).forEach(item => {
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          const imsi = Object.values(item.metric)[0];
          if (typeof imsi === 'string') {
            const [value, unit] = getLabelUnit(parseFloat(item.value![1]));
            if (!(imsi in subscriberMetrics)) {
              subscriberMetrics[imsi] = {} as Metrics;
            }
            subscriberMetrics[imsi][queryType] = `${value}${unit}`;
          }
        });
      } catch (e) {
        enqueueSnackbar?.('failed fetching current usage information', {
          variant: 'error',
        });
      }
    },
  );
  await Promise.all(requests);
  return subscriberMetrics;
}

async function setSubscriberState(params: {
  networkId: NetworkId;
  subscriberMap: Record<string, Subscriber>;
  setSubscriberMap: (subscriberMap: Record<string, Subscriber>) => void;
  setSessionState: (sessionState: Record<string, SubscriberState>) => void;
  key: string;
  value?: MutableSubscriber | Array<MutableSubscriber>;
  newState?: Record<string, Subscriber>;
}) {
  const {
    networkId,
    subscriberMap,
    setSubscriberMap,
    key,
    value,
    newState,
  } = params;
  if (newState) {
    setSubscriberMap(newState);
    return;
  }
  if (Array.isArray(value)) {
    await MagmaAPI.subscribers.lteNetworkIdSubscribersPost({
      networkId,
      subscribers: value,
    });
    const newSubscriberMap: Record<string, Subscriber> = {};
    value.map(newSubscriber => {
      newSubscriberMap[newSubscriber.id] = newSubscriber as Subscriber;
    });
    // TODO[TS-migration] Should newSubscriberMap be spread here?
    // @ts-ignore
    setSubscriberMap({...subscriberMap, newSubscriberMap});
    return;
  }
  if (value != null) {
    if (key in subscriberMap) {
      await MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdPut({
        networkId,
        subscriber: value,
        subscriberId: key,
      });
    } else {
      await MagmaAPI.subscribers.lteNetworkIdSubscribersPost({
        networkId,
        subscribers: [value],
      });
    }
    const subscribers = (await fetchSubscribers({
      networkId: networkId,
      id: key,
    })) as Subscriber;
    if (subscribers) {
      setSubscriberMap({...subscriberMap, [key]: subscribers});
      return;
    }
  } else {
    await MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdDelete({
      networkId,
      subscriberId: key,
    });
    const newSubscriberMap = {...subscriberMap};
    delete newSubscriberMap[key];
    setSubscriberMap(newSubscriberMap);
  }
}

function getGatewaySubscriberMap(sessions: Record<string, SubscriberState>) {
  const gatewayMap: Record<string, Array<string>> = {};
  Object.keys(sessions).forEach(id => {
    const subscriber = sessions[id];
    const gwHardwareId = subscriber?.directory?.location_history?.[0];
    if (
      gwHardwareId !== null &&
      gwHardwareId !== undefined &&
      gwHardwareId !== ''
    ) {
      if (!(gwHardwareId in gatewayMap)) {
        gatewayMap[gwHardwareId] = [];
      }
      gatewayMap[gwHardwareId].push(id);
    }
  });
  return gatewayMap;
}

export function SubscriberContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
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
      await initSubscriberState({
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
          void fetchSubscriberState({networkId, id}).then(sessions => {
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

export default SubscriberContext;
