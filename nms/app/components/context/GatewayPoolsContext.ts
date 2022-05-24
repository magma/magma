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
  CellularGatewayPool,
  CellularGatewayPoolRecord,
  MutableCellularGatewayPool,
} from '../../../generated-ts';
import type {GatewayPoolId} from '../../../shared/types/network';

// add gateway ID to gateway pool records (gateway primary/secondary)
export type GatewayPoolRecordsType = {
  gateway_id: string;
} & CellularGatewayPoolRecord;

export type gatewayPoolsStateType = {
  gatewayPool: CellularGatewayPool;
  gatewayPoolRecords: Array<GatewayPoolRecordsType>;
};

/* GatewayPoolsContextType
state: gateway pool config and associated gateway pool records
setState: POST, PUT, DELETE gateway pool config
updateGatewayPoolRecords: POST, PUT, DELETE gateway pool records
*/
export type GatewayPoolsContextType = {
  state: Record<string, gatewayPoolsStateType>;
  setState: (
    key: GatewayPoolId,
    val?: MutableCellularGatewayPool,
  ) => Promise<void>;
  updateGatewayPoolRecords: (
    key: GatewayPoolId,
    val?: MutableCellularGatewayPool,
    resources?: Array<GatewayPoolRecordsType>,
  ) => Promise<void>;
};

export default React.createContext<GatewayPoolsContextType>(
  {} as GatewayPoolsContextType,
);
