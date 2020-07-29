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
import type {
  apn,
  enodeb,
  lte_gateway,
  rule_id,
  subscriber,
} from '@fbcnms/magma-api';

import CellWifiIcon from '@material-ui/icons/CellWifi';
import KPITray from '../../components/KPITray';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';

type Props = {
  lteGatwayResp: ?{[string]: lte_gateway},
  enb: ?{[string]: enodeb},
  subscriber: ?{[string]: subscriber},
  policyRules: ?Array<rule_id>,
  apns: ?{[string]: apn},
};

export default function NetworkKPI(props: Props) {
  return (
    <KPITray
      data={[
        {
          icon: CellWifiIcon,
          category: 'Gateways',
          value: props.lteGatwayResp
            ? Object.keys(props.lteGatwayResp).length
            : 0,
        },
        {
          icon: SettingsInputAntennaIcon,
          category: 'eNodeBs',
          value: props.enb ? Object.keys(props.enb).length : 0,
        },
        {
          icon: PeopleIcon,
          category: 'Subscribers',
          value: props.subscriber ? Object.keys(props.subscriber).length : 0,
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
      ]}
    />
  );
}
