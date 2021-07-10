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

import CellWifiIcon from '@material-ui/icons/CellWifi';
import FEGGateway from './FEGEquipmentGateway';
import React from 'react';
import TopBar from '../../components/TopBar';

import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

/**
 * Returns the full equipment dashboard of the federation network.
 * It consists of an internal equipment dashboard component to display
 * the useful information about the federation gateways.
 */
function FEGEquipmentDashboard() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <Switch>
        <Route
          path={relativePath('/overview')}
          component={EquipmentDashboardInternal}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}
/**
 * It consists of a top bar to navigate and a federation gateway component
 * to provide information about the federation gateways.
 */
function EquipmentDashboardInternal() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Equipment"
        tabs={[
          {
            label: 'Federation Gateways',
            to: '/gateway',
            icon: CellWifiIcon,
          },
        ]}
      />
      <Switch>
        <Route path={relativePath('/gateway')} component={FEGGateway} />
        <Redirect to={relativeUrl('/gateway')} />
      </Switch>
    </>
  );
}

export default FEGEquipmentDashboard;
