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
import ApnOverview from './ApnOverview';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import PolicyOverview from './PolicyOverview';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import TopBar from '../../components/TopBar';

import {ApnJsonConfig} from './ApnOverview';
import {PolicyJsonConfig} from './PolicyOverview';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

export default function TrafficDashboard() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Traffic"
        tabs={[
          {
            label: 'Policies',
            to: '/policy',
            icon: LibraryBooksIcon,
          },
          {
            label: 'APNs',
            to: '/apn',
            icon: RssFeedIcon,
          },
        ]}
      />

      <Switch>
        <Route
          path={relativePath('/policy/:policyId/json')}
          component={PolicyJsonConfig}
        />
        <Route
          path={relativePath('/policy/json')}
          component={PolicyJsonConfig}
        />
        <Route
          path={relativePath('/apn/:apnId/json')}
          component={ApnJsonConfig}
        />
        <Route path={relativePath('/apn/json')} component={ApnJsonConfig} />
        <Route path={relativePath('/policy')} component={PolicyOverview} />
        <Route path={relativePath('/apn')} component={ApnOverview} />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}
