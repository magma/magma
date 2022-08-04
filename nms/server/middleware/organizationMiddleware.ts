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
 */

import Sequelize from 'sequelize';

import {Organization, jsonArrayContains} from '../../shared/sequelize_models';
import type {NextFunction, Request, Response} from 'express';
import type {OrganizationModel} from '../../shared/sequelize_models/models/organization';

async function getOrganizationFromHost(
  host: string,
): Promise<OrganizationModel | null | undefined> {
  const org = await Organization.findOne({
    where: jsonArrayContains('customDomains', host),
  });
  if (org) {
    return org;
  }
  const subdomain = host.split('.')[0];
  return await Organization.findOne({
    where: Sequelize.where(
      Sequelize.fn('lower', Sequelize.col('name')),
      Sequelize.fn('lower', subdomain),
    ),
  });
}

export async function getOrganization(
  req: Request,
): Promise<OrganizationModel> {
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

export default function organizationMiddleware() {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      req.organization = () => getOrganization(req);
      next();
    } catch (err) {
      res.status(404).send();
      next(err);
    }
  };
}
