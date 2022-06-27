/*
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

import type {ExpressResponse} from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSRequest} from '../auth/access';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import asyncHandler from '../util/asyncHandler';
import express from 'express';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AuditLogEntry, User} from '../../shared/sequelize_models';

const MAX_AUDITLOG_ROWS = 5000;
const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();
router.get(
  '/auditlog/async',
  asyncHandler(async (req: FBCNMSRequest, res: ExpressResponse) => {
    const organization = await req.organization();
    //$FlowFixMe[incompatible-call]
    const data = await AuditLogEntry.findAll({
      where: {organization: organization.name},
      limit: MAX_AUDITLOG_ROWS,
    });

    // cleaner way would be to define association.
    // will do that post db migration release
    const userMap = {};
    const allUsers = await User.findAll();
    allUsers.forEach(item => {
      userMap[item.id] = item.email;
    });
    const userLogData = data.map(item => ({
      id: item.id,
      item: item,
      status: item.status,
      objectId: item.objectId,
      objectType: item.objectType,
      mutationType: item.mutationType,
      mutationData: item.mutationData,
      actingUserId: item.actingUserId,
      actingUserEmail: userMap[item.actingUserId] ?? 'undefined',
      createdAt: item.createdAt,
      updatedAt: item.updatedAt,
      url: item.url,
      ipAddress: item.ipAddress,
    }));
    res.status(200).send(userLogData);
  }),
);

export default router;
