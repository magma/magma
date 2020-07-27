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
 * @flow
 * @format
 */
import * as React from 'react';
import ToggleButton from '@material-ui/lab/ToggleButton';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  button: {
    background: theme.palette.background.default,
    color: 'black',
    borderRight: '1px solid #ddd',
    borderRadius: '4px',
    height: '30px',
    border: 0,
  },
  selected: {
    color: theme.palette.blue60,
  },
  notSelected: {
    color: theme.palette.black,
  },
}));

type Props = {
  onClick: () => void,
  icon: React.Node,
  isSelected?: boolean,
};

const MapButton = (props: Props) => {
  const {onClick, isSelected, icon} = props;
  const classes = useStyles();
  return (
    <ToggleButton value={1} className={classes.button} onClick={onClick}>
      <span className={isSelected ? classes.selected : classes.notSelected}>
        {icon}
      </span>
    </ToggleButton>
  );
};

export default MapButton;
