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

import typeof SvgIcon from '@material-ui/core/@@SvgIcon';
import type {ButtonProps} from '../Button';
import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import SelectMenu from './SelectMenu';
import classNames from 'classnames';
import emptyFunction from '../../../../../fbc_js_core/util/emptyFunction';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  menu: {
    margin: '8px 0px',
  },
  menuDockRight: {
    position: 'absolute',
    right: '0',
  },
}));

export type PopoverMenuProps<TValue> = $ReadOnly<{|
  className?: string,
  menuClassName?: string,
  menuDockRight?: boolean,
  children: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange?: (value: TValue) => void | (() => void),
  leftIcon?: SvgIcon,
  rightIcon?: SvgIcon,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  onVisibilityChange?: (isVisible: boolean) => void,
  ...ButtonProps,
|}>;

const PopoverMenu = <TValue>(props: PopoverMenuProps<TValue>) => {
  const {
    className,
    menuClassName,
    children,
    leftIcon,
    rightIcon,
    menuDockRight,
    onChange,
    variant,
    skin,
    disabled,
    onVisibilityChange,
    tooltip,
    useEllipsis,
    ...selectMenuProps
  } = props;
  const classes = useStyles();
  return (
    <BasePopoverTrigger
      onVisibilityChange={onVisibilityChange}
      popover={
        <SelectMenu
          {...selectMenuProps}
          onChange={onChange || emptyFunction}
          size="normal"
          className={classNames(classes.menu, menuClassName, {
            [classes.menuDockRight]: menuDockRight,
          })}
        />
      }>
      {(onShow, _onHide, contextRef) => (
        <Button
          onClick={onShow}
          ref={contextRef}
          variant={variant}
          skin={skin || 'regular'}
          disabled={disabled}
          className={className}
          leftIcon={leftIcon}
          rightIcon={rightIcon}
          tooltip={tooltip}
          useEllipsis={useEllipsis}>
          {children}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default PopoverMenu;
