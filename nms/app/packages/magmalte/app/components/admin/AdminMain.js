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

import nullthrows from '@fbcnms/util/nullthrows';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

type Props = {
  navItems: () => React.Node,
  navRoutes: () => React.Node,
};

export default function AdminMain(props: Props) {
  const classes = useStyles();
  const {tabs, user, ssoEnabled} = useContext(AppContext);

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={props.navItems()}
        projects={getProjectLinks(tabs, user)}
        user={nullthrows(user)}
        showSettings={!ssoEnabled}
      />
      <AppContent>{props.navRoutes()}</AppContent>
    </div>
  );
}
