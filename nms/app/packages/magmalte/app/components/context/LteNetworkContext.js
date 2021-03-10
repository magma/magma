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
import type {UpdateNetworkProps as FegLteUpdateNetworkProps} from '../../state/feg_lte/NetworkState';
import type {UpdateNetworkProps as LteUpdateNetworkProps} from '../../state/lte/NetworkState';
import type {feg_lte_network, lte_network} from '@fbcnms/magma-api';

import React from 'react';

export type LteNetworkContextType = {
  state: $Shape<lte_network & feg_lte_network>,
  updateNetworks: (
    props: $Shape<LteUpdateNetworkProps & FegLteUpdateNetworkProps>,
  ) => Promise<void>,
};

export default React.createContext<LteNetworkContextType>({});
