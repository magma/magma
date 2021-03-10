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
'use strict';
import type {
  policy_id,
  policy_qos_profile,
  policy_rule,
  rating_group,
} from '@fbcnms/magma-api';

import React from 'react';

export type PolicyContextType = {
  state: {[string]: policy_rule},
  ratingGroups: {[string]: rating_group},
  setRatingGroups: (key: string, val?: rating_group) => Promise<void>,
  qosProfiles: {[string]: policy_qos_profile},
  setQosProfiles: (key: string, val?: policy_qos_profile) => Promise<void>,
  setState: (
    key: policy_id,
    val?: policy_rule,
    isNetworkWide?: boolean,
  ) => Promise<void>,
};

export default React.createContext<PolicyContextType>({});
