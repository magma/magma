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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
import React from 'react';

type OverviewProps = {
  name?: string,
  networkIds?: Array<string>,
};

/**
 * Organization basic information
 */
export default function OrganizationSummary(props: OverviewProps) {
  const {name, networkIds} = props;
  const orgName = window.CONFIG.appData.user.tenant;

  const kpiData: DataRows[] = [
    [
      {
        category: 'Organization Name',
        value: name || '',
      },
    ],
    [
      {
        category: 'Accessible Networks',
        value: [...(networkIds || [])].join(', ') || '-',
      },
    ],
    [
      {
        category: 'Link to Organization Portal',
        isLink: true,
        value: `${window.location.origin.replace(orgName, name)}/nms`,
      },
    ],
  ];
  return <DataGrid data={kpiData} />;
}
