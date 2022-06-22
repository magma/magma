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
import AccountSettings from './AccountSettings';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Admin from './admin/Admin';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AppContent from './layout/AppContent';
import AppSideBar from './AppSideBar';
import {Route, Routes} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

export default function IndexWithoutNetwork() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <AppSideBar items={[]} />
      <AppContent>
        <Routes>
          <Route path="/admin/*" element={<Admin />} />
          <Route path="/settings/*" element={<AccountSettings />} />
        </Routes>
      </AppContent>
    </div>
  );
}
