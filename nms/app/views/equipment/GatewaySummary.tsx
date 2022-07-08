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
 */
import type {DataRows} from '../../components/DataGrid';
import type {LteGateway} from '../../../generated-ts';

import DataGrid from '../../components/DataGrid';
import React from 'react';

export default function GatewaySummary({gwInfo}: {gwInfo: LteGateway}) {
  const version = gwInfo.status?.platform_info?.packages?.[0]?.version;

  const data: Array<DataRows> = [
    [
      {
        value: gwInfo.description,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: gwInfo.id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Hardware UUID',
        value: gwInfo.device?.hardware_id || '-',
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Version',
        value: version ?? 'null',
        statusCircle: false,
      },
    ],
  ];

  return <DataGrid data={data} />;
}
