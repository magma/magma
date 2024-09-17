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
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import Button from '@mui/material/Button';
import Menu, {MenuProps} from '@mui/material/Menu';
import React from 'react';
import withStyles from '@mui/styles/withStyles';
import {colors, shadows} from '../theme/default';

const StyledMenu = withStyles({
  paper: {
    border: `1px solid ${colors.button.lightOutline}`,
    boxShadow: shadows.DP3,
  },
})((props: MenuProps) => (
  <Menu
    elevation={0}
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

type Props = {
  label: string;
  children: React.ReactNode;
  'data-testid'?: string;
  size?: 'small' | 'medium' | 'large';
  className?: string;
};

export default function MenuButton(props: Props) {
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement>();
  const [minWidth, setMinWidth] = React.useState<number>();

  const onButtonClick = (event: React.SyntheticEvent<HTMLButtonElement>) => {
    setMinWidth(event.currentTarget.getBoundingClientRect().width);
    setAnchorEl(event.currentTarget);
  };
  const onClose = () => {
    setAnchorEl(undefined);
    setMinWidth(undefined);
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
