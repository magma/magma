/*
 * Copyright 2022 The Magma Authors.
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

import {TokenSet} from 'openid-client';
import type {OrganizationModel} from './shared/sequelize_models/models/organization';
import type {UserModel} from './shared/sequelize_models/models/user';

declare global {
  namespace Express {
    type User = UserModel;

    interface Request {
      csrfToken: () => string; // from csrf
      body: object; // from bodyParser
      session?: {oidc?: {tokenSet: TokenSet}};
      organization?: () => Promise<OrganizationModel>;
      logIn: (
        user: UserModel,
        callback: (err?: Error | null | undefined) => void,
      ) => void;
      logOut: () => void;
      logout: () => void;
      user: UserModel;
      isAuthenticated: () => boolean;
      isUnauthenticated: () => boolean;
      access: {loginUrl: string};
    }
  }
}
