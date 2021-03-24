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
import type {
  cellular_gateway_pool,
  cellular_gateway_pool_record,
  gateway_pool_id,
  mutable_cellular_gateway_pool,
} from '@fbcnms/magma-api';

import React from 'react';

export type GatewayPoolRecordsType = {
  gateway_id: string,
} & cellular_gateway_pool_record;
export type gatewayPoolsStateType = {
  gatewayPool: cellular_gateway_pool,
  gatewayPoolRecords: Array<GatewayPoolRecordsType>,
};
export type GatewayPoolsContextType = {
  state: {[string]: gatewayPoolsStateType},
  setState: (
    key: gateway_pool_id,
    val?: mutable_cellular_gateway_pool,
    resources?: Array<GatewayPoolRecordsType>,
  ) => Promise<void>,
  updateGatewayPoolResources: (
    key: gateway_pool_id,
    val?: mutable_cellular_gateway_pool,
    resources?: Array<GatewayPoolRecordsType>,
  ) => Promise<void>,
};

export default React.createContext<GatewayPoolsContextType>({});
