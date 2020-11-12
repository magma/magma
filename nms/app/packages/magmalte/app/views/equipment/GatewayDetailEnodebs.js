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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import EnodebContext from '../../components/context/EnodebContext';
import Link from '@material-ui/core/Link';
import React from 'react';

import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useContext, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

type EnodebRowType = {
  name: string,
  id: string,
  health: string,
};

export default function GatewayDetailEnodebs({gwInfo}: {gwInfo: lte_gateway}) {
  const {history, match} = useRouter();
  const enbCtx = useContext(EnodebContext);
  const enbInfo =
    gwInfo.connected_enodeb_serials?.reduce(
      (enbs: {[string]: EnodebInfo}, serial: string) => {
        if (enbCtx.state.enbInfo[serial] != null) {
          enbs[serial] = enbCtx.state.enbInfo[serial];
        }
        return enbs;
      },
      {},
    ) || {};
  const [currRow, setCurrRow] = useState<EnodebRowType>({});

  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        health: isEnodebHealthy(enbInf) ? 'Good' : 'Bad',
        ipAddress: enbInf.enb_state.ip_address ?? 'Not Available',
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
        {
          title: 'Serial Number',
          field: 'id',
          render: currRow => (
            <Link
              variant="body2"
              component="button"
              onClick={() => {
                history.push(
                  match.url.replace(
                    `gateway/${gwInfo.id}`,
                    `enodeb/${currRow.id}`,
                  ),
                );
              }}>
              {currRow.id}
            </Link>
          ),
        },
        {title: 'Health', field: 'health'},
        {title: 'IP Address', field: 'ipAddress'},
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
