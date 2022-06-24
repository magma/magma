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
import ApnContext from '../context/ApnContext';
import CbsdContext from '../context/CbsdContext';
import EnodebContext from '../context/EnodebContext';
import GatewayContext from '../context/GatewayContext';
import GatewayPoolsContext from '../context/GatewayPoolsContext';
import GatewayTierContext from '../context/GatewayTierContext';
import InitSubscriberState from '../../state/lte/SubscriberState';
import LoadingFiller from '../LoadingFiller';
import LteNetworkContext from '../context/LteNetworkContext';
import NetworkContext from '../context/NetworkContext';
import PolicyContext from '../context/PolicyContext';
import SubscriberContext from '../context/SubscriberContext';
import TraceContext from '../context/TraceContext';
import type {
  Apn,
  BaseNameRecord,
  CallTrace,
  CallTraceConfig,
  FegLteNetwork,
  FegNetwork,
  LteGateway,
  LteNetwork,
  MutableCallTrace,
  MutableCbsd,
  MutableSubscriber,
  NetworkRanConfigs,
  PaginatedCbsds,
  PolicyQosProfile,
  PolicyRule,
  RatingGroup,
  Tier,
} from '../../../generated-ts';
import type {EnodebInfo} from './EnodebUtils';
import type {EnodebState} from '../context/EnodebContext';
import type {gatewayPoolsStateType} from '../context/GatewayPoolsContext';

import * as cbsdState from '../../state/lte/CbsdState';
import MagmaAPI from '../../../api/MagmaAPI';
import {
  FEG_LTE,
  LTE,
  NetworkId,
  SubscriberId,
} from '../../../shared/types/network';
import {
  InitEnodeState,
  InitGatewayPoolState,
  InitTierState,
  SetEnodebState,
  SetGatewayPoolsState,
  SetGatewayState,
  SetTierState,
  UpdateGateway,
  UpdateGatewayPoolRecords,
  UpdateGatewayProps,
} from '../../state/lte/EquipmentState';
import {InitTraceState, SetCallTraceState} from '../../state/TraceState';
import {SetApnState} from '../../state/lte/ApnState';
import {
  SetBaseNameState,
  SetPolicyState,
  SetQosProfileState,
  SetRatingGroupState,
} from '../../state/PolicyState';
import {UpdateNetworkState as UpdateFegLteNetworkState} from '../../state/feg_lte/NetworkState';
import {UpdateNetworkState as UpdateFegNetworkState} from '../../state/feg/NetworkState';
import {UpdateNetworkState as UpdateLteNetworkState} from '../../state/lte/NetworkState';
import {
  getGatewaySubscriberMap,
  setSubscriberState,
} from '../../state/lte/SubscriberState';
import {useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

type Props = {
  networkId: NetworkId;
  networkType: string;
  children: React.ReactNode;
};

export function GatewayContextProvider(props: Props) {
  const {networkId} = props;
  const [lteGateways, setLteGateways] = useState<Record<string, LteGateway>>(
    {},
  );
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        const lteGateways = (
          await MagmaAPI.lteGateways.lteNetworkIdGatewaysGet({
            networkId,
          })
        ).data;
        setLteGateways(lteGateways);
      } catch (e) {
        enqueueSnackbar?.('failed fetching gateway information', {
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
    <GatewayContext.Provider
      value={{
        state: lteGateways,
        setState: (key, value?, newState?) => {
          return SetGatewayState({
            lteGateways,
            setLteGateways,
            networkId,
            key,
            value,
            newState,
          });
        },
        updateGateway: props =>
          UpdateGateway({
            networkId,
            setLteGateways,
            ...props,
          } as UpdateGatewayProps),
      }}>
      {props.children}
    </GatewayContext.Provider>
  );
}

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
    refetch();
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
              throw e;
            })
            .then(() => {
              refetch();
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
              throw e;
            })
            .then(() => {
              refetch();
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
              refetch();
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
              refetch();
            });
        },
      }}>
      {children}
    </CbsdContext.Provider>
  );
}

