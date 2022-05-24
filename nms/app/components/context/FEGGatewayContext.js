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

// $FlowFixMe migrated to typescript
import type {FederationGatewayHealthStatus} from '../../components/GatewayUtils';
import type {
  federation_gateway,
  gateway_id,
  mutable_federation_gateway,
} from '../../../generated/MagmaAPIBindings';

import React from 'react';

export type FEGGatewayContextType = {
  state: {[string]: federation_gateway},
  setState: (
    key: gateway_id,
    val?: mutable_federation_gateway,
    newState?: {[string]: federation_gateway},
  ) => Promise<void>,
  health: {[gateway_id]: FederationGatewayHealthStatus},
  activeFegGatewayId: gateway_id,
};

export default React.createContext<FEGGatewayContextType>({});
