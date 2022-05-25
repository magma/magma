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
import type {UpdateGatewayProps} from '../../state/lte/EquipmentState';
import type {
  gateway_id,
  lte_gateway,
  mutable_lte_gateway,
} from '../../../generated/MagmaAPIBindings';

import React from 'react';

export type GatewayContextType = {
  state: {[string]: lte_gateway},
  setState: (
    key: gateway_id,
    val?: mutable_lte_gateway,
    newState?: {[string]: lte_gateway},
  ) => Promise<void>,
  updateGateway: (props: $Shape<UpdateGatewayProps>) => Promise<void>,
};

export default React.createContext<GatewayContextType>({});
