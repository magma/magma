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
import type {FegLteNetwork, LteNetwork} from '../../../generated-ts';
import type {UpdateNetworkProps as FegLteUpdateNetworkProps} from '../../state/feg_lte/NetworkState';
import type {UpdateNetworkProps as LteUpdateNetworkProps} from '../../state/lte/NetworkState';

import React from 'react';

export type UpdateNetworkContextProps = Partial<
  LteUpdateNetworkProps & FegLteUpdateNetworkProps
>;

export type LteNetworkContextType = {
  state: Partial<LteNetwork & FegLteNetwork>;
  updateNetworks: (props: UpdateNetworkContextProps) => Promise<void>;
};

export default React.createContext<LteNetworkContextType>(
  {} as LteNetworkContextType,
);
