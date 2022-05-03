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
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import Button from '@material-ui/core/Button';
import Menu from '@material-ui/core/Menu';
import React from 'react';
import {colors, shadows} from '../theme/default';
import {withStyles} from '@material-ui/core/styles';

import type {Node} from 'react';

const StyledMenu = withStyles({
  paper: {
    border: `1px solid ${colors.button.lightOutline}`,
    boxShadow: shadows.DP3,
  },
})(props => (
  <Menu
    elevation={0}
    getContentAnchorEl={null}
    anchorOrigin={{
      vertical: 'bottom',
      horizontal: 'center',
    }}
    transformOrigin={{
      vertical: 'top',
      horizontal: 'center',
    }}
    {...props}
  />
));

type Props = {|
  label: string,
  children: Node,
  'data-testid'?: string,
  size?: 'small' | 'medium' | 'large',
  className?: string,
|};

export default function MenuButton(props: Props) {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [minWidth, setMinWidth] = React.useState(null);

  const onButtonClick = event => {
    setMinWidth(event.currentTarget.getBoundingClientRect().width);
    setAnchorEl(event.currentTarget);
  };
  const onClose = () => {
    setAnchorEl(null);
    setMinWidth(null);
  };

  const {children, label, ...passthroughProps} = props;
  return (
    <div>
      <Button
        variant="contained"
        color="primary"
        onClick={onButtonClick}
        endIcon={<ArrowDropDownIcon />}
        {...passthroughProps}>
        {label}
      </Button>
      <StyledMenu
        PaperProps={{style: {minWidth}}}
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={onClose}>
        {children}
      </StyledMenu>
    </div>
  );
}
