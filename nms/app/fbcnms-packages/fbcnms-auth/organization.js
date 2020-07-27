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
 * @flow
 * @format
 */

import type {FBCNMSMiddleWareRequest} from '@fbcnms/express-middleware';

export async function injectOrganizationParams<T: {[string]: any}>(
  req: FBCNMSMiddleWareRequest,
  params: T,
): Promise<T & {organization?: string}> {
  if (req.organization) {
    const organization = await req.organization();
    return {
      ...params,
      organization: organization.name,
    };
  }
  return params;
}
