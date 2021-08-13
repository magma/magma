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
  network_id,
  policy_id,
  policy_qos_profile,
  policy_rule,
  rating_group,
} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

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

/* SetQosProfileState
SetQosProfileState
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
