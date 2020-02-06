/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PermissionHandlingProps} from '@fbcnms/ui/components/design-system/Form/FormAction';

import * as React from 'react';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import PopoverMenu from '@fbcnms/ui/components/design-system/Select/PopoverMenu';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo} from 'react';

export type MenuOption = {
  onClick: () => void,
  caption: React.Node,
  ...PermissionHandlingProps,
};

type Props = {
  options: Array<MenuOption>,
  menuIcon?: React.Node,
  className?: ?string,
};

const useStyles = makeStyles({
  menu: {
    width: 'auto',
  },
  menuButton: {
    minWidth: 'unset',
    paddingLeft: 0,
    paddingRight: 0,
  },
  icon: {
    color: symphony.palette.D400,
  },
  disabled: {
    opacity: 0.5,
    cursor: 'default',
  },
});

const TableRowOptionsButton = (props: Props) => {
  const {options, menuIcon, className} = props;
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
        label: opt.caption,
        value: optIndex,
        ignorePermissions: opt.ignorePermissions,
        hideWhenDisabled: opt.hideWhenDisabled,
      }))}
      onChange={handleOptionClick}
      menuClassName={classes.menu}
      className={classes.menuButton}>
      {menuIcon ?? (
        <MoreVertIcon
          className={classNames(className, {
            [classes.icon]: isEnabled,
          })}
        />
      )}
    </PopoverMenu>
  );
};

export default TableRowOptionsButton;
