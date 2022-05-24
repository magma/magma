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

import type {
  MutableSubscriber,
  Subscriber,
  SubscriberState,
} from '../../../generated-ts';

import React from 'react';
import {GatewayId, SubscriberId} from '../../../shared/types/network';
import {SubscriberForbiddenNetworkTypesEnum} from '../../../generated-ts';

export type Metrics = {
  currentUsage: string;
  dailyAvg: string;
};

/* SubscriberContextType
state: paginated subscribers
sessionState: paginated subscribers session state
metrics: subscriber metrics
gwSubscriberMap: gateway subscriber map
totalCount: total count of subscribers
setState: POST, PUT, DELETE subscriber
*/
export type SubscriberContextType = {
  state: Record<string, Subscriber>;
  sessionState: Record<string, SubscriberState>;
  forbiddenNetworkTypes: Record<
    string,
    Array<SubscriberForbiddenNetworkTypesEnum>
  >;
  metrics?: Record<string, Metrics>;
  gwSubscriberMap: Record<GatewayId, Array<SubscriberId>>;
  totalCount: number;
  setState?: (
    key: string,
    val?: MutableSubscriber | Array<MutableSubscriber>,
    newState?: Record<string, Subscriber>,
    newSessionState?: Record<string, SubscriberState>,
  ) => Promise<void>;
};

export default React.createContext<SubscriberContextType>(
  {} as SubscriberContextType,
);
