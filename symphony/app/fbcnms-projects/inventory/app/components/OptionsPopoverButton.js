/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';

import * as React from 'react';
import PopoverMenu from '@fbcnms/ui/components/design-system/Select/PopoverMenu';
import classNames from 'classnames';
import {ThreeDotsVerticalIcon} from '@fbcnms/ui/components/design-system/Icons';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo} from 'react';

export type MenuOption = {|
  onClick: () => void,
  caption: React.Node,
  ...PermissionHandlingProps,
|};

type Props = {|
  options: Array<MenuOption>,
  menuIcon?: React.Node,
  className?: ?string,
  popoverMenuClassName?: ?string,
  onVisibilityChange?: (isVisible: boolean) => void,
|};

const useStyles = makeStyles(() => ({
  menu: {
    width: 'auto',
  },
  menuButton: {
    minWidth: 'unset',
    paddingLeft: 0,
    paddingRight: 0,
  },
  icon: {
    padding: '4px',
    backgroundColor: 'white',
    borderRadius: '100%',
    cursor: 'pointer',
  },
  disabled: {
    opacity: 0.5,
    cursor: 'default',
  },
}));

const OptionsPopoverButton = (props: Props) => {
  const {
    options,
    menuIcon,
    className,
    popoverMenuClassName,
    onVisibilityChange,
  } = props;
  const classes = useStyles();
  const isEnabled = useMemo(() => props.options.length > 0, [
    props.options.length,
  ]);
  const handleOptionClick = useCallback(
    optIndex => {
      options[optIndex].onClick();
    },
    [options],
  );
  return (
    <PopoverMenu
      disabled={!isEnabled}
      variant="text"
      menuDockRight={true}
      options={options.map((opt, optIndex) => ({
        key: optIndex + '',
        label: opt.caption,
        value: optIndex,
        ignorePermissions: opt.ignorePermissions,
        hideWhenDisabled: opt.hideWhenDisabled,
      }))}
      onChange={handleOptionClick}
      menuClassName={classes.menu}
      className={classNames(classes.menuButton, popoverMenuClassName)}
      onVisibilityChange={onVisibilityChange}>
      {menuIcon ?? (
        <ThreeDotsVerticalIcon
          className={classNames(className, {
            [classes.icon]: isEnabled,
          })}
        />
      )}
    </PopoverMenu>
  );
};

export default OptionsPopoverButton;
