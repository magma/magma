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

import MagmaAPI from '../../api/MagmaAPI';
import type {
  BaseNameRecord,
  PolicyQosProfile,
  PolicyRule,
  RatingGroup,
} from '../../generated-ts';
import type {NetworkId, PolicyId} from '../../shared/types/network';

type Props = {
  networkId: NetworkId;
  policies: Record<string, PolicyRule>;
  setPolicies: (arg0: Record<string, PolicyRule>) => void;
  key: PolicyId;
  value?: PolicyRule;
};

export async function SetPolicyState(props: Props) {
  const {networkId, policies, setPolicies, key, value} = props;

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

type BaseNameProps = {
  networkId: NetworkId;
  baseNames: Record<string, BaseNameRecord>;
  setBaseNames: (arg0: Record<string, BaseNameRecord>) => void;
  key: string;
  // base name id
  value?: BaseNameRecord;
};

export async function SetBaseNameState(props: BaseNameProps) {
  const {networkId, baseNames, setBaseNames, key, value} = props;

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

type QosProfileProps = {
  networkId: NetworkId;
  qosProfiles: Record<string, PolicyQosProfile>;
  setQosProfiles: (arg0: Record<string, PolicyQosProfile>) => void;
  key: string;
  value?: PolicyQosProfile;
};

/* SetQosProfileState
SetQosProfileState
if key and value are passed in,
if key is not present, a new profile is created (POST)
if key is present, existing profile is updated (PUT)
if value is not present, the profile is deleted (DELETE)
*/
export async function SetQosProfileState(props: QosProfileProps) {
  const {networkId, qosProfiles, setQosProfiles, key, value} = props;

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

type RatingGroupProps = {
  networkId: NetworkId;
  ratingGroups: Record<string, RatingGroup>;
  setRatingGroups: (arg0: Record<string, RatingGroup>) => void;
  key: string;
  value?: RatingGroup;
};

/* SetRatingGroupState
SetRatingGroupState
if key and value are passed in,
if key is not present, a new profile is created (POST)
if key is present, existing profile is updated (PUT)
if value is not present, the profile is deleted (DELETE)
*/
export async function SetRatingGroupState(props: RatingGroupProps) {
  const {networkId, ratingGroups, setRatingGroups, key, value} = props;

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
