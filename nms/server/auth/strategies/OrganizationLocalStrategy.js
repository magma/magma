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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSMiddleWareRequest} from '../../../server/middleware';
// $FlowFixMe migrated to typescript
import type {UserModel} from '../../../shared/sequelize_models/models/user';

import bcrypt from 'bcryptjs';
import {Strategy as LocalStrategy} from 'passport-local';
// $FlowFixMe migrated to typescript
import {getUserFromRequest} from '../util';

export default function () {
  return new LocalStrategy(
    {
      usernameField: 'email',
      passwordField: 'password',
      passReqToCallback: true,
    },
    validateUser,
  );
}

export async function validateUser(
  req: FBCNMSMiddleWareRequest,
  email: string,
  password: string,
  done: (?Error, UserModel | ?boolean, ?{message: string}) => void,
) {
  try {
    const user = await getUserFromRequest(req, email);
    if (!user) {
      return done(null, false, {
        message: 'Username or password invalid!',
      });
    }

    if (await bcrypt.compare(password, user.password)) {
      done(null, user);
    } else {
      done(null, false, {message: 'Invalid username or password!'});
    }
  } catch (error) {
    done(error);
  }
}
