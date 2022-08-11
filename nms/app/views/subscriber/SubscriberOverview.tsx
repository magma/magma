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
 */
import PeopleIcon from '@mui/icons-material/People';
import React from 'react';
import SettingsIcon from '@mui/icons-material/Settings';
import SubscriberDetail from './SubscriberDetail';
import SubscriberStateTable from './SubscriberStateTable';
import SubscriberTable from './SubscriberTable';
import TopBar from '../../components/TopBar';
import {Navigate, Route, Routes} from 'react-router-dom';

const TITLE = 'Subscribers';

export default function SubscriberDashboard() {
  return (
    <Routes>
      <Route
        path="/overview/config/:subscriberId/*"
        element={<SubscriberDetail />}
      />
      <Route
        path="/overview/sessions/:subscriberId/*"
        element={<SubscriberDetail />}
      />

      <Route path="/overview/*" element={<SubscribersOverview />} />
      <Route index element={<Navigate to="overview" replace />} />
    </Routes>
  );
}

export function SubscribersOverview() {
  return (
    <>
      <TopBar
        header={TITLE}
        tabs={[
          {
            label: 'Config',
            to: 'config',
            icon: SettingsIcon,
          },
          {
            label: 'Sessions',
            to: 'sessions',
            icon: PeopleIcon,
          },
        ]}
      />
      <Routes>
        <Route path="/config" element={<SubscriberTable />} />
        <Route path="/sessions" element={<SubscriberStateTable />} />
        <Route index element={<Navigate to="config" replace />} />
      </Routes>
    </>
  );
}
