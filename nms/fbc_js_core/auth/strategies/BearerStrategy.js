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
 * @flow
 * @format
 */

import type {FBCNMSMiddleWareRequest} from '../../../fbc_js_core/express_middleware';

import {Strategy} from 'passport-http-bearer';
import {clientFromRequest} from '../oidc/client';

export default function BearerTokenStrategy() {
  return new Strategy({passReqToCallback: true}, verify);
}

type TokenUser = {
  email: string,
  organization: string,
};

const verify = async (
  req: FBCNMSMiddleWareRequest,
  token: string,
  done: (?Error, ?TokenUser | ?boolean) => void,
) => {
  try {
    const user = await authenticateToken(token, req);
    if (!user) {
      throw new Error('Invalid token!');
    }
    return done(null, user);
  } catch (e) {
    done(e);
  }
};

const authenticateToken = async (
  accessToken: string,
  req: FBCNMSMiddleWareRequest,
): Promise<TokenUser> => {
  const org = await req.organization();
  const client = await clientFromRequest(req);
  const user = await client.userinfo(accessToken);
  return {...user, organization: org.name};
};
