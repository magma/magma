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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ActionQuery} from '../../components/ActionTable';
import type {EnqueueSnackbarOptions} from 'notistack';
import type {Metrics} from '../../components/context/SubscriberContext';
import type {SubscriberContextType} from '../../components/context/SubscriberContext';
import type {SubscriberRowType} from '../../views/subscriber/SubscriberOverview';
import type {
  mutable_subscriber,
  mutable_subscribers,
  network_id,
  subscriber,
  subscriber_state,
} from '../../../generated/MagmaAPIBindings';

import MagmaV1API from '../../../generated/WebClient';

import {
  DEFAULT_PAGE_SIZE,
  getLabelUnit,
} from '../../views/subscriber/SubscriberUtils';

type FetchProps = {
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
  networkId: string,
  id?: string,
  subscriberMap?: {[string]: subscriber},
  sessionState?: {[string]: subscriber_state},
  token?: string,
  pageSize?: number,
};

type InitSubscriberStateProps = {
  networkId: network_id,
  setSubscriberMap: ({[string]: subscriber}) => void,
  setSessionState: ({[string]: subscriber_state}) => void,
  setSubscriberMetrics: ({[string]: Metrics}) => void,
  setTotalCount: number => void,
  enqueueSnackbar?: (
    msg: string,
    cfg: EnqueueSnackbarOptions,
  ) => ?(string | number),
};

export async function FetchSubscribers(props: FetchProps) {
  const {networkId, enqueueSnackbar, id, token, pageSize} = props;
  if (id !== null && id !== undefined) {
    try {
      return await MagmaV1API.getLteByNetworkIdSubscribersBySubscriberId({
        networkId,
        subscriberId: id,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber information', {
        variant: 'error',
      });
    }
  } else {
    try {
      return await MagmaV1API.getLteByNetworkIdSubscribers({
        networkId,
        pageSize: pageSize ?? DEFAULT_PAGE_SIZE,
        pageToken: token ?? '',
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber information', {
        variant: 'error',
      });
    }
  }
}

