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

export const ACTION = {
  PERMIT: 'PERMIT',
  DENY: 'DENY',
};

export const DIRECTION = {
  UPLINK: 'UPLINK',
  DOWNLINK: 'DOWNLINK',
};

export const PROTOCOL = {
  IPPROTO_IP: 'IPPROTO_IP',
  IPPROTO_UDP: 'IPPROTO_UDP',
  IPPROTO_TCP: 'IPPROTO_TCP',
  IPPROTO_ICMP: 'IPPROTO_ICMP',
};
