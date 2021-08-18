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
 * @flow strict-local
 * @format
 */

import type {$AxiosError} from 'axios';

import React from 'react';

import {useSnackbar} from '@fbcnms/ui/hooks';

export default function NetworkError({error}: {error: $AxiosError<string>}) {
  let errorMessage = error.message;
  if (error.response && error.response.status >= 400) {
    errorMessage = error.response?.statusText;
  }
  useSnackbar(
    'Unable to communicate with magma controller: ' + errorMessage,
    {variant: 'error'},
    !!error,
  );
  return <div />;
}
