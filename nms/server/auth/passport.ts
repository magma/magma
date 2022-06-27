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

import passport from 'passport';
import {User} from '../../shared/sequelize_models';

import type {FBCNMSMiddleWareRequest} from '../../server/middleware';
import type {UserModel} from '../../shared/sequelize_models/models/user';

type OutputRequest<T> = {
  logIn: (user: T, callback: (err?: Error | null | undefined) => void) => void;
  logOut: () => void;
  logout: () => void;
  user: T;
  isAuthenticated: () => boolean;
  isUnauthenticated: () => boolean;
} & FBCNMSMiddleWareRequest;
export type FBCNMSPassportRequest = OutputRequest<UserModel>;

function use() {
  passport.serializeUser((user, done) => {
    done(null, user.id);
  });

  passport.deserializeUser<number>((id, done) => {
    User.findByPk(id)
      .then(user => done(null, user))
      .catch(error => done(error));
  });
}

export default {
  use,
};