export async function FetchSubscriberState(props: FetchProps) {
  const {networkId, enqueueSnackbar, id} = props;
  if (id !== null && id !== undefined) {
    try {
      return await MagmaV1API.getLteByNetworkIdSubscriberStateBySubscriberId({
        networkId,
        subscriberId: id,
      });
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  } else {
    try {
      return await MagmaV1API.getLteByNetworkIdSubscriberState({
        networkId,
      });
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
    setTotalCount,
    enqueueSnackbar,
  } = props;
  const subscriberResponse = await FetchSubscribers({
    networkId,
    enqueueSnackbar,
  });
  if (subscriberResponse) {
    setSubscriberMap(subscriberResponse.subscribers);
    setTotalCount(subscriberResponse.total_count);
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
  setSessionState: ({[string]: subscriber_state}) => void,
  key: string,
  value?: mutable_subscriber | mutable_subscribers,
  newState?: {[string]: subscriber},
  newSessionState?: {[string]: subscriber_state},
};

export async function setSubscriberState(props: SubscriberStateProps) {
  const {
    networkId,
    subscriberMap,
    setSubscriberMap,
    setSessionState,
    key,
    value,
    newState,
    newSessionState,
  } = props;
  if (newState) {
    setSubscriberMap(newState);
    return;
  }
  if (newSessionState) {
    // $FlowIgnore
    setSessionState(newSessionState.sessionState);
    return;
  }
  if (Array.isArray(value)) {
    await MagmaV1API.postLteByNetworkIdSubscribers({
      networkId,
      subscribers: value,
    });
    const newSubscriberMap = {};
    value.map(newSubscriber => {
      newSubscriberMap[newSubscriber.id] = newSubscriber;
    });
    setSubscriberMap({...subscriberMap, newSubscriberMap});
    return;
  }
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
        subscribers: [value],
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

export function getGatewaySubscriberMap(sessions: {
  [string]: subscriber_state,
}) {
  const gatewayMap = {};
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

export type SubscriberQueryType = {
  networkId: string,
  query: ActionQuery,
  maxPageRowCount: number,
  setMaxPageRowCount: number => void,
  pageSize: number,
  tokenList: Array<string>,
  setTokenList: (Array<string>) => void,
  ctx: SubscriberContextType,
  subscriberMetrics?: {[string]: Metrics},
  deleteTable: boolean,
};
/**
 * Used with material-table remote data feature to get paginated subscribers.
 * Returns a promise holding subscriber rows data, the current page and the subscribers total count.
 *
 * @param {string} networkId ID of the network.
 * @param {ActionQuery} query Subscriber query holding page number, page size, total count, order and filters.
 * @param {number} pageSize Size of subscriber page. (default is 10)
 * @param {Array<string>} tokenList List of page tokens used to get next/previous page.
 * @param {(Array<string>) => void} setTokenList Set token list.
 * @param {SubscriberContextType} ctx Subscriber context to set subscriber state.
 * @param {{[string]: Metrics}} subscriberMetrics Metrics used for subscriber Current Usage and Daily Average.
 * @param {boolean} deleteTable Add more fields to subscriber if set to true
 * @return Promise holding subscriber rows data, the current page and the totalCount.
 */
export async function handleSubscriberQuery(props: SubscriberQueryType) {
  const {
    networkId,
    query,
    pageSize,
    maxPageRowCount,
    setMaxPageRowCount,
    tokenList,
    setTokenList,
    ctx,
    subscriberMetrics,
    deleteTable,
  } = props;
  return new Promise(async (resolve, reject) => {
    try {
      // search subscriber by IMSI
      let subscriberSearch = {};
      const search = query.search;
      if (search.startsWith('IMSI') && search.length > 9) {
        const searchedSubscriber = await FetchSubscribers({
          networkId,
          id: search,
        });
        const metrics = subscriberMetrics?.[`${search}`];
        if (searchedSubscriber) {
          subscriberSearch = {
            name: searchedSubscriber.id,
            imsi: searchedSubscriber.id,
            service: searchedSubscriber.lte?.state || '',
            currentUsage: metrics?.currentUsage ?? '0',
            dailyAvg: metrics?.dailyAvg ?? '0',
            lastReportedTime:
              searchedSubscriber.monitoring?.icmp?.last_reported_time === 0
                ? new Date(
                    searchedSubscriber.monitoring?.icmp?.last_reported_time,
                  )
                : '-',
          };
        }
      }

      const page =
        maxPageRowCount < query.page * query.pageSize
          ? maxPageRowCount / query.pageSize
          : query.page;
      const subscriberResponse = await FetchSubscribers({
        networkId,
        token: tokenList[page] ?? tokenList[tokenList.length - 1],
        pageSize,
      });

      const newTokenList = tokenList;
      // add next page token in token list to get next subscriber page.
      let totalCount = 0;
      if (subscriberResponse) {
        if (!newTokenList.includes(subscriberResponse.next_page_token)) {
          newTokenList.push(subscriberResponse.next_page_token);
        }
        totalCount = subscriberResponse.total_count;
        setMaxPageRowCount(totalCount);
        setTokenList([...newTokenList]);
        // set subscriber state with current subscriber rows.
        if (!deleteTable) {
          ctx.setState?.('', undefined, subscriberResponse.subscribers);
        }
      }
      const tableData: Array<SubscriberRowType> = subscriberResponse
        ? Object.keys(subscriberResponse.subscribers).map((imsi: string) => {
            const subscriberInfo = subscriberResponse.subscribers[imsi] || {};
            const metrics = subscriberMetrics?.[`${imsi}`];
            // Additional fields displayed in subscriber delete dialog
            const deleteSubscriber = !deleteTable
              ? {}
              : {
                  authKey: subscriberInfo.lte.auth_key,
                  authOpc: subscriberInfo.lte.auth_opc,
                  dataPlan: subscriberInfo.lte.sub_profile,
                  apns: subscriberInfo.active_apns,
                  policies: subscriberInfo.active_policies,
                  state:
                    subscriberInfo.lte?.state === 'ACTIVE'
                      ? 'ACTIVE'
                      : 'INACTIVE',
                };
            const subscriber = {
              name: subscriberInfo.name ?? imsi,
              imsi: imsi,
              service: subscriberInfo.lte?.state || '',
              currentUsage: metrics?.currentUsage ?? '0',
              dailyAvg: metrics?.dailyAvg ?? '0',
              lastReportedTime:
                subscriberInfo.monitoring?.icmp?.last_reported_time === 0
                  ? new Date(
                      subscriberInfo.monitoring?.icmp?.last_reported_time,
                    )
                  : '-',
            };
            return {...subscriber, ...deleteSubscriber};
          })
        : [];
      resolve({
        data:
          search.startsWith('IMSI') && search.length > 9
            ? [subscriberSearch]
            : tableData,
        page: page,
        totalCount: totalCount,
      });
    } catch (e) {
      reject(e?.message ?? 'error retrieving subscribers');
    }
  });
}
