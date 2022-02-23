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

import type {
  base_name_record,
  network_id,
  policy_id,
  policy_qos_profile,
  policy_rule,
  rating_group,
} from '../../generated/MagmaAPIBindings';

import MagmaV1API from '../../generated/WebClient';

type Props = {
  networkId: network_id,
  policies: {[string]: policy_rule},
  setPolicies: ({[string]: policy_rule}) => void,
  key: policy_id,
  value?: policy_rule,
};

export async function SetPolicyState(props: Props) {
  const {networkId, policies, setPolicies, key, value} = props;
  if (value != null) {
    if (!(key in policies)) {
      await MagmaV1API.postNetworksByNetworkIdPoliciesRules({
        networkId: networkId,
        policyRule: value,
      });
    } else {
      await MagmaV1API.putNetworksByNetworkIdPoliciesRulesByRuleId({
        networkId: networkId,
        ruleId: key,
        policyRule: value,
      });
    }
    // eslint-disable-next-line max-len
    const policyRule = await MagmaV1API.getNetworksByNetworkIdPoliciesRulesByRuleId(
      {
        networkId: networkId,
        ruleId: key,
      },
    );

    if (policyRule) {
      const newPolicies = {...policies, [key]: policyRule};
      setPolicies(newPolicies);
    }
  } else {
    await MagmaV1API.deleteNetworksByNetworkIdPoliciesRulesByRuleId({
      networkId: networkId,
      ruleId: key,
    });
    const newPolicies = {...policies};
    delete newPolicies[key];
    setPolicies(newPolicies);
  }
}

type BaseNameProps = {
  networkId: network_id,
  baseNames: {[string]: base_name_record},
  setBaseNames: ({[string]: base_name_record}) => void,
  key: string, // base name id
  value?: base_name_record,
};

export async function SetBaseNameState(props: BaseNameProps) {
  const {networkId, baseNames, setBaseNames, key, value} = props;
  if (value != null) {
    if (!(key in baseNames)) {
      await MagmaV1API.postNetworksByNetworkIdPoliciesBaseNames({
        networkId: networkId,
        baseNameRecord: value,
      });
    } else {
      await MagmaV1API.putNetworksByNetworkIdPoliciesBaseNamesByBaseName({
        networkId: networkId,
        baseName: key,
        baseNameRecord: value,
      });
    }
    // eslint-disable-next-line max-len
    const baseName = await MagmaV1API.getNetworksByNetworkIdPoliciesBaseNamesByBaseName(
      {
        networkId: networkId,
        baseName: key,
      },
    );

    if (baseName) {
      const newBaseNames = {...baseNames, [key]: baseName};
      setBaseNames(newBaseNames);
    }
  } else {
    await MagmaV1API.deleteNetworksByNetworkIdPoliciesBaseNamesByBaseName({
      networkId: networkId,
      baseName: key,
    });
    const newBaseNames = {...baseNames};
    delete newBaseNames[key];
    setBaseNames(newBaseNames);
  }
}

type QosProfileProps = {
  networkId: network_id,
  qosProfiles: {[string]: policy_qos_profile},
  setQosProfiles: ({[string]: policy_qos_profile}) => void,
  key: string,
  value?: policy_qos_profile,
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
      await MagmaV1API.postLteByNetworkIdPolicyQosProfiles({
        networkId: networkId,
        policy: value,
      });
    } else {
      await MagmaV1API.putLteByNetworkIdPolicyQosProfilesByProfileId({
        networkId: networkId,
        profileId: key,
        profile: value,
      });
    }
    // eslint-disable-next-line max-len
    const qosProfile = await MagmaV1API.getLteByNetworkIdPolicyQosProfilesByProfileId(
      {
        networkId: networkId,
        profileId: key,
      },
    );

    if (qosProfile) {
      const newPolicies = {...qosProfiles, [key]: qosProfile};
      setQosProfiles(newPolicies);
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdPolicyQosProfilesByProfileId({
      networkId: networkId,
      profileId: key,
    });
    const newQosProfiles = {...qosProfiles};
    delete newQosProfiles[key];
    setQosProfiles(newQosProfiles);
  }
}

type RatingGroupProps = {
  networkId: network_id,
  ratingGroups: {[string]: rating_group},
  setRatingGroups: ({[string]: rating_group}) => void,
  key: string,
  value?: rating_group,
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
      await MagmaV1API.postNetworksByNetworkIdRatingGroups({
        networkId: networkId,
        ratingGroup: value,
      });
    } else {
      await MagmaV1API.putNetworksByNetworkIdRatingGroupsByRatingGroupId({
        networkId: networkId,
        ratingGroupId: parseInt(key),
        ratingGroup: value,
      });
    }
    // eslint-disable-next-line max-len
    const ratingGroup = await MagmaV1API.getNetworksByNetworkIdRatingGroupsByRatingGroupId(
      {
        networkId: networkId,
        ratingGroupId: parseInt(key),
      },
    );

    if (ratingGroup) {
      const newRatingGroups = {...ratingGroups, [key]: ratingGroup};
      setRatingGroups(newRatingGroups);
    }
  } else {
    await MagmaV1API.deleteNetworksByNetworkIdRatingGroupsByRatingGroupId({
      networkId: networkId,
      ratingGroupId: parseInt(key),
    });
    const newRatingGroups = {...ratingGroups};
    delete newRatingGroups[key];
    setRatingGroups(newRatingGroups);
  }
}
