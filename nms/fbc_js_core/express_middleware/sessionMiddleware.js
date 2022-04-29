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

import session from 'express-session';

type SessionMiddlewareOptions = {|
  devMode: boolean,
  // $FlowIgnore[value-as-type]
  sessionStore: session.Session,
  sessionToken: string,
|};

export default function sessionMiddleware(
  options: SessionMiddlewareOptions,
  // $FlowIgnore[value-as-type]
): Middleware {
  options.sessionStore.sync();
  return session({
    cookie: {
      secure: !options.devMode,
    },
    // Used to sign the session cookie
    secret: options.sessionToken,
    resave: false,
    saveUninitialized: true,
    store: options.sessionStore,
    unset: 'destroy',
  });
}
