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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {DataRows} from '../../components/DataGrid';
import type {federation_gateway} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
import React from 'react';

/**
 * Returns the federation gateway description, id, hardware uuid, and it's
 * version.
 * @param {federation_gateway} gwInfo The federation gateway being looked at.
 */
export default function GatewaySummary({gwInfo}: {gwInfo: federation_gateway}) {
  const version = gwInfo.status?.platform_info?.packages?.[0]?.version;

  const data: DataRows[] = [
    [
      {
        category: 'Name',
        value: gwInfo.name,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: gwInfo.id,
      },
    ],
    [
      {
        category: 'Hardware UUID',
        value: gwInfo.device?.hardware_id || '-',
      },
    ],
    [
      {
        category: 'Version',
        value: version ?? 'Unknown',
      },
    ],
  ];

  return <DataGrid data={data} />;
}
