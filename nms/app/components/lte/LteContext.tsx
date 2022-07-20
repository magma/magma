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
import {ApnContextProvider} from '../../context/ApnContext';
import {CbsdContextProvider} from '../../context/CbsdContext';
import {EnodebContextProvider} from '../../context/EnodebContext';
import {FEG_LTE, LTE, NetworkId} from '../../../shared/types/network';
import {GatewayContextProvider} from '../../context/GatewayContext';
import {GatewayPoolsContextProvider} from '../../context/GatewayPoolsContext';
import {GatewayTierContextProvider} from '../../context/GatewayTierContext';
import {LteNetworkContextProvider} from '../../context/LteNetworkContext';
import {PolicyProvider} from '../../context/PolicyContext';
import {SubscriberContextProvider} from '../../context/SubscriberContext';
import {TraceContextProvider} from '../../context/TraceContext';

type Props = {
  networkId: NetworkId;
  networkType: string;
  children: React.ReactNode;
};

export function LteContextProvider(props: Props) {
  const {networkId, networkType} = props;
  const lteNetwork = networkType === LTE || networkType === FEG_LTE;
  if (!lteNetwork) {
    return <>{props.children}</>;
  }

  return (
    <LteNetworkContextProvider networkId={networkId}>
      <PolicyProvider networkId={networkId}>
        <ApnContextProvider networkId={networkId}>
          <SubscriberContextProvider networkId={networkId}>
            <GatewayTierContextProvider networkId={networkId}>
              <EnodebContextProvider networkId={networkId}>
                <GatewayContextProvider networkId={networkId}>
                  <GatewayPoolsContextProvider networkId={networkId}>
                    <TraceContextProvider networkId={networkId}>
                      <CbsdContextProvider networkId={networkId}>
                        {props.children}
                      </CbsdContextProvider>
                    </TraceContextProvider>
                  </GatewayPoolsContextProvider>
                </GatewayContextProvider>
              </EnodebContextProvider>
            </GatewayTierContextProvider>
          </SubscriberContextProvider>
        </ApnContextProvider>
      </PolicyProvider>
    </LteNetworkContextProvider>
  );
}
