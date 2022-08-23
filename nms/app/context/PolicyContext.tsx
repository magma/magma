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
import LteNetworkContext from './LteNetworkContext';
import MagmaAPI from '../api/MagmaAPI';
import NetworkContext from './NetworkContext';
import React, {useCallback, useContext, useEffect, useState} from 'react';
import {
  BaseNameRecord,
  FegNetwork,
  PolicyQosProfile,
  PolicyRule,
  RatingGroup,
} from '../../generated';
import {FEG_LTE, NetworkId, PolicyId} from '../../shared/types/network';
import {omit} from 'lodash';
import {updateFegNetworkState} from './FEGNetworkContext';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

export type PolicyContextType = {
  state: Record<string, PolicyRule>;
  baseNames: Record<string, BaseNameRecord>;
  setBaseNames: (key: string, val?: BaseNameRecord) => Promise<void>;
  ratingGroups: Record<string, RatingGroup>;
  setRatingGroups: (key: string, val?: RatingGroup) => Promise<void>;
  qosProfiles: Record<string, PolicyQosProfile>;
  setQosProfiles: (key: string, val?: PolicyQosProfile) => Promise<void>;
  setState: (
    key: PolicyId,
    val?: PolicyRule,
    isNetworkWide?: boolean,
  ) => Promise<void>;
  refetch: () => void;
};
type PolicyProviderProps = {
  networkId: NetworkId;
  children: React.ReactNode;
};

const PolicyContext = React.createContext<PolicyContextType>(
  {} as PolicyContextType,
);

async function setPolicyState(params: {
  networkId: NetworkId;
  policies: Record<string, PolicyRule>;
  setPolicies: (policies: Record<string, PolicyRule>) => void;
  key: PolicyId;
  value?: PolicyRule;
}) {
  const {networkId, policies, setPolicies, key, value} = params;

  if (value != null) {
    if (!(key in policies)) {
      await MagmaAPI.policies.networksNetworkIdPoliciesRulesPost({
        networkId: networkId,
        policyRule: value,
      });
    } else {
      await MagmaAPI.policies.networksNetworkIdPoliciesRulesRuleIdPut({
        networkId: networkId,
        ruleId: key,
        policyRule: value,
      });
    }

    const policyRule = (
      await MagmaAPI.policies.networksNetworkIdPoliciesRulesRuleIdGet({
        networkId: networkId,
        ruleId: key,
      })
    ).data;

    if (policyRule) {
      const newPolicies = {...policies, [key]: policyRule};
      setPolicies(newPolicies);
    }
  } else {
    await MagmaAPI.policies.networksNetworkIdPoliciesRulesRuleIdDelete({
      networkId: networkId,
      ruleId: key,
    });
    const newPolicies = {...policies};
    delete newPolicies[key];
    setPolicies(newPolicies);
  }
}

async function setBaseNameState(params: {
  networkId: NetworkId;
  baseNames: Record<string, BaseNameRecord>;
  setBaseNames: (baseNames: Record<string, BaseNameRecord>) => void;
  key: string;
  // base name id
  value?: BaseNameRecord;
}) {
  const {networkId, baseNames, setBaseNames, key, value} = params;

  if (value != null) {
    if (!(key in baseNames)) {
      await MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesPost({
        networkId: networkId,
        baseNameRecord: value,
      });
    } else {
      await MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesBaseNamePut({
        networkId: networkId,
        baseName: key,
        baseNameRecord: value,
      });
    }

    const baseName = (
      await MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesBaseNameGet({
        networkId: networkId,
        baseName: key,
      })
    ).data;

    if (baseName) {
      const newBaseNames = {...baseNames, [key]: baseName};
      setBaseNames(newBaseNames);
    }
  } else {
    await MagmaAPI.policies.networksNetworkIdPoliciesBaseNamesBaseNameDelete({
      networkId: networkId,
      baseName: key,
    });
    const newBaseNames = {...baseNames};
    delete newBaseNames[key];
    setBaseNames(newBaseNames);
  }
}

/** setQosProfileState
 * if key and value are passed in,
 * if key is not present, a new profile is created (POST)
 * if key is present, existing profile is updated (PUT)
 * if value is not present, the profile is deleted (DELETE)
 */
