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
  gateway_id,
  mutable_subscriber,
  subscriber,
  subscriber_id,
  subscriber_state,
} from '@fbcnms/magma-api';

import React from 'react';

export type Metrics = {
  currentUsage: string,
  dailyAvg: string,
};

export type SubscriberContextType = {
  state: {[string]: subscriber},
  metrics?: {[string]: Metrics},
  gwSubscriberMap: {[gateway_id]: Array<subscriber_id>},
  sessionState: {[string]: subscriber_state},
  setState?: (
    key: string,
    val?: mutable_subscriber,
    newState?: {
      state: {[string]: subscriber},
      sessionState: {[string]: subscriber_state},
    },
  ) => Promise<void>,
};

export default React.createContext<SubscriberContextType>({});
