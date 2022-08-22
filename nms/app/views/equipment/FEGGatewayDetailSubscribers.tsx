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
import FEGSubscriberContext from '../../context/FEGSubscriberContext';
import Link from '@mui/material/Link';
import LoadingFiller from '../../components/LoadingFiller';
import React from 'react';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {fetchSubscribers} from '../../util/SubscriberState';
import {useContext, useEffect, useState} from 'react';
import {useInterval} from '../../hooks';
import {useNavigate, useResolvedPath} from 'react-router-dom';
import type {FederationGateway, Subscriber} from '../../../generated';

/**
 * @property {FederationGateway} gwInfo The Federation gateway being looked at
 * @property {boolean} refresh Boolean telling to autorefresh or not
 */
type FEGGatewayDetailType = {
  gwInfo: FederationGateway;
  refresh: boolean;
};

/**
 * @property {string} name Subscriber name
 * @property {string} id Subscriber id
 * @property {string} service Subscriber service status
 */
type SubscriberRowType = {
  name: string;
  id: string;
  service: string;
};

/**
 * Returns a table of subscribers serviced by the federation gateway.
 *
 * @param {FEGGatewayDetailType} props
 */
export default function GatewayDetailSubscribers(props: FEGGatewayDetailType) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const [subscriberRows, setSubscriberRows] = useState<
    Array<SubscriberRowType>
  >([]);
  const ctx = useContext(FEGSubscriberContext);
  const [isLoading, setIsLoading] = useState(true);
  // Auto refresh every REFRESH_INTERVAL seconds
  useInterval(() => ctx.refetch(), props.refresh ? REFRESH_INTERVAL : null);

  const sessionState = ctx?.sessionState || {};
  const subscriberToNetworkIdMap: Record<string, string> = {};

  Object.keys(sessionState).map(servicedNetworkId => {
    const servicedNetworkSessionState = sessionState[servicedNetworkId] ?? {};
    Object?.keys(servicedNetworkSessionState).map(subscriberImsi => {
      subscriberToNetworkIdMap[subscriberImsi] = servicedNetworkId;
    });
  });
  // get all the subscribers IMSI number serviced by the federation network
  const subscribersImsi = JSON.stringify(Object.keys(subscriberToNetworkIdMap));

  useEffect(() => {
    const fetchSubscribersInfo = async () => {
      const newSubscriberRows: Array<SubscriberRowType> = [];
      //TODO: - @andreilee bulk fetch from a paginated api endpoint
      await Promise.all(
        Object.keys(subscriberToNetworkIdMap).map(async subscriberImsi => {
          const subscriberInfo = (await fetchSubscribers({
            networkId: subscriberToNetworkIdMap[subscriberImsi],
            id: subscriberImsi,
          })) as Subscriber;
          newSubscriberRows.push({
            name: subscriberInfo?.name || subscriberImsi,
            id: subscriberImsi,
            service: subscriberInfo?.lte?.state || '-',
          });
        }),
      );
      setSubscriberRows(newSubscriberRows);
      setIsLoading(false);
    };
    void fetchSubscribersInfo();
    // rerun only when a new subscriber session has been added
  }, [subscribersImsi]); // eslint-disable-line react-hooks/exhaustive-deps

  if (isLoading) {
    return <LoadingFiller />;
  }

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
