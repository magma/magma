/*
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
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

import {getLabelUnit} from '../../views/subscriber/SubscriberUtils';
import type {Metrics} from '../../components/context/SubscriberContext';
import type {
  mutable_subscriber,
  network_id,
  subscriber,
  subscriber_state,
} from '@fbcnms/magma-api';

type FetchProps = {
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
  networkId: string,
  id?: string,
  subscriberMap?: {[string]: subscriber},
  sessionState?: {[string]: subscriber_state},
};
type InitSubscriberStateProps = {
  networkId: network_id,
  setSubscriberMap: ({[string]: subscriber}) => void,
  setSessionState: ({[string]: subscriber_state}) => void,
  setSubscriberMetrics?: ({[string]: Metrics}) => void,
  enqueueSnackbar?: (msg: string, cfg: {}) => ?(string | number),
};
export async function FetchSubscribers(props: FetchProps) {
  const {networkId, enqueueSnackbar, id} = props;
  if (id !== null && id !== undefined) {
    try {
      const subscriber = await MagmaV1API.getLteByNetworkIdSubscribersBySubscriberId(
        {
          networkId,
          subscriberId: id,
        },
      );
      if (subscriber) {
        return subscriber;
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber information', {
        variant: 'error',
      });
      return;
    }
  } else {
    try {
      const subscribers = await MagmaV1API.getLteByNetworkIdSubscribers({
        networkId,
      });
      if (subscribers) {
        return subscribers;
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber information', {
        variant: 'error',
      });
      return;
    }
  }
}

export async function FetchSubscriberState(props: FetchProps) {
  const {networkId, enqueueSnackbar, id} = props;
  if (id !== null && id !== undefined) {
    try {
      const session = await MagmaV1API.getLteByNetworkIdSubscriberStateBySubscriberId(
        {
          networkId,
          subscriberId: id,
        },
      );
      if (session) {
        return session;
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  } else {
    try {
      const sessions = await MagmaV1API.getLteByNetworkIdSubscriberState({
        networkId,
      });
      if (sessions) {
        return sessions;
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  }
}

export async function fetchSubscriberMetrics(props: FetchProps) {
  const {networkId, enqueueSnackbar} = props;
  const subscriberMetrics = {};
  const queries = {
    dailyAvg: 'avg (avg_over_time(ue_reported_usage[24h])) by (IMSI)',
    currentUsage: 'sum (ue_reported_usage) by (IMSI)',
  };

  const requests = Object.keys(queries).map(async (queryType: string) => {
    try {
      const resp = await MagmaV1API.getNetworksByNetworkIdPrometheusQuery({
        networkId,
        query: queries[queryType],
      });

      resp?.data?.result?.filter(Boolean).forEach(item => {
        const imsi = Object.values(item?.metric)?.[0];
        if (typeof imsi === 'string') {
          const [value, unit] = getLabelUnit(parseFloat(item?.value?.[1]));
          if (!(imsi in subscriberMetrics)) {
            subscriberMetrics[imsi] = {};
          }
          subscriberMetrics[imsi][queryType] = `${value}${unit}`;
        }
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching current usage information', {
        variant: 'error',
      });
    }
  });
  await Promise.all(requests);
  return subscriberMetrics;
}
export default async function InitSubscriberState(
  props: InitSubscriberStateProps,
) {
  const {
    networkId,
    setSubscriberMap,
    setSubscriberMetrics,
    setSessionState,
    enqueueSnackbar,
  } = props;
  const subscribers = await FetchSubscribers({networkId, enqueueSnackbar});
  if (subscribers) {
    setSubscriberMap(subscribers);
  }

  const state = await FetchSubscriberState({networkId, enqueueSnackbar});
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

type SubscriberStateProps = {
  networkId: network_id,
  subscriberMap: {[string]: subscriber},
  setSubscriberMap: ({[string]: subscriber}) => void,
  key: string,
  value?: mutable_subscriber,
};

export async function setSubscriberState(props: SubscriberStateProps) {
  const {networkId, subscriberMap, setSubscriberMap, key, value} = props;
  if (value != null) {
    if (key in subscriberMap) {
      await MagmaV1API.putLteByNetworkIdSubscribersBySubscriberId({
        networkId,
        subscriber: value,
        subscriberId: key,
      });
      const subscribers = await FetchSubscribers({
        networkId: networkId,
        id: key,
      });
      if (subscribers) {
        setSubscriberMap({...subscriberMap, [key]: subscribers});
        return;
      }
    } else {
      await MagmaV1API.postLteByNetworkIdSubscribers({
        networkId,
        subscriber: value,
      });
    }
    setSubscriberMap({...subscriberMap, [key]: value});
  } else {
    await MagmaV1API.deleteLteByNetworkIdSubscribersBySubscriberId({
      networkId,
      subscriberId: key,
    });
    const newSubscriberMap = {...subscriberMap};
    delete newSubscriberMap[key];
    setSubscriberMap(newSubscriberMap);
  }
}

export function getSubscriberGatewayMap(subscribers: {[string]: subscriber}) {
  const gatewayMap = {};
  Object.keys(subscribers).forEach(id => {
    const subscriber = subscribers[id];
    const gwHardwareId = subscriber?.state?.directory?.location_history?.[0];
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
