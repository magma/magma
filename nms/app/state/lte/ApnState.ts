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
 */

import MagmaAPI from '../../../api/MagmaAPI';
import {Apn} from '../../../generated-ts';
import {NetworkId} from '../../../shared/types/network';

type Props = {
  networkId: NetworkId;
  apns: Record<string, Apn>;
  setApns: (arg0: Record<string, Apn>) => void;
  key: string;
  value?: Apn;
};

export async function SetApnState(props: Props) {
  const {networkId, apns, setApns, key, value} = props;

  if (value != null) {
    if (!(key in apns)) {
      await MagmaAPI.apns.lteNetworkIdApnsPost({
        networkId: networkId,
        apn: value,
      });
      setApns({...apns, [key]: value});
    } else {
      await MagmaAPI.apns.lteNetworkIdApnsApnNamePut({
        networkId: networkId,
        apnName: key,
        apn: value,
      });
      setApns({...apns, [key]: value});
    }

    const apn = (
      await MagmaAPI.apns.lteNetworkIdApnsApnNameGet({
        networkId: networkId,
        apnName: key,
      })
    ).data;

    if (apn) {
      const newApns = {...apns, [key]: apn};
      setApns(newApns);
    }
  } else {
    await MagmaAPI.apns.lteNetworkIdApnsApnNameDelete({
      networkId: networkId,
      apnName: key,
    });
    const newApns = {...apns};
    delete newApns[key];
    setApns(newApns);
  }
}
