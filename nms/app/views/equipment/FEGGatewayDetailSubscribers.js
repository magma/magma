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

import type {
  federation_gateway,
  subscriber,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
import FEGSubscriberContext from '../../components/context/FEGSubscriberContext';
import Link from '@material-ui/core/Link';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import React from 'react';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {FetchSubscribers} from '../../state/lte/SubscriberState';
import {
  REFRESH_INTERVAL,
  RefreshTypeEnum,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
import {useEffect, useState} from 'react';
import {useNavigate, useParams, useResolvedPath} from 'react-router-dom';

/**
 * @property {federation_gateway} gwInfo The Federation gateway being looked at
 * @property {boolean} refresh Boolean telling to autorefresh or not
 */
type FEGGatewayDetailType = {
  gwInfo: federation_gateway,
  refresh: boolean,
};

/**
 * @property {string} name Subscriber name
 * @property {string} id Subscriber id
 * @property {string} service Subscriber service status
 */
type SubscriberRowType = {
  name: string,
  id: string,
  service: string,
};

/**
 * Returns a table of subscribers serviced by the federation gateway.
 *
 * @param {FEGGatewayDetailType} props
 */
export default function GatewayDetailSubscribers(props: FEGGatewayDetailType) {
  const resolvedPath = useResolvedPath('');
  const navigate = useNavigate();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const [subscriberRows, setSubscriberRows] = useState<
    Array<SubscriberRowType>,
  >([]);
  const [isLoading, setIsLoading] = useState(true);
  // Auto refresh every REFRESH_INTERVAL seconds
  const ctx = useRefreshingContext({
    context: FEGSubscriberContext,
    networkId: networkId,
    type: RefreshTypeEnum.FEG_SUBSCRIBER,
    interval: REFRESH_INTERVAL,
    refresh: props.refresh,
  });
  const sessionState = ctx?.sessionState || {};
  const subscriberToNetworkIdMap = {};

  Object.keys(sessionState).map(servicedNetworkId => {
    // $FlowIgnore[prop-missing] because refresh context returns other things too like state, enbInfo and each have their own property
    const servicedNetworkSessionState = sessionState[servicedNetworkId] ?? {};
    Object?.keys(servicedNetworkSessionState).map(subscriberImsi => {
      subscriberToNetworkIdMap[subscriberImsi] = servicedNetworkId;
    });
  });
  // get all the subscribers IMSI number serviced by the federation network
  const subscribersImsi = JSON.stringify(Object.keys(subscriberToNetworkIdMap));

  useEffect(() => {
    const fetchSubscribersInfo = async () => {
      const newSubscriberRows = [];
      //TODO: - @andreilee bulk fetch from a paginated api endpoint
      await Promise.all(
        Object.keys(subscriberToNetworkIdMap).map(async subscriberImsi => {
          // $FlowIgnore because it can be called with different values when getting paginated subscribers
          const subscriberInfo: subscriber = await FetchSubscribers({
            networkId: subscriberToNetworkIdMap[subscriberImsi],
            id: subscriberImsi,
          });
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
    fetchSubscribersInfo();
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
