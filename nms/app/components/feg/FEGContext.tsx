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

import * as React from 'react';
import {FEGGatewayContextProvider} from '../../context/FEGGatewayContext';
import {FEGNetworkContextProvider} from '../../context/FEGNetworkContext';
import {FEGSubscriberContextProvider} from '../../context/FEGSubscriberContext';
import {GatewayTierContextProvider} from '../../context/GatewayTierContext';
import {NetworkId} from '../../../shared/types/network';

/**
 * A context provider for federation networks. It is used in sharing
 * information like the network information or the gateways information.
 * @param {object} props contains the network id and its type
 */
export function FEGContextProvider(props: {
  networkId: NetworkId;
  children: React.ReactNode;
}) {
  const {networkId} = props;

  return (
    <FEGNetworkContextProvider networkId={networkId}>
      <FEGSubscriberContextProvider networkId={networkId}>
        <GatewayTierContextProvider networkId={networkId}>
          <FEGGatewayContextProvider networkId={networkId}>
            {props.children}
          </FEGGatewayContextProvider>
        </GatewayTierContextProvider>
      </FEGSubscriberContextProvider>
    </FEGNetworkContextProvider>
  );
}
