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
import MagmaAPI from '../api/MagmaAPI';
import {DEFAULT_PAGE_SIZE} from '../views/subscriber/SubscriberUtils';
import {EnqueueSnackbar} from '../hooks/useSnackbar';
import {NetworkId} from '../../shared/types/network';
import {SubscriberState} from '../../generated';
import {getServicedAccessNetworks} from '../components/FEGServicingAccessGatewayKPIs';
import type {ActionQuery} from '../components/ActionTable';
import type {
  Metrics,
  SubscriberContextType,
} from '../context/SubscriberContext';
import type {PaginatedSubscribers, Subscriber} from '../../generated';

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

type FetchParams = {
  enqueueSnackbar?: EnqueueSnackbar;
  networkId: string;
  id?: string;
  subscriberMap?: Record<string, Subscriber>;
  sessionState?: Record<string, SubscriberState>;
  token?: string;
  pageSize?: number;
};

/**
 * Props passed when fetching subscriber state.
 *
 * @param {NetworkId} networkId Id of the federation network.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error.
 */
type FetchFegSubscriberStateParams = {
  networkId: NetworkId;
  enqueueSnackbar?: EnqueueSnackbar;
};
type FegSubscriberState = Record<NetworkId, Record<string, SubscriberState>>;

export async function fetchSubscribers(params: FetchParams) {
  const {networkId, enqueueSnackbar, id, token, pageSize} = params;
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
      const searchedSubscriber = (await fetchSubscribers({
        networkId,
        id: search,
      })) as Subscriber;
      const metrics = subscriberMetrics?.[`${search}`];
      if (searchedSubscriber) {
        subscriberSearch = {
          name: searchedSubscriber.name || '',
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
    const subscriberResponse = (await fetchSubscribers({
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

export async function fetchSubscriberState(params: FetchParams) {
  const {networkId, enqueueSnackbar, id} = params;
  if (id !== null && id !== undefined) {
    try {
      const subscriber = await MagmaAPI.subscribers.lteNetworkIdSubscriberStateSubscriberIdGet(
        {
          networkId,
          subscriberId: id,
        },
      );
      if (subscriber) {
        return {[id]: subscriber.data};
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  } else {
    try {
      const subscribers = await MagmaAPI.subscribers.lteNetworkIdSubscriberStateGet(
        {
          networkId,
        },
      );
      if (subscribers) {
        return subscribers.data;
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching subscriber state', {
        variant: 'error',
      });
      return;
    }
  }
}

/**
 * Fetches and returns the subscriber session state of all the serviced
 * federated lte networks under by this federation network.
 *
 * @param {FetchFegSubscriberStateParams} props an object containing the network id and snackbar to display error.
 * @returns {{[string]:{[string]: subscriber_state}}} returns an object containing the serviced
 *   network ids mapped to the each of their subscriber state. It returns an empty object and
 *   displays any error encountered on the snackbar when it fails to fetch the session state.
 */
export async function fetchFegSubscriberState(
  props: FetchFegSubscriberStateParams,
): Promise<FegSubscriberState> {
  const {networkId, enqueueSnackbar} = props;
  const servicedAccessNetworks = await getServicedAccessNetworks(
    networkId,
    enqueueSnackbar,
  );
  const sessionState: FegSubscriberState = {};
  for (const servicedAccessNetwork of servicedAccessNetworks) {
    const servicedAccessNetworkId = servicedAccessNetwork.id;
    const state = await fetchSubscriberState({
      networkId: servicedAccessNetworkId,
      enqueueSnackbar,
    });
    // group session states under their network id
    sessionState[servicedAccessNetworkId] = state as Record<
      string,
      SubscriberState
    >;
  }
  return sessionState;
}
