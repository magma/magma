/**
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

import express from 'express';
import request from 'supertest';
import router from '../routes';
import {AccessRoles} from '../../../shared/roles';
import {Organization, User} from '../../../shared/sequelize_models';

describe('Organization routes', () => {
  const app = express().use('', router);

  beforeEach(async () => {
    await Organization.sync();
    await User.sync();
  });

  afterEach(async () => {
    await Organization.drop();
    await User.drop();
  });

  function createUser(email: string, organizationName: string) {
    return User.create({
      email,
      role: AccessRoles.USER,
      organization: organizationName,
      password: '1234',
    });
  }

  it('deleting an organization also deletes its users', async () => {
    const organization = await Organization.create({name: 'test'});
    const otherOrganization = await Organization.create({name: 'other'});

    await createUser('user1@magma.test', organization.name);
    await createUser('user2@magma.test', organization.name);
    const user3 = await createUser('user3@magma.test', otherOrganization.name);

    await request(app)
      .delete(`/organization/async/${organization.id}`)
      .expect(200);

    const organizations = await Organization.findAll();
    expect(organizations.map(org => org.id)).toEqual([otherOrganization.id]);
    const users = await User.findAll();
    expect(users.map(user => user.email)).toEqual([user3.email]);
  });

  it('deleting an unknown organization is a noop', async () => {
    await request(app).delete(`/organization/async/4711`).expect(200);
  });
});