export function EnodebContextProvider(props: Props) {
  const {networkId} = props;
  const [enbInfo, setEnbInfo] = useState<Record<string, EnodebInfo>>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<NetworkRanConfigs>(
    {} as NetworkRanConfigs,
  );
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        const [lteRanConfigsResp] = await Promise.allSettled([
          MagmaAPI.lteNetworks.lteNetworkIdCellularRanGet({networkId}),
          InitEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
        ]);
        if (lteRanConfigsResp.status === 'fulfilled') {
          setLteRanConfigs(lteRanConfigsResp.value.data);
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching enodeb information', {
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
    <EnodebContext.Provider
      value={{
        state: {enbInfo},
        lteRanConfigs: lteRanConfigs,
        setState: (key: string, value?, newState?: EnodebState) => {
          return SetEnodebState({
            enbInfo,
            setEnbInfo,
            networkId,
            key,
            value,
            newState,
          });
        },
      }}>
      {props.children}
    </EnodebContext.Provider>
  );
}

export function TraceContextProvider(props: Props) {
  const {networkId} = props;
  const [traceMap, setTraceMap] = useState<Record<string, CallTrace>>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchLteState = async () => {
      if (networkId == null) {
        return;
      }
      await InitTraceState({
        networkId,
        setTraceMap,
        enqueueSnackbar,
      }),
        setIsLoading(false);
    };
    void fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <TraceContext.Provider
      value={{
        state: traceMap,
        setState: (key: string, value?: MutableCallTrace | CallTraceConfig) =>
          SetCallTraceState({
            networkId,
            callTraces: traceMap,
            setCallTraces: setTraceMap,
            key,
            value,
          }),
      }}>
      {props.children}
    </TraceContext.Provider>
  );
}

export function SubscriberContextProvider(props: Props) {
  const {networkId} = props;
  const [subscriberMap, setSubscriberMap] = useState({});
  const [sessionState, setSessionState] = useState({});
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
      }),
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
          newSessionState?,
        ) =>
          setSubscriberState({
            networkId,
            subscriberMap,
            setSubscriberMap,
            setSessionState,
            key,
            value,
            newState,
            newSessionState,
          }),
      }}>
      {props.children}
    </SubscriberContext.Provider>
  );
}

export function GatewayPoolsContextProvider(props: Props) {
  const {networkId} = props;
  const [isLoading, setIsLoading] = useState(true);
  const [gatewayPools, setGatewayPools] = useState<
    Record<string, gatewayPoolsStateType>
  >({});
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        await InitGatewayPoolState({
          enqueueSnackbar,
          networkId,
          setGatewayPools,
        });
      } catch (e) {
        enqueueSnackbar?.('failed fetching gateway pool information', {
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
    <GatewayPoolsContext.Provider
      value={{
        state: gatewayPools,
        setState: (key, value?) =>
          SetGatewayPoolsState({
            gatewayPools,
            setGatewayPools,
            networkId,
            key,
            value,
          }),
        updateGatewayPoolRecords: (key, value?, resources?) =>
          UpdateGatewayPoolRecords({
            gatewayPools,
            setGatewayPools,
            networkId,
            key,
            value,
            resources,
          }),
      }}>
      {props.children}
    </GatewayPoolsContext.Provider>
  );
}

export function GatewayTierContextProvider(props: Props) {
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
          InitTierState({networkId, setTiers, enqueueSnackbar}),
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
          SetTierState({tiers, setTiers, networkId, key, value}),
      }}>
      {props.children}
    </GatewayTierContext.Provider>
  );
}

