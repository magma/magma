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
import type {GatewayId} from '../../../shared/types/network';
import type {LteGateway, MutableLteGateway} from '../../../generated-ts';
import type {UpdateGatewayProps} from '../../state/lte/EquipmentState';

import React from 'react';

export type GatewayContextType = {
  state: Record<string, LteGateway>;
  setState: (
    key: GatewayId,
    val?: MutableLteGateway,
    newState?: Record<string, LteGateway>,
  ) => Promise<void>;
  updateGateway: (props: Partial<UpdateGatewayProps>) => Promise<void>;
};

export default React.createContext<GatewayContextType>(
  {} as GatewayContextType,
);
