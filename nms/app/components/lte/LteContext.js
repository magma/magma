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
 *
 * @flow strict-local
 * @format
 */
import * as React from 'react';
// $FlowFixMe migrated to typescript
import ApnContext from '../context/ApnContext';
import CbsdContext from '../context/CbsdContext';
// $FlowFixMe migrated to typescript
import EnodebContext from '../context/EnodebContext';
import GatewayContext from '../context/GatewayContext';
import GatewayPoolsContext from '../context/GatewayPoolsContext';
import GatewayTierContext from '../context/GatewayTierContext';
import InitSubscriberState from '../../state/lte/SubscriberState';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import LteNetworkContext from '../context/LteNetworkContext';
import MagmaV1API from '../../../generated/WebClient';
// $FlowFixMe migrated to typescript
import NetworkContext from '../context/NetworkContext';
import PolicyContext from '../context/PolicyContext';
import SubscriberContext from '../context/SubscriberContext';
import TraceContext from '../context/TraceContext';
// $FlowFixMe migrated to typescript
import type {EnodebInfo} from './EnodebUtils';
// $FlowFixMe migrated to typescript
import type {EnodebState} from '../context/EnodebContext';
import type {
  apn,
  base_name,
  base_name_record,
  call_trace,
  call_trace_config,
  feg_lte_network,
  feg_network,
  lte_gateway,
  lte_network,
  mutable_call_trace,
  mutable_cbsd,
  mutable_subscriber,
  mutable_subscribers,
  network_id,
  network_ran_configs,
  network_type,
  paginated_cbsds,
  policy_qos_profile,
  policy_rule,
  rating_group,
  subscriber_id,
  tier,
} from '../../../generated/MagmaAPIBindings';
import type {gatewayPoolsStateType} from '../context/GatewayPoolsContext';

import * as cbsdState from '../../state/lte/CbsdState';
// $FlowFixMe migrated to typescript
import {FEG_LTE, LTE} from '../../../shared/types/network';
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
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

type Props = {
  networkId: network_id,
  networkType: network_type,
  children: React.Node,
};

export function GatewayContextProvider(props: Props) {
  const {networkId} = props;
  const [lteGateways, setLteGateways] = useState<{[string]: lte_gateway}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        const lteGateways = await MagmaV1API.getLteByNetworkIdGateways({
          networkId,
        });
        setLteGateways(lteGateways);
      } catch (e) {
        enqueueSnackbar?.('failed fetching gateway information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
          UpdateGateway({networkId, setLteGateways, ...props}),
      }}>
      {props.children}
    </GatewayContext.Provider>
  );
}