export function PolicyProvider(props: Props) {
  const {networkId} = props;
  const networkCtx = useContext(NetworkContext);
  const lteNetworkCtx = useContext(LteNetworkContext);
  const [policies, setPolicies] = useState<Record<string, PolicyRule>>({});
  const [baseNames, setBaseNames] = useState<Record<string, BaseNameRecord>>(
    {},
  );
  const [qosProfiles, setQosProfiles] = useState<
    Record<string, PolicyQosProfile>
  >({});
  const [ratingGroups, setRatingGroups] = useState<Record<string, RatingGroup>>(
    {},
  );
  const [fegNetwork, setFegNetwork] = useState<FegNetwork>({} as FegNetwork);
  const [fegPolicies, setFegPolicies] = useState<Record<string, PolicyRule>>(
    {},
  );
  const [isLoading, setIsLoading] = useState(true);
  const networkType = networkCtx.networkType;
  const enqueueSnackbar = useEnqueueSnackbar();
  let fegNetworkId: string | undefined;
  if (networkType === FEG_LTE) {
    fegNetworkId = lteNetworkCtx.state?.federation?.feg_network_id;
  }

  useEffect(() => {
    const fetchState = async () => {
      try {
        setPolicies(
          (
            await MagmaAPI.policies.networksNetworkIdPoliciesRulesviewfullGet({
              networkId,
            })
          ).data,
        );

        // Base Names
        const baseNameIDs: Array<string> = (
          await MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesGet({
            networkId,
          })
        ).data;
        const baseNameRecords = await Promise.all(
          baseNameIDs.map(baseNameID =>
            MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesBaseNameGet({
              networkId,
              baseName: baseNameID,
            }),
          ),
        );
        const newBaseNames: Record<string, BaseNameRecord> = {};
        baseNameRecords.map(({data: record}) => {
          newBaseNames[record.name] = record;
        });
        setBaseNames(newBaseNames);

        setRatingGroups(
          // TODO[TS-migration] What is the actual type here?
          ((
            await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsGet({
              networkId,
            })
          ).data as unknown) as Record<string, RatingGroup>,
        );
        setQosProfiles(
          (
            await MagmaAPI.policies.lteNetworkIdPolicyQosProfilesGet({
              networkId,
            })
          ).data,
        );
        if (fegNetworkId) {
          setFegNetwork(
            (
              await MagmaAPI.federationNetworks.fegNetworkIdGet({
                networkId: fegNetworkId,
              })
            ).data,
          );
          setFegPolicies(
            (
              await MagmaAPI.policies.networksNetworkIdPoliciesRulesviewfullGet(
                {
                  networkId: fegNetworkId,
                },
              )
            ).data,
          );
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching policy information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    void fetchState();
  }, [networkId, fegNetworkId, networkType, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }
  return (
    <PolicyContext.Provider
      value={{
        state: policies,
        ratingGroups: ratingGroups,
        setRatingGroups: async (key, value) => {
          await SetRatingGroupState({
            networkId,
            ratingGroups,
            setRatingGroups,
            key,
            value,
          });
        },

        baseNames: baseNames,
        setBaseNames: async (key, value) => {
          await SetBaseNameState({
            networkId,
            baseNames,
            setBaseNames,
            key,
            value,
          });
        },

        qosProfiles: qosProfiles,
        setQosProfiles: async (key, value) => {
          await SetQosProfileState({
            networkId,
            qosProfiles,
            setQosProfiles,
            key,
            value,
          });
        },
        setState: async (key, value?, isNetworkWide?) => {
          if (networkType === FEG_LTE) {
            const fegNetworkID =
              lteNetworkCtx.state?.federation?.feg_network_id;
            await SetPolicyState({
              policies,
              setPolicies,
              networkId,
              key,
              value,
            });

            // duplicate the policy on feg_network as well
            if (fegNetworkID) {
              await SetPolicyState({
                policies: fegPolicies,
                setPolicies: setFegPolicies,
                networkId: fegNetworkID,
                key,
                value,
              });
            }
          } else {
            await SetPolicyState({
              policies,
              setPolicies,
              networkId,
              key,
              value,
            });
          }
          if (isNetworkWide === true) {
            // we only support isNetworkWide rules now(and not basenames)
            let ruleNames = [];
            let fegRuleNames = [];

            if (value != null) {
              ruleNames =
                lteNetworkCtx.state?.subscriber_config
                  ?.network_wide_rule_names ?? [];
              fegRuleNames =
                fegNetwork.subscriber_config?.network_wide_rule_names ?? [];

              // update subscriber config if necessary
              if (!ruleNames.includes(key)) {
                ruleNames.push(key);
                void lteNetworkCtx.updateNetworks({
                  networkId,
                  subscriberConfig: {
                    network_wide_base_names:
                      lteNetworkCtx.state?.subscriber_config
                        ?.network_wide_base_names,
                    network_wide_rule_names: ruleNames,
                  },
                });
              }

              if (!fegRuleNames.includes(key)) {
                fegRuleNames.push(key);
                if (networkType === FEG_LTE && fegNetwork) {
                  void UpdateFegNetworkState({
                    networkId: fegNetwork.id,
                    subscriberConfig: {
                      network_wide_base_names:
                        fegNetwork.subscriber_config?.network_wide_base_names,
                      network_wide_rule_names: fegRuleNames,
                    },
                    setFegNetwork,
                    refreshState: true,
                  });
                }
              }
            }
          } else {
            // delete network wide rules for the key if present
            let ruleNames = [];
            let fegRuleNames = [];
            const oldRuleNames =
              lteNetworkCtx.state?.subscriber_config?.network_wide_rule_names ??
              [];
            const oldFegRuleNames =
              fegNetwork.subscriber_config?.network_wide_rule_names ?? [];

            if (oldRuleNames.includes(key)) {
              ruleNames = oldRuleNames.filter(function (ruleId) {
                return ruleId !== key;
              });
              void lteNetworkCtx.updateNetworks({
                networkId,
                subscriberConfig: {
                  network_wide_base_names:
                    lteNetworkCtx.state?.subscriber_config
                      ?.network_wide_base_names,
                  network_wide_rule_names: ruleNames,
                },
              });
            }

            // if we have old feg rul
            if (oldFegRuleNames.includes(key)) {
              fegRuleNames = oldFegRuleNames.filter(function (ruleId) {
                return ruleId !== key;
              });

              if (networkType === FEG_LTE && fegNetwork) {
                void UpdateFegNetworkState({
                  networkId: fegNetwork.id,
                  subscriberConfig: {
                    network_wide_base_names:
                      fegNetwork.subscriber_config?.network_wide_base_names,
                    network_wide_rule_names: fegRuleNames,
                  },
                  setFegNetwork,
                  refreshState: true,
                });
              }
            }
          }
        },
      }}>
      {props.children}
    </PolicyContext.Provider>
  );
}

export function ApnProvider(props: Props) {
  const {networkId} = props;
  const [apns, setApns] = useState<Record<string, Apn>>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        setApns(
          (
            await MagmaAPI.apns.lteNetworkIdApnsGet({
              networkId,
            })
          ).data,
        );
      } catch (e) {
        enqueueSnackbar?.('failed fetching APN information', {
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
    <ApnContext.Provider
      value={{
        state: apns,
        setState: (key, value?) => {
          return SetApnState({
            apns,
            setApns,
            networkId,
            key,
            value,
          });
        },
      }}>
      {props.children}
    </ApnContext.Provider>
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
    return <>props.children</>;
  }

  return (
    <LteNetworkContextProvider {...{networkId, networkType}}>
      <PolicyProvider {...{networkId, networkType}}>
        <ApnProvider {...{networkId, networkType}}>
          <SubscriberContextProvider {...{networkId, networkType}}>
            <GatewayTierContextProvider {...{networkId, networkType}}>
              <EnodebContextProvider {...{networkId, networkType}}>
                <GatewayContextProvider {...{networkId, networkType}}>
                  <GatewayPoolsContextProvider {...{networkId, networkType}}>
                    <TraceContextProvider {...{networkId, networkType}}>
                      <CbsdContextProvider {...{networkId, networkType}}>
                        {props.children}
                      </CbsdContextProvider>
                    </TraceContextProvider>
                  </GatewayPoolsContextProvider>
                </GatewayContextProvider>
              </EnodebContextProvider>
            </GatewayTierContextProvider>
          </SubscriberContextProvider>
        </ApnProvider>
      </PolicyProvider>
    </LteNetworkContextProvider>
  );
}
