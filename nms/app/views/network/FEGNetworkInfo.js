/*
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
import type {feg_network} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
import React from 'react';

type Props = {
  fegNetwork: $Shape<feg_network>,
};

/**
 * Returns information about the federation network.
 * @param {object} props: has a property called fegNetwork that has
 * information about the federation network.
 */
export default function NetworkInfo(props: Props) {
  const kpiData: DataRows[] = [
    [
      {
        category: 'ID',
        value: props.fegNetwork.id,
      },
    ],
    [
      {
        category: 'Name',
        value: props.fegNetwork.name,
      },
    ],
    [
      {
        category: 'Description',
        value: props.fegNetwork.description || '-',
      },
    ],
    [
      {
        category: 'Served Federated LTE Network IDs',
        value: props.fegNetwork?.federation?.served_network_ids?.join() || '-',
        tooltip:
          'List of Federated LTE Network IDs serviced under this federation network',
      },
    ],
    [
      {
        category: 'Served Virtual Federation Network IDs',
        value: props.fegNetwork?.federation?.served_nh_ids?.join() || '-',
        tooltip:
          'List of Neutral Host (or Virtual) Federation Networks IDs serviced under this federation network',
      },
    ],
  ];
  return <DataGrid data={kpiData} testID="feg_info" />;
}
