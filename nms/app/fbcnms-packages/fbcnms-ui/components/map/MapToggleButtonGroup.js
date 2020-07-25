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
import ToggleButtonGroup from '@material-ui/core/Typography';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  toggleGroup: {
    boxShadow: '0px 0px 0.5px 0.5px grey',
    borderRadius: '4px',
  },
}));

type Props = {children: ?React.Node};

const MapToggleButtonGroup = (props: Props) => {
  const classes = useStyles();

  return (
    <ToggleButtonGroup className={classes.toggleGroup}>
      {props.children}
    </ToggleButtonGroup>
  );
};

export default MapToggleButtonGroup;
