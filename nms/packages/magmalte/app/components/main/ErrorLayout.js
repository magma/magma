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
import AppContent from '../layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';

import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '../Settings';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

export default function ErrorLayout({children}: {children: React.Node}) {
  const classes = useStyles();
  const {user, tabs, ssoEnabled} = React.useContext(AppContext);

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={[]}
        secondaryItems={[]}
        projects={getProjectLinks(tabs, user)}
        showSettings={shouldShowSettings({
          isSuperUser: user.isSuperUser,
          ssoEnabled,
        })}
        user={user}
      />
      <AppContent>{children}</AppContent>
    </div>
  );
}
