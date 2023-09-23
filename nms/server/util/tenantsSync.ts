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

import OrchestratorAPI from '../api/OrchestratorAPI';
import axios from 'axios';
import {OrganizationModel} from '../../shared/sequelize_models/models/organization';

export async function syncOrganizationWithOrc8rTenant(
  organization: OrganizationModel,
): Promise<void> {
  let orc8rTenant;
  try {
    orc8rTenant = (
      await OrchestratorAPI.tenants.tenantsTenantIdGet({
        tenantId: organization.id,
      })
    ).data;
  } catch (error) {
    // Ignore "not found" since there is no guarantee NMS and Orc8r are in sync
    rethrowUnlessNotFoundError(error);
  }

  if (orc8rTenant) {
    await OrchestratorAPI.tenants.tenantsTenantIdPut({
      tenant: {
        id: organization.id,
        name: organization.name,
        networks: organization.networkIDs,
      },
      tenantId: organization.id,
    });
  } else {
    await OrchestratorAPI.tenants.tenantsPost({
      tenant: {
        id: organization.id,
        name: organization.name,
        networks: organization.networkIDs,
      },
    });
  }
}

export function rethrowUnlessNotFoundError(error: unknown) {
  if (!(axios.isAxiosError(error) && error?.response?.status === 404)) {
    throw error;
  }
}
