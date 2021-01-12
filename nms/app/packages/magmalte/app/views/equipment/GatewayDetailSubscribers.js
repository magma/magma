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
import type {lte_gateway} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import Link from '@material-ui/core/Link';
import React from 'react';
import SubscriberContext from '../../components/context/SubscriberContext';

import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

type SubscriberRowType = {
  id: string,
  service: string,
};

export default function GatewayDetailSubscribers({
  gwInfo,
}: {
  gwInfo: lte_gateway,
}) {
  const {history, match} = useRouter();
  const subscriberCtx = useContext(SubscriberContext);
  const gwSubscriberMap =
    subscriberCtx.gwSubscriberMap[gwInfo?.device?.hardware_id] || [];

  const subscriberRows: Array<SubscriberRowType> = gwSubscriberMap.map(
    (serialNum: string) => {
      const subscriberInfo = subscriberCtx.state[serialNum];
      return {
        name: subscriberInfo?.name || serialNum,
        id: serialNum,
        service: subscriberInfo.lte.state,
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
                history.push(
                  match.url.replace(
                    `equipment/overview/gateway/${gwInfo.id}`,
                    `subscribers/overview/${currRow.id}`,
                  ),
                );
              }}>
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
