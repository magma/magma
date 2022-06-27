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

import type {ExpressResponse, NextFunction} from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSRequest} from '../auth/access';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import asyncHandler from '../util/asyncHandler';

export const hostOrgMiddleware = asyncHandler(
  async (req: FBCNMSRequest, res: ExpressResponse, next: NextFunction) => {
    if (req.organization) {
      const organization = await req.organization();
      if (organization.isHostOrg) {
        next();
        return;
      }
    }

    return res.redirect(req.access?.loginUrl || '/user/login');
  },
);
