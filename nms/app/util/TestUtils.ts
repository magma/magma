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

import MagmaAPI from '../../api/MagmaAPI';
import {AxiosResponse} from 'axios';

type Response<
  SERVICE extends object,
  API_METHOD extends keyof SERVICE
> = SERVICE[API_METHOD] extends (
  ...args: Array<any>
) => Promise<AxiosResponse<infer DATA>>
  ? DATA
  : never;

export function mockAPI<
  SERVICE extends typeof MagmaAPI[keyof typeof MagmaAPI],
  API_METHOD extends jest.FunctionPropertyNames<SERVICE>
>(service: SERVICE, apiMethod: API_METHOD): jest.SpyInstance;
export function mockAPI<
  SERVICE extends typeof MagmaAPI[keyof typeof MagmaAPI],
  API_METHOD extends jest.FunctionPropertyNames<SERVICE>
>(
  service: SERVICE,
  apiMethod: API_METHOD,
  data: Response<SERVICE, API_METHOD>,
): jest.SpyInstance<AxiosResponse<{data: Response<SERVICE, API_METHOD>}>>;
export function mockAPI(service: any, apiMethod: any, data?: any) {
  return jest.spyOn(service, apiMethod).mockResolvedValue({
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    data,
  } as any);
}

export function mockAPIOnce<
  SERVICE extends typeof MagmaAPI[keyof typeof MagmaAPI],
  API_METHOD extends jest.FunctionPropertyNames<SERVICE>
>(
  service: SERVICE,
  apiMethod: API_METHOD,
  data: Response<SERVICE, API_METHOD>,
): jest.SpyInstance<AxiosResponse<{data: Response<SERVICE, API_METHOD>}>> {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
  return jest.spyOn(service, apiMethod).mockResolvedValueOnce({
    data,
  } as any);
}
