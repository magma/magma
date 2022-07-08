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

import type {FederationGatewayHealthStatus} from '../GatewayUtils';

import React from 'react';
import {
  FederationGateway,
  MutableFederationGateway,
} from '../../../generated-ts';
import {GatewayId} from '../../../shared/types/network';

export type FEGGatewayContextType = {
  state: Record<string, FederationGateway>;
  setState: (
    key: GatewayId,
    val?: MutableFederationGateway,
    newState?: Record<string, FederationGateway>,
  ) => Promise<void>;
  health: Record<GatewayId, FederationGatewayHealthStatus>;
  activeFegGatewayId: GatewayId;
};

export default React.createContext<FEGGatewayContextType>(
  {} as FEGGatewayContextType,
);
