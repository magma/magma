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

import axios, {AxiosError, AxiosResponse} from 'axios';

export function isAxiosErrorResponse<T = {message?: string}>(
  error: any,
): error is AxiosError<T> & {response: AxiosResponse<T>} {
  return axios.isAxiosError(error) && !!error.response;
}

export function getErrorMessage(
  error: unknown,
  fallbackMessage = 'Unknown Error',
): string {
  let errorMessage;
  if (isAxiosErrorResponse(error)) {
    errorMessage = error.response?.data?.message ?? error.message;
  } else if (error instanceof Error) {
    errorMessage = error.message;
  }
  return errorMessage || fallbackMessage;
}
