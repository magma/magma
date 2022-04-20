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

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../../../fbc_js_core/auth/access';

import asyncHandler from '../../../fbc_js_core/util/asyncHandler';
import express from 'express';

import {AuditLogEntry} from '../../../fbc_js_core/sequelize_models';

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();
router.get(
  '/auditlog/async',
  asyncHandler(async (req: FBCNMSRequest, res: ExpressResponse) => {
    const organization = await req.organization();
    const data = await AuditLogEntry.findAll({
      where: {organization: organization.name},
      limit: 20,
    });
    res.status(200).send(data);
  }),
);

export default router;
