/*
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
import type {apn, network_id} from '@fbcnms/magma-api';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

type Props = {
  networkId: network_id,
  apns: {[string]: apn},
  setApns: ({[string]: apn}) => void,
  key: string,
  value?: apn,
};

export async function SetApnState(props: Props) {
  const {networkId, apns, setApns, key, value} = props;
  if (value != null) {
    if (!(key in apns)) {
      await MagmaV1API.postLteByNetworkIdApns({
        networkId: networkId,
        apn: value,
      });
      setApns({...apns, [key]: value});
    } else {
      await MagmaV1API.putLteByNetworkIdApnsByApnName({
        networkId: networkId,
        apnName: key,
        apn: value,
      });
      setApns({...apns, [key]: value});
    }
    const apn = await MagmaV1API.getLteByNetworkIdApnsByApnName({
      networkId: networkId,
      apnName: key,
    });
    if (apn) {
      const newApns = {...apns, [key]: apn};
      setApns(newApns);
    }
  } else {
    await MagmaV1API.deleteLteByNetworkIdApnsByApnName({
      networkId: networkId,
      apnName: key,
    });
    const newApns = {...apns};
    delete newApns[key];
    setApns(newApns);
  }
}
