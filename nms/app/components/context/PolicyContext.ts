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
import React from 'react';
import type {
  BaseNameRecord,
  PolicyQosProfile,
  PolicyRule,
  RatingGroup,
} from '../../../generated-ts';
import type {PolicyId} from '../../../shared/types/network';

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
};

export default React.createContext<PolicyContextType>({} as PolicyContextType);
