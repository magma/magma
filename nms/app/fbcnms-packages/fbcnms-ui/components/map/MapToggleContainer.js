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
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  toggleContainer: {
    background: theme.palette.background.default,
    padding: 0,
    border: 0,
    borderStyle: 'solid',
    borderRadius: '4px',
    display: 'inline-block',
  },
}));

type Props = {children: ?React.Node};

const MapToggleContainer = (props: Props) => {
  const classes = useStyles();
  return <div className={classes.toggleContainer}>{props.children}</div>;
};

export default MapToggleContainer;
