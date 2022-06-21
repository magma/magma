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
 */
import type {ActionQuery} from '../../components/ActionTable';
import type {Metrics} from '../../components/context/SubscriberContext';
import type {
  MutableSubscriber,
  PaginatedSubscribers,
  Subscriber,
  SubscriberState,
} from '../../../generated-ts';
import type {OptionsObject} from 'notistack';
import type {SubscriberContextType} from '../../components/context/SubscriberContext';

import MagmaAPI from '../../../api/MagmaAPI';

import {
  DEFAULT_PAGE_SIZE,
  getLabelUnit,
} from '../../views/subscriber/SubscriberUtils';
import {NetworkId} from '../../../shared/types/network';

export type SubscriberRowType = {
  name: string;
  imsi: string;
  activeApns?: string;
  ipAddresses?: string;
  activeSessions?: number;
  service: string;
  currentUsage: string;
  dailyAvg: string;
  lastReportedTime: Date | string;
};

type FetchProps = {
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
  networkId: string;
  id?: string;
  subscriberMap?: Record<string, Subscriber>;
  sessionState?: Record<string, SubscriberState>;
  token?: string;
  pageSize?: number;
};

type InitSubscriberStateProps = {
  networkId: NetworkId;
  setSubscriberMap: (subscriberMap: Record<string, Subscriber>) => void;
  setSessionState: (sessionState: Record<string, SubscriberState>) => void;
  setSubscriberMetrics: (subscriberMetrics: Record<string, Metrics>) => void;
  setTotalCount: (count: number) => void;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

export async function FetchSubscribers(props: FetchProps) {
  const {networkId, enqueueSnackbar, id, token, pageSize} = props;
  if (id !== null && id !== undefined) {
    try {
      return (
        await MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdGet({
          networkId,
          subscriberId: id,
        })
      ).data;
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber information', {
        variant: 'error',
      });
    }
  } else {
    try {
      return (
        await MagmaAPI.subscribers.lteNetworkIdSubscribersGet({
          networkId,
          pageSize: pageSize ?? DEFAULT_PAGE_SIZE,
          pageToken: token ?? '',
        })
      ).data;
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
      return (
        await MagmaAPI.subscribers.lteNetworkIdSubscriberStateSubscriberIdGet({
          networkId,
          subscriberId: id,
        })
      ).data;
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  } else {
    try {
      return (
        await MagmaAPI.subscribers.lteNetworkIdSubscriberStateGet({
          networkId,
        })
      ).data;
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
  const subscriberResponse = (await FetchSubscribers({
    networkId,
    enqueueSnackbar,
  })) as PaginatedSubscribers;
  if (subscriberResponse) {
    setSubscriberMap(subscriberResponse.subscribers);
    setTotalCount(subscriberResponse.total_count);
  }

  const state = (await FetchSubscriberState({
    networkId,
    enqueueSnackbar,
  })) as Record<string, SubscriberState>;
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
  networkId: NetworkId;
  subscriberMap: Record<string, Subscriber>;
  setSubscriberMap: (subscriberMap: Record<string, Subscriber>) => void;
  setSessionState: (sessionState: Record<string, SubscriberState>) => void;
  key: string;
  value?: MutableSubscriber | Array<MutableSubscriber>;
  newState?: Record<string, Subscriber>;
  newSessionState?: Record<string, SubscriberState>;
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
    // TODO[TS-migration] is the type of newSessionState broken?
    setSessionState(
      ((newSessionState as unknown) as {
        sessionState: Record<string, SubscriberState>;
      }).sessionState,
    );
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
      const subscribers = (await FetchSubscribers({
        networkId: networkId,
        id: key,
      })) as Subscriber;
      if (subscribers) {
        setSubscriberMap({...subscriberMap, [key]: subscribers});
        return;
      }
    } else {
      await MagmaAPI.subscribers.lteNetworkIdSubscribersPost({
        networkId,
        subscribers: [value],
      });
    }
    setSubscriberMap({...subscriberMap, [key]: value as Subscriber});
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

export function getGatewaySubscriberMap(
  sessions: Record<string, SubscriberState>,
) {
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

export type SubscriberQueryType = {
  networkId: string;
  query: ActionQuery;
  maxPageRowCount: number;
  setMaxPageRowCount: (rowCount: number) => void;
  pageSize: number;
  tokenList: Array<string>;
  setTokenList: (tokens: Array<string>) => void;
  ctx: SubscriberContextType;
  subscriberMetrics?: Record<string, Metrics>;
  deleteTable: boolean;
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
export async function handleSubscriberQuery(
  props: SubscriberQueryType,
): Promise<{
  data: Array<SubscriberRowType>;
  page: number;
  totalCount: number;
}> {
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
  try {
    // search subscriber by IMSI
    let subscriberSearch = {} as SubscriberRowType;
    const search = query.search;
    if (search.startsWith('IMSI') && search.length > 9) {
      const searchedSubscriber = (await FetchSubscribers({
        networkId,
        id: search,
      })) as Subscriber;
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
    const subscriberResponse = (await FetchSubscribers({
      networkId,
      token: tokenList[page] ?? tokenList[tokenList.length - 1],
      pageSize,
    })) as PaginatedSubscribers;

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
        await ctx.setState?.('', undefined, subscriberResponse.subscribers);
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
                ? new Date(subscriberInfo.monitoring?.icmp?.last_reported_time)
                : '-',
          };
          return {...subscriber, ...deleteSubscriber};
        })
      : [];
    return {
      data:
        search.startsWith('IMSI') && search.length > 9
          ? [subscriberSearch]
          : tableData,
      page: page,
      totalCount: totalCount,
    };
  } catch (e) {
    if (e instanceof Error) {
      throw e;
    }
    throw new Error('error retrieving subscribers');
  }
}
