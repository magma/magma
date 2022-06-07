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

export const ACTION = {
  PERMIT: 'PERMIT',
  DENY: 'DENY',
} as const;

export const APP_NAME = {
  NO_APP_NAME: 'NO_APP_NAME',
  FACEBOOK: 'FACEBOOK',
  FACEBOOK_MESSENGER: 'FACEBOOK_MESSENGER',
  INSTAGRAM: 'INSTAGRAM',
  YOUTUBE: 'YOUTUBE',
  GOOGLE: 'GOOGLE',
  GMAIL: 'GMAIL',
  GOOGLE_DOCS: 'GOOGLE DOCS',
  NETFLIX: 'NETFLIX',
  APPLE: 'APPLE',
  MICROSOFT: 'MICROSOFT',
  REDDIT: 'REDDIT',
  WHATSAPP: 'WHATSAPP',
  GOOGLE_PLAY: 'GOOGLE_PLAY',
  APPSTORE: 'APPSTORE',
  AMAZON: 'AMAZON',
  WECHAT: 'WECHAT',
  TIKTOK: 'TIKTOK',
  TWITTER: 'TWITTER',
  WIKIPEDIA: 'WIKIPEDIA',
  GOOGLE_MAPS: 'GOOGLE_MAPS',
  YAHOO: 'YAHOO',
  IMO: 'IMO',
} as const;

export const APP_SERVICE_TYPE = {
  NO_SERVICE_TYPE: 'NO_SERVICE_TYPE',
  CHAT: 'CHAT',
  AUDIO: 'AUDIO',
  VIDEO: 'VIDEO',
} as const;

export const DIRECTION = {
  UPLINK: 'UPLINK',
  DOWNLINK: 'DOWNLINK',
} as const;

export const PROTOCOL = {
  IPPROTO_IP: 'IPPROTO_IP',
  IPPROTO_UDP: 'IPPROTO_UDP',
  IPPROTO_TCP: 'IPPROTO_TCP',
  IPPROTO_ICMP: 'IPPROTO_ICMP',
} as const;
