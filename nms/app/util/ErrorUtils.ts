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

import {AxiosError, AxiosResponse} from 'axios';

function isAxiosError<T>(error: any): error is AxiosError<T> {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-return,@typescript-eslint/no-unsafe-member-access
  return 'isAxiosError' in error && error.isAxiosError;
}

export function isAxiosErrorResponse<T = {message?: string}>(
  error: any,
): error is AxiosError<T> & {response: AxiosResponse<T>} {
  return isAxiosError(error) && !!error.response;
}

export function getErrorMessage(
  error: unknown,
  fallbackMessage = 'Unknown Error',
): string {
  let errorMessage;
  if (isAxiosError<{message?: string}>(error)) {
    errorMessage = error.response?.data?.message ?? error.message;
  }
  if (error instanceof Error) {
    errorMessage = error.message;
  }
  return errorMessage || fallbackMessage;
}
