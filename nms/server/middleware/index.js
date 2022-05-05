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

export {default as organizationMiddleware} from './organizationMiddleware';
export {default as appMiddleware} from './appMiddleware';
export {default as csrfMiddleware} from './csrfMiddleware';
// $FlowFixMe migrated to typescript
export {default as sessionMiddleware} from './sessionMiddleware';
export {default as webpackSmartMiddleware} from './webpackSmartMiddleware';

import type {OrganizationMiddlewareRequest} from './organizationMiddleware';

export type FBCNMSMiddleWareRequest = {
  csrfToken: () => string, // from csrf
  body: Object, // from bodyParser
  session: any,
} & OrganizationMiddlewareRequest;