async function setQosProfileState(params: {
  networkId: NetworkId;
  qosProfiles: Record<string, PolicyQosProfile>;
  setQosProfiles: (qosProfiles: Record<string, PolicyQosProfile>) => void;
  key: string;
  value?: PolicyQosProfile;
}) {
  const {networkId, qosProfiles, setQosProfiles, key, value} = params;

  if (value != null) {
    if (!(key in qosProfiles)) {
      await MagmaAPI.policies.lteNetworkIdPolicyQosProfilesPost({
        networkId: networkId,
        policy: value,
      });
    } else {
      await MagmaAPI.policies.lteNetworkIdPolicyQosProfilesProfileIdPut({
        networkId: networkId,
        profileId: key,
        profile: value,
      });
    }

    const qosProfile = (
      await MagmaAPI.policies.lteNetworkIdPolicyQosProfilesProfileIdGet({
        networkId: networkId,
        profileId: key,
      })
    ).data;

    if (qosProfile) {
      const newPolicies = {...qosProfiles, [key]: qosProfile};
      setQosProfiles(newPolicies);
    }
  } else {
    await MagmaAPI.policies.lteNetworkIdPolicyQosProfilesProfileIdDelete({
      networkId: networkId,
      profileId: key,
    });
    const newQosProfiles = {...qosProfiles};
    delete newQosProfiles[key];
    setQosProfiles(newQosProfiles);
  }
}

/* setRatingGroupState
 * if key and value are passed in,
 * if key is not present, a new profile is created (POST)
 * if key is present, existing profile is updated (PUT)
 * if value is not present, the profile is deleted (DELETE)
 */
async function setRatingGroupState(params: {
  networkId: NetworkId;
  ratingGroups: Record<string, RatingGroup>;
  setRatingGroups: (ratingGroups: Record<string, RatingGroup>) => void;
  key: string;
  value?: RatingGroup;
}) {
  const {networkId, ratingGroups, setRatingGroups, key, value} = params;

  if (value != null) {
    if (!(key in ratingGroups)) {
      await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsPost({
        networkId: networkId,
        ratingGroup: value,
      });
    } else {
      await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsRatingGroupIdPut(
        {
          networkId: networkId,
          ratingGroupId: parseInt(key),
          ratingGroup: value,
        },
      );
    }

    const ratingGroup = (
      await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsRatingGroupIdGet(
        {
          networkId: networkId,
          ratingGroupId: parseInt(key),
        },
      )
    ).data;

    if (ratingGroup) {
      const newRatingGroups = {...ratingGroups, [key]: ratingGroup};
      setRatingGroups(newRatingGroups);
    }
  } else {
    await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsRatingGroupIdDelete(
      {
        networkId: networkId,
        ratingGroupId: parseInt(key),
      },
    );
    const newRatingGroups = {...ratingGroups};
    delete newRatingGroups[key];
    setRatingGroups(newRatingGroups);
  }
}

export function PolicyProvider(props: PolicyProviderProps) {
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

  const fetchState = useCallback(async () => {
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
        (
          await MagmaAPI.ratingGroups.networksNetworkIdRatingGroupsGet({
            networkId,
          })
        ).data,
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
            await MagmaAPI.policies.networksNetworkIdPoliciesRulesviewfullGet({
              networkId: fegNetworkId,
            })
          ).data,
        );
      }
    } catch (e) {
      enqueueSnackbar?.('failed fetching policy information', {
        variant: 'error',
      });
    }
    setIsLoading(false);
  }, [networkId, fegNetworkId, enqueueSnackbar]);

  useEffect(() => void fetchState(), [fetchState, networkType]);

  if (isLoading) {
    return <LoadingFiller />;
  }
  return (
    <PolicyContext.Provider
      value={{
        state: policies,
        ratingGroups: ratingGroups,
        setRatingGroups: async (key, value) => {
          await setRatingGroupState({
            networkId,
            ratingGroups,
            setRatingGroups,
            key,
            value,
          });
        },

        baseNames: baseNames,
        setBaseNames: async (key, value) => {
          await setBaseNameState({
            networkId,
            baseNames,
            setBaseNames,
            key,
            value,
          });
        },

        qosProfiles: qosProfiles,
        setQosProfiles: async (key, value) => {
          await setQosProfileState({
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
            await setPolicyState({
              policies,
              setPolicies,
              networkId,
              key,
              value,
            });

            // duplicate the policy on feg_network as well
            if (fegNetworkID) {
              await setPolicyState({
                policies: fegPolicies,
                setPolicies: setFegPolicies,
                networkId: fegNetworkID,
                key,
                value: omit(value, 'qos_profile'),
              });
            }
          } else {
            await setPolicyState({
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
                  void updateFegNetworkState({
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
                void updateFegNetworkState({
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
        refetch: () => void fetchState(),
      }}>
      {props.children}
    </PolicyContext.Provider>
  );
}

export default PolicyContext;
