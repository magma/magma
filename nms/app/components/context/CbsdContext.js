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
 *
 * @flow strict-local
 * @format
 */
import React from 'react';

import type {cbsd, mutable_cbsd} from '../../../generated/MagmaAPIBindings';

export type CbsdContextType = {
  state: {
    isLoading: boolean,
    cbsds: cbsd[],
    totalCount: number,
    page: number,
    pageSize: number,
  },
  setPaginationOptions: ({page: number, pageSize: number}) => void,
  refetch: () => Promise<void>,
  create: (newCbsd: mutable_cbsd) => Promise<void>,
  update: (id: number, cbsd: mutable_cbsd) => Promise<void>,
  deregister: (id: number) => Promise<void>,
  remove: (id: number) => Promise<void>,
};

export default React.createContext<CbsdContextType>({});
