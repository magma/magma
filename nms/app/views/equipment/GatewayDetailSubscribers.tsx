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
import ActionTable from '../../components/ActionTable';
import Link from '@mui/material/Link';
import React from 'react';
import SubscriberContext from '../../context/SubscriberContext';

import {REFRESH_INTERVAL} from '../../context/AppContext';
import {useContext} from 'react';
import {useInterval} from '../../hooks';
import {useNavigate, useResolvedPath} from 'react-router-dom';
import type {GatewayDetailType} from './GatewayDetailMain';

type SubscriberRowType = {
  id: string;
  service: string;
};

export default function GatewayDetailSubscribers(props: GatewayDetailType) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const subscriberCtx = useContext(SubscriberContext);
  useInterval(
    () => subscriberCtx.refetchSessionState(),
    props.refresh ? REFRESH_INTERVAL : null,
  );
  const hardware_id = props.gwInfo?.device?.hardware_id;
  const gwSubscriberMap = hardware_id
    ? subscriberCtx.gwSubscriberMap[hardware_id] || []
    : [];

  const subscriberRows: Array<SubscriberRowType> = gwSubscriberMap.map(
    (serialNum: string) => {
      const subscriberInfo = subscriberCtx.state?.[serialNum];
      return {
        name: subscriberInfo?.name || serialNum,
        id: serialNum,
        service: subscriberInfo?.lte.state || '-',
      };
    },
  );

  return (
    <ActionTable
      title=""
      data={subscriberRows}
      columns={[
        {title: 'Name', field: 'name'},
        {
          title: 'Subscriber ID',
          field: 'id',
          render: currRow => (
            <Link
              variant="body2"
              component="button"
              onClick={() => {
                navigate(
                  resolvedPath.pathname.replace(
                    `equipment/overview/gateway/${props.gwInfo.id}`,
                    `subscribers/overview/${currRow.id}`,
                  ),
                );
              }}
              underline="hover">
              {currRow.id}
            </Link>
          ),
        },
        {title: 'Service', field: 'service'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [10],
        toolbar: false,
      }}
    />
  );
}
