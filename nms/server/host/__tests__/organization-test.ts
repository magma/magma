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

import OrchestratorAPI from '../../api/OrchestratorAPI';
import appMiddleware from '../../middleware/appMiddleware';
import express from 'express';
import request from 'supertest';
import router from '../routes';
import {AccessRoles} from '../../../shared/roles';
import {Organization, User} from '../../../shared/sequelize_models';

jest.mock('../../api/OrchestratorAPI', () => ({
  tenants: {
    tenantsTenantIdDelete: jest.fn(),
    tenantsTenantIdGet: jest.fn(),
    tenantsTenantIdPut: jest.fn(),
    tenantsPost: jest.fn(),
  },
}));

const newNetworks = ['aaa', 'bbb', 'ccc'];

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

  it('a new organization can be created in NMS and orc8r', async () => {
    const orgName = 'new-org';

    let newOrganization = await Organization.findOne({where: {name: orgName}});
    expect(newOrganization).toBeNull();

    // Tenant does not exist in orc8r yet
    const mockedTenantIdGet = OrchestratorAPI.tenants
      .tenantsTenantIdGet as jest.Mock<any>;
    mockedTenantIdGet.mockImplementation(() =>
      Promise.reject({isAxiosError: true, response: {status: 404}}),
    );

    await request(app)
      .post('/organization/async')
      .send({name: orgName, networkIDs: ['test'], customDomains: []})
      .expect(200);

    newOrganization = await Organization.findOne({where: {name: orgName}});

    expect(newOrganization).not.toBeNull();
    expect(newOrganization!.name).toBe(orgName);
    expect(newOrganization!.networkIDs).toEqual(['test']);
    expect(OrchestratorAPI.tenants.tenantsPost).toBeCalledWith({
      tenant: {
        id: newOrganization?.id,
        name: newOrganization?.name,
        networks: newOrganization?.networkIDs,
      },
    });
  });

  it('an existing organization can updated in NMS and orc8r', async () => {
    const organization = await Organization.create({name: 'test'});
    expect(organization.networkIDs).toEqual([]);

    // Tenant already exists in orc8r
    const mockedTenantIdGet = OrchestratorAPI.tenants
      .tenantsTenantIdGet as jest.Mock<any>;
    mockedTenantIdGet.mockImplementation(() =>
      Promise.resolve({data: {id: organization.id}}),
    );

    await request(app)
      .put(`/organization/async/${organization.name}`)
      .send({name: 'test', networkIDs: newNetworks})
      .expect(200);

    const updatedOrganization = await Organization.findOne({
      where: {id: organization.id},
    });

    expect(updatedOrganization!.networkIDs).toEqual(newNetworks);
    expect(OrchestratorAPI.tenants.tenantsTenantIdPut).toBeCalledWith({
      tenant: {
        id: organization.id,
        name: organization.name,
        networks: newNetworks,
      },
      tenantId: organization.id,
    });
  });

  it('updating an organization in NMS creates it in orc8r if it does not exist', async () => {
    const organization = await Organization.create({name: 'test'});
    expect(organization.networkIDs).toEqual([]);

    // Tenant does not exist in orc8r yet
    const mockedTenantIdGet = OrchestratorAPI.tenants
      .tenantsTenantIdGet as jest.Mock<any>;
    mockedTenantIdGet.mockImplementation(() =>
      Promise.reject({isAxiosError: true, response: {status: 404}}),
    );

    await request(app)
      .put(`/organization/async/${organization.name}`)
      .send({name: 'test', networkIDs: newNetworks})
      .expect(200);

    const updatedOrganization = await Organization.findOne({
      where: {id: organization.id},
    });

    expect(updatedOrganization!.networkIDs).toEqual(newNetworks);
    expect(OrchestratorAPI.tenants.tenantsPost).toBeCalledWith({
      tenant: {
        id: organization?.id,
        name: organization?.name,
        networks: newNetworks,
      },
    });
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
    expect(OrchestratorAPI.tenants.tenantsTenantIdDelete).toBeCalledWith({
      tenantId: organization.id,
    });
  });

  it('deleting an unknown organization forwards the deletion to orc8r', async () => {
    await request(app).delete(`/organization/async/4711`).expect(200);
    expect(OrchestratorAPI.tenants.tenantsTenantIdDelete).toBeCalledWith({
      tenantId: 4711,
    });
  });
});
