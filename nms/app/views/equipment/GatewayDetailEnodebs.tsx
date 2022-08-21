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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {GatewayDetailType} from './GatewayDetailMain';

import ActionTable from '../../components/ActionTable';
import EnodebContext from '../../context/EnodebContext';
import Link from '@mui/material/Link';
import React, {useContext} from 'react';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useInterval} from '../../hooks';
import {useNavigate, useResolvedPath} from 'react-router-dom';
import {useState} from 'react';

type EnodebRowType = {
  name: string;
  id: string;
  health: string;
  mmeConnected: string;
};

export default function GatewayDetailEnodebs(props: GatewayDetailType) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const enodebContext = useContext(EnodebContext);

  useInterval(
    () => enodebContext.refetch(),
    props.refresh ? REFRESH_INTERVAL : null,
  );

  const enbInfo =
    props.gwInfo.connected_enodeb_serials?.reduce(
      (enbs: Record<string, EnodebInfo>, serial: string) => {
        if (enodebContext.state.enbInfo[serial] != null) {
          enbs[serial] = enodebContext.state.enbInfo[serial];
        }
        return enbs;
      },
      {},
    ) || {};
  const [currRow, setCurrRow] = useState<EnodebRowType>({} as EnodebRowType);

  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      const isEnbManaged = enbInf.enb?.enodeb_config?.config_type === 'MANAGED';
      return {
        health: isEnbManaged ? (isEnodebHealthy(enbInf) ? 'Good' : 'Bad') : '-',
        mmeConnected: enbInf.enb_state?.mme_connected
          ? 'Connected'
          : 'Disconnected',
        ipAddress: enbInf.enb_state.ip_address ?? '-',
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
                navigate(
                  resolvedPath.pathname.replace(
                    `gateway/${props.gwInfo.id}`,
                    `enodeb/${currRow.id}`,
                  ),
                );
              }}
              underline="hover">
              {currRow.id}
            </Link>
          ),
        },
        {title: 'Health', field: 'health'},
        {title: 'MME', field: 'mmeConnected'},
        {title: 'IP Address', field: 'ipAddress'},
      ]}
      handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
      menuItems={[
        {
          name: 'View',
          handleFunc: () => {
            navigate(
              resolvedPath.pathname.replace(
                `gateway/${props.gwInfo.id}`,
                `enodeb/${currRow.id}`,
              ),
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
