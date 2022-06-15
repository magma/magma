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

import * as React from 'react';
import AssignmentIcon from '@material-ui/icons/Assignment';
import AuditLog from './AuditLog';
import Networks from './Networks';
import PeopleIcon from '@material-ui/icons/People';
import SignalCellularAlt from '@material-ui/icons/SignalCellularAlt';
import TopBar from '../TopBar';
import UsersSettings from '../UsersSettings';
import {Navigate, Route, Routes} from 'react-router-dom';

const TITLE = 'Administration';

export default function Admin() {
  const tabs = [
    {
      label: 'Users',
      to: 'users',
      icon: PeopleIcon,
    },
    {
      label: 'Audit Log',
      to: 'audit_log',
      icon: AssignmentIcon,
    },
    {
      label: 'Networks',
      to: 'networks',
      icon: SignalCellularAlt,
    },
  ];

  return (
    <>
      <TopBar header={TITLE} tabs={tabs} />
      <Routes>
        <Route path="/users" element={<UsersSettings />} />
        <Route path="/audit_log" element={<AuditLog />} />
        <Route path="/networks/*" element={<Networks />} />
        <Route index element={<Navigate to="users" replace />} />
      </Routes>
    </>
  );
}
