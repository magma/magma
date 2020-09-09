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
import type {DataRows} from '../../components/DataGrid';
import type {apn, rule_id} from '@fbcnms/magma-api';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataGrid from '../../components/DataGrid';
import EnodebContext from '../../components/context/EnodebContext';
import GatewayContext from '../../components/context/GatewayContext';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import SubscriberContext from '../../components/context/SubscriberContext';

import {useContext} from 'react';

type Props = {
  policyRules: ?Array<rule_id>,
  apns: ?{[string]: apn},
};

export default function NetworkKPI(props: Props) {
  const gwCtx = useContext(GatewayContext);
  const enbCtx = useContext(EnodebContext);
  const subscriberCtx = useContext(SubscriberContext);

  const kpiData: DataRows[] = [
    [
      {
        icon: CellWifiIcon,
        category: 'Gateways',
        value: Object.keys(gwCtx.state).length,
      },
      {
        icon: SettingsInputAntennaIcon,
        category: 'eNodeBs',
        value: Object.keys(enbCtx.state.enbInfo).length,
      },
      {
        icon: PeopleIcon,
        category: 'Subscribers',
        value: Object.keys(subscriberCtx.state).length,
      },
      {
        icon: LibraryBooksIcon,
        category: 'Policies',
        value: props.policyRules ? props.policyRules.length : 0,
      },
      {
        icon: RssFeedIcon,
        category: 'APNs',
        value: props.apns ? Object.keys(props.apns).length : 0,
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}
