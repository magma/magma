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

import useAxios, {ExtendedConfig, Result} from '../../hooks/useAxios';
import {AxiosResponse} from 'axios';

export function mockUseAxios(responses: {[url: string]: {data: any}}) {
  (useAxios as jest.Mock).mockImplementation(
    (config: ExtendedConfig<any>): Result<any> => {
      const response = responses[config.url!];
      if ('onResponse' in config) {
        if (response) {
          config.onResponse!(response as AxiosResponse<any>);
          delete responses[config.url!];
        }
        return {} as Result<any>; // Either onResponse or the response result should be used but not both
      } else {
        return {isLoading: false, response} as Result<any>;
      }
    },
  );
}
