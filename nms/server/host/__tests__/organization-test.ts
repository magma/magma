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

import appMiddleware from '../../middleware/appMiddleware';
import express from 'express';
import request from 'supertest';
import router from '../routes';
import {AccessRoles} from '../../../shared/roles';
import {Organization, User} from '../../../shared/sequelize_models';

describe('Organization routes', () => {
  const app = express().use(appMiddleware()).use('', router);

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

  it('a new organization can be created', async () => {
    const orgName = 'new-org';

    let newOrganization = await Organization.findOne({where: {name: orgName}});
    expect(newOrganization).toBeNull();

    await request(app)
      .post('/organization/async')
      .send({name: orgName, networkIDs: ['test'], customDomains: []})
      .expect(200);

    newOrganization = await Organization.findOne({where: {name: orgName}});

    expect(newOrganization).not.toBeNull();
    expect(newOrganization!.name).toBe(orgName);
    expect(newOrganization!.networkIDs).toEqual(['test']);
  });

  it('an organization can updated', async () => {
    const organization = await Organization.create({name: 'test'});
    expect(organization.networkIDs).toEqual([]);

    await request(app)
      .put(`/organization/async/${organization.name}`)
      .send({name: 'test', networkIDs: ['aaa', 'bbb', 'ccc']})
      .expect(200);

    const updatedOrganization = await Organization.findOne({
      where: {id: organization.id},
    });

    expect(updatedOrganization!.networkIDs).toEqual(['aaa', 'bbb', 'ccc']);
  });

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
