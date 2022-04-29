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

import Sequelize from 'sequelize';

import {
  Organization,
  jsonArrayContains,
} from '../../fbc_js_core/sequelize_models';
import type {ExpressRequest, ExpressResponse, NextFunction} from 'express';
import type {OrganizationType} from '../../fbc_js_core/sequelize_models/models/organization';

async function getOrganizationFromHost(
  host: string,
): Promise<?OrganizationType> {
  const org = await Organization.findOne({
    where: jsonArrayContains('customDomains', host),
  });
  if (org) {
    return org;
  }
  const subdomain = host.split('.')[0];
  return await Organization.findOne({
    where: {
      name: Sequelize.where(
        Sequelize.fn('lower', Sequelize.col('name')),
        Sequelize.fn('lower', subdomain),
      ),
    },
  });
}

export async function getOrganization(req: {
  get(field: string): string | void,
}): Promise<OrganizationType> {
  for (const header of ['host', 'x-forwarded-host']) {
    const host = req.get(header);
    if (host != null && host !== '') {
      const org = await getOrganizationFromHost(host);
      if (org) {
        return org;
      }
    }
  }
  throw new Error('Invalid organization!');
}

// We don't depend on organization to be there in other modules.
export type OrganizationRequestPart = {
  organization: () => Promise<OrganizationType>,
};
export type OrganizationMiddlewareRequest = ExpressRequest &
  $Shape<OrganizationRequestPart>;

export default function organizationMiddleware() {
  return (
    req: OrganizationMiddlewareRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) => {
    try {
      req.organization = () => getOrganization(req);
      next();
    } catch (err) {
      res.status(404).send();
      next(err);
    }
  };
}
