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
import {gray7} from '@fbcnms/ui/theme/colors';
import {makeStyles} from '@material-ui/styles';

type Props = {
  children: React.Node,
};

const useStyles = makeStyles(() => ({
  content: {
    flexGrow: 1,
    height: '100vh',
    overflow: 'auto',
    overflowX: 'hidden',
    backgroundColor: gray7,
  },
}));

export default function AppContent(props: Props) {
  const classes = useStyles();
  return <main className={classes.content}>{props.children}</main>;
}
