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
import type {EnodebInfo} from '../lte/EnodebUtils';
import type {NetworkRanConfigs} from '../../../generated-ts';

import React from 'react';

export type EnodebState = {
  enbInfo: Record<string, EnodebInfo>;
};

export type EnodebContextType = {
  state: EnodebState;
  lteRanConfigs?: NetworkRanConfigs;
  setState: (
    key: string,
    val?: EnodebInfo,
    newState?: EnodebState,
  ) => Promise<void>;
};

export default React.createContext<EnodebContextType>({} as EnodebContextType);
