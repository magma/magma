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
 */

import type {OrganizationMiddlewareRequest} from '../../middleware/organizationMiddleware';

import {Client, Issuer} from 'openid-client';

const _clientCache: Record<string, Client> = {};

export async function clientFromRequest(req: OrganizationMiddlewareRequest) {
  if (!req.organization) {
    throw new Error('Must be using organization');
  }

  const {
    name,
    ssoOidcClientID,
    ssoOidcClientSecret,
    ssoOidcConfigurationURL,
  } = await req.organization();

  if (_clientCache[name]) {
    return _clientCache[name];
  }

  const issuer = await Issuer.discover(ssoOidcConfigurationURL);
  _clientCache[name] = new issuer.Client({
    client_id: ssoOidcClientID,
    client_secret: ssoOidcClientSecret,
    response_types: ['code'],
  });
  return _clientCache[name];
}
