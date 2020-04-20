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
import ActionTable from '../../components/ActionTable';
import React from 'react';

import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '@fbcnms/magma-api';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

type EnodebRowType = {
  name: string,
  id: string,
  health: string,
};

export default function GatewayDetailEnodebs({
  gwInfo,
  enbInfo,
}: {
  gwInfo: lte_gateway,
  enbInfo: {[string]: EnodebInfo},
}) {
  const {history, match} = useRouter();
  const [currRow, setCurrRow] = useState<EnodebRowType>({});

  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        health: isEnodebHealthy(enbInf) ? 'Good' : 'Bad',
        name: enbInf.enb.name,
        id: serialNum,
      };
    },
  );

  return (
    <ActionTable
      title=""
      data={enbRows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Serial Number', field: 'id'},
        {title: 'Health', field: 'health'},
      ]}
      handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
      menuItems={[
        {
          name: 'View',
          handleFunc: () => {
            history.push(
              match.url.replace(`gateway/${gwInfo.id}`, `enodeb/${currRow.id}`),
            );
          },
        },
        {name: 'Edit'},
        {name: 'Remove'},
        {name: 'Deactivate'},
        {name: 'Reboot'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5],
        toolbar: false,
      }}
    />
  );
}
