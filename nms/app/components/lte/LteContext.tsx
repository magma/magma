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
import CbsdContext from '../context/CbsdContext';
import InitSubscriberState, {
  FetchSubscriberState,
  getGatewaySubscriberMap,
  setSubscriberState,
} from '../../state/lte/SubscriberState';
import LoadingFiller from '../LoadingFiller';
import LteNetworkContext from '../context/LteNetworkContext';
import NetworkContext from '../context/NetworkContext';
import SubscriberContext from '../context/SubscriberContext';
import {ApnContextProvider} from '../context/ApnContext';
import {EnodebContextProvider} from '../context/EnodebContext';
import {GatewayContextProvider} from '../context/GatewayContext';
import {GatewayPoolsContextProvider} from '../context/GatewayPoolsContext';
import {GatewayTierContextProvider} from '../context/GatewayTierContext';
import {PolicyProvider} from '../context/PolicyContext';
import {TraceContextProvider} from '../context/TraceContext';
import {useCallback, useContext, useEffect, useMemo, useState} from 'react';
import type {
  FegLteNetwork,
  LteNetwork,
  MutableCbsd,
  MutableSubscriber,
  PaginatedCbsds,
  Subscriber,
  SubscriberState,
} from '../../../generated';

import * as cbsdState from '../../state/lte/CbsdState';
import MagmaAPI from '../../api/MagmaAPI';
import {
  FEG_LTE,
  LTE,
  NetworkId,
  SubscriberId,
} from '../../../shared/types/network';
import {UpdateNetworkState as UpdateFegLteNetworkState} from '../../state/feg_lte/NetworkState';
import {UpdateNetworkState as UpdateLteNetworkState} from '../../state/lte/NetworkState';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

type Props = {
  networkId: NetworkId;
  networkType: string;
  children: React.ReactNode;
};

export function CbsdContextProvider({networkId, children}: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();

  const [isLoading, setIsLoading] = useState(false);
  const [fetchResponse, setFetchResponse] = useState<PaginatedCbsds>({
    cbsds: [],
    total_count: 0,
  });
  const [paginationOptions, setPaginationOptions] = useState<{
    page: number;
    pageSize: number;
  }>({
    page: 0,
    pageSize: 10,
  });

  const refetch = useCallback(() => {
    return cbsdState.fetch({
      networkId,
      page: paginationOptions.page,
      pageSize: paginationOptions.pageSize,
      setIsLoading,
      setFetchResponse,
      enqueueSnackbar,
    });
  }, [
    networkId,
    paginationOptions.page,
    paginationOptions.pageSize,
    setIsLoading,
    setFetchResponse,
    enqueueSnackbar,
  ]);

  useEffect(() => {
    void refetch();
  }, [refetch, paginationOptions.page, paginationOptions.pageSize]);

  const state = useMemo(() => {
    return {
      isLoading,
      cbsds: fetchResponse.cbsds,
      totalCount: fetchResponse.total_count,
      page: paginationOptions.page,
      pageSize: paginationOptions.pageSize,
    };
  }, [
    isLoading,
    fetchResponse.cbsds,
    fetchResponse.total_count,
    paginationOptions.page,
    paginationOptions.pageSize,
  ]);

  return (
    <CbsdContext.Provider
      value={{
        state,
        setPaginationOptions,
        refetch,
        create: (newCbsd: MutableCbsd) => {
          return cbsdState
            .create({
              networkId,
              newCbsd,
            })
            .catch(e => {
              enqueueSnackbar?.('failed to create CBSD', {
                variant: 'error',
              });
              throw e as Error;
            })
            .then(() => {
              void refetch();
            });
        },
        update: (id: number, cbsd: MutableCbsd) => {
          return cbsdState
            .update({
              networkId,
              id,
              cbsd,
            })
            .catch(e => {
              enqueueSnackbar?.('failed to update CBSD', {
                variant: 'error',
              });
              throw e as Error;
            })
            .then(() => {
              void refetch();
            });
        },
        deregister: (id: number) => {
          return cbsdState
            .deregister({
              networkId,
              id,
            })
            .catch(() => {
              enqueueSnackbar?.('failed to deregister CBSD', {
                variant: 'error',
              });
            })
            .then(() => {
              void refetch();
            });
        },
        remove: (id: number) => {
          return cbsdState
            .remove({
              networkId,
              cbsdId: id,
            })
            .catch(() => {
              enqueueSnackbar?.('failed to remove CBSD', {
                variant: 'error',
              });
            })
            .then(() => {
              void refetch();
            });
        },
      }}>
      {children}
    </CbsdContext.Provider>
  );
}

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

export function LteNetworkContextProvider(props: Props) {
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
            return UpdateFegLteNetworkState({
              networkId,
              setLteNetwork,
              refreshState,
              ...props,
            });
          } else {
            return UpdateLteNetworkState({
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
