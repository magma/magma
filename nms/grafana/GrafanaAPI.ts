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

import axios, {AxiosError, AxiosRequestConfig, AxiosResponse} from 'axios';

import type {
  AddOrgUserResponse,
  CreateDashboardResponse,
  CreateDatasourceResponse,
  CreateOrgResponse,
  CreateUserResponse,
  Dashboard,
  DeleteOrgResponse,
  GetDatasourcesResponse,
  GetHealthResponse,
  GetOrgUsersResponse,
  GetUserResponse,
  OrgUser,
  Organization,
  PostDatasource,
  StarDashboardResponse,
  User,
} from './GrafanaAPIType';

type GrafanaPromise<T> = Promise<GrafanaResponse<T>>;

export type GrafanaResponse<T> = {
  status: number;
  data: T;
};

export type GrafanaClient = {
  getUser: (loginOrEmail: string) => GrafanaPromise<GetUserResponse>;
  createUser: (user: User) => GrafanaPromise<CreateUserResponse>;

  getOrg: (orgName: string) => GrafanaPromise<Organization>;
  addOrg: (orgName: string) => GrafanaPromise<CreateOrgResponse>;
  deleteOrg: (orgID: number) => GrafanaPromise<DeleteOrgResponse>;
  addUserToOrg: (
    orgID: number,
    user: OrgUser,
  ) => GrafanaPromise<AddOrgUserResponse>;
  getUsersInOrg: (orgID: number) => GrafanaPromise<GetOrgUsersResponse>;

  createDatasource: (
    ds: PostDatasource,
    orgID: number,
  ) => GrafanaPromise<CreateDatasourceResponse>;
  updateDatasource: (
    dsID: number,
    orgID: number,
    ds: PostDatasource,
  ) => GrafanaPromise<CreateDatasourceResponse>;
  getDatasources: (orgID: number) => GrafanaPromise<GetDatasourcesResponse>;

  createDashboard: (
    db: Dashboard,
    orgID: number,
  ) => GrafanaPromise<CreateDashboardResponse>;

  starDashboard: (
    dbID: number,
    orgID: number,
    username: string,
  ) => GrafanaPromise<StarDashboardResponse>;

  getHealth: () => GrafanaPromise<GetHealthResponse>;
};

async function request<Data>(req: AxiosRequestConfig): GrafanaPromise<Data> {
  try {
    const res = (await axios(req)) as AxiosResponse<Data>;
    return {status: res.status, data: res.data};
  } catch (error) {
    return {
      status: (error as AxiosError).response?.status || 0,
      data:
        (error as AxiosError<Data>).response?.data ||
        (({
          message: 'unknown error',
        } as unknown) as Data),
    };
  }
}

const client = (
  apiURL: string,
  constHeaders: {[key: string]: string},
): GrafanaClient => ({
  async getUser(loginOrEmail: string): GrafanaPromise<GetUserResponse> {
    return await request({
      url: apiURL + `/api/users/lookup`,
      params: {loginOrEmail: loginOrEmail},
      method: 'GET',
      headers: constHeaders,
    });
  },

  async createUser(user: User): GrafanaPromise<CreateUserResponse> {
    return await request({
      url: apiURL + `/api/admin/users`,
      method: 'POST',
      data: user,
      headers: constHeaders,
    });
  },

  async getOrg(orgName: string): GrafanaPromise<Organization> {
    return await request({
      url: apiURL + `/api/orgs/name/${orgName}`,
      method: 'GET',
      headers: constHeaders,
    });
  },

  async addOrg(orgName: string): GrafanaPromise<CreateOrgResponse> {
    return await request({
      url: apiURL + '/api/orgs',
      method: 'POST',
      data: {name: orgName},
      headers: constHeaders,
    });
  },

  async deleteOrg(orgID: number): GrafanaPromise<DeleteOrgResponse> {
    return await request({
      url: apiURL + `/api/orgs/${orgID}`,
      method: 'DELETE',
      headers: constHeaders,
    });
  },

  async addUserToOrg(
    orgID: number,
    user: OrgUser,
  ): GrafanaPromise<AddOrgUserResponse> {
    return await request({
      url: apiURL + `/api/orgs/${orgID}/users`,
      method: 'POST',
      data: user,
      headers: constHeaders,
    });
  },

  async getUsersInOrg(orgID: number): GrafanaPromise<GetOrgUsersResponse> {
    return await request({
      url: apiURL + `/api/orgs/${orgID}/users`,
      method: 'GET',
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async createDatasource(
    ds: PostDatasource,
    orgId: number,
  ): GrafanaPromise<CreateDatasourceResponse> {
    return await request({
      url: apiURL + `/api/datasources`,
      method: 'POST',
      data: ds,
      headers: {
        ...constHeaders,
        'X-Grafana-Org-Id': orgId.toString(),
        'Content-Type': 'application/json',
      },
    });
  },

  async updateDatasource(dsID: number, orgID: number, ds: PostDatasource) {
    return await request({
      url: apiURL + `/api/datasources/${dsID}`,
      method: 'PUT',
      data: ds,
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async getDatasources(orgID: number): GrafanaPromise<GetDatasourcesResponse> {
    return await request({
      url: apiURL + `/api/datasources`,
      method: 'GET',
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async createDashboard(
    db: Dashboard,
    orgID: number,
  ): GrafanaPromise<CreateDashboardResponse> {
    return await request({
      url: apiURL + `/api/dashboards/db/`,
      method: 'POST',
      data: db,
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async starDashboard(
    dbID: number,
    orgID: number,
    username: string,
  ): GrafanaPromise<StarDashboardResponse> {
    return await request({
      url: apiURL + `/api/user/stars/dashboard/${dbID}`,
      method: 'POST',
      headers: {
        'X-WEBAUTH-USER': username,
        'X-Grafana-Org-Id': orgID.toString(),
      },
    });
  },

  async getHealth(): GrafanaPromise<GetHealthResponse> {
    return await request({
      url: apiURL + `/api/health`,
      method: 'GET',
      headers: constHeaders,
    });
  },
});

export default client;
