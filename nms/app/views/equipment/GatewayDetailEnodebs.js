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
// $FlowFixMe migrated to typescript
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {GatewayDetailType} from './GatewayDetailMain';

import ActionTable from '../../components/ActionTable';
import EnodebContext from '../../components/context/EnodebContext';
import Link from '@material-ui/core/Link';
import React from 'react';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
// $FlowFixMe migrated to typescript
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {useNavigate, useParams, useResolvedPath} from 'react-router-dom';
import {useState} from 'react';

type EnodebRowType = {
  name: string,
  id: string,
  health: string,
  mmeConnected: string,
};

export default function GatewayDetailEnodebs(props: GatewayDetailType) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  // Auto refresh  every 30 seconds
  const enbState = useRefreshingContext({
    context: EnodebContext,
    networkId: networkId,
    type: 'enodeb',
    interval: REFRESH_INTERVAL,
    refresh: props.refresh || false,
  });
  const enbInfo =
    props.gwInfo.connected_enodeb_serials?.reduce(
      (enbs: {[string]: EnodebInfo}, serial: string) => {
        // $FlowIgnore
        if (enbState?.enbInfo?.[serial] != null) {
          // $FlowIgnore
          enbs[serial] = enbState.enbInfo?.[serial];
        }
        return enbs;
      },
      {},
    ) || {};
  const [currRow, setCurrRow] = useState<EnodebRowType>({});

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
              }}>
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
