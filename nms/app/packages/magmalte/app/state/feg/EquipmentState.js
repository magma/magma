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
    federation_gateway,
    mutable_federation_gateway,
    network_id,
  } from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
/**************************** Gateway State **********************************/
type GatewayStateProps = {
    networkId: network_id,
    fegGateways: {[string]: federation_gateway},
    setFegGateways: ({[string]: federation_gateway}) => void,
    key: gateway_id,
    value?: mutable_federation_gateway,
    newState?: {[string]: federation_gateway},
};

export async function SetGatewayState(props: GatewayStateProps) {
    const {networkId, fegGateways, setFegGateways, key, value, newState} = props;
    if (newState) {
        setFegGateways(newState);
        return;
    }
    if (value != null) {
        if (!(key in fegGateways)) {
        await MagmaV1API.postFegByNetworkIdGateways({
            networkId: networkId,
            gateway: value,
        });
        setFegGateways({...fegGateways, [key]: value});
        } else {
        await MagmaV1API.putFegByNetworkIdGatewaysByGatewayId({
            networkId: networkId,
            gatewayId: key,
            gateway: value,
        });
        setFegGateways({...fegGateways, [key]: value});
        }
    } else {
        await MagmaV1API.deleteFegByNetworkIdGatewaysByGatewayId({
        networkId: networkId,
        gatewayId: key,
        });
        const newFegGateways = {...fegGateways};
        delete newFegGateways[key];
        setFegGateways(newFegGateways);
    }
}