export function CbsdContextProvider({networkId, children}: Props) {
  const enqueueSnackbar = useEnqueueSnackbar();

  const [isLoading, setIsLoading] = useState(false);
  const [fetchResponse, setFetchResponse] = useState<paginated_cbsds>({
    cbsds: [],
    total_count: 0,
  });
  const [paginationOptions, setPaginationOptions] = useState<{
    page: number,
    pageSize: number,
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
        create: (newCbsd: mutable_cbsd) => {
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
        update: (id: number, cbsd: mutable_cbsd) => {
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
  const [enbInfo, setEnbInfo] = useState<{[string]: EnodebInfo}>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();
  useEffect(() => {
    const fetchState = async () => {
      try {
        if (networkId == null) {
          return;
        }
        const [lteRanConfigsResp] = await Promise.allSettled([
          MagmaV1API.getLteByNetworkIdCellularRan({networkId}),
          InitEnodeState({networkId, setEnbInfo, enqueueSnackbar}),
        ]);
        if (lteRanConfigsResp.value) {
          setLteRanConfigs(lteRanConfigsResp.value);
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching enodeb information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
        setLteRanConfigs: lteRanConfigs => setLteRanConfigs(lteRanConfigs),
      }}>
      {props.children}
    </EnodebContext.Provider>
  );
}

export function TraceContextProvider(props: Props) {
  const {networkId} = props;
  const [traceMap, setTraceMap] = useState<{[string]: call_trace}>({});
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
    fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <TraceContext.Provider
      value={{
        state: traceMap,
        setState: (
          key: string,
          value?: mutable_call_trace | call_trace_config,
        ) =>
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
  const [forbiddenNetworkTypes, setForbiddenNetworkTypes] = useState({});
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
        setForbiddenNetworkTypes,
        setSubscriberMetrics,
        setSessionState,
        setTotalCount,
        enqueueSnackbar,
      }),
        setIsLoading(false);
    };
    fetchLteState();
  }, [networkId, enqueueSnackbar]);

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <SubscriberContext.Provider
      value={{
        forbiddenNetworkTypes: forbiddenNetworkTypes,
        state: subscriberMap,
        metrics: subscriberMetrics,
        sessionState: sessionState,
        totalCount: totalCount,
        gwSubscriberMap: getGatewaySubscriberMap(sessionState),
        setState: (
          key: subscriber_id,
          value?: mutable_subscriber | mutable_subscribers,
          newState?,
          newSessionState?,
        ) =>
          setSubscriberState({
            networkId,
            subscriberMap,
            setSubscriberMap,
            setSessionState,
            setForbiddenNetworkTypes,
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
  const [gatewayPools, setGatewayPools] = useState<{
    [string]: gatewayPoolsStateType,
  }>({});
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
    fetchState();
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
  const [tiers, setTiers] = useState<{[string]: tier}>({});
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
          MagmaV1API.getChannelsByChannelId({channelId: 'stable'}),
          InitTierState({networkId, setTiers, enqueueSnackbar}),
        ]);
        if (stableChannelResp.value) {
          setSupportedVersions(
            stableChannelResp.value.supported_versions.reverse(),
          );
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching tier information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
  const [policies, setPolicies] = useState<{[string]: policy_rule}>({});
  const [baseNames, setBaseNames] = useState<{
    [string]: base_name_record,
  }>({});
  const [qosProfiles, setQosProfiles] = useState<{
    [string]: policy_qos_profile,
  }>({});
  const [ratingGroups, setRatingGroups] = useState<{[string]: rating_group}>(
    {},
  );
  const [fegNetwork, setFegNetwork] = useState<feg_network>({});
  const [fegPolicies, setFegPolicies] = useState<{[string]: policy_rule}>({});
  const [isLoading, setIsLoading] = useState(true);
  const networkType = networkCtx.networkType;
  const enqueueSnackbar = useEnqueueSnackbar();
  let fegNetworkId = '';
  if (networkType === FEG_LTE) {
    fegNetworkId = lteNetworkCtx.state?.federation.feg_network_id;
  }

  useEffect(() => {
    const fetchState = async () => {
      try {
        setPolicies(
          await MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull({
            networkId,
          }),
        );
        // Base Names
        // eslint-disable-next-line max-len
        const baseNameIDs: Array<base_name> = await MagmaV1API.getNetworksByNetworkIdPoliciesBaseNames(
          {
            networkId,
          },
        );
        const baseNameRecords: Array<base_name_record> = await Promise.all(
          baseNameIDs.map(baseNameID =>
            MagmaV1API.getNetworksByNetworkIdPoliciesBaseNamesByBaseName({
              networkId,
              baseName: baseNameID,
            }),
          ),
        );
        const newBaseNames: {[string]: base_name_record} = {};
        baseNameRecords.map(record => {
          newBaseNames[record.name] = record;
        });
        setBaseNames(newBaseNames);

        setRatingGroups(
          // $FlowIgnore
          await MagmaV1API.getNetworksByNetworkIdRatingGroups({networkId}),
        );
        setQosProfiles(
          await MagmaV1API.getLteByNetworkIdPolicyQosProfiles({networkId}),
        );
        if (fegNetworkId != null && fegNetworkId !== '') {
          setFegNetwork(
            await MagmaV1API.getFegByNetworkId({networkId: fegNetworkId}),
          );
          setFegPolicies(
            await MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull({
              networkId: fegNetworkId,
            }),
          );
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching policy information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
            const fegNetworkID = lteNetworkCtx.state?.federation.feg_network_id;
            await SetPolicyState({
              policies,
              setPolicies,
              networkId,
              key,
              value,
            });

            // duplicate the policy on feg_network as well
            if (fegNetworkID != null) {
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
                lteNetworkCtx.updateNetworks({
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
                  UpdateFegNetworkState({
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
              ruleNames = oldRuleNames.filter(function (
                ruleId,
                _unused0,
                _unused1,
              ) {
                return ruleId !== key;
              });
              lteNetworkCtx.updateNetworks({
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
              fegRuleNames = oldFegRuleNames.filter(function (
                ruleId,
                _unused0,
                _unused1,
              ) {
                return ruleId !== key;
              });

              if (networkType === FEG_LTE && fegNetwork) {
                UpdateFegNetworkState({
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
  const [apns, setApns] = useState<{[string]: apn}>({});
  const [isLoading, setIsLoading] = useState(true);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    const fetchState = async () => {
      try {
        setApns(
          await MagmaV1API.getLteByNetworkIdApns({
            networkId,
          }),
        );
      } catch (e) {
        enqueueSnackbar?.('failed fetching APN information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
    $Shape<lte_network & feg_lte_network>,
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
            MagmaV1API.getFegLteByNetworkId({networkId}),
            MagmaV1API.getFegLteByNetworkIdSubscriberConfig({networkId}),
          ]);
          if (fegLteResp.value) {
            let subscriber_config = {};
            if (fegLteSubscriberConfigResp.value) {
              subscriber_config = fegLteSubscriberConfigResp.value;
            }
            setLteNetwork({...fegLteResp.value, subscriber_config});
          }
        } else {
          setLteNetwork(await MagmaV1API.getLteByNetworkId({networkId}));
        }
      } catch (e) {
        enqueueSnackbar?.('failed fetching network information', {
          variant: 'error',
        });
      }
      setIsLoading(false);
    };
    fetchState();
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
    return props.children;
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
