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

import * as React from 'react';
import AssignmentIcon from '@material-ui/icons/Assignment';
import AuditLog from './AuditLog';
import Networks from './Networks';
import PeopleIcon from '@material-ui/icons/People';
import SignalCellularAlt from '@material-ui/icons/SignalCellularAlt';
import TopBar from '../TopBar';
import UsersSettings from '../UsersSettings';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '../../../fbc_js_core/ui/hooks';

const TITLE = 'Administration';

export default function Admin() {
  const {relativeUrl, relativePath} = useRouter();
  console.log({relativeUrl: relativeUrl(''), relativePath: relativePath('')});

  const tabs = [
    {
      label: 'Users',
      to: '/users',
      icon: PeopleIcon,
    },
    {
      label: 'Audit Log',
      to: '/audit_log',
      icon: AssignmentIcon,
    },
    {
      label: 'Networks',
      to: '/networks',
      icon: SignalCellularAlt,
    },
  ];

  return (
    <>
      <TopBar header={TITLE} tabs={tabs} />
      <Switch>
        <Route path={relativeUrl('/users')} component={UsersSettings} />
        <Route path={relativeUrl('/audit_log')} component={AuditLog} />
        <Route path={relativeUrl('/networks')} component={Networks} />
        <Redirect to={relativeUrl('/users')} />
      </Switch>
    </>
  );
}
