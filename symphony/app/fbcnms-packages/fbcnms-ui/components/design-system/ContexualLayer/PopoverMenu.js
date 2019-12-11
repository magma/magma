/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import BasePopoverTrigger from './BasePopoverTrigger';
import Button from '../Button';
import SelectMenu from './SelectMenu';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  menu: {
    margin: '8px 0px',
  },
});

type Props<TValue> = {
  className?: string,
  children: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange: (value: TValue) => void | (() => void),
};

const PopoverMenu = <TValue>({children, ...selectMenuProps}: Props<TValue>) => {
  const classes = useStyles();
  return (
    <BasePopoverTrigger
      popover={
        <SelectMenu
          {...selectMenuProps}
          size="normal"
          className={classes.menu}
        />
      }>
      {(onShow, contextRef) => (
        <Button onClick={onShow} ref={contextRef} variant="text">
          {children}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default PopoverMenu;
