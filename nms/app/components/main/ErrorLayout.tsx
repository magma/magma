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
import AppContent from '../layout/AppContent';
import AppSideBar from '../AppSideBar';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    display: 'flex',
  },
});

export default function ErrorLayout({children}: {children: React.ReactNode}) {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <AppSideBar items={[]} />
      <AppContent>{children}</AppContent>
    </div>
  );
}
